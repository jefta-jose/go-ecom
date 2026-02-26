package products

import (
	"context"

	repo "github.com/jefta-jose/go-ecom/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateProduct(ctx context.Context, tempProduct CreateProductParams) (repo.Product, error)
	ListProducts(ctx context.Context) ([]repo.Product, error)
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

type svc struct {
	repo repo.Querier
}

func (s *svc) CreateProduct(ctx context.Context, tempProduct CreateProductParams) (repo.Product, error) {
	params := repo.CreateProductParams{
		Name:           tempProduct.Name,
		PriceInCenters: tempProduct.PriceInCenters,
		Quantity:       tempProduct.Quantity,
	}
	return s.repo.CreateProduct(ctx, params)
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}
