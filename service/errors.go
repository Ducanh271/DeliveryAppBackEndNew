package service

import "errors"

var (
	ErrEmailInUse    = errors.New("Email already in use")
	ErrEmailNotFound = errors.New("Email not found")
	ErrInvalidEmail  = errors.New("Invalid email")

	ErrUserNotFound = errors.New("User not found")
	ErrUserInactive = errors.New("User is inactive")
	ErrUserBanned   = errors.New("User is banned")

	ErrInvalidResetToken    = errors.New("Invalid or expired reset token")
	ErrTokenPurposeMismatch = errors.New("Invalid token purpose")

	ErrInvalidRefreshToken = errors.New("Invalid refresh token")
	ErrRefreshTokenExpired = errors.New("Refresh token is expired")

	ErrInvalidOTP  = errors.New("Invalid OTP")
	ErrOTPExpired  = errors.New("Expired OTP")
	ErrOTPNotFound = errors.New("OTP not found")

	ErrInvalidPassword = errors.New("Invalid password")

	ErrProductNotFound = errors.New("Product not found")
	ErrNoImages        = errors.New("At least one image is required")

	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderNotPending    = errors.New("order is not in pending status")
	ErrNotYourOrder       = errors.New("this order does not belong to you")
	ErrInvalidReceiver    = errors.New("Receiver is not relative with this order")
	ErrOrderNotOwned      = errors.New("you do not have permission to manage this order")
	ErrMaxOrdersReached   = errors.New("you have reached the maximum number of orders")
	ErrCannotCancel       = errors.New("cannot cancel this order (not pending)")
	ErrOrderNotProcessing = errors.New("order is not available for shipping")

	ErrUserCannotReview = errors.New("you can only review products you have purchased and received")
	ErrReviewExists     = errors.New("you have already reviewed this product for this order")
)
