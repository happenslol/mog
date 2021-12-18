package codegen

import (
	"fmt"
	"strings"
)

func findIDType(strct Strct) Field {
	for _, fld := range strct.Fields {
		// If one of the existing fields is the index field, we don't
		// need any extra imports since it will already be imported
		// for the field helper methods
		if fld.BsonKey == "_id" {
			return fld
		}
	}

	// If we drop through, there's no explicit ID key and
	// we use primitive.ObjectID, which also needs an import.
	return Field{
		Type: "primitive.ObjectID",
		Imports: map[string]struct{}{
			"go.mongodb.org/mongo-driver/bson/primitive": {}}}
}

func splitPackageUID(uid string) (pkg, strct string, err error) {
	parts := strings.Split(uid, ".")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("Invalid struct identifier: %s", uid)
	}

	pkg = strings.Join(parts[:len(parts)-1], ".")
	strct = parts[len(parts)-1]
	return
}
