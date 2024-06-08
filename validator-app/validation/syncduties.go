package validation

type SyncDutiesResponse struct {
	// PublicValidatorKeys is a list of public keys of validators that had sync committee duties for the specified slot.
	PublicValidatorKeys []string
}

func GetSyncDuties(slot uint64) (*SyncDutiesResponse, error) {

}
