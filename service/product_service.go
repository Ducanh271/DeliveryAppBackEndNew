package service

import (
	"context"
	"database/sql"
	"errors"
	"example.com/delivery-app/dto"
	"example.com/delivery-app/infrastructure/storage"
	"example.com/delivery-app/models"
	"example.com/delivery-app/repository"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"mime/multipart"
	"time"
)

type ProductService struct {
	productRepo repository.ProductRepository
	imageSvc    storage.ImageStorageService
	uow         repository.UnitOfWork
}

func NewProductService(
	productRepo repository.ProductRepository,
	imageSvc storage.ImageStorageService,
	uow repository.UnitOfWork,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		imageSvc:    imageSvc,
		uow:         uow,
	}
}

// === LOGIC TỪ CreateNewProductHandler ===
func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest, files []*multipart.FileHeader) (*models.Product, []string, error) {

	if len(files) == 0 {
		return nil, nil, ErrNoImages
	}

	// BƯỚC 1: Upload ảnh lên Cloudinary TRƯỚC (An toàn hơn)
	// (Nếu upload fail, không cần rollback CSDL)
	type UploadedImage struct {
		URL      string
		PublicID string
	}
	var uploadedImages []UploadedImage
	var urls []string // Mảng URL để trả về

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, nil, fmt.Errorf("lỗi mở file: %w", err)
		}

		url, publicID, err := s.imageSvc.UploadProductImage(ctx, file, fileHeader.Filename)
		file.Close() // Đóng file ngay sau khi upload
		if err != nil {
			// (Lý tưởng: nên gọi DeleteImage cho các ảnh đã lỡ upload)
			return nil, nil, fmt.Errorf("lỗi upload ảnh: %w", err)
		}

		uploadedImages = append(uploadedImages, UploadedImage{URL: url, PublicID: publicID})
		urls = append(urls, url)
	}

	// BƯỚC 2: Lưu vào CSDL trong một Transaction (UoW)
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		QtyInitial:  req.QtyInitial,
		QtySold:     req.QtySold,
		CreatedAt:   time.Now(),
	}

	err := s.uow.Execute(func(repoProvider func(repoType any) any) error {
		// (Giả sử UoW của bạn có thể cung cấp ProductRepository)
		prodRepo := repoProvider((*repository.ProductRepository)(nil)).(repository.ProductRepository)
		// 2a. Tạo Product
		productID, err := prodRepo.Create(product)
		if err != nil {
			return err
		}
		product.ID = productID

		// 2b. Lưu và link ảnh
		isFirst := true
		for _, img := range uploadedImages {
			// 1. Tạo ảnh trong bảng `Images`
			imageID, err := prodRepo.CreateImage(img.URL, img.PublicID)
			if err != nil {
				return err
			}

			// 2. Link ảnh vào bảng `ProductImages`
			if err := prodRepo.LinkImageToProduct(productID, imageID, isFirst); err != nil {
				return err
			}
			isFirst = false
		}
		return nil
	})

	if err != nil {
		// (Lý tưởng: có 1 job để xóa ảnh mồ côi trên Cloudinary)
		return nil, nil, fmt.Errorf("lỗi khi lưu vào CSDL: %w", err)
	}

	return product, urls, nil
}

// === LOGIC TỪ DeleteProductHandler ===
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {

	// BƯỚC 1: Lấy thông tin ảnh (cả Image IDs và Public IDs) TRƯỚC
	images, err := s.productRepo.GetImagesByProductID(id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("lỗi khi lấy thông tin ảnh: %w", err)
	}

	// BƯỚC 2: Xóa khỏi CSDL (trong Transaction)
	err = s.uow.Execute(func(repoProvider func(any) any) error {
		prodRepo := repoProvider((*repository.ProductRepository)(nil)).(repository.ProductRepository)

		var imageIDs []int64
		for _, img := range images {
			imageIDs = append(imageIDs, img.ID)
		}

		// 2a. Xóa Product (Bảng Products)
		// 	   -> Tự động CASCADE xóa các link trong `ProductImages`
		if err := prodRepo.DeleteProduct(id); err != nil {
			return err // Trả về lỗi gốc (vd: sql.ErrNoRows)
		}

		// 2b. Xóa Images (Bảng Images)
		// 	   -> Tự động CASCADE xóa các link (còn lại, nếu có) trong `ProductImages`
		if len(imageIDs) > 0 {
			if err := prodRepo.DeleteImages(imageIDs); err != nil {
				return fmt.Errorf("lỗi xóa images: %w", err)
			}
		}
		return nil // Commit
	})

	if err != nil {
		return err // Trả về lỗi CSDL (handler sẽ check ErrNoRows)
	}

	// BƯỚC 3: Xóa ảnh khỏi Cloudinary (Sau khi DB đã thành công)
	go func() {
		log.Printf("Bắt đầu dọn dẹp %d ảnh cho product %d", len(images), id)
		for _, img := range images {
			if img.PublicID != "" {
				if err := s.imageSvc.DeleteImage(context.Background(), img.PublicID); err != nil {
					log.Printf("Lỗi xóa ảnh mồ côi (orphan) %s: %v", img.PublicID, err)
				}
			}
		}
	}()

	return nil
}

// === LOGIC TỪ GetProductByIDHandler (TỐI ƯU HÓA) ===
func (s *ProductService) GetProductDetail(id int64) (*dto.ProductDetailResponse, error) {
	var product *models.Product
	var images []models.ProductImage
	var rating struct {
		Avg   float64
		Count int64
	}
	// Dùng errgroup để chạy 3 query song song
	g, _ := errgroup.WithContext(context.Background())

	// 1. Lấy Product
	g.Go(func() error {
		var err error
		product, err = s.productRepo.GetByID(id)
		if err != nil {
			return err
		} // Trả về lỗi (vd: sql.ErrNoRows)
		return nil
	})

	// 2. Lấy Ảnh
	g.Go(func() error {
		var err error
		images, err = s.productRepo.GetImagesByProductID(id)
		return err // (Không sao nếu lỗi, vẫn trả về product)
	})

	// 3. Lấy Rating
	g.Go(func() error {
		var err error
		rating.Avg, rating.Count, err = s.productRepo.GetAverageRating(id)
		return err // (Không sao nếu lỗi)
	})

	// Đợi tất cả hoàn thành
	if err := g.Wait(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err // Lỗi CSDL
	}

	// 4. Mapping
	var imageResponses []dto.ProductImageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, dto.ProductImageResponse{
			ID:     img.ID,
			URL:    img.URL,
			IsMain: img.IsMain,
		})
	}

	return &dto.ProductDetailResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		QtyInitial:  product.QtyInitial,
		QtySold:     product.QtySold,
		CreatedAt:   product.CreatedAt,
		Images:      imageResponses,
		AvgRate:     rating.Avg,
		ReviewCount: int(rating.Count), // Chuyển int64 sang int
	}, nil
}

// === LOGIC TỪ GetProductsHandler/SearchHandler (GIẢI QUYẾT N+1) ===
// (Hàm này dùng cho cả GetAll và Search)
func (s *ProductService) GetProducts(keyword string, page, limit int64) (*dto.ProductListResponse, error) {

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	var products []models.Product
	var total int64
	var err error

	// BƯỚC 1: Repo chỉ tìm Product (2 query)
	if keyword == "" {
		products, total, err = s.productRepo.GetAll(page, limit)
	} else {
		products, total, err = s.productRepo.Search(keyword, page, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy product list: %w", err)
	}

	// Nếu không có sản phẩm, trả về rỗng
	if len(products) == 0 {
		return &dto.ProductListResponse{
			Products:   []dto.ProductResponse{},
			Pagination: dto.Pagination{Total: 0, Page: int(page), Limit: int(limit), TotalPages: 0},
		}, nil
	}

	// Lấy list IDs
	productIDs := make([]int64, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
	}

	// BƯỚC 2: Repo lấy TẤT CẢ ảnh (1 query)
	imagesMap, err := s.productRepo.GetImagesByProductIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy images: %w", err)
	}

	// BƯỚC 3: Repo lấy TẤT CẢ rating (1 query)
	ratingsMap, err := s.productRepo.GetAverageRatingsByProductIDs(productIDs)
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy ratings: %w", err)
	}

	// BƯỚC 4: Service "lắp ráp" (mapping)
	var productResponses []dto.ProductResponse
	for _, p := range products {

		// Lấy ảnh từ map
		var modelImages []models.ProductImage
		if imgs, ok := imagesMap[p.ID]; ok {
			modelImages = imgs
		}

		// Lấy rating từ map
		rating := struct {
			Avg   float64
			Count int64
		}{Avg: 0, Count: 0}
		if r, ok := ratingsMap[p.ID]; ok {
			rating = r
		}

		productResponses = append(productResponses, dto.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			QtyInitial:  p.QtyInitial,
			QtySold:     p.QtySold,
			CreatedAt:   p.CreatedAt,
			Images:      modelImages,
			AvgRate:     rating.Avg,
			ReviewCount: rating.Count,
		})
	}

	// Trả về DTO cuối cùng
	totalPages := (total + limit - 1) / limit
	return &dto.ProductListResponse{
		Products: productResponses,
		Pagination: dto.Pagination{
			Total:      int(total),
			Page:       int(page),
			Limit:      int(limit),
			TotalPages: int(totalPages),
		},
	}, nil
}

// === LOGIC TỪ GetNumberOfProductHandler ===
func (s *ProductService) GetTotalProducts() (int64, error) {
	return s.productRepo.GetTotalCount()
} // === LOGIC TỪ CreateNewProductHandler ===
