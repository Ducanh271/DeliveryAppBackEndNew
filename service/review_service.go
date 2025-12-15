package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"example.com/delivery-app/dto"
	"example.com/delivery-app/infrastructure/storage"
	"example.com/delivery-app/models"
	"example.com/delivery-app/repository"
)

type ReviewService struct {
	reviewRepo repository.ReviewRepository
	imageSvc   storage.ImageStorageService
	uow        repository.UnitOfWork
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	imageSvc storage.ImageStorageService,
	uow repository.UnitOfWork,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		imageSvc:   imageSvc,
		uow:        uow,
	}
}

// === 1. CREATE REVIEW ===
func (s *ReviewService) CreateReview(ctx context.Context, userID int64, req *dto.CreateReviewRequest, files []*multipart.FileHeader) (*models.Review, error) {
	// 1. Check quyền review (đã mua và nhận hàng chưa)
	canReview, err := s.reviewRepo.CanUserReview(userID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to check review permission: %w", err)
	}
	if !canReview {
		return nil, ErrUserCannotReview
	}

	// 2. Check đã review chưa (tránh spam)
	exists, err := s.reviewRepo.CheckReviewExists(userID, req.OrderID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if exists {
		return nil, ErrReviewExists
	}

	// 3. Upload ảnh (nếu có) lên Cloudinary TRƯỚC
	type UploadedImage struct {
		URL      string
		PublicID string
	}
	var uploadedImages []UploadedImage

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}

		// Folder "reviews" trên Cloudinary
		// Tạo tên file unique hơn để tránh trùng lặp nếu cần
		url, publicID, err := s.imageSvc.UploadProductImage(ctx, file, fmt.Sprintf("review_%d_%s", userID, fileHeader.Filename))
		file.Close() // Đóng file ngay sau khi upload xong
		if err != nil {
			// (Lý tưởng: nên rollback các ảnh đã upload nếu có lỗi xảy ra ở ảnh sau)
			return nil, fmt.Errorf("failed to upload review image: %w", err)
		}
		uploadedImages = append(uploadedImages, UploadedImage{URL: url, PublicID: publicID})
	}

	// 4. Lưu vào DB (Transaction UoW)
	review := &models.Review{
		ProductID: req.ProductID,
		UserID:    userID,
		OrderID:   req.OrderID,
		Rate:      req.Rate,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	err = s.uow.Execute(func(factory func(any) any) error {
		// Lấy repo từ factory (ép kiểu cẩn thận)
		repo := factory((*repository.ReviewRepository)(nil)).(repository.ReviewRepository)

		// 4a. Tạo Review
		reviewID, err := repo.Create(review)
		if err != nil {
			return err
		}
		review.ID = reviewID

		// 4b. Tạo và Link ảnh
		for _, img := range uploadedImages {
			// Insert vào bảng Images
			imageID, err := repo.CreateImage(img.URL, img.PublicID)
			if err != nil {
				return err
			}
			// Link vào bảng ReviewImages
			if err := repo.LinkImageToReview(reviewID, imageID); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	return review, nil
}

// === 2. GET REVIEWS BY PRODUCT ===
func (s *ReviewService) GetReviewsByProduct(productID int64, page, limit int) (*dto.ReviewListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// 1. Lấy danh sách review (chưa có ảnh)
	reviews, total, err := s.reviewRepo.GetByProductID(productID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	if len(reviews) == 0 {
		return &dto.ReviewListResponse{
			Reviews:    []dto.ReviewResponse{},
			Pagination: dto.Pagination{Total: total, Page: page, Limit: limit, TotalPages: 0},
		}, nil
	}

	// 2. Lấy danh sách ID review để query ảnh (Batch Query)
	var reviewIDs []int64
	for _, r := range reviews {
		reviewIDs = append(reviewIDs, r.ID)
	}

	// 3. Query ảnh 1 lần duy nhất cho tất cả review
	imagesMap, err := s.reviewRepo.GetImagesByReviewIDs(reviewIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get review images: %w", err)
	}

	// 4. Mapping (Lắp ráp Review + Images)
	var responseList []dto.ReviewResponse
	for _, r := range reviews {
		// Lấy ảnh từ map
		var imgRes []dto.ReviewImageResponse
		if imgs, ok := imagesMap[r.ID]; ok {
			for _, img := range imgs {
				imgRes = append(imgRes, dto.ReviewImageResponse{
					ID:  img.ID,
					URL: img.URL,
				})
			}
		}
		// Nếu không có ảnh, khởi tạo mảng rỗng để JSON trả về [] thay vì null
		if imgRes == nil {
			imgRes = []dto.ReviewImageResponse{}
		}

		responseList = append(responseList, dto.ReviewResponse{
			ID:        r.ID,
			UserID:    r.UserID,
			UserName:  r.UserName, // Lấy từ join bảng Users
			Rate:      r.Rate,
			Content:   r.Content,
			CreatedAt: r.CreatedAt,
			Images:    imgRes,
		})
	}

	totalPages := (total + limit - 1) / limit
	return &dto.ReviewListResponse{
		Reviews: responseList,
		Pagination: dto.Pagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
}
