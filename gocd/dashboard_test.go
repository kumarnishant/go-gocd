package gocd

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestDashboardService_Get(t *testing.T) {
	type fields struct {
		client *Client
		log    *logrus.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Dashboard test",
			fields: fields{
				client: (&Configuration{
					Server:       "http://aa2e0b202494a11e9aaa80270a75f15a-574269451.us-east-2.elb.amazonaws.com/go/",
					SkipSslCheck: true,
				}).Client(),
				log: logrus.StandardLogger(),
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//ds := &DashboardService{
			//	client: tt.fields.client,
			//	log:    tt.fields.log,
			//}
			tt.fields.client.Dashboard.Get(tt.args.ctx, "")
		})
	}
}
