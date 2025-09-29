package keeper

import (
	"context"
	"errors"

	"flstorage/x/fedstoraging/types"

	"cosmossdk.io/collections"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Port.Set(ctx, genState.PortId); err != nil {
		return err
	}
	for _, elem := range genState.StoredFileMap {
		if err := k.StoredFile.Set(ctx, elem.OriginalHash, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.DataAccessPermissionMap {
		if err := k.DataAccessPermission.Set(ctx, elem.PermissionId, elem); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	genesis.PortId, err = k.Port.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}
	if err := k.StoredFile.Walk(ctx, nil, func(_ string, val types.StoredFile) (stop bool, err error) {
		genesis.StoredFileMap = append(genesis.StoredFileMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.DataAccessPermission.Walk(ctx, nil, func(_ string, val types.DataAccessPermission) (stop bool, err error) {
		genesis.DataAccessPermissionMap = append(genesis.DataAccessPermissionMap, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
