package templates

import "text/template"

const tmplRaw = `package [[ .PackageName ]]

import (
	"reflect"
	"context"

	mogutil "github.com/happenslol/mog/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	[[ range $imp := .Imports ]]
	"[[ $imp ]]"
	[[ end ]]
)

// Fieldname methods

[[ range $strct := .Structs ]]
type _[[ $strct.Name ]] struct{}

var [[ $strct.Name ]] = new(_[[ $strct.Name ]])

[[ range $fld := $strct.Fields ]]
func (_[[ $strct.Name ]]) [[ $fld.Name ]](filter interface{}) bson.D {
	return bson.D{{
		Key: "[[ $fld.BsonKey ]]",
		Value: filter}}
}
[[ end ]]
[[ end ]]

// Collections

[[ range $name, $col := .Collections ]]
type [[ $col.CollectionType ]] struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func New[[ $col.CollectionType ]](db *mongo.Database,
	options ...*options.CollectionOptions) *[[ $col.CollectionType ]] {
	col := db.Collection("[[ $name ]]", options...)
	return &[[ $col.CollectionType ]]{db, col}
}

func (c *[[ $col.CollectionType ]]) EnsureIndices(
	idxs []mongo.IndexModel, timeout time.Duration,
) error {
	return mogutil.EnsureIndices(c.collection, idxs, timeout)
}

func (c *[[ $col.CollectionType ]]) Clone(opts ...*options.CollectionOptions) (*[[ $col.CollectionType ]], error) {
	cloned, err := c.collection.Clone(opts...)
	if err != nil {
		return nil, err
	}

	return &[[ $col.CollectionType ]]{c.database, cloned}, nil
}

func (c *[[ $col.CollectionType ]]) Name() string {
	return "[[ $name ]]"
}

func (c *[[ $col.CollectionType ]]) Database() *mongo.Database {
	return c.database
}

func (c *[[ $col.CollectionType ]]) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return c.collection.BulkWrite(ctx, models, opts...)
}

func (c *[[ $col.CollectionType ]]) InsertOne(ctx context.Context, document *[[ $col.ModelType ]],
	opts ...*options.InsertOneOptions) ([[ $col.IDType ]], error) {
	result, err := c.collection.InsertOne(ctx, document, opts...)
	if err != nil {
		return "", err
	}

	return result.InsertedID.([[ $col.IDType ]]), nil
}

func (c *[[ $col.CollectionType ]]) InsertMany(ctx context.Context, documents []*[[ $col.ModelType ]],
	opts ...*options.InsertManyOptions) ([][[ $col.IDType ]], error) {
	interfaceDocs := make([]interface{}, len(documents))
	for i, d := range documents {
		interfaceDocs[i] = d
	}

	result, err := c.collection.InsertMany(ctx, interfaceDocs, opts...)
	if err != nil {
		return nil, err
	}

	inserted := make([][[ $col.IDType ]], len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		inserted[i] = id.([[ $col.IDType ]])
	}

	return inserted, nil
}

func (c *[[ $col.CollectionType ]]) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	result, err := c.collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (c *[[ $col.CollectionType ]]) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (int64, error) {
	result, err := c.collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (c *[[ $col.CollectionType ]]) UpdateByID(ctx context.Context, id [[ $col.IDType ]], update interface{},
	opts ...*options.UpdateOptions) ([[ $col.IDType ]], error) {
	result, err := c.collection.UpdateByID(ctx, id, update, opts...)
	if err != nil {
		return "", err
	}

	return result.UpsertedID.([[ $col.IDType ]]), nil
}

func (c *[[ $col.CollectionType ]]) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) ([[ $col.IDType ]], error) {
	result, err := c.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return "", err
	}

	return result.UpsertedID.([[ $col.IDType ]]), nil
}

func (c *[[ $col.CollectionType ]]) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (c *[[ $col.CollectionType ]]) ReplaceOne(ctx context.Context, filter interface{},
	replacement *[[ $col.ModelType ]], opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (c *[[ $col.CollectionType ]]) Aggregate(ctx context.Context, pipeline interface{},
	v []interface{}, opts ...*options.AggregateOptions) ([]interface{}, error) {
	cur, err := c.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	results := v
	arrayType := reflect.TypeOf(v).Elem()
	for cur.Next(ctx) {
		elem := reflect.New(arrayType).Interface()
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (c *[[ $col.CollectionType ]]) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	return c.collection.CountDocuments(ctx, filter, opts...)
}

func (c *[[ $col.CollectionType ]]) EstimatedDocumentCount(ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return c.collection.EstimatedDocumentCount(ctx, opts...)
}

func (c *[[ $col.CollectionType ]]) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions) ([]interface{}, error) {
	return c.collection.Distinct(ctx, fieldName, filter, opts...)
}

func (c *[[ $col.CollectionType ]]) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) ([]*[[ $col.ModelType ]], error) {
	results := []*[[ $col.ModelType ]]{}
	cur, err := c.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem [[ $col.ModelType ]]
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	return results, nil
}

func (c *[[ $col.CollectionType ]]) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) (*[[ $col.ModelType ]], error) {
	var result [[ $col.ModelType ]]

	if err := c.collection.
		FindOne(ctx, filter).
		Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *[[ $col.CollectionType ]]) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) (*[[ $col.ModelType ]], error) {
	var result [[ $col.ModelType ]]

	if err := c.collection.
		FindOneAndDelete(ctx, filter, opts...).
		Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *[[ $col.CollectionType ]]) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement *[[ $col.ModelType ]], opts ...*options.FindOneAndReplaceOptions,
) (*[[ $col.ModelType ]], error) {
	var result [[ $col.ModelType ]]

	if err := c.collection.
		FindOneAndReplace(ctx, filter, replacement, opts...).
		Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *[[ $col.CollectionType ]]) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions,
) (*[[ $col.ModelType ]], error) {
	var result [[ $col.ModelType ]]

	if err := c.collection.
		FindOneAndUpdate(ctx, filter, update, opts...).
		Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *[[ $col.CollectionType ]]) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return c.collection.Watch(ctx, pipeline, opts...)
}

func (c *[[ $col.CollectionType ]]) Indexes() mongo.IndexView {
	return c.collection.Indexes()
}

func (c *[[ $col.CollectionType ]]) Drop(ctx context.Context) error {
	return c.collection.Drop(ctx)
}
[[ end ]]
`

// We need to change the delimiters since we use {{}} in things like
// bson.D arrays
var Tmpl = template.Must(
	template.New("mog-codegen").Delims("[[", "]]").Parse(tmplRaw))