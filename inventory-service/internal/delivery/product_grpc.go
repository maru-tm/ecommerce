package delivery

import (
	"context"
	"inventory-service/internal/domain"
	"inventory-service/internal/proto"
	"inventory-service/internal/usecase"
	"log"
)

type ProductServiceServer struct {
	proto.UnimplementedProductServiceServer
	useCase usecase.ProductUseCase
}

func NewProductServiceServer(useCase usecase.ProductUseCase) *ProductServiceServer {
	log.Println("Initializing ProductServiceServer")
	return &ProductServiceServer{
		useCase: useCase,
	}
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *proto.Product) (*proto.Product, error) {
	log.Printf("CreateProduct called with request: %v", req)

	product := &domain.Product{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Category:    domain.Category{ID: req.GetCategory().GetId(), Name: req.GetCategory().GetName()},
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
		Description: req.GetDescription(),
	}

	log.Printf("Product created: %v", product)

	err := s.useCase.CreateProduct(product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}

	log.Printf("Product successfully created: %v", product)

	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Category:    &proto.Category{Id: product.Category.ID, Name: product.Category.Name},
		Price:       product.Price,
		Stock:       int32(product.Stock),
		Description: product.Description,
	}, nil
}

func (s *ProductServiceServer) GetProductByID(ctx context.Context, req *proto.ProductId) (*proto.Product, error) {
	log.Printf("GetProductByID called with product ID: %v", req.GetId())

	product, err := s.useCase.GetProductByID(req.GetId())
	if err != nil {
		log.Printf("Error getting product by ID: %v", err)
		return nil, err
	}

	log.Printf("Product found: %v", product)

	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Category:    &proto.Category{Id: product.Category.ID, Name: product.Category.Name},
		Price:       product.Price,
		Stock:       int32(product.Stock),
		Description: product.Description,
	}, nil
}

func (s *ProductServiceServer) ListProducts(ctx context.Context, req *proto.Empty) (*proto.ProductList, error) {
	log.Println("ListProducts called")

	products, err := s.useCase.ListProducts()
	if err != nil {
		log.Printf("Error listing products: %v", err)
		return nil, err
	}

	log.Printf("Found %d products", len(products))

	var productList []*proto.Product
	for _, p := range products {
		productList = append(productList, &proto.Product{
			Id:          p.ID,
			Name:        p.Name,
			Category:    &proto.Category{Id: p.Category.ID, Name: p.Category.Name},
			Price:       p.Price,
			Stock:       int32(p.Stock),
			Description: p.Description,
		})
	}

	log.Printf("Returning product list with %d products", len(productList))

	return &proto.ProductList{Products: productList}, nil
}

func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *proto.Product) (*proto.Product, error) {
	log.Printf("UpdateProduct called with request: %v", req)

	product := &domain.Product{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Category:    domain.Category{ID: req.GetCategory().GetId(), Name: req.GetCategory().GetName()},
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
		Description: req.GetDescription(),
	}

	log.Printf("Product to update: %v", product)

	err := s.useCase.UpdateProduct(req.GetId(), product)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return nil, err
	}

	log.Printf("Product successfully updated: %v", product)

	return &proto.Product{
		Id:          product.ID,
		Name:        product.Name,
		Category:    &proto.Category{Id: product.Category.ID, Name: product.Category.Name},
		Price:       product.Price,
		Stock:       int32(product.Stock),
		Description: product.Description,
	}, nil
}

func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *proto.ProductId) (*proto.Empty, error) {
	log.Printf("DeleteProduct called with product ID: %v", req.GetId())

	err := s.useCase.DeleteProduct(req.GetId())
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return nil, err
	}

	log.Printf("Product successfully deleted with ID: %v", req.GetId())

	return &proto.Empty{}, nil
}
