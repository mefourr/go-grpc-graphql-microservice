package order

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"log"
)

type Repository interface {
	Close() error
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func (r *PostgresRepository) Close() error {
	log.Println("Closing database connection.")
	err := r.db.Close()
	if err != nil {
		log.Printf("Error closing database: %v\n", err)
		return err
	}
	log.Println("Database connection closed.")
	return nil
}

func (r *PostgresRepository) PutOrder(ctx context.Context, o Order) error {
	log.Printf("Putting order with ID: %s\n", o.ID)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction for order %s: %v\n", o.ID, err)
		return err
	}

	defer func() {
		if err != nil {
			log.Printf("Rolling back transaction for order %s due to error: %v\n", o.ID, err)
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction for order %s: %v\n", o.ID, err)
		} else {
			log.Printf("Order %s committed successfully.\n", o.ID)
		}
	}()

	_, err = r.db.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountId,
		o.TotalPrice,
	)
	if err != nil {
		log.Printf("Error inserting order %s: %v\n", o.ID, err)
		return err
	}

	start, _ := tx.PrepareContext(
		ctx,
		pq.CopyIn("order_products", "order_id", "product_id", "quantity"),
	)
	for _, prod := range o.Products {
		_, err = start.ExecContext(ctx, o.ID, prod.ID, prod.Quantity)
		if err != nil {
			log.Printf("Error inserting product for order %s: %v\n", o.ID, err)
			return err
		}
	}
	_, err = start.ExecContext(ctx)
	if err != nil {
		log.Printf("Error finalizing product insertion for order %s: %v\n", o.ID, err)
		return err
	}
	log.Printf("Order %s and its products inserted successfully.\n", o.ID)
	return start.Close()
}

func (r *PostgresRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	log.Printf("Fetching orders for account %s\n", accountId)
	query :=
		`SELECT 
    		o.id,
    		o.created_at,
    		o.account_id,
    		o.total_price::money::numeric::float8,
			opr.product_id,
			opr.quantity
		FROM orders o
    	JOIN order_products opr ON o.id = opr.order_id
        WHERE o.account_id = $1
        ORDER BY o.id`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		accountId,
	)
	if err != nil {
		log.Printf("Error executing query for account %s: %v\n", accountId, err)
		return nil, err
	}
	defer rows.Close()

	order := &Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	var orders []Order
	var products []OrderedProduct

	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountId,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			log.Printf("Error scanning row for account %s: %v\n", accountId, err)
			return nil, err
		}

		log.Printf("Processing order ID: %s, Product ID: %s, Quantity: %d\n", order.ID, orderedProduct.ID, orderedProduct.Quantity)

		if lastOrder.ID != "" && lastOrder.ID != order.ID {
			orders = append(orders, Order{
				ID:         lastOrder.ID,
				AccountId:  lastOrder.AccountId,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   products,
			})
			products = nil
		}

		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})

		*lastOrder = *order
	}

	// Final append for the last order
	orders = append(orders, Order{
		ID:         lastOrder.ID,
		AccountId:  lastOrder.AccountId,
		CreatedAt:  lastOrder.CreatedAt,
		TotalPrice: lastOrder.TotalPrice,
		Products:   products,
	})

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows for account %s: %v\n", accountId, err)
		return nil, err
	}

	log.Printf("Fetched %d orders for account %s\n", len(orders), accountId)
	return orders, nil
}

func (r *PostgresRepository) Ping() error {
	log.Println("Pinging database.")
	err := r.db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v\n", err)
		return err
	}
	log.Println("Database ping successful.")
	return nil
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	log.Printf("Connecting to database at %s\n", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Printf("Error opening database connection: %v\n", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v\n", err)
		return nil, err
	}

	log.Println("Database connection successful.")
	return &PostgresRepository{db: db}, nil
}
