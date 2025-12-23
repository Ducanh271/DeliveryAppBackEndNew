package repository

import (
	"database/sql"
	"fmt"
)

type UnitOfWork interface {
	Execute(fn func(repoProvider func(repoType any) any) error) error
}

type sqlUnitOfWork struct {
	db *sql.DB
}

// NewUnitOfWork là hàm khởi tạo UoW
func NewUnitOfWork(db *sql.DB) UnitOfWork {
	return &sqlUnitOfWork{
		db: db,
	}
}
func (uow *sqlUnitOfWork) Execute(fn func(repoProvider func(repoType any) any) error) error {
	// 1. Bắt đầu transaction
	tx, err := uow.db.Begin()
	if err != nil {
		return fmt.Errorf("không thể bắt đầu transaction: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 2. Đây là "nhà máy" (factory)
	// Nó sẽ tạo repo dựa trên "kiểu" (type) bạn yêu cầu
	repoProvider := func(repoType any) any {
		switch repoType {
		case (*UserRepository)(nil): // Nếu hỏi "UserRepository"
			return NewUserRepository(tx) // Trả về repo với transaction

		case (*ProductRepository)(nil): // Nếu hỏi "ProductRepository"
			return NewProductRepository(tx) // Trả về repo với transaction

		case (*RefreshTokenRepository)(nil): // (Thêm repo khác nếu cần)
			return NewRefreshTokenRepository(tx)
		case (*OrderRepository)(nil):
			return NewOrderRepository(tx)
		case (*ReviewRepository)(nil):
			return NewReviewRepository(tx)
		default:
			return nil
		}
	}

	// 3. Chạy logic của service, truyền "nhà máy" vào
	err = fn(repoProvider)
	if err != nil {
		return err // Lỗi, sẽ Rollback
	}

	// 4. Commit
	committed = true
	return tx.Commit()
}

// func (uow *sqlUnitOfWork) Execute(fn func(UserRepository) error) error {
// 	tx, err := uow.db.Begin()
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", err)
// 	}
// 	commited := false
// 	defer func() {
// 		if !commited {
// 			tx.Rollback()
// 		}
// 	}()
//
// 	repo := NewUserRepository(tx)
//
// 	if err := fn(repo); err != nil {
// 		return fmt.Errorf("failed to execute function in unit of work: %w", err)
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}
// 	commited = true
// 	return nil
// }
