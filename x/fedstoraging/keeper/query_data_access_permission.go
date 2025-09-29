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

func (q queryServer) ListDataAccessPermission(ctx context.Context, req *types.QueryAllDataAccessPermissionRequest) (*types.QueryAllDataAccessPermissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	dataAccessPermissions, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.DataAccessPermission,
		req.Pagination,
		func(_ string, value types.DataAccessPermission) (types.DataAccessPermission, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDataAccessPermissionResponse{DataAccessPermission: dataAccessPermissions, Pagination: pageRes}, nil
}

func (q queryServer) GetDataAccessPermission(ctx context.Context, req *types.QueryGetDataAccessPermissionRequest) (*types.QueryGetDataAccessPermissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.DataAccessPermission.Get(ctx, req.PermissionId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetDataAccessPermissionResponse{DataAccessPermission: val}, nil
}
