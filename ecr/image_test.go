package ecr

import (
	"github.com/Taimee/ecr-lifecycle/ecs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"reflect"
	"testing"
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
