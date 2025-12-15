package storage

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// cloudinaryService là implementation của ImageStorageService
type cloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService là hàm khởi tạo, nhận client Cloudinary
func NewCloudinaryService(cld *cloudinary.Cloudinary) ImageStorageService {
	return &cloudinaryService{cld: cld}
}

// Triển khai hàm Upload
func (s *cloudinaryService) UploadProductImage(ctx context.Context, file multipart.File, fileName string) (string, string, error) {
	uploadResult, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   "product", // Bạn có thể lấy folder từ config
	})
	if err != nil {
		return "", "", err
	}
	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

// Triển khai hàm Delete
func (s *cloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
