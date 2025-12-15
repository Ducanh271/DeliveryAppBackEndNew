package service

type EmailService interface {
	SendOTP(toEmail string, otp string) error
	SendPasswordResetEmail(toEmail string, resetLink string) error
}
