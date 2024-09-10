package usecase

import (
	"context"
	"diplom/infrastructure"

	"go.uber.org/zap"
)

type BlockService struct {
	logger           *zap.Logger
	blockchainClient infrastructure.BlockchainClient
}

func NewBlockService(blockchainClient infrastructure.BlockchainClient) *BlockService {
	return &BlockService{
		blockchainClient: blockchainClient,
	}
}

func (b BlockService) GetLatestBlockNumber(ctx context.Context) (uint64, error) {

	blockNumber, err := b.blockchainClient.GetBlockNumber(ctx)
	if err != nil {
		b.logger.Error("The infrastructure layer returned an error", zap.Error(err))
		return 0, err
	}

	b.blockchainClient.GetPriceActive(ctx)

	return blockNumber, nil
}
