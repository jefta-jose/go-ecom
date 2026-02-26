package products

import (
	"log"
	"net/http"

	"github.com/jefta-jose/go-ecom/internal/json"
)

// This just means: “Every handler has a reference to the service it can call.”
// The handler does not implement the Service interface itself.
// It just uses the service to get things done.
type handler struct {
	service Service
}

// this is dependency injection in action
// NewHandler is a function that takes a Service as an argument and returns a pointer to a handler struct.
// 
func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var tempProduct CreateProductParams
	if err := json.Read(r, &tempProduct); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), tempProduct)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdProduct)
}
