package utils

import (
	"fmt"
	"testing"
)

func TestGetFreePort(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "test-1",
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port, err := GetFreePort()
			fmt.Printf("port:%d, err:%v\n", port, err)
		})
	}
}
