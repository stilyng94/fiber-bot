package repos

import (
	"context"

	"github.com/fiber-bot/models"
	"github.com/uptrace/bun"
)

type HandleAddProductPayload struct {
	Title       string `json:"title" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=1"`
	ImageUrl    string `json:"imageUrl" validate:"required,url"`
	Description string `json:"description" validate:"required,gt=1"`
}

type ProductRepo interface {
	CreateProduct(ctx context.Context, payload HandleAddProductPayload) (*models.ProductModel, error)
	UpdateProduct(ctx context.Context, ID int, status models.ProductStatus) (*models.ProductModel, error)
	DeleteProduct(ctx context.Context, ID int) error
	GetProduct(ctx context.Context, ID int) (*models.ProductModel, error)
	GetProducts(ctx context.Context) (*models.Paginate[models.ProductModel], error)
}

type ProductRepoImpl struct {
	db *bun.DB
}

// GetProducts implements ProductRepo.
func (u *ProductRepoImpl) GetProducts(ctx context.Context) (*models.Paginate[models.ProductModel], error) {
	products := []models.ProductModel{}
	count, err := u.db.NewSelect().Model(&products).Where("status = ?", "available").Limit(10).ScanAndCount(ctx)
	if err != nil {
		return nil, err
	}
	return &models.Paginate[models.ProductModel]{
		Edges: products,
		Page:  1, Limit: 10, Total: count,
	}, nil
}

// DeleteProduct implements ProductRepo.
func (u *ProductRepoImpl) DeleteProduct(ctx context.Context, ID int) error {
	product := &models.ProductModel{ID: int64(ID)}
	_, err := u.db.NewDelete().Model(product).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetProduct implements ProductRepo.
func (u *ProductRepoImpl) GetProduct(ctx context.Context, ID int) (*models.ProductModel, error) {
	var product models.ProductModel
	err := u.db.NewSelect().Model(&product).Where("id = ?", int64(ID)).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// UpdateProduct implements ProductRepo.
func (u *ProductRepoImpl) UpdateProduct(ctx context.Context, ID int, status models.ProductStatus) (*models.ProductModel, error) {
	product := &models.ProductModel{ID: int64(ID), Status: status}
	_, err := u.db.NewUpdate().Model(product).Column("status").WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// CreateProduct implements ProductRepo.
func (u *ProductRepoImpl) CreateProduct(ctx context.Context, payload HandleAddProductPayload) (*models.ProductModel, error) {

	product := &models.ProductModel{Title: payload.Title, ImageUrl: payload.ImageUrl, Price: payload.Price, Description: payload.Description}
	_, err := u.db.NewInsert().Model(product).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func NewProductRepo(db *bun.DB) ProductRepo {
	return &ProductRepoImpl{db: db}
}
