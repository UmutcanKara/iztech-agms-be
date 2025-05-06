package auth

import (
	"context"
	"errors"
	"fmt"
	"iztech-agms/util"
	"strconv"
	"strings"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(5) * time.Second,
	}
}

func (s *service) getUserByUsername(ctx context.Context, username string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.Repository.getUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *service) getUsersByQuery(ctx context.Context, q map[string][]string) ([]*GetUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	queryTemp := `SELECT id, user_name FROM "USERS"`
	queryTemp = fmt.Sprintf("%s WHERE ", queryTemp)
	for key, values := range q {
		for index, value := range values {
			if index > 0 {
				queryTemp = fmt.Sprintf("%s OR ", queryTemp)
			}
			queryTemp = fmt.Sprintf("%s %s='%s'", queryTemp, key, value)
		}
		queryTemp = fmt.Sprintf("%s AND", queryTemp)
	}
	queryTemp = strings.TrimSuffix(queryTemp, "WHERE ")
	queryTemp = strings.TrimSuffix(queryTemp, "AND")
	queryTemp = fmt.Sprintf("%s ORDER BY id ASC", queryTemp)

	// Get matching users
	println(queryTemp)
	users, err := s.Repository.getUsersByQuery(ctx, queryTemp)
	if err != nil {
		return nil, err
	}
	res := make([]*GetUserRes, len(users))
	for i, v := range users {
		res[i] = &GetUserRes{
			UserName: v.UserName,
		}
	}
	return res, nil
}
func (s *service) getUsers(ctx context.Context, req *GetUsersReq) ([]*GetUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	resUser, err := s.Repository.getUserByUsername(ctx, req.UserName)
	if err != nil {
		return nil, err
	}
	err = s.compareEncrypted(req.Token, resUser.Session)
	if err != nil {
		return nil, err
	}

	resUsers, err := s.Repository.getUsers(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*GetUserRes, len(resUsers))
	for i, v := range resUsers {
		res[i] = &GetUserRes{
			UserName: v.UserName,
		}
	}
	return res, nil
}

func (s *service) createUsers(ctx context.Context, req *CreateUsersReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users := make([]*CreateUserReq, len(req.Users))

	for i, user := range req.Users {
		hashedPwd, err := util.HashPassword(user.Password)
		if err != nil {
			return err
		}
		users[i] = &CreateUserReq{
			UserName: user.UserName,
			Password: hashedPwd,
			Role:     user.Role,
		}
	}

	err := s.Repository.createUsers(ctx, users)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) register(ctx context.Context, user *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	//err := s.compareInputs(user.Password, user.PasswordConfirm)
	//if err != nil {
	//	return nil, err
	//}

	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		UserName: user.UserName,
		Password: hashedPwd,
	}

	r, err := s.Repository.createUser(ctx, u)
	if err != nil {
		return nil, err
	}
	res := &CreateUserRes{
		ID:       strconv.Itoa(r.ID),
		UserName: r.UserName,
	}
	return res, nil
}

func (s *service) login(ctx context.Context, u *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user := &User{
		UserName: u.UserName,
		Password: u.Password,
	}
	res, err := s.Repository.getUserByUsername(ctx, user.UserName)
	if err != nil {
		return nil, err
	}

	if err = s.compareEncrypted(res.Password, user.Password); err != nil {
		return nil, err
	}
	unencrypted := strconv.Itoa(res.ID) + time.Now().String()
	token, err := s.hashToken(unencrypted)
	if err != nil {
		return nil, err
	}

	err = s.Repository.setToken(ctx, user.UserName, unencrypted)
	if err != nil {
		return nil, err
	}

	return &LoginUserRes{ID: strconv.Itoa(res.ID), UserName: res.UserName, Token: token}, nil
}

func (s *service) changePassword(ctx context.Context, u *ChangePasswordReq, token string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	if err := s.compareInputs(u.NewPassword, u.NewPasswordConfirm); err != nil {
		return errors.New("incorrect")
	}
	if err := s.compareInputs(u.Password, u.NewPassword); err == nil {
		return errors.New("invalid")
	}
	user, err := s.Repository.getUserByUsername(ctx, u.UserName)
	if err != nil {
		return err
	}
	//println(user.Password, user.Session)
	err = s.compareEncrypted(user.Password, u.Password)
	if err != nil {
		return err
	}
	err = s.compareEncrypted(token, user.Session)
	if err != nil {
		return err
	}
	hashed, err := util.HashPassword(u.NewPassword)
	if err != nil {
		return err
	}
	u.Password = hashed

	err = s.Repository.changePassword(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) hashToken(token string) (string, error) {
	hash, err := util.HashPassword(token)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s *service) compareEncrypted(hashedStr string, str string) error {
	err := util.CheckPassword(str, hashedStr)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) compareInputs(b string, a string) error {
	ok := strings.Compare(a, b)
	if ok != 0 {
		return errors.New("inputs do not match")
	}
	return nil
}
