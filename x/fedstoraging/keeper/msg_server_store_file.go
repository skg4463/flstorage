package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"flstorage/x/fedstoraging/types"
)

func (k msgServer) StoreFile(goCtx context.Context, msg *types.MsgStoreFile) (*types.MsgStoreFileResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // StoredFile 객체 생성.
    // Creator 필드는 ignite scaffold map이 자동으로 추가해 줌.
	var storedFile = types.StoredFile{
		Creator:      msg.Creator,
		OriginalHash: msg.OriginalHash,
		Tag:          msg.Tag,
		ShardHashes:  msg.ShardHashes,
	}

    // 최신 collections 문법을 사용하여 상태에 저장.
	err := k.StoredFile.Set(ctx, storedFile.OriginalHash, storedFile)
	if err != nil {
		return nil, err
	}

	return &types.MsgStoreFileResponse{}, nil
}