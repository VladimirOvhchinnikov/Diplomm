package infrastructure

import "context"

type BlockchainClient interface {
	Connect(ctx context.Context) error
	GetBlockNumber(ctx context.Context) (uint64, error)
	GetPriceActive(ctx context.Context) (string, error)
}
