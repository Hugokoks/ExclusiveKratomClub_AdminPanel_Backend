package limiter

const (
	KiB = 1 << 10 // 1024
	MiB = 1 << 20 // 1,048,576
)

var OrderFilterLimit = map[string]int{

	"id":             40,
	"firstName":      50,
	"lastName":       50,
	"email":          120,
	"address":        120,
	"paymentMethod":  25,
	"deliveryMethod": 20,
	"status":         10,
	"sortBy":         15,
	"sortOrder":      6,
}
