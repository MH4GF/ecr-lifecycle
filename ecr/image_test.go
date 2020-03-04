package ecr

import (
	"github.com/MH4GF/ecr-lifecycle/ecs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"reflect"
	"testing"
	"time"
)

func TestImage_Uris(t *testing.T) {
	r := ecr.Repository{
		RepositoryUri: aws.String("012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge"),
	}

	cases := []struct {
		input    Image
		expected []string
	}{
		{
			input: Image{
				Detail: &ecr.ImageDetail{
					ImageTags: []*string{
						aws.String("latest"),
					},
					ImageDigest: aws.String("sha256:example"),
				},
			},
			expected: []string{
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge@sha256:example",
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:latest",
			},
		},
		{
			input: Image{
				Detail: &ecr.ImageDetail{
					ImageTags: []*string{
						aws.String("latest"),
						aws.String("prod"),
					},
					ImageDigest: aws.String("sha256:example"),
				},
			},
			expected: []string{
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge@sha256:example",
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:latest",
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:prod",
			},
		},
		{
			input: Image{
				Detail: &ecr.ImageDetail{
					ImageTags:   []*string{},
					ImageDigest: aws.String("sha256:example"),
				},
			},
			expected: []string{
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge@sha256:example",
			},
		},
	}

	for _, c := range cases {
		if !reflect.DeepEqual(c.input.Uris(r), c.expected) {
			t.Errorf("Expected to contain %v, but not.", c.expected)
		}
	}
}

func TestImage_IsImageUsedRunningTasks(t *testing.T) {
	r := ecr.Repository{
		RepositoryUri: aws.String("012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge"),
	}
	i := Image{
		Detail: &ecr.ImageDetail{
			ImageTags: []*string{
				aws.String("latest"),
			},
			ImageDigest: aws.String("sha256:example"),
		},
	}

	cases := []struct {
		input    []ecs.Task
		expected bool
	}{
		{
			input:    []ecs.Task{},
			expected: false,
		},
		{
			input: []ecs.Task{
				{
					Image: "012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:latest",
				},
			},
			expected: true,
		},
		{
			input: []ecs.Task{
				{
					Image: "012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:different-tag",
				},
			},
			expected: false,
		},
		{
			input: []ecs.Task{
				{
					Image: "012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:latest",
				},
				{
					Image: "012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:different-tag",
				},
			},
			expected: true,
		},
	}

	for _, c := range cases {
		if i.IsImageUsedRunningTasks(c.input, r) != c.expected {
			t.Errorf("Expected value to be %v, but not.", c.expected)
		}
	}
}

func strToTime(s string) *time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return &t
}

func Test_sortImages(t *testing.T) {
	type args struct {
		images []*Image
	}
	tests := []struct {
		name string
		args args
		want []*Image
	}{
		{
			name: "作成日が最新のものから降順でソートされる",
			args: struct {
				images []*Image
			}{
				images: []*Image{
					{
						Detail: &ecr.ImageDetail{
							ImagePushedAt: strToTime("2020-02-01"),
						},
					},
					{
						Detail: &ecr.ImageDetail{
							ImagePushedAt: strToTime("2020-01-31"),
						},
					},
					{
						Detail: &ecr.ImageDetail{
							ImagePushedAt: strToTime("2020-02-02"),
						},
					},
				},
			},
			want: []*Image{
				{
					Detail: &ecr.ImageDetail{
						ImagePushedAt: strToTime("2020-02-02"),
					},
				},
				{
					Detail: &ecr.ImageDetail{
						ImagePushedAt: strToTime("2020-02-01"),
					},
				},
				{
					Detail: &ecr.ImageDetail{
						ImagePushedAt: strToTime("2020-01-31"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortImages(tt.args.images); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortImages() = %v, want %v", got, tt.want)
			}
		})
	}
}
