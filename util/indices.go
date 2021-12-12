package util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndices(c *mongo.Collection,
	idxs []mongo.IndexModel, timeout time.Duration,
) error {
	if timeout == 0 {
		timeout = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	iv := c.Indexes()

	cur, err := iv.List(ctx)
	if err != nil {
		return err
	}

	existing := []bson.M{}
	for cur.Next(ctx) {
		var idx bson.M
		err = cur.Decode(&idx)
		if err != nil {
			return err
		}

		existing = append(existing, idx)
	}

	opts := options.CreateIndexesOptions{}
	for _, idx := range idxs {
		if keys, ok := idx.Keys.(primitive.M); ok && len(keys) > 1 {
			return fmt.Errorf("Expected 'bson.D' as compound index type, got '%T'", idx.Keys)
		}

		nameParts := GetNameParts(idx)

		name := strings.Join(nameParts, "_")
		if idx.Options == nil {
			idx.Options = &options.IndexOptions{}
		}

		idx.Options.Name = &name

		exists := false
		for _, ex := range existing {
			if name == ex["name"] {
				exists = true
				break
			}
		}

		if !exists {
			if _, err := iv.CreateOne(ctx, idx, &opts); err != nil {
				return fmt.Errorf("Failed to create index: %w", err)
			}
		}
	}

	if err = ctx.Err(); err != nil {
		return fmt.Errorf("Time out while creating indexes: %w", err)
	}

	return nil
}

func GetNameParts(idx mongo.IndexModel) []string {
	nameParts := []string{}
	switch keys := idx.Keys.(type) {
	case primitive.M:
		for k, v := range keys {
			nameParts = append(nameParts, fmt.Sprintf("%v_%v", k, v))
		}
	case primitive.D:
		for _, v := range keys {
			nameParts = append(nameParts, fmt.Sprintf("%v_%v", v.Key, v.Value))
		}
	}
	return nameParts
}
