package validation

type BlockRewardSlot struct {
	// Status describes Whether the slot contains a block produced by a MEV relay or a vanilla block (built internally in the validator node).
	Status string `json:"status"`
	// Reward describes The amount of reward the node operator/validator received for including the block in that slot (in GWEI).
	Reward float64 `json:"reward"`
}

func GetBlockRewardSlot(slot string) (*BlockRewardSlot, error) {

}
