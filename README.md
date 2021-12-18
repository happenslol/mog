# mog

`mog` is a **Mo**ngoDB code **g**enerator library that can be used to generate helpers for building queries using the official [mongodb go driver](https://github.com/mongodb/mongo-go-driver).

## Usage

Generate a default config file:

```bash
go run github.com/happenslol/mog init
```

Add your models and primitives to the config:

```yaml
primitives:
- github.com/happenslol/mog/testdata.CustomID
- github.com/happenslol/mog/testdata.PageCount

collections:
  authors: github.com/happenslol/mog/testdata.Author
  books: github.com/happenslol/mog/testdata.Author
```

Generate code:

```bash
go run github.com/happenslol/mog
```

Import the generated code and the query builder helpers:

```go
import (
  . "github.com/happenslol/mog/filter"
  . "github.com/happenslol/mog/modelgen"
  cols "github.com/happenslol/mog/colgen"
)
```

Write your queries:

```go
col := cols.NewAuthorsCollection(db)

authors, err := col.Find(ctx, And(
  Author.Age(Gt(50)),
  Author.Age(Lt(60)),
  Author.Books(Book.Pages(Gte(1000))),
))
```

This will generate the following query:

```json
{
  "$and": [
    {"age": {"$gt": 50}},
    {"age": {"$lt": 60}},
    {"books.pages": { "$gte": 1000 }},
  ]
}
```

The helper methods all take and return `bson.D` (or their respective type in case there is only 1 type allowed), and can be combined and nested in any order. `mog` respects `bson` tags, will correctly infer any struct and ID types (including embedded structs) and will generate query helpers for any nested structs that are not marked primitives.

The generated collection code simply copies the interface from the `mongo.Collection` type and reimplements all methods while casting the inputs and results to your given type. The filters implement all mongodb operators listed [here](https://docs.mongodb.com/manual/reference/operator/query/).
