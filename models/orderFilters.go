package models

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
