package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/tendermint/rpc/client/mock"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

type MockClient struct {
	mock.Client
	heights []time.Time
}

func (client MockClient) Block(height *int64) (*ctypes.ResultBlock, error) {
	return &ctypes.ResultBlock{
		BlockMeta: &types.BlockMeta{
			Header: types.Header{
				Time: client.heights[*height],
			},
		},
	}, nil
}

func TestFindBlockHeightWithDate(t *testing.T) {
	startTime := time.Date(2019, 01, 01, 0, 0, 0, 0, time.Local)
	client := MockClient{
		heights: []time.Time{
			startTime,
			startTime.Add(20 * time.Second),
			startTime.Add(100 * time.Second),
			startTime.Add(200 * time.Second),
			startTime.Add(300 * time.Second)},
	}

	tests := []struct {
		startTime time.Time
		expectedHeight int64
		error bool
	}{
		{startTime, 0, true},
		{startTime.Add(20 * time.Second), int64(1), false},
		{startTime.Add(50 * time.Second), int64(2), false},
		{startTime.Add(150 * time.Second), int64(3), false},
		{startTime.Add(250 * time.Second), int64(4), false},
		{startTime.Add(350 * time.Second), 0, true},
	}

	for _, test := range tests {
		height, err := findBlockHeightFromDate(test.startTime, client, 4)
		if test.error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedHeight, *height)
		}
	}
}

