package auth

import (
	"context"
	"database/sql"
)

type User struct {
	ID              int            `json:"id"`
	UserName        string         `json:"user_name"`
	Password        string         `json:"password"`
	PwdResetToken   sql.NullString `json:"pwd_reset_token"`
	PwdResetExpires sql.NullTime   `json:"pwd_reset_expires"`
	Session         string
}

type Repository interface {
	// GET Reqs
	getUserByUsername(ctx context.Context, username string) (*User, error)
	getUsersByQuery(ctx context.Context, query string) ([]*User, error)
	getUsers(ctx context.Context) ([]*User, error)

	// POST Reqs
	createUser(ctx context.Context, user *User) (*User, error)
	setToken(ctx context.Context, uname string, token string) error

	//// PUT Reqs
	changePassword(ctx context.Context, user *ChangePasswordReq) error
	//
	//// DELETE Reqs
	//deleteUser(ctx context.Context, user *User) error
}

type Service interface {
	// GET Reqs
	getUsersByQuery(ctx context.Context, q map[string][]string) ([]*GetUserRes, error)
	getUsers(ctx context.Context, user *GetUsersReq) ([]*GetUserRes, error)
	getUserByUsername(ctx context.Context, username string) (*User, error)

	// POST Reqs
	createUser(ctx context.Context, user *CreateUserReq) (*CreateUserRes, error)
	register(ctx context.Context, user *CreateUserReq) (*CreateUserRes, error)
	login(ctx context.Context, u *LoginUserReq) (*LoginUserRes, error)

	//// PUT Reqs
	changePassword(ctx context.Context, u *ChangePasswordReq, token string) error

	//// DELETE Reqs
	//deleteUser(ctx context.Context, u *DeleteUserReq) error
	// UTILS

	compareEncrypted(hashedStr string, str string) error
	compareInputs(b string, a string) error
	hashToken(token string) (string, error)
}

type GetUsersReq struct {
	UserName string `json:"user_name" db:"user_name"`
	Token    string `json:"token" db:"token"`
}

type GetUserRes struct {
	UserName string `json:"user_name"`
	Role     string `json:"role"`
}

type GetUserByUsernameReq struct {
	UserName string `json:"user_name" db:"user_name"`
}

type CreateUserReq struct {
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
}

type LoginUserReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Token    string `json:"token"`
}

type ChangePasswordReq struct {
	UserName           string `json:"user_name"`
	Password           string `json:"password"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`
}

type ChangePasswordRes struct {
	UserName string `json:"user_name"`
}

type DeleteUserReq struct {
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	Confirmation string `json:"confirmation"`
}
