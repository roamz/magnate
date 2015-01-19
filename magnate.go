package magnate

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func Init(c Client) error {
	indexers := []collectionIndexer{
		Marker{},
	}

	var err error

	for _, collection := range indexers {
		for _, index := range collection.Indexes() {
			if err = c.C(collection).EnsureIndex(index); err != nil {
				return err
			}
		}
	}

	return nil
}

type Namer interface {
	CollectionName() string
}

type collectionIndexer interface {
	Namer
	Indexes() []mgo.Index
}

type Client struct {
	*mgo.Database
}

func (c Client) C(n Namer) *mgo.Collection {
	return c.Database.C(n.CollectionName())
}

type Marker struct {
	ID      bson.ObjectId `bson:"_id"`
	Number  int           `bson:"number"`
	Label   string        `bson:"label"`
	Partial bool          `bson:"partial,omitempty"`
}

func (m Marker) CollectionName() string {
	return "magnate_migrations"
}

func (m Marker) Indexes() []mgo.Index {
	return []mgo.Index{
		{
			Key:    []string{"number"},
			Unique: true,
		},
	}
}
