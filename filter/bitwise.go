package filter

import "go.mongodb.org/mongo-driver/bson"

func BitsAllClear(v interface{}) bson.D {
	return bson.D{{Key: "$bitsAllClear", Value: v}}
}

func BitsAllSet(v interface{}) bson.D {
	return bson.D{{Key: "$bitsAllSet", Value: v}}
}

func BitsAnyClear(v interface{}) bson.D {
	return bson.D{{Key: "$bitsAnyClear", Value: v}}
}

func BitsAnySet(v interface{}) bson.D {
	return bson.D{{Key: "$bitsAnySet", Value: v}}
}
