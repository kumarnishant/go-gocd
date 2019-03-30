package gocd

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetPipelineService_Get(t *testing.T) {
	type fields struct {
		client *Client
		log    *logrus.Logger
	}
	type args struct {
		ctx          context.Context
		pipelineName string
		instanceId   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Pipeline test",
			fields: fields{
				client: (&Configuration{
					Server:       "http://aa2e0b202494a11e9aaa80270a75f15a-574269451.us-east-2.elb.amazonaws.com/go/",
					SkipSslCheck: true,
				}).Client(),
				log: logrus.StandardLogger(),
			},
			args: args{
				ctx:          context.Background(),
				pipelineName: "test-da-CD-4",
				instanceId:   12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &GetPipelineService{
				client: tt.fields.client,
				log:    tt.fields.log,
			}
			ps.Get(tt.args.ctx, tt.args.pipelineName, tt.args.instanceId)

		})
	}
}
