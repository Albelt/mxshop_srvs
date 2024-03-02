package model

import (
	"fmt"
	"testing"
)

func TestDecodePassword(t *testing.T) {
	type args struct {
		origin string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name:  "test-1",
			args:  args{origin: "123456"},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt, hash := encodePassword(tt.args.origin)
			fmt.Printf("salt:%s\n", salt)
			fmt.Printf("hash:%s\n", hash)

			ok := verifyPassword("12345", salt, hash)
			fmt.Printf("verify:%v\n", ok)
		})
	}
}
