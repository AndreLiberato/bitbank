package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AndreLiberato/bitbank/internal/domain"
	_ "modernc.org/sqlite"
)

type AccountRepository interface {
	FindAll() ([]domain.Account, error)
	FindByNumber(number string) (*domain.Account, error)
	Save(account domain.Account) error
	Update(account domain.Account) error
	Exists(number string) (bool, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository() (AccountRepository, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".bitbank")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create data directory: %w", err)
	}
	db, err := sql.Open("sqlite", filepath.Join(dir, "bitbank.db"))
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &accountRepository{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		number  TEXT PRIMARY KEY,
		balance REAL NOT NULL DEFAULT 0
	)`)
	if err != nil {
		return fmt.Errorf("could not create accounts table: %w", err)
	}
	return nil
}

func (r *accountRepository) FindAll() ([]domain.Account, error) {
	rows, err := r.db.Query(`SELECT number, balance FROM accounts`)
	if err != nil {
		return nil, fmt.Errorf("could not query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.Number, &a.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *accountRepository) FindByNumber(number string) (*domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRow(`SELECT number, balance FROM accounts WHERE number = ?`, number).
		Scan(&a.Number, &a.Balance)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not query account: %w", err)
	}
	return &a, nil
}

func (r *accountRepository) Save(account domain.Account) error {
	_, err := r.db.Exec(`INSERT INTO accounts (number, balance) VALUES (?, ?)`,
		account.Number, account.Balance)
	if err != nil {
		return fmt.Errorf("could not save account: %w", err)
	}
	return nil
}

func (r *accountRepository) Update(account domain.Account) error {
	result, err := r.db.Exec(`UPDATE accounts SET balance = ? WHERE number = ?`,
		account.Balance, account.Number)
	if err != nil {
		return fmt.Errorf("could not update account: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account %s not found", account.Number)
	}
	return nil
}

func (r *accountRepository) Exists(number string) (bool, error) {
	account, err := r.FindByNumber(number)
	if err != nil {
		return false, err
	}
	return account != nil, nil
}
