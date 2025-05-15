// domain/category.go
package domain

type Category struct {
	ID   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Product struct {
	ID          string   `bson:"_id,omitempty" json:"id"`
	Name        string   `bson:"name" json:"name"`
	Category    Category `bson:"category" json:"category"`
	Price       float64  `bson:"price" json:"price"`
	Stock       int      `bson:"stock" json:"stock"`
	Description string   `bson:"description" json:"description"`
}
