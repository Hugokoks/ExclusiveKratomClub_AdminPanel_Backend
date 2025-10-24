package db

import (
	"AdminPanelAPI/apperrors"
	"context"
	"errors"

	ekc_db "github.com/Hugokoks/kratomclub-go-common/db"
	"github.com/jackc/pgx/v5"
)



type OrderDetail struct {
	ID               int64
	EkcID           string
	Status           string
	CreatedAt        string
	CustomerFirstName string
	CustomerLastName  string
	CustomerEmail     string
	CustomerPhone     string
	DeliveryAddress  string
	PaymentMethod    string
	DeliveryMethod   string
	DiscountCzk     float64
	ShippingCzk     float64
	SubtotalCzk     float64
	PaymentFeeCzk   float64
	TotalCzk       float64
	WeightGrams    int
	Note           string
	Items		  []OrderItemDetail
}
type OrderItemDetail struct {
	ProductName   string
	VariantLabel string
	ImageSrc     string
	UnitPrice    float64
	Quantity     int
	TotalPrice   float64
}



func SelectOrderDetail(ctx context.Context, orderID string)(*OrderDetail,error){
	query := `
		select 
			o.id,
			o.number,
			o.status,
			TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at, 
			o.customer_first_name,
			o.customer_last_name,
			o.customer_email,
			o.customer_phone,
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
			o.discount_czk,
			o.shipping_czk,
			o.subtotal_czk,
			o.payment_fee_czk,
			o.total_czk,
			o.weight_grams,
			o.note
		from orders o where o.number=$1
	`
	row := ekc_db.Pool.QueryRow(ctx, query, orderID) 

	var o OrderDetail
	err := row.Scan(
		&o.ID,
		&o.EkcID,
		&o.Status,
		&o.CreatedAt,
		&o.CustomerFirstName,
		&o.CustomerLastName,
		&o.CustomerEmail,
		&o.CustomerPhone,
		&o.DeliveryAddress,
		&o.PaymentMethod,
		&o.DeliveryMethod,
		&o.DiscountCzk,
		&o.ShippingCzk,
		&o.SubtotalCzk,
		&o.PaymentFeeCzk,
		&o.TotalCzk,
		&o.WeightGrams,
		&o.Note,
	)

	if err != nil {
		
		////order neexistuje podle order nubmer
		if errors.Is(err,pgx.ErrNoRows){

			return nil, apperrors.ErrOrdersNotFound
		}
		/////neocekavana db chyba
		return nil, err
	}

	query = `
	SELECT p.product_name,pv.variant_label,pv.img_src,oi.unit_price,oi.quantity,oi.total_price FROM order_items oi 
		join products p on oi.product_id = p.id
		join product_variants pv on pv.id = oi.variant_id
	WHERE oi.order_id = $1
	`
	rows, err := ekc_db.Pool.Query(ctx, query, o.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItemDetail
	for rows.Next() {
		var item OrderItemDetail
		if err := rows.Scan(
			&item.ProductName,
			&item.VariantLabel,
			&item.ImageSrc,
			&item.UnitPrice,
			&item.Quantity,
			&item.TotalPrice,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	o.Items = items
	return &o, nil
}