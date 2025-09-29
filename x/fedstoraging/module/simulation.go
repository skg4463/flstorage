package fedstoraging

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	fedstoragingsimulation "flstorage/x/fedstoraging/simulation"
	"flstorage/x/fedstoraging/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	fedstoragingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&fedstoragingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgStoreFile          = "op_weight_msg_fedstoraging"
		defaultWeightMsgStoreFile int = 100
	)

	var weightMsgStoreFile int
	simState.AppParams.GetOrGenerate(opWeightMsgStoreFile, &weightMsgStoreFile, nil,
		func(_ *rand.Rand) {
			weightMsgStoreFile = defaultWeightMsgStoreFile
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStoreFile,
		fedstoragingsimulation.SimulateMsgStoreFile(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRequestDataAccess          = "op_weight_msg_fedstoraging"
		defaultWeightMsgRequestDataAccess int = 100
	)

	var weightMsgRequestDataAccess int
	simState.AppParams.GetOrGenerate(opWeightMsgRequestDataAccess, &weightMsgRequestDataAccess, nil,
		func(_ *rand.Rand) {
			weightMsgRequestDataAccess = defaultWeightMsgRequestDataAccess
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestDataAccess,
		fedstoragingsimulation.SimulateMsgRequestDataAccess(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
