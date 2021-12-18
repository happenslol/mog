package filter

import "go.mongodb.org/mongo-driver/bson"

func All(v ...interface{}) bson.D {
	return bson.D{{Key: "$all", Value: v}}
}

func ElemMatch(v interface{}) bson.D {
	return bson.D{{Key: "$elemMatch", Value: v}}
}

func Size(size uint64) bson.D {
	return bson.D{{Key: "$size", Value: size}}
}
