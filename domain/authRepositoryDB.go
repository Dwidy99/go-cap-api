package domain

import (
	"capi/errs"
	"capi/logger"

	"github.com/jmoiron/sqlx"
)

type AuthRepositoryDB struct {
	db *sqlx.DB
}

func NewAuthRepositoryDB(db *sqlx.DB) AuthRepositoryDB {
	return AuthRepositoryDB{db}
}

func (d AuthRepositoryDB) FindBy(username, password string) (*Login, *errs.AppErr){
	query := `SELECT username, u.customer_id, role, string_agg(a.account_id::varchar(255), ',') as account_numbers FROM users u 
					LEFT JOIN accounts a ON a.customer_id = u.customer_id
					WHERE username = $1 and password = $2 GROUP BY username, a.customer_id, role`
	
	var login Login
	err := d.db.Get(&login, query, username, password)
	if err != nil {
		if err != nil {
			return nil, errs.NewAuthenticationError("Invalid Credential")
		} else {
			logger.Error("Error while verifying login request to Database")
			return nil, errs.NewUnexpectedError("Unexpected Database Error") 
		}
	}

	return &login, nil
}