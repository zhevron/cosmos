package api

const (
	HEADER_CONSISTENCY_LEVEL    = "x-ms-consistency-level"
	HEADER_CONTENT_TYPE         = "Content-Type"
	HEADER_CONTINUATION         = "x-ms-continuation"
	HEADER_DATE                 = "x-ms-date"
	HEADER_IS_QUERY             = "x-ms-documentdb-isquery"
	HEADER_IS_UPSERT            = "x-ms-documentdb-is-upsert"
	HEADER_MAX_ITEM_COUNT       = "x-ms-max-item-count"
	HEADER_PARTITION_KEY        = "x-ms-documentdb-partitionkey"
	HEADER_QUERY_CROSSPARTITION = "x-ms-documentdb-query-enablecrosspartition"
	HEADER_QUERY_METRICS        = "x-ms-documentdb-populatequerymetrics"
	HEADER_OFFER_AUTOPILOT      = "x-ms-cosmos-offer-autopilot-settings"
	HEADER_OFFER_THROUGHPUT     = "x-ms-offer-throughput"
	HEADER_REQUEST_CHARGE       = "x-ms-request-charge"
	HEADER_RESOURCE_QUOTA       = "x-ms-resource-quota"
	HEADER_RESOURCE_USAGE       = "x-ms-resource-usage"
	HEADER_RETRY_AFTER          = "retry-after-ms"
	HEADER_SESSION_TOKEN        = "x-ms-session-token" // nolint:gosec
	HEADER_VERSION              = "x-ms-version"
	PARTITION_KEY_VERSION       = 2
	TIME_FORMAT                 = "Mon, 02 Jan 2006 15:04:05 GMT"
)

type BaseModel struct {
	RID  string   `json:"_rid"`
	Etag string   `json:"_etag"`
	Ts   DateTime `json:"_ts"`
}
