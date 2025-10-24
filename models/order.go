package models

type OrderPatchPayload struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}
type OrderFilters struct {
	ID             string `form:"id"`
	FirstName      string `form:"firstName"`
	LastName       string `form:"lastName"`
	Email          string `form:"email"`
	Address        string `form:"address"`
	PaymentMethod  string `form:"paymentMethod"`
	DeliveryMethod string `form:"deliveryMethod"`
	Status         string `form:"status"`
	SortBy         string `form:"sortBy"`
	SortOrder      string `form:"sortOrder"`
}

type Order struct {
	ID              string  `json:"id"`
	CreatedAt       string  `json:"createdAt"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	Email           string  `json:"email"`
	DeliveryAddress string  `json:"deliveryAddress"`
	PaymentMethod   string  `json:"paymentMethod"`
	DeliveryMethod  string  `json:"deliveryMethod"`
	TotalPrice      float64 `json:"totalPrice"`
	ItemCount       int     `json:"itemCount"`
	Status          string  `json:"status"`
}