package keeper

import (
	"context"

	"flstorage/x/fedstoraging/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) StoreFile(ctx context.Context, msg *types.MsgStoreFile) (*types.MsgStoreFileResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgStoreFileResponse{}, nil
}
