package filter

import "go.mongodb.org/mongo-driver/bson"

func Not(v interface{}) bson.D {
	return bson.D{{Key: "$not", Value: v}}
}

func And(v ...interface{}) bson.D {
	return bson.D{{Key: "$and", Value: v}}
}

func Or(v ...interface{}) bson.D {
	return bson.D{{Key: "$or", Value: v}}
}

func Nor(v ...interface{}) bson.D {
	return bson.D{{Key: "$nor", Value: v}}
}
