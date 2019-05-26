package mongo

import (
	"context"
	"fmt"

	"github.com/getaceres/payment-demo/payment"
	"github.com/getaceres/payment-demo/persistence"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	paymentCollectionName = "payments"
)

type MongoPayment struct {
	ID      string          `json:"_id" bson:"_id"`
	Payment payment.Payment `json:"payment"`
}

type MongoPaymentRepository struct {
	collection                  *mongo.Collection
	defaultFindAndUpdateOptions *options.FindOneAndUpdateOptions
}

func NewMongoPaymentRepository(connectionURI, database string) (*MongoPaymentRepository, error) {
	var result MongoPaymentRepository
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionURI), nil)
	if err != nil {
		return &result, fmt.Errorf("Error creating MongoDB client: %s", err.Error())
	}
	result.collection = client.Database(database).Collection(paymentCollectionName)
	result.defaultFindAndUpdateOptions = options.FindOneAndUpdate().SetReturnDocument(options.After)
	return &result, nil
}

func (m *MongoPaymentRepository) findAndDo(ID string, function func(bson.M, *payment.Payment) *mongo.SingleResult, toUpdate *payment.Payment, action string) (payment.Payment, error) {
	var result MongoPayment
	err := function(bson.M{"_id": ID}, toUpdate).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result.Payment, persistence.NotFoundError{
				ElementType: persistence.PaymentElementType,
				ID:          ID,
			}
		}
		return result.Payment, fmt.Errorf("Error %s payment %s: %s", action, ID, err.Error())
	}
	return result.Payment, nil
}

func (m *MongoPaymentRepository) AddPayment(pay payment.Payment) (payment.Payment, error) {
	pay.ID = uuid.New().String()
	result, err := m.collection.InsertOne(context.Background(), MongoPayment{
		ID:      pay.ID,
		Payment: pay,
	})
	if err != nil {
		return payment.Payment{}, fmt.Errorf("Error saving payment: %s", err.Error())
	}
	pay.ID = result.InsertedID.(string)
	return pay, nil
}

func (m *MongoPaymentRepository) UpdatePayment(pay payment.Payment) (payment.Payment, error) {
	return m.findAndDo(pay.ID, func(filter bson.M, pay *payment.Payment) *mongo.SingleResult {
		return m.collection.FindOneAndReplace(context.Background(), filter, MongoPayment{
			ID:      pay.ID,
			Payment: *pay,
		}, options.FindOneAndReplace().SetReturnDocument(options.After))
	}, &pay, "updating")
}

func (m *MongoPaymentRepository) DeletePayment(id string) (payment.Payment, error) {
	return m.findAndDo(id, func(filter bson.M, pay *payment.Payment) *mongo.SingleResult {
		return m.collection.FindOneAndDelete(context.Background(), filter)
	}, nil, "deleting")
}

func (m *MongoPaymentRepository) GetPayment(id string) (payment.Payment, error) {
	return m.findAndDo(id, func(filter bson.M, pay *payment.Payment) *mongo.SingleResult {
		return m.collection.FindOne(context.Background(), filter)
	}, nil, "getting")
}

func (m *MongoPaymentRepository) GetPayments(filter map[string]string) ([]payment.Payment, error) {
	result := make([]payment.Payment, 0)
	cursor, err := m.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return result, fmt.Errorf("Error getting all payments: %s", err.Error())
	}

	var decoded MongoPayment
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&decoded)
		if err != nil {
			return result, fmt.Errorf("Error decoding result: %s", err.Error())
		}
		result = append(result, decoded.Payment)
	}
	return result, nil
}
