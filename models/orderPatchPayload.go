package models

type OrderPatchPayload struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}
