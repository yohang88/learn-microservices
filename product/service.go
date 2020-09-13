package main

func NewService() Service {
	return service{}
}

type service struct{}

type Service interface {
	GetHealthCheck() (string, error)
	GetProduct(Id string) (Product, error)
}

type Product struct {
	Id 			string `json:"id"`
	Name 		string `json:"name"`
	Description string `json:"description"`
}

func (service) GetHealthCheck() (string, error) {
	return "ok", nil
}

func (service) GetProduct(Id string) (Product, error) {
	return Product{
		Id: Id,
		Name: "Test Product Name",
		Description: "Test Product Description",
	}, nil
}
