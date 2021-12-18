package filter

import "go.mongodb.org/mongo-driver/bson"

func GeoIntersects(v interface{}) bson.D {
	return bson.D{{Key: "$geoIntersects", Value: v}}
}

func GeoWithin(v interface{}) bson.D {
	return bson.D{{Key: "$geoWithin", Value: v}}
}

func Near(v interface{}) bson.D {
	return bson.D{{Key: "$near", Value: v}}
}

func NearSphere(v interface{}) bson.D {
	return bson.D{{Key: "$nearSphere", Value: v}}
}

func Geometry(geometry interface{}) bson.D {
	return bson.D{{Key: "$geometry", Value: geometry}}
}

func MinDistance(v float64) bson.D {
	return bson.D{{Key: "$geometry", Value: v}}
}

func MaxDistance(v float64) bson.D {
	return bson.D{{Key: "$geometry", Value: v}}
}
