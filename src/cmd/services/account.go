package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tdb/ms-payment/src/cmd/dao"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountInput struct {
	Name       string         `json:"name"`
	PixKey     string         `json:"pix_key"`
	CreditCard dao.CreditCard `json:"credit_card"`
}

type AccountService interface {
	CreateAccount(account AccountInput) error
}

type accountService struct {
	accountDao dao.AccountDao
}

func NewAccountService(client *mongo.Client) AccountService {
	return &accountService{
		accountDao: dao.NewAccountDao(client),
	}
}

func (s *accountService) CreateAccount(account AccountInput) error {
	if account.PixKey == "" && account.CreditCard.Number == "" {
		return errors.New("pix key or credit card number is required")
	}

	if account.Name == "" {
		return errors.New("name is required")
	}

	accountModel := dao.Account{
		Name:       account.Name,
		Number:     uuid.New().String(),
		PixKey:     account.PixKey,
		CreditCard: &account.CreditCard,
		Balance:    0,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	return s.accountDao.CreateAccount(&accountModel)
}
