package models

// Product represents the product collection schema
type Product struct {
	ID    string  `json:"id" bson:"_id,omitempty"`
	Name  string  `json:"name" bson:"name"`
	Price float64 `json:"price" bson:"price"`
	Stock float64 `json:"stock" bson:"stock"`
}
