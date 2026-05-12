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
		balance REAL NOT NULL DEFAULT 0,
		account_type TEXT NOT NULL DEFAULT 'simples',
		points INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		return fmt.Errorf("could not create accounts table: %w", err)
	}
	if err := addColumnIfMissing(db, "accounts", "account_type", "TEXT NOT NULL DEFAULT 'simples'"); err != nil {
		return err
	}
	if err := addColumnIfMissing(db, "accounts", "points", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	return nil
}

func addColumnIfMissing(db *sql.DB, table, column, definition string) error {
	exists, err := columnExists(db, table, column)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition))
	if err != nil {
		return fmt.Errorf("could not add column %s: %w", column, err)
	}
	return nil
}

func columnExists(db *sql.DB, table, column string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false, fmt.Errorf("could not inspect table %s: %w", table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid, notNull, pk int
		var name, columnType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}

func (r *accountRepository) FindAll() ([]domain.Account, error) {
	rows, err := r.db.Query(`SELECT number, balance, account_type, points FROM accounts`)
	if err != nil {
		return nil, fmt.Errorf("could not query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.Number, &a.Balance, &a.Type, &a.Points); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *accountRepository) FindByNumber(number string) (*domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRow(`SELECT number, balance, account_type, points FROM accounts WHERE number = ?`, number).
		Scan(&a.Number, &a.Balance, &a.Type, &a.Points)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not query account: %w", err)
	}
	return &a, nil
}

func (r *accountRepository) Save(account domain.Account) error {
	_, err := r.db.Exec(`INSERT INTO accounts (number, balance, account_type, points) VALUES (?, ?, ?, ?)`,
		account.Number, account.Balance, account.Type, account.Points)
	if err != nil {
		return fmt.Errorf("could not save account: %w", err)
	}
	return nil
}

func (r *accountRepository) Update(account domain.Account) error {
	result, err := r.db.Exec(`UPDATE accounts SET balance = ?, account_type = ?, points = ? WHERE number = ?`,
		account.Balance, account.Type, account.Points, account.Number)
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
