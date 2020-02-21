package cosmos

import (
	"context"

	"github.com/zhevron/cosmos/api"
)

type Database struct {
	api.Database

	client *Client
}

func (d Database) Collection(ctx context.Context, id string) (*Collection, error) {
	var collection api.Collection
	if _, err := d.client.get(ctx, createCollectionLink(d.ID, id), &collection, nil); err != nil {
		return nil, err
	}

	return &Collection{
		Collection: collection,
		database:   &d,
	}, nil
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
	}

	return collections, nil
}

// TODO: Database.CreateCollection
// TODO: Database.ReplaceCollection
// TODO: Database.DeleteCollection

func (d Database) Client() *Client {
	return d.client
}
