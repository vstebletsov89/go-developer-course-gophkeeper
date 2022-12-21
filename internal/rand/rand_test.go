package rand

import (
	"reflect"
	"testing"
)

func TestGenerateRandom(t *testing.T) {
	tests := []struct {
		name   string
		size   int
		result int
	}{
		{
			name:   "generate random with 1 byte length",
			size:   1,
			result: 1,
		},
		{
			name:   "generate random with 10 byte length",
			size:   10,
			result: 10,
		},
		{
			name:   "generate random with 100 byte length",
			size:   100,
			result: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateRandom(tt.size); !reflect.DeepEqual(len(got), tt.result) {
				t.Errorf("GenerateRandom() = %v, want %v", got, tt.result)
			}
		})
	}
}
