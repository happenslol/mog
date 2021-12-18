package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func Exists(v bool) bson.D {
	return bson.D{{Key: "$exists", Value: v}}
}

func Type(v ...bsontype.Type) bson.D {
	// TODO: Does this actually make a difference?
	if len(v) == 1 {
		return bson.D{{Key: "$type", Value: v[0]}}
	}

	return bson.D{{Key: "$type", Value: v}}
}
