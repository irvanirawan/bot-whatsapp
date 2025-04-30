package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	ColumnID         = "id"
	ColumnName       = "name"
	ColumnToken      = "token"
	ColumnWebhook    = "webhook"
	ColumnJID        = "jid"
	ColumnQRCode     = "qrcode"
	ColumnConnected  = "connected"
	ColumnExpiration = "expiration"
	ColumnEvents     = "events"
	ColumnTimeouts   = "timeouts"
)

type DatabaseHandler struct {
	DB *sql.DB
}

func (dbHandler *DatabaseHandler) Initialize() error {
	var err error
	// Initialize SQLite database
	dbHandler.DB, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		return err
	}

	// Create users table if not exists
	sqlStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS users (
		%s INTEGER NOT NULL PRIMARY KEY,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL,
		%s TEXT NOT NULL DEFAULT "", 
		%s TEXT NOT NULL DEFAULT "", 
		%s TEXT NOT NULL DEFAULT "", 
		%s INTEGER, 
		%s INTEGER, 
		%s TEXT NOT NULL DEFAULT "All", 
		%s TEXT
	);`,
		ColumnID, ColumnName, ColumnToken, ColumnWebhook, ColumnJID, ColumnQRCode,
		ColumnConnected, ColumnExpiration, ColumnEvents, ColumnTimeouts)
	_, err = dbHandler.DB.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Insert default user if not exists
	sqlStmt = fmt.Sprintf(`INSERT INTO users(%s, %s, %s, %s, %s, %s, %s, %s)
				SELECT 'irvanbotwa', 'irvanbotwa', 'irvanbotwa', 'irvanbotwa', 'irvanbotwa', 0, 0, 'All'
				WHERE NOT EXISTS (SELECT 1 FROM users WHERE %s = 'irvanbotwa')`,
		ColumnName, ColumnToken, ColumnWebhook, ColumnJID, ColumnQRCode, ColumnConnected, ColumnExpiration, ColumnEvents, ColumnName)
	_, err = dbHandler.DB.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error inserting default user: %v", err)
	}
	return nil
}

func (dbHandler *DatabaseHandler) GetQrCode() (string, error) {
	var qrCode string
	sqlStmt := fmt.Sprintf("SELECT %s FROM users WHERE %s = ? LIMIT 1", ColumnQRCode, ColumnID)
	rows, err := dbHandler.DB.Query(sqlStmt, 1)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&qrCode); err != nil {
			return "", err
		}
	}
	return qrCode, nil
}

func (dbHandler *DatabaseHandler) UpdateConnectionStatus(id int, status int) error {
	sqlStmt := fmt.Sprintf("UPDATE users SET %s = ? WHERE %s = ?", ColumnConnected, ColumnID)
	_, err := dbHandler.DB.Exec(sqlStmt, status, id)
	if err != nil {
		return fmt.Errorf("error updating connection status: %v", err)
	}
	return nil
}

func (dbHandler *DatabaseHandler) UpdateQrCode(qrCode string) error {
	sqlStmt := fmt.Sprintf("UPDATE users SET %s = ? WHERE %s = 1", ColumnQRCode, ColumnID)
	_, err := dbHandler.DB.Exec(sqlStmt, qrCode)
	if err != nil {
		return fmt.Errorf("error saving QR code: %v", err)
	}
	return nil
}
