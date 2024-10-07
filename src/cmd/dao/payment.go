package dao

import (
	"context"

	"github.com/tdb/ms-payment/src/cmd/enum"
	"go.mongodb.org/mongo-driver/mongo"
)

type Payment struct {
	ID            string             `json:"id" bson:"_id"`
	Payer         string             `json:"payer" bson:"payer"`
	Amount        float64            `json:"amount" bson:"amount"`
	PaymentMethod enum.PaymentMethod `json:"payment_method" bson:"payment_method"`
	CreatedAt     int64              `json:"created_at" bson:"created_at"`
}

type PaymentDao interface {
	RegisterPayment(payment *Payment) error
}

type paymentDao = baseDao

func NewPaymentDao(client *mongo.Client) PaymentDao {
	return &paymentDao{
		client: client,
	}
}

const PAYMENT_COLLECTION = "payments"

func (dao *paymentDao) RegisterPayment(payment *Payment) error {
	ctx := context.TODO()
	collection := dao.client.Database(DATABASE_NAME).Collection(PAYMENT_COLLECTION)

	_, err := collection.InsertOne(ctx, payment)
	if err != nil {
		return err
	}

	return nil
}
