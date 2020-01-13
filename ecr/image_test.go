package ecr

import (
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
				},
			},
			expected: []string{"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:latest"},
		},
		{
			input: Image{
				Detail: &ecr.ImageDetail{
					ImageTags: []*string{
						aws.String("latest"),
						aws.String("prod"),
					},
				},
			},
			expected: []string{
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:latest",
				"012345678910.dkr.ecr.ap-northeast-1.amazonaws.com/hoge:prod",
			},
		},
		{
			input: Image{
				Detail: &ecr.ImageDetail{
					ImageTags: []*string{},
				},
			},
			expected: []string{},
		},
	}

	for _, c := range cases {
		if !reflect.DeepEqual(c.input.Uris(r), c.expected) {
			t.Errorf("Expected to contain %v, but not.", c.expected)
		}
	}
}
