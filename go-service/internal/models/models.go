package models

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type Address struct {
	City    string `json:"city" binding:"required,min=2,max=64"`
	Street  string `json:"street" binding:"required,min=3,max=128"`
	ZipCode string `json:"zip_code" binding:"required,len=6,numeric"`
}

type Item struct {
	Name     string  `json:"name" binding:"required,min=2,max=64"`
	Quantity int     `json:"quantity" binding:"required,min=1,max=100"`
	Price    float64 `json:"price" binding:"required,gt=0"`
}

type Metadata struct {
	Priority string   `json:"priority" binding:"required,oneof=low medium high"`
	Tags     []string `json:"tags" binding:"required,min=1,max=5,dive,required,min=2,max=20"`
}

type ProcessRequest struct {
	RequestID string   `json:"request_id" binding:"required,uuid4"`
	Customer  string   `json:"customer" binding:"required,min=3,max=64"`
	Address   Address  `json:"address" binding:"required"`
	Items     []Item   `json:"items" binding:"required,min=1,max=10,dive"`
	Metadata  Metadata `json:"metadata" binding:"required"`
}

type ProcessResponse struct {
	RequestID   string   `json:"request_id"`
	ApprovedBy  string   `json:"approved_by"`
	ItemsCount  int      `json:"items_count"`
	TotalAmount float64  `json:"total_amount"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
}
