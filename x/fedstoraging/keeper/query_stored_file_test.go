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

func createNStoredFile(keeper keeper.Keeper, ctx context.Context, n int) []types.StoredFile {
	items := make([]types.StoredFile, n)
	for i := range items {
		items[i].OriginalHash = strconv.Itoa(i)
		items[i].Tag = strconv.Itoa(i)
		items[i].ShardHashes = []string{`abc` + strconv.Itoa(i), `xyz` + strconv.Itoa(i)}
		_ = keeper.StoredFile.Set(ctx, items[i].OriginalHash, items[i])
	}
	return items
}

func TestStoredFileQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNStoredFile(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetStoredFileRequest
		response *types.QueryGetStoredFileResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetStoredFileRequest{
				OriginalHash: msgs[0].OriginalHash,
			},
			response: &types.QueryGetStoredFileResponse{StoredFile: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetStoredFileRequest{
				OriginalHash: msgs[1].OriginalHash,
			},
			response: &types.QueryGetStoredFileResponse{StoredFile: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetStoredFileRequest{
				OriginalHash: strconv.Itoa(100000),
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
			response, err := qs.GetStoredFile(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestStoredFileQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNStoredFile(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllStoredFileRequest {
		return &types.QueryAllStoredFileRequest{
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
			resp, err := qs.ListStoredFile(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredFile), step)
			require.Subset(t, msgs, resp.StoredFile)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListStoredFile(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.StoredFile), step)
			require.Subset(t, msgs, resp.StoredFile)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListStoredFile(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.StoredFile)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListStoredFile(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
