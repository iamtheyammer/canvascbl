package orautil

import "testing"

func TestBuildLimitExpression(t *testing.T) {
	type args struct {
		limit uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "limit_0",
			args: struct{ limit uint64 }{limit: uint64(0)},
			want: "FETCH NEXT 0 ROWS",
		},
		{
			name: "limit_1000",
			args: struct{ limit uint64 }{limit: uint64(1000)},
			want: "FETCH NEXT 1000 ROWS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildLimitExpression(tt.args.limit); got != tt.want {
				t.Errorf("BuildLimitExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
