package filter

import "go.mongodb.org/mongo-driver/bson"

func Eq(v interface{}) bson.D {
	return bson.D{{Key: "$eq", Value: v}}
}

func Gt(v interface{}) bson.D {
	return bson.D{{Key: "$gt", Value: v}}
}

func Gte(v interface{}) bson.D {
	return bson.D{{Key: "$gte", Value: v}}
}

func In(v ...interface{}) bson.D {
	return bson.D{{Key: "$in", Value: v}}
}

func Lt(v interface{}) bson.D {
	return bson.D{{Key: "$lt", Value: v}}
}

func Lte(v interface{}) bson.D {
	return bson.D{{Key: "$lte", Value: v}}
}

func Ne(v interface{}) bson.D {
	return bson.D{{Key: "$ne", Value: v}}
}

func Nin(v ...interface{}) bson.D {
	return bson.D{{Key: "$nin", Value: v}}
}
