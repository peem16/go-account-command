package repository

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type accountEventRepositoryDB struct {
	*mongo.Collection
}

func NewAccountEventRepositoryDB(db *mongo.Database) AccountEventRepository {
	collection := db.Collection("accountevents")
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

	return &accountEventRepositoryDB{Collection: collection}
}

func (r *accountEventRepositoryDB) CreateEvent(a AccountEvent) error {
	_, err := r.Collection.InsertOne(context.Background(), a)
	return err
}

func (r *accountEventRepositoryDB) ClearAccount() error {
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

func (r *accountEventRepositoryDB) FindOneAccountEventByID(id string, types []string) (*AccountEvent, error) {

	accountEvent := AccountEvent{}

	filter := bson.D{primitive.E{Key: "accountID", Value: id}}
	err := r.Collection.FindOne(context.Background(), filter).Decode(&accountEvent)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &accountEvent, nil
}

func (r *accountEventRepositoryDB) GetAccountEventByID(id string) ([]AccountEvent, error) {
	accountEvent := []AccountEvent{}

	filter := bson.D{primitive.E{Key: "accountID", Value: id}}

	cursor, err := r.Collection.Find(context.Background(), filter)

	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &accountEvent)

	return accountEvent, err
}

func (r *accountEventRepositoryDB) GetAccountEvents(id []string) ([]AccountEvent, error) {
	accountEvent := []AccountEvent{}

	filter := bson.M{"accountID": bson.M{"$in": id}}
	opts := options.Find().SetSort(bson.D{{"type", 1}})

	cursor, err := r.Collection.Find(context.Background(), filter, opts)

	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &accountEvent)

	return accountEvent, err
}
