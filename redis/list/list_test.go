package data

import (
	"context"
	"golang-code-learn/redis/common"
	"testing"

	"github.com/go-redis/redis/v8"
)

func Test_listPopPush(t *testing.T) {
	type args struct {
		ctx context.Context
		rc  *redis.Client
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{context.Background(), common.NewClient()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listPopPush(tt.args.ctx, tt.args.rc)
		})
	}
}

func Test_listPopPushWithLua(t *testing.T) {
	type args struct {
		ctx context.Context
		rc  *redis.Client
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{context.Background(), common.NewClient()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listPopPushWithLua(tt.args.ctx, tt.args.rc)
		})
	}
}

func BenchmarkListPopPush(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listPopPush(context.Background(), common.NewClient())
	}
}

func BenchmarkListPopPushWithLua(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listPopPushWithLua(context.Background(), common.NewClient())
	}
}
