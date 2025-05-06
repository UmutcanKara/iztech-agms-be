package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/georgysavva/scany/v2/sqlscan"
	"strings"
	"time"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func (r repository) getUserByUsername(ctx context.Context, username string) (*User, error) {
	var u []*User

	query := `SELECT id, password, user_name, session FROM "USERS" WHERE user_name = $1`

	err := sqlscan.Select(ctx, r.db, &u, query, username)
	if err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, sql.ErrNoRows
	}
	return u[0], nil
}

func (r repository) getUsersByQuery(ctx context.Context, query string) ([]*User, error) {
	var u []*User

	err := sqlscan.Select(ctx, r.db, &u, query)
	if err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (r repository) getUsers(ctx context.Context) ([]*User, error) {
	var u []*User
	query := `SELECT id, user_name, session FROM "USERS"`

	err := sqlscan.Select(ctx, r.db, &u, query)
	if err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (r repository) createUser(ctx context.Context, user *User) (*User, error) {
	var uID int

	subEndDate := time.Now().AddDate(1, 0, 0)

	query := `INSERT INTO "USERS" (user_name, password, role) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, user.UserName, user.Password, user.Role, subEndDate).Scan(&uID)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	user.ID = uID
	return user, nil
}

func (r repository) createUsers(ctx context.Context, user []*CreateUserReq) error {

	values := ""
	for _, u := range user {
		values = fmt.Sprintf("%s, (%s, %s, %s)", values, u.UserName, u.Password, u.Role)
	}
	values = strings.TrimPrefix(values, ", ")
	query := `INSERT INTO "USERS" (user_name, password, role) 
				SELECT "USERS".user_name,
				       "USERS".password,
				       "USERS".role
				FROM (VALUES ( $1 ))
				    `
	err := r.db.QueryRowContext(ctx, query, values)
	if err != nil {
		println(err.Err())
		return err.Err()
	}
	return nil
}

func (r repository) setToken(ctx context.Context, uname string, token string) error {

	query := `UPDATE "USERS" SET session = $1 WHERE user_name = $2`

	_, err := r.db.ExecContext(ctx, query, token, uname)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) changePassword(ctx context.Context, u *ChangePasswordReq) error {
	query := `UPDATE "USERS" SET password = $1 WHERE user_name = $2`
	_, err := r.db.ExecContext(ctx, query, u.Password, u.UserName)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}
