package domain

import (
	"capi/errs"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func NewMock() (*sqlx.DB, sqlmock.Sqlmock){
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatalf("An Error '%s' was not Expected when Opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return sqlxDB, mock
}

func TestCustomerRepositoryDB_FindAll(t *testing.T) {
	type args struct {
		status string
	}

	tests := []struct {
		name string
		args args
		want []Customer
		wantErr *errs.AppErr
	}{
		{
			"Success Get All Data Customers",
			args{""},
			[]Customer{
				{"1", "User1", "Tangerang", "15540", "2012-01-01", "1"},
				{"2", "User2", "Jakarta", "15641", "2022-02-12", "1"},
				{"3", "User3", "Bandung", "15542", "2021-06-12", "1"},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T)  {
			db, mock := NewMock()
			repo := NewCustomerRepositoryDB(db)

			rows := mock.NewRows([]string{"customer_id", "name", "city", "zipcode", "date_of_birth", "status"}).AddRow("1", "User1", "Tangerang", "15540", "2012-01-01", "1").AddRow("2", "User2", "Jakarta", "15641", "2022-02-12", "1")

			mock.ExpectQuery(`select \* form customers`).WillReturnRows(rows)
			got, got1 := repo.FindAll(tt.args.status)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomerRepositoryDB.FindAll() got = %v, want %v", got, tt.want)

			}
			if !reflect.DeepEqual(got1, tt.wantErr) {
				t.Errorf("CustomerRepositoryDB.FindAll() got1 = %v, want %v", got1, tt.wantErr)
			}
		})
	}
}

func TestCustomerRepositoryDB_FindByID(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}

	type args struct {
		id string
	}

	tests := []struct {
		name string
		fields fields
		args args
		want *Customer
		want1 *errs.AppErr
	}{

	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			s := CustomerRepositoryDB{
				db: tt.fields.db,
			}
			got, got1 := s.FindByID(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomerRepositoryDB.FindByID() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CustomerRepositoryDB.FindByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}