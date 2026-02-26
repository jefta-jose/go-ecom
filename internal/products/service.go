package products

import (
	"context"

	repo "github.com/jefta-jose/go-ecom/internal/adapters/postgresql/sqlc"
)

// Service is an interface that defines the methods for managing products. 
// It includes methods for creating a product and listing products.
type Service interface {
	CreateProduct(ctx context.Context, tempProduct CreateProductParams) (repo.Product, error)
	ListProducts(ctx context.Context) ([]repo.Product, error)
}

// NewService is a constructor function that takes a repo.Querier and returns a new instance of the Service interface.
// It creates a new svc struct, which implements the Service interface, and initializes it with the provided repository.
func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

// svc is a struct that implements the Service interface. It has a single field, repo, which is of type repo.Querier.
type svc struct {
	repo repo.Querier
}

// CreateProduct takes a context and a CreateProductParams struct, converts the parameters to the format expected by the repository, and then calls the CreateProduct method of the repository to create a new product in the database. It returns the created product or an error if something goes wrong.
func (s *svc) CreateProduct(ctx context.Context, tempProduct CreateProductParams) (repo.Product, error) {
	params := repo.CreateProductParams{
		Name:           tempProduct.Name,
		PriceInCenters: tempProduct.PriceInCenters,
		Quantity:       tempProduct.Quantity,
	}
	return s.repo.CreateProduct(ctx, params)
}

// ListProducts takes a context and calls the ListProducts method of the repository to retrieve a 
// list of products from the database. It returns a slice of products or an error if something goes wrong.
func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}
