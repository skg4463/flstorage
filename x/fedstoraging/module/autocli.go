package fedstoraging

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"flstorage/x/fedstoraging/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListStoredFile",
					Use:       "list-stored-file",
					Short:     "List all StoredFile",
				},
				{
					RpcMethod:      "GetStoredFile",
					Use:            "get-stored-file [id]",
					Short:          "Gets a StoredFile",
					Alias:          []string{"show-stored-file"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "original_hash"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
