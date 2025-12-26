package utils

import (
	"github.com/microcosm-cc/bluemonday"
	"sync"
)

var (
	// Dùng singleton để không phải khởi tạo policy nhiều lần gây tốn tài nguyên
	strictPolicy *bluemonday.Policy
	ugcPolicy    *bluemonday.Policy
	once         sync.Once
)

func initPolicies() {
	// StrictPolicy: Loại bỏ TOÀN BỘ thẻ HTML. Dùng cho Tên, Tiêu đề, Mã đơn hàng.
	// Ví dụ: "Hello <b>World</b>" -> "Hello World"
	strictPolicy = bluemonday.StrictPolicy()

	// UGCPolicy (User Generated Content): Cho phép các thẻ HTML an toàn (b, i, u, p...).
	// Loại bỏ script, iframe, object. Dùng cho bài viết blog, mô tả sản phẩm dài.
	ugcPolicy = bluemonday.UGCPolicy()
}

// SanitizeText: Dành cho các trường văn bản thuần túy (Tên, Chat, Note)
func SanitizeText(input string) string {
	once.Do(initPolicies)
	return strictPolicy.Sanitize(input)
}

// SanitizeHTML: Dành cho các trường cho phép định dạng (Mô tả sản phẩm rich text)
func SanitizeHTML(input string) string {
	once.Do(initPolicies)
	return ugcPolicy.Sanitize(input)
}
