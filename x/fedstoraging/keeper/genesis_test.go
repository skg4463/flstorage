package keeper_test

import (
	"testing"

	"flstorage/x/fedstoraging/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:        types.DefaultParams(),
		PortId:        types.PortID,
		StoredFileMap: []types.StoredFile{{OriginalHash: "0"}, {OriginalHash: "1"}}, DataAccessPermissionMap: []types.DataAccessPermission{{PermissionId: "0"}, {PermissionId: "1"}}}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, genesisState.PortId, got.PortId)
	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.StoredFileMap, got.StoredFileMap)
	require.EqualExportedValues(t, genesisState.DataAccessPermissionMap, got.DataAccessPermissionMap)

}
