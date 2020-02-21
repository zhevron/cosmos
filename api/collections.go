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

type ListCollectionsResponse struct {
	DocumentCollections []Collection `json:"DocumentCollections"`
}

// TODO: CreateCollectionRequest ()
