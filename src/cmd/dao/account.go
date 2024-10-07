package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	ID         string      `json:"id" bson:"_id"`
	Number     string      `json:"number" bson:"number"`
	Name       string      `json:"name" bson:"name"`
	Balance    float64     `json:"balance" bson:"balance"`
	PixKey     string      `json:"pix_key,omitempty" bson:"pix_key,omitempty"`
	CreditCard *CreditCard `json:"credit_card,omitempty" bson:"credit_card,omitempty"`
	CreatedAt  int64       `json:"created_at" bson:"created_at"`
	UpdatedAt  int64       `json:"updated_at" bson:"updated_at"`
}

type CreditCard struct {
	Number          string `json:"number" bson:"number"`
	CVV             string `json:"cvv" bson:"cvv"`
	Holder          string `json:"holder" bson:"holder"`
	ExpirationMonth int    `json:"expiration_month" bson:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year" bson:"expiration_year"`
}

type AccountDao interface {
	CreateAccountIndexes()
	FindAccountByPixKey(pixKey string) (*Account, error)
	FindAccountByCreditCard(creditCardNumber string) (*Account, error)
	CheckBalance(accountId string) (float64, error)
	UpdateBalance(accountId string, amount float64) error
	CreateAccount(account *Account) error
}

type accountDao = baseDao

func NewAccountDao(client *mongo.Client) AccountDao {
	return &accountDao{
		client: client,
	}
}

const ACCOUNT_COLLECTION = "accounts"

func (dao *accountDao) CreateAccountIndexes() {
	ctx := context.TODO()
	collection := dao.client.Database(DATABASE_NAME).Collection(ACCOUNT_COLLECTION)

	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"pix_key": 1},
	})

	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"credit_card.number": 1},
	})
}

func (dao *accountDao) FindAccountByPixKey(pixKey string) (*Account, error) {
	ctx := context.TODO()
	collection := dao.client.Database(DATABASE_NAME).Collection(ACCOUNT_COLLECTION)

	var account Account

	err := collection.FindOne(ctx, map[string]string{"pix_key": pixKey}).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (dao *accountDao) FindAccountByCreditCard(creditCardNumber string) (*Account, error) {
	ctx := context.TODO()
	collection := dao.client.Database(DATABASE_NAME).Collection(ACCOUNT_COLLECTION)

	var account Account

	err := collection.FindOne(ctx, map[string]string{"credit_card.number": creditCardNumber}).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (dao *accountDao) CheckBalance(accountId string) (float64, error) {
	ctx := context.TODO()
	collection := dao.client.Database(DATABASE_NAME).Collection(ACCOUNT_COLLECTION)

	var account Account

	err := collection.FindOne(ctx, map[string]string{"_id": accountId}).Decode(&account)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

func (dao *accountDao) UpdateBalance(accountId string, amount float64) error {
	ctx := context.TODO()
	collection := dao.client.Database(DATABASE_NAME).Collection(ACCOUNT_COLLECTION)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": accountId}, bson.M{"$set": bson.M{"balance": amount}})
	if err != nil {
		return err
	}

	return nil
}

func (dao *accountDao) CreateAccount(account *Account) error {
	ctx := context.TODO()
	collection := dao.client.Database(DATABASE_NAME).Collection(ACCOUNT_COLLECTION)

	_, err := collection.InsertOne(ctx, account)
	if err != nil {
		return err
	}

	return nil
}
