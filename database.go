package cosmos

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/patrickmn/go-cache"

	"github.com/zhevron/cosmos/api"
)

type Database struct {
	api.Database

	client *Client
	cache  *cache.Cache
}

func (d Database) ListCollections(ctx context.Context) ([]*Collection, error) {
	span, ctx := d.startSpan(ctx, "cosmos.ListCollections")
	defer span.Finish()

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

func (d Database) GetCollection(ctx context.Context, id string) (*Collection, error) {
	span, ctx := d.startSpan(ctx, "cosmos.GetCollection")
	defer span.Finish()

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

// TODO: Database.CreateCollection
// TODO: Database.ReplaceCollection

func (d Database) DeleteCollection(ctx context.Context, id string) error {
	span, ctx := d.startSpan(ctx, "cosmos.DeleteCollection")
	defer span.Finish()

	_, err := d.client.delete(ctx, createCollectionLink(d.ID, id), nil)
	return err
}

func (d Database) Client() *Client {
	return d.client
}

func (d Database) startSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := d.client.startSpan(ctx, operationName)
	span.SetTag("db.instance", d.ID)

	return span, ctx
}
