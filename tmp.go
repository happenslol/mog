package main

import (
	"context"

	"github.com/happenslol/mog/testdata"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersCollection struct {
	name       string
	database   *mongo.Database
	collection *mongo.Collection
}

func NewUsersCollection(db *mongo.Database,
	options ...*options.CollectionOptions) *UsersCollection {
	col := db.Collection("users", options...)
	return &UsersCollection{"users", db, col}
}

func (c *UsersCollection) Clone(opts ...*options.CollectionOptions) (*UsersCollection, error) {
	cloned, err := c.collection.Clone(opts...)
	if err != nil {
		return nil, err
	}

	return &UsersCollection{c.name, c.database, cloned}, nil
}

func (c *UsersCollection) Name() string {
	return c.name
}

func (c *UsersCollection) Database() *mongo.Database {
	return c.database
}

func (c *UsersCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return nil, nil
}

func (c *UsersCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (string, error) {
	return "", nil
}

func (c *UsersCollection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (string, error) {
	return "", nil
}

func (c *UsersCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	return 0, nil
}

func (c *UsersCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	return 0, nil
}

func (c *UsersCollection) UpdateByID(ctx context.Context, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (c *UsersCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (c *UsersCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (c *UsersCollection) ReplaceOne(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (c *UsersCollection) Aggregate(ctx context.Context, pipeline interface{},
	v interface{}, opts ...*options.AggregateOptions) error {
	return nil
}

func (c *UsersCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	return 0, nil
}

func (c *UsersCollection) EstimatedDocumentCount(ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return 0, nil
}

func (c *UsersCollection) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions) ([]*testdata.Author, error) {
	return nil, nil
}

func (c *UsersCollection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) ([]*testdata.Author, error) {
	return nil, nil
}

func (c *UsersCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *testdata.Author {
	return nil
}

func (c *UsersCollection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *testdata.Author {
	return nil
}

func (c *UsersCollection) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *testdata.Author {
	return nil
}

func (c *UsersCollection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *testdata.Author {
	return nil
}

func (c *UsersCollection) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return nil, nil
}

func (c *UsersCollection) Indexes() mongo.IndexView {
	return c.collection.Indexes()
}

func (c *UsersCollection) Drop(ctx context.Context) error {
	return c.collection.Drop(ctx)
}
