package filter

import "go.mongodb.org/mongo-driver/bson"

func Expr(v interface{}) bson.D {
	return bson.D{{Key: "$expr", Value: v}}
}

func JSONSchema(v map[string]interface{}) bson.D {
	// TODO: Is there a better type for this?
	return bson.D{{Key: "$jsonSchema", Value: v}}
}

func Mod(div, rem int64) bson.D {
	return bson.D{{Key: "$mod", Value: []int64{div, rem}}}
}

func Regex(pattern string) bson.D {
	return bson.D{{Key: "$regex", Value: pattern}}
}

func RegexOpt(pattern, options string) bson.D {
	return bson.D{
		{Key: "$regex", Value: pattern},
		{Key: "$options", Value: options}}
}

func Text(
	search, language string,
	caseSensitive, diacriticSensitive bool,
) bson.D {
	return bson.D{{
		Key: "$text",
		Value: bson.D{
			{Key: "$search", Value: search},
			{Key: "$language", Value: language},
			{Key: "$caseSensitive", Value: caseSensitive},
			{Key: "$diacriticSensitive", Value: diacriticSensitive},
		}}}
}
