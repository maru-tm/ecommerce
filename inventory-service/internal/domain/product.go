// domain/product.go
package domain

type ProductRepository interface {
	CreateProduct(product *Product) error
	GetProductByID(id string) (*Product, error)
	GetProductByName(name string) (*Product, error)
	ListProducts() ([]Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id string) error
	CheckStock(productID string) (bool, error)
	SearchProducts(keyword string, categoryID string) ([]Product, error)
}

type ProductUseCase interface {
	CreateProduct(product *Product) error
	GetProductByID(id string) (*Product, error)
	ListProducts() ([]Product, error)
	UpdateProduct(id string, product *Product) error
	DecreaseStock(productID string, quantity int) error
	DeleteProduct(id string) error
	CheckStock(productID string, quantity int) (bool, error)
	SearchProducts(query string, categoryID string) ([]Product, error)
	UpdateProductStock(productID string, quantity int) error
}
