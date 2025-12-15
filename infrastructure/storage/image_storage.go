package storage

import (
	"context"
	"mime/multipart"
)

// ImageStorageService là hợp đồng cho bất kỳ dịch vụ lưu trữ hình ảnh nào
type ImageStorageService interface {
	// UploadProductImage nhận file và trả về URL và PublicID
	UploadProductImage(ctx context.Context, file multipart.File, fileName string) (string, string, error)

	// DeleteImage xóa ảnh dựa trên PublicID
	DeleteImage(ctx context.Context, publicID string) error
}
