package products

import (
	"context"

	repo "github.com/jefta-jose/go-ecom/internal/adapters/postgresql/sqlc"
)

// Service is an interface that defines the methods for managing products. 
type Service interface {
	CreateProduct(ctx context.Context, tempProduct CreateProductParams) (repo.Product, error)
	ListProducts(ctx context.Context) ([]repo.Product, error)
	
}

// this struct contains a connection to the database
type svc struct {
	repo repo.Querier
}

// this service constructor that takes in the repo argument and returns a new Service Instance by implementing the svc struct.
func NewService(repo repo.Querier) Service {
	// it returns the implementation of the svc struct (db connection)
	return &svc{repo: repo}
}

// ***************************************************************************
// in order to implement interfaces in Go we need to create the methods
// through the svc struct
// ***************************************************************************

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
