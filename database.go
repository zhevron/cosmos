package cosmos

import (
	"context"

	"github.com/patrickmn/go-cache"

	"github.com/zhevron/cosmos/api"
)

type Database struct {
	api.Database

	client *Client
	cache  *cache.Cache
}

func (d Database) Collection(ctx context.Context, id string) (*Collection, error) {
	if collection, found := d.cache.Get(id); found {
		return collection.(*Collection), nil
	}

	var coll api.Collection
	if _, err := d.client.get(ctx, createCollectionLink(d.ID, id), &coll, nil); err != nil {
		return nil, err
	}

	collection := &Collection{
		Collection: coll,
		database:   &d,
	}
	d.cache.Set(coll.ID, collection, cache.DefaultExpiration)

	return collection, nil
}

func (d Database) ListCollections(ctx context.Context) ([]*Collection, error) {
	var res api.ListCollectionsResponse
	if _, err := d.client.get(ctx, createCollectionLink(d.ID, ""), &res, nil); err != nil {
		return nil, err
	}

	collections := make([]*Collection, len(res.DocumentCollections))
	for i, c := range res.DocumentCollections {
		collections[i] = &Collection{
			Collection: c,
			database:   &d,
		}
		d.cache.Set(c.ID, collections[i], cache.DefaultExpiration)
	}

	return collections, nil
}

// TODO: Database.CreateCollection
// TODO: Database.ReplaceCollection
// TODO: Database.DeleteCollection

func (d Database) Client() *Client {
	return d.client
}
