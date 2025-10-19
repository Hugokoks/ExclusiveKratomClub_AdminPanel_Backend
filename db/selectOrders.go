package db

import (
	"AdminPanelAPI/apperrors"
	"AdminPanelAPI/models"
	"context"
	"fmt"
	"strings"

	ekc_db "github.com/Hugokoks/kratomclub-go-common/db"
)

// Struktura, která odpovídá tomu, co posíláme na frontend
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

// Funkce pro bezpečné sestavení a spuštění dotazu
func SelectOrders(ctx context.Context, filters models.OrderFilters) ([]Order, error) {
	// Základní dotaz už je bez JOINu
	query := `
		SELECT 
			o.number, 
			TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at, 
			o.customer_first_name, 
			o.customer_last_name, 
			o.customer_email, 
			CASE 
				WHEN o.delivery_method LIKE '%home%' THEN o.address_street || ', ' || o.address_city
				ELSE o.pickup_name
			END as delivery_address,
			CASE 
                WHEN POSITION('-' IN o.payment_method) > 0 THEN SPLIT_PART(o.payment_method, '-', 2)
                ELSE o.payment_method
            END as payment_method,
            CASE 
                WHEN POSITION('-' IN o.delivery_method) > 0 THEN SPLIT_PART(o.delivery_method, '-', 2)
                ELSE o.delivery_method
            END as delivery_method,
			
			o.total_czk,
			(SELECT SUM(oi.quantity) FROM order_items oi WHERE oi.order_id = o.id) as item_count,
			o.status
		FROM orders o
		WHERE 1=1
	`

	args := []interface{}{}
	argId := 1

	// --- Dynamické přidávání podmínek pro filtry ---
	if filters.ID != "" {
		query += fmt.Sprintf(" AND o.number ILIKE $%d", argId)
		args = append(args, "%"+filters.ID+"%")
		argId++
	}
	if filters.FirstName != "" {
		query += fmt.Sprintf(" AND o.customer_first_name ILIKE $%d", argId)
		args = append(args, "%"+filters.FirstName+"%")
		argId++
	}
	if filters.LastName != "" {
		query += fmt.Sprintf(" AND o.customer_last_name ILIKE $%d", argId)
		args = append(args, "%"+filters.LastName+"%")
		argId++
	}
	if filters.Email != "" {
		query += fmt.Sprintf(" AND o.customer_email ILIKE $%d", argId)
		args = append(args, "%"+filters.Email+"%")
		argId++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND o.status = $%d", argId)
		args = append(args, filters.Status)
		argId++
	}
	if filters.PaymentMethod != "" {
		query += fmt.Sprintf(" AND o.payment_method = $%d", argId)
		args = append(args, filters.PaymentMethod)
		argId++
	}
	if filters.DeliveryMethod != "" {
		query += fmt.Sprintf(" AND o.delivery_method = $%d", argId)
		args = append(args, filters.DeliveryMethod)
		argId++
	}

	// Speciální logika pro adresu
	if filters.Address != "" {
		addressPattern := "%" + filters.Address + "%"
		query += fmt.Sprintf(" AND (o.address_street ILIKE $%d OR o.address_city ILIKE $%d OR o.pickup_name ILIKE $%d)", argId, argId+1, argId+2)
		args = append(args, addressPattern, addressPattern, addressPattern)
		argId += 3
	}

	// --- BEZPEČNÉ PŘIDÁNÍ ŘAZENÍ (ORDER BY) ---
	allowedSortBy := map[string]string{
		"date":      "o.created_at",
		"price":     "o.total_czk",
		"itemCount": "item_count", // Toto funguje díky aliasu
	}
	sortByColumn, ok := allowedSortBy[filters.SortBy]
	if !ok {
		sortByColumn = "o.created_at"
	}

	sortOrder := "DESC"
	if strings.ToUpper(filters.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortByColumn, sortOrder)

	rows, err := ekc_db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID, &o.CreatedAt, &o.FirstName, &o.LastName, &o.Email, &o.DeliveryAddress,
			&o.PaymentMethod, &o.DeliveryMethod, &o.TotalPrice, &o.ItemCount, &o.Status,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if len(orders) == 0 {
		return nil, apperrors.ErrOrdersNotFound
	}

	return orders, nil
}
