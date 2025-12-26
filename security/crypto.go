package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt mã hóa plainText bằng key (AES-GCM)
func Encrypt(plainText string, keyString string) (string, error) {
	// Key phải đủ 32 byte (cho AES-256)
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Sử dụng GCM (Galois/Counter Mode) để vừa mã hóa vừa xác thực
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Tạo nonce (số ngẫu nhiên duy nhất)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Mã hóa: Seal(dst, nonce, plaintext, additionalData)
	// Ta gộp nonce vào đầu chuỗi kết quả để dùng khi giải mã
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	// Trả về dạng Base64 để lưu vào DB (dạng text)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt giải mã cipherText (Base64) về plainText
func Decrypt(cipherTextBase64 string, keyString string) (string, error) {
	key := []byte(keyString)
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Tách nonce và nội dung mã hóa thực sự
	nonce, cipherBytes := cipherText[:nonceSize], cipherText[nonceSize:]

	// Giải mã
	plainText, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
