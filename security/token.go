package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	// Tạo số ngẫu nhiên từ 0 đến 999999 (6 chữ số)
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %v", err)
	}
	// Định dạng số thành chuỗi 6 chữ số, thêm số 0 nếu cần
	return fmt.Sprintf("%06d", n), nil
}
