package infrastructure

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type PolygonTestNet struct {
	rpc    string
	logger *zap.Logger

	client *ethclient.Client
}

func NewTestNet(rpc string, logger *zap.Logger) *PolygonTestNet {

	return &PolygonTestNet{
		rpc:    rpc,
		logger: logger,
	}
}

func (t *PolygonTestNet) Connect(ctx context.Context) error {

	if t.client != nil {
		t.logger.Info("The client is already running")
		return nil
	}

	client, err := ethclient.Dial(t.rpc)
	if err != nil {
		t.logger.Error("The client was unable to connect to the testnet network.", zap.Error(err))
		return err
	}

	t.client = client
	t.logger.Info("Successfully connected to the testnet network")
	return nil
}
func (t *PolygonTestNet) GetBlockNumber(ctx context.Context) (uint64, error) {

	if t.client == nil {
		t.logger.Info("Network request not possible. client not initialized. Connecting...")
		t.Connect(ctx)
	}

	blocknumber, err := t.client.BlockNumber(ctx)
	if err != nil {
		t.logger.Error("Failed to obtain block number", zap.Error(err))
		return 0, err
	}

	t.logger.Info("Successfully retrieved block number", zap.Uint64("blockNumber", blocknumber))
	return blocknumber, nil
}

func (t *PolygonTestNet) GetPriceActive(ctx context.Context) (string, error) {

	return "", nil

}
