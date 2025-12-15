package notification

import (
	"example.com/delivery-app/config"
	"example.com/delivery-app/service"
	"net/smtp"
)

type smtpService struct {
	config config.EmailConfig
}

func NewEmailService(cfg config.EmailConfig) service.EmailService {
	return &smtpService{config: cfg}
}

func (s *smtpService) SendOTP(toEmail string, otp string) error {
	from := s.config.From
	password := s.config.Password
	smtpHost := s.config.Host
	smtpPort := s.config.Port

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject: Your OTP Code\n"
	body := "Your OTP code is: " + otp + "\n"
	message := []byte(subject + "\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return err
	}
	return nil
}
func (s *smtpService) SendPasswordResetEmail(toEmail string, resetLink string) error {
	return nil
}
