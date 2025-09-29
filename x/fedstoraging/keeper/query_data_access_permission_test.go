package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"flstorage/x/fedstoraging/keeper"
	"flstorage/x/fedstoraging/types"
)

func createNDataAccessPermission(keeper keeper.Keeper, ctx context.Context, n int) []types.DataAccessPermission {
	items := make([]types.DataAccessPermission, n)
	for i := range items {
		items[i].PermissionId = strconv.Itoa(i)
		items[i].Granted = true
		_ = keeper.DataAccessPermission.Set(ctx, items[i].PermissionId, items[i])
	}
	return items
}

func TestDataAccessPermissionQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNDataAccessPermission(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetDataAccessPermissionRequest
		response *types.QueryGetDataAccessPermissionResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDataAccessPermissionRequest{
				PermissionId: msgs[0].PermissionId,
			},
			response: &types.QueryGetDataAccessPermissionResponse{DataAccessPermission: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetDataAccessPermissionRequest{
				PermissionId: msgs[1].PermissionId,
			},
			response: &types.QueryGetDataAccessPermissionResponse{DataAccessPermission: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDataAccessPermissionRequest{
				PermissionId: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetDataAccessPermission(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestDataAccessPermissionQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNDataAccessPermission(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDataAccessPermissionRequest {
		return &types.QueryAllDataAccessPermissionRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListDataAccessPermission(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DataAccessPermission), step)
			require.Subset(t, msgs, resp.DataAccessPermission)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListDataAccessPermission(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DataAccessPermission), step)
			require.Subset(t, msgs, resp.DataAccessPermission)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListDataAccessPermission(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.DataAccessPermission)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListDataAccessPermission(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
