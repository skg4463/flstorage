package types

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		PortId: PortID, StoredFileMap: []StoredFile{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	storedFileIndexMap := make(map[string]struct{})

	for _, elem := range gs.StoredFileMap {
		index := fmt.Sprint(elem.OriginalHash)
		if _, ok := storedFileIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for storedFile")
		}
		storedFileIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
