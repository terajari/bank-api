package repository

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/terajari/bank-api/model"
)

func TestCreateAccount(t *testing.T) {
	type args struct {
		ctx     context.Context
		account model.Accounts
	}

	test := []struct {
		name    string
		args    args
		actual  func(sqlmock.Sqlmock)
		want    model.Accounts
		wantErr bool
	}{
		{
			name: "success create account",
			args: args{
				ctx:     context.TODO(),
				account: model.Accounts{ID: "testID", Owner: "testOwner", Balance: 20000, Currency: "IDR"},
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (id, owner, balance, currency) VALUES ($1, $2, $3, $4) RETURNING id, owner, balance, currency, created_at")).
					WithArgs("testID", "testOwner", 20000, "IDR").
					WillReturnRows(s.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
						AddRow("testID", "testOwner", 20000, "IDR", time.Time{}))
			},
			want:    model.Accounts{ID: "testID", Owner: "testOwner", Balance: 20000, Currency: "IDR"},
			wantErr: false,
		},
		{
			name: "failed create account",
			args: args{
				ctx:     context.TODO(),
				account: model.Accounts{ID: "testID", Owner: "testOwner", Balance: 20000, Currency: "IDR"},
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (id, owner, balance, currency) VALUES ($1, $2, $3, $4) RETURNING id, owner, balance, currency, created_at")).
					WithArgs("testID", "testOwner", 20000, "IDR").
					WillReturnError(errors.New("failed"))
			},
			want:    model.Accounts{},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.actual(mock)

			r := NewAccountsRepository(sqlx.NewDb(db, "sqlmock"))
			got, err := r.Create(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	test := []struct {
		name    string
		args    args
		actual  func(sqlmock.Sqlmock)
		want    model.Accounts
		wantErr bool
	}{
		{
			name: "success to get account by id",
			args: args{
				ctx: context.TODO(),
				id:  "testID",
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT id, owner, balance, currency FROM accounts WHERE id = $1 LIMIT 1")).
					WithArgs("testID").
					WillReturnRows(s.NewRows([]string{"id", "owner", "balance", "currency"}).
						AddRow("testID", "testOwner", 20000, "IDR"))
			},
			want:    model.Accounts{ID: "testID", Owner: "testOwner", Balance: 20000, Currency: "IDR"},
			wantErr: false,
		},
		{
			name: "failed to get account by id",
			args: args{
				ctx: context.TODO(),
				id:  "testID",
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT id, owner, balance, currency FROM accounts WHERE id = $1")).
					WithArgs("testID").
					WillReturnError(errors.New("failed"))
			},
			want:    model.Accounts{},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.actual(mock)

			r := NewAccountsRepository(sqlx.NewDb(db, "sqlmock"))
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	type args struct {
		ctx     context.Context
		account model.Accounts
	}

	test := []struct {
		name    string
		args    args
		actual  func(sqlmock.Sqlmock)
		want    model.Accounts
		wantErr bool
	}{
		{
			name: "success to update account",
			args: args{
				ctx:     context.TODO(),
				account: model.Accounts{ID: "testID", Owner: "testOwner", Balance: 50000, Currency: "IDR"},
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("UPDATE accounts SET balance = $2 WHERE id = $1 RETURNING id, owner, balance, currency, created_at")).
					WithArgs("testID", 50000).
					WillReturnRows(s.NewRows([]string{"id", "owner", "balance", "currency", "created_at"}).
						AddRow("testID", "testOwner", 50000, "IDR", time.Time{}))
			},
			want:    model.Accounts{ID: "testID", Owner: "testOwner", Balance: 50000, Currency: "IDR"},
			wantErr: false,
		},

		{
			name: "failed to update account",
			args: args{
				ctx:     context.TODO(),
				account: model.Accounts{ID: "testID", Owner: "testOwner", Balance: 50000, Currency: "IDR"},
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("UPDATE accounts SET balance = $2 WHERE id = $1 RETURNING id, owner, balance, currency, created_at")).
					WithArgs("testID", 50000).
					WillReturnError(errors.New("failed"))
			},
			want:    model.Accounts{},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.actual(mock)

			r := NewAccountsRepository(sqlx.NewDb(db, "sqlmock"))
			updatedAcc, err := r.Update(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if updatedAcc != tt.want {
				t.Errorf("Update() got = %v, want %v", updatedAcc, tt.want)
			}
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	test := []struct {
		name    string
		args    args
		actual  func(sqlmock.Sqlmock)
		want    model.Accounts
		wantErr bool
	}{
		{
			name: "success to delete account",
			args: args{
				ctx: context.TODO(),
				id:  "testID",
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectExec(regexp.QuoteMeta("DELETE FROM accounts WHERE id = $1")).
					WithArgs("testID").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "failed to delete account",
			args: args{
				ctx: context.TODO(),
				id:  "testID",
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectExec(regexp.QuoteMeta("DELETE FROM accounts WHERE id = $1")).
					WithArgs("testID").
					WillReturnError(errors.New("failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.actual(mock)

			r := NewAccountsRepository(sqlx.NewDb(db, "sqlmock"))

			err = r.Delete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestListAccounts(t *testing.T) {
	type args struct {
		ctx    context.Context
		owner  string
		limit  int
		offset int
	}

	test := []struct {
		name    string
		args    args
		actual  func(sqlmock.Sqlmock)
		want    []model.Accounts
		wantErr bool
	}{
		{
			name: "success to list accounts",
			args: args{
				ctx:    context.TODO(),
				owner:  "testOwner",
				limit:  10,
				offset: 0,
			},
			actual: func(s sqlmock.Sqlmock) {
				rows := s.NewRows([]string{"id", "owner", "balance", "currency"}).
					AddRow("testId", "testOwner", 50000, "IDR").
					AddRow("testId", "testOwner", 4, "USD")

				s.ExpectQuery(regexp.QuoteMeta("SELECT id, owner, balance, currency FROM accounts WHERE owner = $1 LIMIT $2 OFFSET $3 ORDER BY id")).
					WillReturnRows(rows)
			},
			want: []model.Accounts{{
				ID:       "testId",
				Owner:    "testOwner",
				Balance:  50000,
				Currency: "IDR",
			},
				{
					ID:       "testId",
					Owner:    "testOwner",
					Balance:  4,
					Currency: "USD",
				},
			},
			wantErr: false,
		},
		{
			name: "failed to list accounts",
			args: args{
				ctx:    context.TODO(),
				owner:  "testOwner",
				limit:  10,
				offset: 0,
			},
			actual: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT id, owner, balance, currency FROM accounts WHERE owner = $1 LIMIT $2 OFFSET $3 ORDER BY id")).
					WillReturnError(errors.New("failed"))
			},
			want:    []model.Accounts{},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.actual(mock)

			r := NewAccountsRepository(sqlx.NewDb(db, "sqlmock"))

			accounts, err := r.List(tt.args.ctx, tt.args.owner, tt.args.limit, tt.args.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(accounts, tt.want) {
				t.Errorf("List() got = %v, want %v", accounts, tt.want)
				return
			}
		})
	}
}
