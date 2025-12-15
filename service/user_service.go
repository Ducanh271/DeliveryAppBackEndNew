package service

import (
	"database/sql"
	"errors"
	"example.com/delivery-app/dto"
	"example.com/delivery-app/models"
	"example.com/delivery-app/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// UserService chứa logic nghiệp vụ quản lý user
type UserService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository // Cần cho BanUser
	uow              repository.UnitOfWork             // Cần cho CreateShipper
}

// NewUserService là hàm khởi tạo
func NewUserService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	uow repository.UnitOfWork,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		uow:              uow,
	}
}

// === LOGIC TỪ CreateShipper ===
func (s *UserService) CreateShipper(req *dto.CreateShipperRequest) (*models.User, error) {
	// 1. Kiểm tra email (giống SignUp)
	exists, err := s.userRepo.CheckEmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra email: %w", err) // Lỗi 500
	}
	if exists {
		return nil, ErrEmailInUse // Lỗi 409
	}

	// 2. Validate mật khẩu (nếu cần)

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("lỗi hash password: %w", err) // Lỗi 500
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Address:  req.Address,
		Role:     "shipper",
	}

	// 4. Dùng Unit of Work để tạo user và set status
	err = s.uow.Execute(func(repoProvider func(repoType any) any) error {
		repo := repoProvider((*repository.UserRepository)(nil)).(repository.UserRepository)
		insertedID, err := repo.CreateUser(user)
		if err != nil {
			return fmt.Errorf("lỗi tạo user: %w", err)
		}
		user.ID = insertedID

		// Set status=1 (active)
		if err := repo.UpdateStatusUser(user.Email, 1); err != nil {
			return fmt.Errorf("lỗi cập nhật status shipper: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err // Trả về lỗi từ UoW
	}
	return user, nil
}

// === LOGIC TỪ ProfileHandler ===
func (s *UserService) GetUserProfile(userID int64) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || user == nil {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("lỗi lấy user profile: %w", err)
	}

	// Mapping: Chuyển models.User -> dto.UserProfileResponse
	// để không làm lộ Password, OTP...
	return &dto.UserProfileResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Phone:   user.Phone,
		Address: user.Address,
		Role:    user.Role,
	}, nil
}

func (s *UserService) BanUser(userID int64) error {
	// Dùng UoW (hoặc transaction) để đảm bảo 2 việc
	// (Hoặc bạn có thể bỏ UoW nếu không quá quan trọng)

	// 1. Xóa refresh token (lấy từ logic cũ)
	if err := s.refreshTokenRepo.DeleteByUserID(userID); err != nil {
		return fmt.Errorf("lỗi xóa refresh token: %w", err)
	}

	// 2. Cập nhật status
	if err := s.userRepo.UpdateStatusUserByUserID(userID, 2); err != nil { // 2 = Banned
		return fmt.Errorf("lỗi ban user: %w", err)
	}
	return nil
}

// === LOGIC TỪ UnBanUserAccountHandler ===
func (s *UserService) UnbanUser(userID int64) error {
	// Cập nhật status
	if err := s.userRepo.UpdateStatusUserByUserID(userID, 1); err != nil { // 1 = Active
		return fmt.Errorf("lỗi unban user: %w", err)
	}
	return nil
}

// === LOGIC TỪ GetAllCustomers/ShippersHandler ===
func (s *UserService) GetUsersByRole(role string, page, limit int) (*dto.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := s.userRepo.GetAllUserWithType(role, page, limit)
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy user list: %w", err)
	}

	// Mapping
	var userResponses []dto.UserProfileResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserProfileResponse{
			ID:      user.ID,
			Email:   user.Email,
			Name:    user.Name,
			Phone:   user.Phone,
			Address: user.Address,
			Role:    user.Role,
			Status:  user.Status,
		})
	}

	totalPages := (total + limit - 1) / limit
	return &dto.UserListResponse{
		Users: userResponses,
		Pagination: dto.Pagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
}

// === LOGIC TỪ GetNumberOfCustomer/ShipperHandler ===
func (s *UserService) GetDashboardStats() (*dto.DashboardStatsResponse, error) {
	numCustomers, err := s.userRepo.GetNumberOfCustomer()
	if err != nil {
		return nil, err
	}
	numShippers, err := s.userRepo.GetNumberOfShipper()
	if err != nil {
		return nil, err
	}

	return &dto.DashboardStatsResponse{
		TotalCustomers: numCustomers,
		TotalShippers:  numShippers,
	}, nil
}
