package services

import (
	"errors"
	"time"

	"github.com/tdb/ms-payment/src/cmd/dao"
	"github.com/tdb/ms-payment/src/cmd/enum"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentInput struct {
	Amount        float64            `json:"amount"`
	PaymentMethod enum.PaymentMethod `json:"payment_method"`
	PixKey        string             `json:"pix_key,omitempty"`
	CreditCard    *dao.CreditCard    `json:"credit_card,omitempty"`
}

type PaymentService interface {
	RegisterPayment(paymentInput PaymentInput) error
}

type paymentService struct {
	accountDao dao.AccountDao
	paymentDao dao.PaymentDao
}

func NewPaymentService(client *mongo.Client) PaymentService {
	return &paymentService{
		accountDao: dao.NewAccountDao(client),
		paymentDao: dao.NewPaymentDao(client),
	}
}

func (s *paymentService) RegisterPayment(payment PaymentInput) error {
	switch payment.PaymentMethod {
	case enum.PIX:
		if payment.PixKey == "" {
			return errors.New("pix key is required")
		}

		account, err := s.accountDao.FindAccountByPixKey(payment.PixKey)
		if err != nil {
			return errors.New("account not found")
		}

		balance, err := s.accountDao.CheckBalance(account.ID)
		if err != nil {
			return err
		}

		if balance < payment.Amount {
			return errors.New("insufficient balance")
		}

		err = s.paymentDao.RegisterPayment(&dao.Payment{
			Payer:         account.Number,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			CreatedAt:     time.Now().Unix(),
		})

		if err != nil {
			return err
		}

		err = s.accountDao.UpdateBalance(account.ID, balance-payment.Amount)
		if err != nil {
			return err
		}

		return nil
	case enum.CREDIT_CARD:
		if payment.CreditCard == nil {
			return errors.New("credit card is required")
		}

		account, err := s.accountDao.FindAccountByCreditCard(payment.CreditCard.Number)
		if err != nil {
			return err
		}

		balance, err := s.accountDao.CheckBalance(account.ID)
		if err != nil {
			return err
		}

		if balance < payment.Amount {
			return errors.New("insufficient balance")
		}

		err = s.paymentDao.RegisterPayment(&dao.Payment{
			Payer:         account.Number,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			CreatedAt:     time.Now().Unix(),
		})

		if err != nil {
			return err
		}

		err = s.accountDao.UpdateBalance(account.ID, balance-payment.Amount)
		if err != nil {
			return err
		}

		return nil
	case enum.CASH:
		return nil
	default:
		return errors.New("invalid payment method")
	}
}
