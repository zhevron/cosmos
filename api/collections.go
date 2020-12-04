package api

type IndexMode string

const (
	IndexModeConsistent IndexMode = "Consistent"
	IndexModeLazy       IndexMode = "Lazy"
)

type IndexKind string

const (
	IndexKindHash    IndexKind = "Hash"
	IndexKindRange   IndexKind = "Range"
	IndexKindSpatial IndexKind = "Spatial"
)

type IndexDataType string

const (
	IndexDataTypeString     IndexDataType = "String"
	IndexDataTypeNumber     IndexDataType = "Number"
	IndexDataTypePoint      IndexDataType = "Point"
	IndexDataTypePolygon    IndexDataType = "Polygon"
	IndexDataTypeLineString IndexDataType = "LineString"
)

type Collection struct {
	BaseModel

	ID             string         `json:"id"`
	IndexingPolicy IndexingPolicy `json:"indexingPolicy"`
}

type IndexingPolicy struct {
	Automatic     bool      `json:"automatic"`
	IndexingMode  IndexMode `json:"indexingMode"`
	IncludedPaths []struct {
		Path    string `json:"path"`
		Indexes []struct {
			DataType  IndexDataType `json:"dataType"`
			Kind      IndexKind     `json:"kind"`
			Precision int32         `json:"precision"`
		} `json:"indexes"`
	} `json:"includedPaths"`
	ExcludedPaths []struct {
		Path string `json:"path"`
	} `json:"excludedPaths"`
}

type PartitionKeyKind string

const (
	PartitionKeyKindHash PartitionKeyKind = "Hash"
)

type PartitionKey struct {
	Paths   []string         `json:"paths"`
	Kind    PartitionKeyKind `json:"kind"`
	Version int              `json:"version"`
}

type AutopilotSettings struct {
	MaxThroughput int `json:"maxThroughput"`
}

type ListCollectionsResponse struct {
	DocumentCollections []Collection `json:"DocumentCollections"`
}

type CreateCollectionRequest struct {
	ID             string         `json:"id"`
	PartitionKey   PartitionKey   `json:"partitionKey"`
	IndexingPolicy IndexingPolicy `json:"indexingPolicy,omitempty"`
}

type ReplaceCollectionRequest struct {
	ID             string         `json:"id"`
	IndexingPolicy IndexingPolicy `json:"indexingPolicy,omitempty"`
}
