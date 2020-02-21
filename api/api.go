package api

const (
	HEADER_CONSISTENCY_LEVEL    = "x-ms-consistency-level"
	HEADER_CONTENT_TYPE         = "Content-Type"
	HEADER_CONTINUATION         = "x-ms-continuation"
	HEADER_DATE                 = "x-ms-date"
	HEADER_IS_QUERY             = "x-ms-documentdb-isquery"
	HEADER_MAX_ITEM_COUNT       = "x-ms-max-item-count"
	HEADER_PARTITION_KEY        = "x-ms-documentdb-partitionkey"
	HEADER_QUERY_CROSSPARTITION = "x-ms-documentdb-query-enablecrosspartition"
	HEADER_SESSION_TOKEN        = "x-ms-session-token" // nolint:gosec
	HEADER_VERSION              = "x-ms-version"
	TIME_FORMAT                 = "Mon, 02 Jan 2006 15:04:05 GMT"
)

type BaseModel struct {
	RID  string `json:"_rid"`
	Etag string `json:"_etag"`
	Ts   int64  `json:"_ts"`
}
