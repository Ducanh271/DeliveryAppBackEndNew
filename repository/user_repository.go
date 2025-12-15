package repository

import (
	"database/sql"
	"errors"
	"example.com/delivery-app/models"
	"time"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error) //we can relplace any with interface{}
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

type UserRepository interface {
	// Function don't need tx
	CheckEmailExists(email string) (bool, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID int64) (*models.User, error)
	GetNumberOfCustomer() (int64, error)
	GetNumberOfShipper() (int64, error)
	UpdateStatusUserByUserID(userID int64, status int) error
	ClearOTP(userID int64) error
	SetResetOTP(email, otp string, expiry time.Time) error
	UpdatePasswordByEmail(email, hashed string) error
	GetAllUserWithType(role string, page, limit int) ([]models.User, int, error)
	ClearResetOTP(userID int64) error
	// tx functions
	UpdateStatusUser(email string, status int) error
	UpdateOTP(userEmail string, otp string, expiry time.Time) error
	CreateUser(user *models.User) (int64, error)
	// WithTx returns a new UserRepository that uses the given transaction.
	WithTx(tx *sql.Tx) UserRepository
}

type userRepo struct {
	db DBTX
}

func NewUserRepository(db DBTX) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) WithTx(tx *sql.Tx) UserRepository {
	return &userRepo{db: tx}
}

func (r *userRepo) CheckEmailExists(email string) (bool, error) {
	var exist int
	query := "select 1 from users where email = ? limit 1"
	err := r.db.QueryRow(query, email).Scan(&exist)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil

}

func (r *userRepo) CreateUser(user *models.User) (int64, error) {
	query := "insert into users (name, email, password, phone, address, role, created_at) values (?,?, ?, ?, ?, ?,?)"
	user.CreatedAt = time.Now()

	result, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.Address, user.Role, user.CreatedAt)
	if err != nil {
		return 0, err
	}
	// ...
	id, _ := result.LastInsertId()
	return id, nil
}
func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	query := "select id, name, email, password, phone, address, role, created_at, otp_code, otp_expires_at, status, reset_otp, reset_otp_expires_at from users where email = ?"
	row := r.db.QueryRow(query, email)
	var user models.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &user.CreatedAt, &user.OTPCode, &user.OTPExpiresAt, &user.Status, &user.ResetOTP, &user.ResetOTPExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetUserByID(userID int64) (*models.User, error) {
	query := "select id, name, email, password, phone, address, role, created_at from users where id = ?"
	row := r.db.QueryRow(query, userID)
	var user models.User
	var createdAtstr string

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &createdAtstr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtstr)
	return &user, nil
}

func (r *userRepo) GetNumberOfCustomer() (int64, error) {
	query := "select count(*) from users where role = 'customer'"
	var num int64
	err := r.db.QueryRow(query).Scan(&num)
	if err != nil {
		return 0, err
	}
	return num, nil
}
func (r *userRepo) GetNumberOfShipper() (int64, error) {
	query := "select count(*) from users where role = 'shipper'"
	var num int64
	err := r.db.QueryRow(query).Scan(&num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// tx
func (r *userRepo) UpdateOTP(userEmail string, otp string, expiry time.Time) error {
	query := `update users set otp_code = ?, otp_expires_at = ? where email = ?`
	_, err := r.db.Exec(query, otp, expiry, userEmail)
	return err
}
func (r *userRepo) UpdateStatusUserByUserID(userID int64, status int) error {
	updateQuery := `UPDATE users SET status = ? WHERE id = ?`
	_, err := r.db.Exec(updateQuery, status, userID)
	return err
}

// func UpdateStatusUserTx(tx *sql.Tx, email string, status int) error {
// 	updateQuery := `UPDATE users SET status = ? WHERE email = ?`
// 	_, err := tx.Exec(updateQuery, status, email)
// 	return err

func (r *userRepo) UpdateStatusUser(email string, status int) error {
	updateQuery := `UPDATE users SET status = ? WHERE email = ?`
	_, err := r.db.Exec(updateQuery, status, email)
	return err
}
func (r *userRepo) ClearOTP(userID int64) error {
	_, err := r.db.Exec(`UPDATE users SET otp_code = NULL, otp_num = 0, otp_expires_at = NULL WHERE id = ?`, userID)
	return err
}

func (r *userRepo) SetResetOTP(email, otp string, expiry time.Time) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET reset_otp = ?, reset_otp_expires_at = ?, reset_otp_num = 5
		WHERE email = ?`, otp, expiry, email)
	return err
}

func (r *userRepo) UpdatePasswordByEmail(email, hashed string) error {
	_, err := r.db.Exec(`UPDATE users SET password = ? WHERE email = ?`, hashed, email)
	return err
}

// func get customer or shipper for admin
func (r *userRepo) GetAllUserWithType(role string, page, limit int) ([]models.User, int, error) {
	offset := (page - 1) * limit
	var total int
	err := r.db.QueryRow("select count(*) from users where role = ?", role).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := "select id, name, email, phone, address, role, status from users where role = ? order by id limit ? offset ?"
	rows, err := r.db.Query(query, role, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.Role,
			&user.Status,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return users, total, nil

}
func (r *userRepo) ClearResetOTP(userID int64) error {
	_, err := r.db.Exec(`UPDATE users SET reset_otp = NULL, reset_otp_num = 0, reset_otp_expires_at = NULL WHERE id = ?`, userID)
	return err
}
