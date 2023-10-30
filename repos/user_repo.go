package repos

import (
	"context"

	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/models"
	"github.com/uptrace/bun"
)

type UserRepo interface {
	CreateUser(ctx context.Context, payload internal.DecodeTelegramHashPayload) (*models.UserModel, error)
	UpdateUser(ctx context.Context, ID int, name string) (*models.UserModel, error)
	DeleteUser(ctx context.Context, ID int) error
	GetUser(ctx context.Context, ID int) (*models.UserModel, error)
	GetUsers(ctx context.Context) (*models.Paginate[models.UserModel], error)
}

type UserRepoImpl struct {
	db *bun.DB
}

// GetUsers implements UserRepo.
func (u *UserRepoImpl) GetUsers(ctx context.Context) (*models.Paginate[models.UserModel], error) {
	users := []models.UserModel{}
	count, err := u.db.NewSelect().Model(&users).Limit(10).ScanAndCount(ctx)
	if err != nil {
		return nil, err
	}
	return &models.Paginate[models.UserModel]{
		Edges: users,
		Page:  1, Limit: 10, Total: count,
	}, nil
}

// DeleteUser implements UserRepo.
func (u *UserRepoImpl) DeleteUser(ctx context.Context, ID int) error {
	user := &models.UserModel{ID: int64(ID)}
	_, err := u.db.NewDelete().Model(user).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetUser implements UserRepo.
func (u *UserRepoImpl) GetUser(ctx context.Context, ID int) (*models.UserModel, error) {
	var user models.UserModel
	err := u.db.NewSelect().Model(&user).Where("id = ?", int64(ID)).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser implements UserRepo.
func (u *UserRepoImpl) UpdateUser(ctx context.Context, ID int, name string) (*models.UserModel, error) {
	user := &models.UserModel{ID: int64(ID), FirstName: name}
	_, err := u.db.NewUpdate().Model(user).Column("name").WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser implements UserRepo.
func (u *UserRepoImpl) CreateUser(ctx context.Context, payload internal.DecodeTelegramHashPayload) (*models.UserModel, error) {
	user := &models.UserModel{FirstName: payload.FirstName, TelegramID: payload.TelegramID,
		LastName: payload.LastName, Username: payload.Username,
	}

	err := u.db.NewSelect().Model(user).Where("telegram_id = ?", payload.TelegramID).Scan(ctx)
	if err != nil {
		// create
		_, err = u.db.NewInsert().Model(user).Exec(ctx)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	//update
	_, err = u.db.NewUpdate().Model(user).Column("first_name", "last_name", "username", "chat_id", "chat_type").Where("telegram_id = ?", payload.TelegramID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &UserRepoImpl{db: db}
}
