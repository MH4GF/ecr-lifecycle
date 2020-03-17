package main

import (
	"github.com/MH4GF/ecr-lifecycle/ecr"
	"github.com/MH4GF/ecr-lifecycle/ecs"
	"testing"
)

func TestConfig_validate(t *testing.T) {
	type fields struct {
		EcsAssumeRoleArns  []string
		Region             string
		Keep               int
		EcrClient          ecr.Client
		Repositories       []ecr.Repository
		EcsAllRunningTasks []ecs.Task
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "正常系",
			fields: fields{
				EcsAssumeRoleArns: []string{
					"arn:aws:iam::12345678901:role/fuga",
				},
				Region:             "ap-northeast-1",
				Keep:               50,
				EcrClient:          ecr.Client{},
				Repositories:       nil,
				EcsAllRunningTasks: nil,
			},
			wantErr: false,
		},
		{
			name: "EcsAssumeRoleArnsが空sliceの場合",
			fields: fields{
				EcsAssumeRoleArns:  []string{},
				Region:             "ap-northeast-1",
				Keep:               50,
				EcrClient:          ecr.Client{},
				Repositories:       nil,
				EcsAllRunningTasks: nil,
			},
			wantErr: true,
		},
		{
			name: "EcsAssumeRoleArnsのいづれかが20文字以下の場合",
			fields: fields{
				EcsAssumeRoleArns: []string{
					"arn:aws:iam::12345678901:role/fuga",
					"1234567890123456789",
				},
				Region:             "ap-northeast-1",
				Keep:               50,
				EcrClient:          ecr.Client{},
				Repositories:       nil,
				EcsAllRunningTasks: nil,
			},
			wantErr: true,
		},
		{
			name: "Regionが空文字の場合",
			fields: fields{
				EcsAssumeRoleArns: []string{
					"arn:aws:iam::12345678901:role/fuga",
				},
				Region:             "",
				Keep:               50,
				EcrClient:          ecr.Client{},
				Repositories:       nil,
				EcsAllRunningTasks: nil,
			},
			wantErr: true,
		},
		{
			name: "Keepが1以下の場合",
			fields: fields{
				EcsAssumeRoleArns: []string{
					"arn:aws:iam::12345678901:role/fuga",
				},
				Region:             "ap-northeast-1",
				Keep:               0,
				EcrClient:          ecr.Client{},
				Repositories:       nil,
				EcsAllRunningTasks: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				EcsAssumeRoleArns:  tt.fields.EcsAssumeRoleArns,
				Region:             tt.fields.Region,
				Keep:               tt.fields.Keep,
				EcrClient:          tt.fields.EcrClient,
				Repositories:       tt.fields.Repositories,
				EcsAllRunningTasks: tt.fields.EcsAllRunningTasks,
			}
			if err := c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
