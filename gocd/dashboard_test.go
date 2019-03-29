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
					Server:       "http://ac328a6a13cc011e9b7ea02981a66812-1369645122.ap-southeast-1.elb.amazonaws.com/go/",
					Username:     "admin",
					Password:     "asdfuQWDsqas123qw#",
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
			tt.fields.client.Dashboard.Get(tt.args.ctx)
		})
	}
}
