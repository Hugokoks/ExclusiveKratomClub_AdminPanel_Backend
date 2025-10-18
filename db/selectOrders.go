package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"AdminPanelAPI/apperrors"
	"AdminPanelAPI/models"

	ekc_db "github.com/Hugokoks/kratomclub-go-common/db"
)

type Order struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	DeliveryAddress string    `json:"deliveryAddress"`
	PaymentMethod   string    `json:"paymentMethod"`
	DeliveryMethod  string    `json:"deliveryMethod"`
	TotalPrice      float64   `json:"totalPrice"`
	ItemCount       int       `json:"itemCount"`
	Status          string    `json:"status"`
}

func SelectOrders(ctx context.Context, filters models.OrderFilters) ([]Order, error) {

	query := `
		SELECT 
			o.number, o_d.created_at, o.customer_first_name, o.customer_last_name, o.customer_email, 
			CASE 
				WHEN o.delivery_method = 'packeta-home' THEN o.address_street || ', ' || o.address_city
				ELSE o_d.pickup_name
			END as delivery_address,
			o.payment_method, o.delivery_method, o_d.total_czk, 
			-- Pro každou objednávku (o.id) sečteme hodnoty ze sloupce 'quantity'
			(SELECT SUM(oi.quantity) FROM order_items oi WHERE oi.order_id = o.id) as item_count,
			o.status
		FROM orders o
		LEFT JOIN order_details o_d ON o.id = o_d.order_id
		WHERE 1=1
	`

	args := []interface{}{}

	argId := 1
	// --- Dynamicky přidáváme podmínky pro každý filtr, POKUD není prázdný ---
	if filters.ID != "" {
		// ILIKE je jako LIKE, ale nerozlišuje velká/malá písmena
		query += fmt.Sprintf(" AND o.number ILIKE $%d", argId)
		args = append(args, "%"+filters.ID+"%") // Hledáme text kdekoliv v řetězci
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

	// Speciální logika pro adresu - hledáme ve více sloupcích
	if filters.Address != "" {
		// Hledáme zadaný text buď v ulici, městě, NEBO v názvu výdejního místa
		// Závorky jsou důležité pro správné pořadí operací
		addressPattern := "%" + filters.Address + "%"
		query += fmt.Sprintf(" AND (o.address_street ILIKE $%d OR o.address_city ILIKE $%d OR o_d.pickup_name ILIKE $%d)", argId, argId+1, argId+2)
		args = append(args, addressPattern, addressPattern, addressPattern)
		argId += 3
	}

	// --- BEZPEČNÉ PŘIDÁNÍ ŘAZENÍ (ORDER BY) ---
	allowedSortBy := map[string]string{
		"date":      "o_d.created_at",
		"price":     "o_d.total_czk",
		"itemCount": "item_count",
	}
	sortByColumn, ok := allowedSortBy[filters.SortBy]
	if !ok {
		sortByColumn = "o_d.created_at"
	}

	sortOrder := "DESC"
	if strings.ToUpper(filters.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortByColumn, sortOrder)
	/// Spustíme finální dotaz

	rows, err := ekc_db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Načteme výsledky
	var orders []Order
	for rows.Next() {
		var o Order
		// Pořadí musí přesně odpovídat pořadí v SELECT
		if err := rows.Scan(
			&o.ID, &o.CreatedAt, &o.FirstName, &o.LastName, &o.Email, &o.DeliveryAddress,
			&o.PaymentMethod, &o.DeliveryMethod, &o.TotalPrice, &o.ItemCount, &o.Status,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if len(orders) == 0 {
		// ...vrátíme náš nový, specifický error.
		return nil, apperrors.ErrOrdersNotFound
	}
	return orders, nil
}
