package account

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type Repository interface {
	Close() error
	PutAccount(ctx context.Context, a Account) error
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func (r *PostgresRepository) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) PutAccount(ctx context.Context, a Account) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO account(id, name) VALUES($1, $2)",
		a.ID,
		a.Name,
	)
	return err
}

func (r *PostgresRepository) GetAccountById(ctx context.Context, id string) (*Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT a.id, a.name FROM accounts a WHERE a.id = $1",
		id,
	)

	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}

	return a, nil
}
func (r *PostgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT a.id, a.name FROM accounts a ORDER BY DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		a := &Account{}
		if err := rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}
