package validation

import (
	"errors"
	"math/big"
)

type BlockRewardSlot struct {
	// Status describes Whether the slot contains a block produced by a MEV relay or a vanilla block (built internally in the validator node).
	Status string `json:"status"`
	// Reward describes The amount of reward the node operator/validator received for including the block in that slot (in GWEI).
	Reward float64 `json:"reward"`
}

func GetBlockRewardSlot(slot uint64) (*BlockRewardSlot, error) {
	// Get MEV2 Client
	mev2, errMEV2 := getMEV2RelayBackendClient()
	if errMEV2 != nil {
		return nil, errMEV2
	}

	mev2.EthBlockNumber()

	// Get Web3 Client
	client, errClient := getWeb3BackendClient()
	if errClient != nil {
		return nil, errClient
	}
	// Get Current Block Number in the Endpoint's network
	currentBlockNumber, errBlockNumber := client.Eth.GetBlockNumber()
	if errBlockNumber != nil {
		return nil, errBlockNumber
	}
	// Ensure it's not in the future
	if slot > currentBlockNumber {
		return nil, errors.New(ErrSlotInFuture)
	}

	blockNumber := &big.Int{}
	blockNumber.SetUint64(slot)
	blockInfo, errBlockInfo := client.Eth.GetBlocByNumber(blockNumber, true)
	if errBlockInfo != nil {
		return nil, errBlockInfo
	}
	blockInfo

}
