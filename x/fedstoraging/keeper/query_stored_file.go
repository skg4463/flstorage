package keeper

import (
	"context"
	"errors"

	"flstorage/x/fedstoraging/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListStoredFile(ctx context.Context, req *types.QueryAllStoredFileRequest) (*types.QueryAllStoredFileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storedFiles, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.StoredFile,
		req.Pagination,
		func(_ string, value types.StoredFile) (types.StoredFile, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStoredFileResponse{StoredFile: storedFiles, Pagination: pageRes}, nil
}

func (q queryServer) GetStoredFile(ctx context.Context, req *types.QueryGetStoredFileRequest) (*types.QueryGetStoredFileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.StoredFile.Get(ctx, req.OriginalHash)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetStoredFileResponse{StoredFile: val}, nil
}
