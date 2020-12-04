package cosmos

import (
	"encoding/json"
	"strconv"

	"github.com/zhevron/cosmos/api"
)

type CreateCollectionOption func(*api.CreateCollectionRequest, map[string]string)

func WithPartitionKey(partitionKey api.PartitionKey) CreateCollectionOption {
	return func(req *api.CreateCollectionRequest, headers map[string]string) {
		partitionKey.Version = api.PARTITION_KEY_VERSION
		req.PartitionKey = partitionKey
	}
}

func WithIndexingPolicy(indexingPolicy api.IndexingPolicy) CreateCollectionOption {
	return func(req *api.CreateCollectionRequest, headers map[string]string) {
		req.IndexingPolicy = indexingPolicy
	}
}

func WithThroughput(throughput int) CreateCollectionOption {
	return func(req *api.CreateCollectionRequest, headers map[string]string) {
		headers[api.HEADER_OFFER_THROUGHPUT] = strconv.Itoa(throughput)
		delete(headers, api.HEADER_OFFER_AUTOPILOT)
	}
}

func WithAutopilot(settings api.AutopilotSettings) CreateCollectionOption {
	return func(req *api.CreateCollectionRequest, headers map[string]string) {
		settingsJSON, _ := json.Marshal(settings)
		headers[api.HEADER_OFFER_AUTOPILOT] = string(settingsJSON)
		delete(headers, api.HEADER_OFFER_THROUGHPUT)
	}
}

type ReplaceCollectionOption func(*api.ReplaceCollectionRequest, map[string]string)

func ReplaceIndexingPolicy(indexingPolicy api.IndexingPolicy) ReplaceCollectionOption {
	return func(req *api.ReplaceCollectionRequest, headers map[string]string) {
		req.IndexingPolicy = indexingPolicy
	}
}
