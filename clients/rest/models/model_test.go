package models

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeribitResponse_Unmarshal(t *testing.T) {
	tests := []struct {
		name          string
		jsonData      string
		expectedResp  DeribitResponse[string]
		expectError   bool
		errorContains string
	}{
		{
			name: "valid_response",
			jsonData: `{
				"jsonrpc": "2.0",
				"id": 123,
				"result": "test_result"
			}`,
			expectedResp: DeribitResponse[string]{
				Jsonrpc: "2.0",
				Id:      int64Ptr(123),
				Result:  "test_result",
			},
		},
		{
			name: "missing_id",
			jsonData: `{
				"jsonrpc": "2.0",
				"result": "test_result"
			}`,
			expectedResp: DeribitResponse[string]{
				Jsonrpc: "2.0",
				Result:  "test_result",
			},
		},
		{
			name:          "invalid_json",
			jsonData:      `{invalid`,
			expectError:   true,
			errorContains: "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp DeribitResponse[string]
			err := json.Unmarshal([]byte(tt.jsonData), &resp)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp.Jsonrpc, resp.Jsonrpc)
			assert.Equal(t, tt.expectedResp.Id, resp.Id)
			assert.Equal(t, tt.expectedResp.Result, resp.Result)
		})
	}
}

func TestDirection_Constants(t *testing.T) {
	tests := []struct {
		direction Direction
		expected  string
	}{
		{Buy, "buy"},
		{Sell, "sell"},
	}

	for _, tt := range tests {
		t.Run(string(tt.direction), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.direction))
		})
	}
}

func TestDepth_Constants(t *testing.T) {
	tests := []struct {
		depth    Depth
		expected int
	}{
		{One, 1},
		{Five, 5},
		{Ten, 10},
		{TwentyFive, 25},
		{OneHundred, 100},
		{OneThousand, 1000},
		{TenThousand, 10000},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("depth_%d", tt.expected), func(t *testing.T) {
			assert.Equal(t, tt.expected, int(tt.depth))
		})
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}
