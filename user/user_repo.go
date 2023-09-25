package user

import (
	"context"

	"github.com/uptrace/bun"
)

type UserRepo interface {
	InsertUser(ctx context.Context, payload HandleAddUserPayload) (*UserModel, error)
	UpdateUser(ctx context.Context, ID int, name string) (*UserModel, error)
	DeleteUser(ctx context.Context, ID int) error
	GetUser(ctx context.Context, ID int) (*UserModel, error)
	GetUsers(ctx context.Context) ([]UserModel, error)
}

type UserRepoImpl struct {
	db *bun.DB
}

// GetUsers implements UserRepo.
func (u *UserRepoImpl) GetUsers(ctx context.Context) ([]UserModel, error) {
	users := []UserModel{}
	err := u.db.NewSelect().Model(&users).Limit(10).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteUser implements UserRepo.
func (u *UserRepoImpl) DeleteUser(ctx context.Context, ID int) error {
	user := &UserModel{ID: int64(ID)}
	_, err := u.db.NewDelete().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetUser implements UserRepo.
func (u *UserRepoImpl) GetUser(ctx context.Context, ID int) (*UserModel, error) {
	var user UserModel
	err := u.db.NewSelect().Model(&user).Where("id = ?", int64(ID)).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser implements UserRepo.
func (u *UserRepoImpl) UpdateUser(ctx context.Context, ID int, name string) (*UserModel, error) {
	user := &UserModel{ID: int64(ID), Name: name}
	_, err := u.db.NewUpdate().Model(user).Column("name").WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// InsertUser implements UserRepo.
func (u *UserRepoImpl) InsertUser(ctx context.Context, payload HandleAddUserPayload) (*UserModel, error) {
	user := &UserModel{Name: payload.Name, Email: payload.Email, Password: payload.Password}
	_, err := u.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &UserRepoImpl{db: db}
}
