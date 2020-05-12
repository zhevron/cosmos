package cosmos

import (
	"context"
	"io"
	"io/ioutil"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/zhevron/cosmos/api"
)

type Collection struct {
	api.Collection

	database *Database
}

func (c Collection) ListDocuments(ctx context.Context) (*DocumentIterator, error) {
	span, ctx := c.startCollectionSpan(ctx, "cosmos.ListDocuments")
	defer span.Finish()

	var listResult api.ListDocumentsResponse
	res, err := c.database.Client().get(ctx, createDocumentLink(c.database.ID, c.ID, ""), &listResult, nil)
	if err != nil {
		return nil, err
	}

	return newDocumentIterator(ctx, c.database.Client(), res, nil, listResult), nil
}

func (c Collection) GetDocument(ctx context.Context, partitionKey interface{}, id string, out interface{}) error {
	span, ctx := c.startCollectionSpan(ctx, "cosmos.GetDocument")
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	_, err := c.database.Client().get(ctx, createDocumentLink(c.database.ID, c.ID, id), out, headers)
	return err
}

func (c Collection) CreateDocument(ctx context.Context, partitionKey interface{}, document interface{}, upsert bool) error {
	span, ctx := c.startCollectionSpan(ctx, "cosmos.CreateDocument")
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}
	if upsert {
		headers[api.HEADER_IS_UPSERT] = "True"
	}

	_, err := c.database.Client().post(ctx, createDocumentLink(c.database.ID, c.ID, ""), document, nil, headers)
	return err
}

func (c Collection) ReplaceDocument(ctx context.Context, partitionKey interface{}, document interface{}) error {
	documentID, err := DocumentID(document)
	if err != nil {
		return err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.ReplaceDOcument", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	_, err = c.database.Client().put(ctx, createDocumentLink(c.database.ID, c.ID, documentID), document, nil, headers)
	return err
}

func (c Collection) DeleteDocument(ctx context.Context, partitionKey interface{}, document interface{}) error {
	documentID, err := DocumentID(document)
	if err != nil {
		return err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.DeleteDocument", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	_, err = c.database.Client().delete(ctx, createDocumentLink(c.database.ID, c.ID, documentID), headers)
	return err
}

func (c Collection) QueryDocuments(ctx context.Context, partitionKey interface{}, query string, params ...api.QueryParameter) (*DocumentIterator, error) {
	span, ctx := c.startCollectionSpan(ctx, "cosmos.QueryDocuments")
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_CONTENT_TYPE: "application/query+json",
		api.HEADER_IS_QUERY:     "True",
	}

	if partitionKey == nil {
		headers[api.HEADER_QUERY_CROSSPARTITION] = "True"
	} else {
		headers[api.HEADER_PARTITION_KEY] = makePartitionKeyHeaderValue(partitionKey)
	}

	queryParams := []api.QueryParameter{}
	for _, p := range params {
		if strings.Contains(query, p.Name) {
			queryParams = append(queryParams, p)
		}
	}

	ext.DBStatement.Set(span, query)
	span.SetTag("cosmos.parameters", queryParams)

	apiQuery := &api.Query{
		Query:      query,
		Parameters: queryParams,
	}

	var queryResult api.ListDocumentsResponse
	res, err := c.database.Client().post(ctx, createDocumentLink(c.database.ID, c.ID, ""), apiQuery, &queryResult, headers)
	if err != nil {
		return nil, err
	}

	return newDocumentIterator(ctx, c.database.Client(), res, apiQuery, queryResult), nil
}

func (c Collection) ListAttachments(ctx context.Context, partitionKey interface{}, document interface{}) ([]*Attachment, error) {
	documentID, err := DocumentID(document)
	if err != nil {
		return nil, err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.ListAttachments", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	var res api.ListAttachmentsResponse
	if _, err := c.database.Client().get(ctx, createAttachmentLink(c.database.ID, c.ID, documentID, ""), &res, headers); err != nil {
		return nil, err
	}

	attachments := make([]*Attachment, len(res.Attachments))
	for i, a := range res.Attachments {
		attachments[i] = &Attachment{
			Attachment: a,
		}
	}

	return attachments, nil
}

func (c Collection) GetAttachment(ctx context.Context, partitionKey interface{}, document interface{}, id string) (*Attachment, error) {
	documentID, err := DocumentID(document)
	if err != nil {
		return nil, err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.GetAttachment", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	var attachment api.Attachment
	if _, err := c.database.Client().get(ctx, createAttachmentLink(c.database.ID, c.ID, documentID, id), &attachment, headers); err != nil {
		return nil, err
	}

	return &Attachment{
		Attachment: attachment,
	}, nil
}

func (c Collection) CreateAttachmentFromReader(ctx context.Context, partitionKey interface{}, document interface{}, id string, contentType string, reader io.Reader) (*Attachment, error) {
	documentID, err := DocumentID(document)
	if err != nil {
		return nil, err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.CreateAttachmentFromReader", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
		api.HEADER_CONTENT_TYPE:  contentType,
		"Slug":                   id,
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var attachment api.Attachment
	if _, err := c.database.Client().post(ctx, createAttachmentLink(c.database.ID, c.ID, documentID, ""), content, &attachment, headers); err != nil {
		return nil, err
	}

	return &Attachment{
		Attachment: attachment,
	}, nil
}

func (c Collection) CreateAttachmentFromMedia(ctx context.Context, partitionKey interface{}, document interface{}, id string, contentType string, media string) (*Attachment, error) {
	documentID, err := DocumentID(document)
	if err != nil {
		return nil, err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.CreateAttachmentFromMedia", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	attachment := api.Attachment{
		ID:          id,
		ContentType: contentType,
		Media:       media,
	}

	if _, err := c.database.Client().post(ctx, createAttachmentLink(c.database.ID, c.ID, documentID, ""), attachment, &attachment, headers); err != nil {
		return nil, err
	}

	return &Attachment{
		Attachment: attachment,
	}, nil
}

func (c Collection) ReplaceAttachmentFromReader(ctx context.Context, partitionKey interface{}, document interface{}, id string, contentType string, reader io.Reader) (*Attachment, error) {
	documentID, err := DocumentID(document)
	if err != nil {
		return nil, err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.ReplaceAttachmentFromReader", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
		api.HEADER_CONTENT_TYPE:  contentType,
		"Slug":                   id,
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var attachment api.Attachment
	if _, err := c.database.Client().put(ctx, createAttachmentLink(c.database.ID, c.ID, documentID, id), content, &attachment, headers); err != nil {
		return nil, err
	}

	return &Attachment{
		Attachment: attachment,
	}, nil
}

func (c Collection) ReplaceAttachmentFromMedia(ctx context.Context, partitionKey interface{}, document interface{}, id string, contentType string, media string) (*Attachment, error) {
	documentID, err := DocumentID(document)
	if err != nil {
		return nil, err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.ReplaceAttachmentFromMedia", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	attachment := api.Attachment{
		ID:          id,
		ContentType: contentType,
		Media:       media,
	}

	if _, err := c.database.Client().put(ctx, createAttachmentLink(c.database.ID, c.ID, documentID, id), attachment, &attachment, headers); err != nil {
		return nil, err
	}

	return &Attachment{
		Attachment: attachment,
	}, nil
}

func (c Collection) DeleteAttachment(ctx context.Context, partitionKey interface{}, document interface{}, id string) error {
	documentID, err := DocumentID(document)
	if err != nil {
		return err
	}

	span, ctx := c.startDocumentSpan(ctx, "cosmos.DeleteAttachment", documentID)
	defer span.Finish()

	headers := map[string]string{
		api.HEADER_PARTITION_KEY: makePartitionKeyHeaderValue(partitionKey),
	}

	_, err = c.database.Client().delete(ctx, createAttachmentLink(c.database.ID, c.ID, documentID, id), headers)
	return err
}

func (c Collection) Database() *Database {
	return c.database
}

func (c Collection) startCollectionSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := c.database.startSpan(ctx, operationName)
	span.SetTag("cosmos.collection", c.ID)

	return span, ctx
}

func (c Collection) startDocumentSpan(ctx context.Context, operationName string, documentID string) (opentracing.Span, context.Context) {
	span, ctx := c.startCollectionSpan(ctx, operationName)
	span.SetTag("cosmos.document", documentID)

	return span, ctx
}
