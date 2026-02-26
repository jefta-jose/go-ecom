package products

type CreateProductParams struct {
    Name           string `json:"name"`
    PriceInCenters int32  `json:"priceInCenters"`
    Quantity       int32  `json:"quantity"`
}
