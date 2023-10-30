package models

type Paginate[T OrderModel | UserModel | ProductModel] struct {
	Edges []T `json:"edges"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
