package repository

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BalanceViweRepositoryDB struct {
	*mongo.Collection
}

func NewBalanceViweRepositoryDB(db *mongo.Database) BalanceViewRepository {
	collection := db.Collection("balanceviews")
	mod := mongo.IndexModel{
		Keys: bson.M{
			"accountID": 1, // index in ascending order
		},
	}
	ind, err := collection.Indexes().CreateOne(context.Background(), mod)

	// Check if the CreateOne() method returned any errors
	if err != nil {
		fmt.Println("Indexes().CreateOne() ERROR:", err)
		os.Exit(1) // exit in case of error
	} else {
		// API call returns string of the index name
		fmt.Println("CreateOne() index:", ind)
	}

	return &BalanceViweRepositoryDB{Collection: collection}
}

func (r *BalanceViweRepositoryDB) CreateBalanceView(b BalanceView) error {

	_, err := r.Collection.InsertOne(context.TODO(), b)
	if err != nil {
		return err
	}

	return err
}

func (r *BalanceViweRepositoryDB) UpsertBulkBalanceView(balanceView []BalanceView) error {

	var operations []mongo.WriteModel

	for _, b := range balanceView {

		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(bson.M{"accountID": b.AccountID})
		operation.SetUpdate(bson.D{{"$set", bson.D{
			{"accountID", b.AccountID}, {"balance", b.Balance},
		}}})
		operation.SetUpsert(true)
		operations = append(operations, operation)
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)

	_, err := r.Collection.BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		return err
	}

	return err
}

func (r *BalanceViweRepositoryDB) UpdateBalanceViewByID(string, BalanceView) error {

	return nil
}

func (r *BalanceViweRepositoryDB) GetBalanceViewByID(id string) (*BalanceView, error) {
	b := &BalanceView{}

	filter := bson.M{"accountID": id}

	err := r.Collection.FindOne(context.Background(), filter).Decode(b)

	return b, err
}

func (r *BalanceViweRepositoryDB) ClearBalanceView() error {
	filter := bson.D{}

	deleteResult, err := r.Collection.DeleteMany(context.Background(), filter, nil)
	if deleteResult.DeletedCount == 0 {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
