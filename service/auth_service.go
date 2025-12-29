package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	// "example.com/delivery-app/middleware"
	"github.com/golang-jwt/jwt/v5"

	"fmt"
	"time"

	"example.com/delivery-app/dto"
	"example.com/delivery-app/models"
	"example.com/delivery-app/repository"
	"example.com/delivery-app/security"
	"example.com/delivery-app/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	uow              repository.UnitOfWork
	emailSvc         EmailService
	jwtSecret        string
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	uow repository.UnitOfWork,
	emailSvc EmailService,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		uow:              uow,
		emailSvc:         emailSvc,
		jwtSecret:        jwtSecret,
	}
}

func (s *AuthService) SignUp(req *dto.SignUpRequest) (*models.User, error) {
	if !utils.IsPasswordStrong(req.Password) {
		return nil, errors.New("Mật khẩu quá yếu, cần ít nhất 8 ký tự bao gồm chữ hoa, số và ký tự đặc biệt")
	}
	// Check if email already exists
	exists, err := s.userRepo.CheckEmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, ErrEmailInUse
	}
	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Address:  req.Address,
		Role:     "customer",
	}
	// 3. Unit of work
	err = s.uow.Execute(func(repoProvider func(repoType any) any) error {
		prodRepo := repoProvider((*repository.UserRepository)(nil)).(repository.UserRepository)
		insertedID, err := prodRepo.CreateUser(user)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		user.ID = insertedID

		// create otp
		otp, err := security.GenerateOTP()
		if err != nil {
			return fmt.Errorf("failed to generate OTP: %w", err)
		}
		expiry := time.Now().Add(10 * time.Minute)

		//save otp
		if err := prodRepo.UpdateOTP(user.Email, otp, expiry); err != nil {
			return fmt.Errorf("failed to update user OTP: %w", err)
		}
		user.OTPCode = &otp
		user.OTPExpiresAt = &expiry

		if err := s.emailSvc.SendOTP(user.Email, otp); err != nil {
			return fmt.Errorf("failed to send OTP email: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil

}

// func (s *AuthService) VerifyOTP(email string, otp string) error {
func (s *AuthService) VerifyOTP(email string, otp string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return ErrEmailNotFound
	}
	if user.OTPCode == nil || user.OTPExpiresAt == nil {
		return ErrOTPNotFound
	}
	if *user.OTPCode != otp {
		return ErrInvalidOTP
	}
	if time.Now().After(*user.OTPExpiresAt) {
		return ErrOTPExpired
	}
	// Mark user as verified
	if err := s.userRepo.UpdateStatusUser(user.Email, 1); err != nil {
		return fmt.Errorf("failed to mark user as verified: %w", err)
	}
	return nil
}

// func Login
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Failed to ger user: %v\n", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}
	if user.Status == 2 {
		return nil, ErrUserBanned
	}
	if user.Status == 0 {
		otp, err := security.GenerateOTP()
		if err != nil {
			log.Printf("Failed to generate OTP: %v\n", err)
			return nil, fmt.Errorf("failed to generate OTP: %w", err)
		}
		expiry := time.Now().Add(10 * time.Minute)
		// start transaction
		err = s.uow.Execute(func(repoProvider func(repoType any) any) error {
			repo := repoProvider((*repository.UserRepository)(nil)).(repository.UserRepository)
			// update otp to db
			if err := repo.UpdateOTP(user.Email, otp, expiry); err != nil {
				return fmt.Errorf("failed to update user OTP: %w", err)
			}
			if err := s.emailSvc.SendOTP(user.Email, otp); err != nil {
				return fmt.Errorf("failed to send OTP email: %w", err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return nil, ErrUserInactive
	}
	accessTokenStr, refreshTokenStr, err := s.CreateTokens(user)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Save(user.ID, refreshTokenStr, time.Now().Add(7*24*time.Hour)); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrEmailNotFound
	}
	otp, err := security.GenerateOTP()
	if err != nil {
		return fmt.Errorf("Failed to created otp: %w\n", err)
	}
	expiry := time.Now().Add(10 * time.Minute)

	err = s.uow.Execute(func(repoProvider func(repoType any) any) error {
		repo := repoProvider((*repository.UserRepository)(nil)).(repository.UserRepository)
		if err := repo.SetResetOTP(user.Email, otp, expiry); err != nil {
			return fmt.Errorf("Failed to save reset otp: %w\n", err)
		}
		if err := s.emailSvc.SendOTP(user.Email, otp); err != nil {
			return fmt.Errorf("Failed to send otp: %w\n", err)
		}
		return nil
	})
	return err
}

// === 2. VERIFY RESET OTP ===
func (s *AuthService) VerifyResetOTP(req *dto.VerifyResetOTPRequest) (string, error) {
	// 1. Lấy user
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", ErrUserNotFound
	}

	// 2. Validate OTP
	if user.ResetOTP == nil || *user.ResetOTP != req.OTP {
		return "", ErrInvalidOTP
	}
	if time.Now().After(*user.ResetOTPExpiresAt) {
		return "", ErrOTPExpired
	}

	// 3. Tạo token "đặt lại mật khẩu"
	claims := jwt.MapClaims{
		"email":   user.Email,
		"purpose": "reset_password", // Rất quan trọng
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo reset token: %w", err) // Lỗi 500
	}

	// 4. (Tùy chọn) Xóa OTP sau khi dùng
	_ = s.userRepo.ClearResetOTP(user.ID) // Bỏ qua lỗi, không quá quan trọng

	return tokenString, nil
}

// === 3. RESET PASSWORD ===
func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// 1. Validate token reset
	claims, err := s.validateResetToken(req.Token)
	if err != nil {
		return err // Trả về ErrInvalidResetToken
	}

	// 2. Kiểm tra mục đích (purpose)
	if p, ok := claims["purpose"].(string); !ok || p != "reset_password" {
		return ErrTokenPurposeMismatch
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return ErrInvalidResetToken // Token không có email
	}
	if !utils.IsPasswordStrong(req.NewPassword) {
		return errors.New("Mật khẩu mới quá yếu, cần ít nhất 8 ký tự, bao gồm chữ hoa, chữ số và kí tự đặc biệt")
	}
	// 4. Hash mật khẩu
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("lỗi hash password: %w", err) // Lỗi 500
	}

	// 5. Cập nhật CSDL (1 lệnh, không cần UoW)
	if err := s.userRepo.UpdatePasswordByEmail(email, string(hashed)); err != nil {
		return fmt.Errorf("lỗi cập nhật password: %w", err) // Lỗi 500
	}

	return nil
}

// (Hàm helper private để validate token)
func (s *AuthService) validateResetToken(tokenStr string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidResetToken
	}

	return claims, nil
}

// func Refresh token
func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// 1. Lấy token cũ từ DB
	// (Dùng repo đã được tiêm vào)
	refreshToken, err := s.refreshTokenRepo.GetByToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get refresh token: %w", err) // Lỗi 500
	}
	if refreshToken == nil {
		return nil, ErrInvalidRefreshToken
	}

	// 2. Kiểm tra hạn
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, ErrRefreshTokenExpired // Lỗi nghiệp vụ 401
	}

	// 3. Lấy thông tin user
	user, err := s.userRepo.GetUserByID(refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get user từ token: %w", err) // Lỗi 500
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 4. Tạo bộ token mới (Access + Refresh)
	accessTokenStr, newRefreshTokenStr, err := s.CreateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo tokens: %w", err) // Lỗi 500
	}

	// 5. Cập nhật token cũ bằng token mới (Token Rotation)
	err = s.refreshTokenRepo.UpdateToken(req.RefreshToken, newRefreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("lỗi cập nhật refresh token: %w", err)
	}

	// 6. Trả về DTO
	return &dto.RefreshTokenResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: newRefreshTokenStr,
	}, nil
}

// func create random string for request token
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// func create access token and refresh token
func (s *AuthService) CreateTokens(user *models.User) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(5 * time.Minute).Unix(),
	})
	accessTokenStr, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", errors.New("Can't create access token")
	}
	refreshTokenStr, err := GenerateRefreshToken()
	if err != nil {
		return "", "", errors.New("Can't create refresh token")
	}
	return accessTokenStr, refreshTokenStr, nil
}

func (s *AuthService) LogOut(reToken string) error {
	stt, err := s.refreshTokenRepo.DeleteByToken(reToken)
	if stt == true {
		if err != nil {
			return ErrInvalidRefreshToken
		}
		return nil
	}
	return errors.New("Can't log out")
}
