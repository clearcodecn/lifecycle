package lifecycle

import (
	"context"
	"testing"
	"time"
)

func TestLifeCycle(t *testing.T) {
	tests := []struct {
		name    string
		hooks   []Hook
		wantErr bool
	}{
		{
			name: "1",
			hooks: []Hook{
				{
					OnStart: func(ctx context.Context) error {
						time.Sleep(10 * time.Second)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return nil
					},
				},
				{
					OnStart: func(ctx context.Context) error {
						time.Sleep(12 * time.Second)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return nil
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lf := New()
			lf.Add(tt.hooks...)
			err := lf.Start(context.Background())
			if err != nil && tt.wantErr {
				return
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
}
