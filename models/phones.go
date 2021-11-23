package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	insertPhoneQuery = `
		INSERT INTO phones (number, code, created_at, updated_at, expired_at) VALUES (
			:number,
			:code,
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			:expired_at
		) RETURNING *
	`

	getPhoneByNumberCodeQuery = `
		SELECT * FROM phones WHERE number = $1 AND code = $2 ORDER BY ASC LIMIT 1
	`

	updatePhoneValidateQuery = `
		UPDATE phones SET validate_at = CURRENT_TIMESTAMP AND expired_at = CURRENT_TIMESTAMP WHERE number = :number AND code = :code RETURNING *
	`
)

type Phones struct {
	ID          uint64       `json:"id" db:"id" redis:"id"`
	Number      string       `json:"number" db:"number" redis:"number"`
	Code        string       `json:"code" db:"code" redis:"code"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at" redis:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at" redis:"updated_at"`
	ExpiredAt   time.Time    `json:"expired_at" db:"expired_at" redis:"expired_at"`
	ValidatedAt time.Time    `json:"validated_at" db:"validated_at" redis:"validated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" db:"deleted_at" redis:"deleted_at"`
}

func InsertPhoneDB(number string, code string, tx *sqlx.Tx) (*Phones, error) {
	params := map[string]interface{}{
		"number":     number,
		"code":       code,
		"expired_at": time.Now().Add(3 * time.Minute),
	}

	query, args, _ := sqlx.Named(insertPhoneQuery, params)
	query = tx.Rebind(query)

	var phones Phones

	err := tx.QueryRowx(query, args...).StructScan(&phones)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &phones, nil
}

func GetPhoneByNumberDB(number string, code string, tx *sqlx.Tx) (*Phones, error) {
	var phones Phones
	err := tx.Select(&phones, getPhoneByNumberCodeQuery, number, code)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &phones, nil
}

func UpdatePhoneValidateDB(number string, code string, tx *sqlx.Tx) (*Phones, error) {
	params := map[string]interface{}{
		"number": number,
		"code":   code,
	}

	query, args, _ := sqlx.Named(updatePhoneValidateQuery, params)
	query = tx.Rebind(query)

	var phones Phones

	err := tx.QueryRowx(query, args...).StructScan(&phones)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &phones, nil
}
