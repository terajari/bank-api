package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/terajari/bank-api/dto"
	mockrepo "github.com/terajari/bank-api/mock/repository"
	"github.com/terajari/bank-api/model"
)

func TestGetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx context.Context
		id  string
	}

	test := []struct {
		name    string
		args    args
		actual  func(repo *mockrepo.MockAccountsRepository)
		want    dto.GetAccountResponse
		wantErr bool
	}{
		{
			name: "success get account by id",
			args: args{
				ctx: context.TODO(),
				id:  "testId",
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().Get(context.TODO(), "testId").
					Return(
						model.Accounts{
							ID:       "testId",
							Owner:    "testOwner",
							Balance:  1000,
							Currency: "testCurrency",
						},
						nil,
					)
			},
			want: dto.GetAccountResponse{
				Id:       "testId",
				Owner:    "testOwner",
				Balance:  1000,
				Currency: "testCurrency",
			},
			wantErr: false,
		},

		{
			name: "failed to get account by id",
			args: args{
				ctx: context.TODO(),
				id:  "testId",
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().Get(context.TODO(), "testId").
					Return(
						model.Accounts{},
						fmt.Errorf("error"),
					)
			},
			want:    dto.GetAccountResponse{},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockAccounRepo := mockrepo.NewMockAccountsRepository(ctrl)

			uc := NewAccountsUsecase(mockAccounRepo)
			if tt.actual != nil {
				tt.actual(mockAccounRepo)
			}

			acc, err := uc.GetAccount(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() err: %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(acc, tt.want) {
				t.Errorf("GetAccount() got: %v, want: %v", acc, tt.want)
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx context.Context
		dto dto.UpdateAccountRequest
	}

	test := []struct {
		name    string
		args    args
		actual  func(repo *mockrepo.MockAccountsRepository)
		want    dto.UpdateAccountResponse
		wantErr bool
	}{
		{
			name: "success update account",
			args: args{
				ctx: context.TODO(),
				dto: dto.UpdateAccountRequest{
					Id:      "testId",
					Balance: 2000,
				},
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().Get(context.TODO(), "testId").Return(
					model.Accounts{
						ID:       "testId",
						Owner:    "testOwner",
						Balance:  1000,
						Currency: "testCurrency",
					},
					nil,
				)
				repo.EXPECT().Update(context.TODO(), &model.Accounts{
					ID:      "testId",
					Balance: 2000,
				}).Return(model.Accounts{
					ID:       "testId",
					Owner:    "testOwner",
					Balance:  2000,
					Currency: "testCurrency",
				},
					nil,
				)
			},
			want: dto.UpdateAccountResponse{
				Id:       "testId",
				Owner:    "testOwner",
				Balance:  2000,
				Currency: "testCurrency",
			},
			wantErr: false,
		},
		{
			name: "failed to update account",
			args: args{
				ctx: context.TODO(),
				dto: dto.UpdateAccountRequest{
					Id:      "testId",
					Balance: 2000,
				},
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().Get(context.TODO(), "testId").Return(
					model.Accounts{
						ID:       "testId",
						Owner:    "testOwner",
						Balance:  1000,
						Currency: "testCurrency",
					},
					nil,
				)
				repo.EXPECT().Update(context.TODO(), &model.Accounts{
					ID:      "testId",
					Balance: 2000,
				}).Return(model.Accounts{},
					fmt.Errorf("error"),
				)
			},
			want:    dto.UpdateAccountResponse{},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockAccounRepo := mockrepo.NewMockAccountsRepository(ctrl)

			uc := NewAccountsUsecase(mockAccounRepo)
			if tt.actual != nil {
				tt.actual(mockAccounRepo)
			}

			updatedAcc, err := uc.UpdateAccount(tt.args.ctx, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccount() err: %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(updatedAcc, tt.want) {
				t.Errorf("UpdateAccount() got %v, want %v", updatedAcc, t)
			}
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx context.Context
		id  string
	}

	test := []struct {
		name    string
		args    args
		actual  func(repo *mockrepo.MockAccountsRepository)
		wantErr bool
	}{
		{
			name: "success delete account",
			args: args{
				ctx: context.TODO(),
				id:  "testId",
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().Get(context.TODO(), "testId").Return(
					model.Accounts{
						ID:       "testId",
						Owner:    "testOwner",
						Balance:  1000,
						Currency: "testCurrency",
					},
					nil,
				)
				repo.EXPECT().Delete(context.TODO(), "testId").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed to delete account",
			args: args{
				ctx: context.TODO(),
				id:  "testId",
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().Get(context.TODO(), "testId").Return(
					model.Accounts{
						ID:       "testId",
						Owner:    "testOwner",
						Balance:  1000,
						Currency: "testCurrency",
					},
					nil,
				)
				repo.EXPECT().Delete(context.TODO(), "testId").Return(fmt.Errorf("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockAccounRepo := mockrepo.NewMockAccountsRepository(ctrl)

			uc := NewAccountsUsecase(mockAccounRepo)
			if tt.actual != nil {
				tt.actual(mockAccounRepo)
			}

			err := uc.DeleteAccount(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAccount() err: %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx context.Context
		dto dto.ListAccountsRequest
	}

	test := []struct {
		name    string
		args    args
		actual  func(repo *mockrepo.MockAccountsRepository)
		want    []dto.GetAccountResponse
		wantErr bool
	}{
		{
			name: "success list accounts",
			args: args{
				ctx: context.TODO(),
				dto: dto.ListAccountsRequest{
					Owner: "testOwner",
					Page:  1,
					Size:  10,
				},
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().List(context.TODO(), "testOwner", 10, 0).Return(
					[]model.Accounts{
						{
							ID:       "testId",
							Owner:    "testOwner",
							Balance:  1000,
							Currency: "testCurrency1",
						},
						{
							ID:       "testId",
							Owner:    "testOwner",
							Balance:  1000,
							Currency: "testCurrency2",
						},
					},
					nil,
				)
			},
			want: []dto.GetAccountResponse{
				{
					Id:       "testId",
					Owner:    "testOwner",
					Balance:  1000,
					Currency: "testCurrency1",
				},
				{
					Id:       "testId",
					Owner:    "testOwner",
					Balance:  1000,
					Currency: "testCurrency2",
				},
			},
			wantErr: false,
		},
		{
			name: "failed to get list accounts",
			args: args{
				ctx: context.TODO(),
				dto: dto.ListAccountsRequest{
					Owner: "testOwner",
					Page:  1,
					Size:  10,
				},
			},
			actual: func(repo *mockrepo.MockAccountsRepository) {
				repo.EXPECT().List(context.TODO(), "testOwner", 10, 0).Return(
					[]model.Accounts{},
					fmt.Errorf("error"),
				)
			},
			want:    []dto.GetAccountResponse{},
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			mockAccounRepo := mockrepo.NewMockAccountsRepository(ctrl)

			uc := NewAccountsUsecase(mockAccounRepo)
			if tt.actual != nil {
				tt.actual(mockAccounRepo)
			}

			accounts, err := uc.ListAccounts(tt.args.ctx, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAccounts() err: %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(accounts, tt.want) {
				t.Errorf("ListAccounts() got %v, want %v", accounts, tt.want)
			}
		})
	}
}
