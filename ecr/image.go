package ecr

import (
	"fmt"
	"github.com/MH4GF/ecr-lifecycle/ecs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"sort"
)

// Image ... Store original ecr.Image
type Image struct {
	Detail *ecr.ImageDetail
}

// Uris ... uriの配列を返す。持っているタグごとのuriと、ImageDigestを用いたuriの両方を返す
func (i *Image) Uris(r ecr.Repository) []string {
	uris := make([]string, 0)

	// 012345678910.dkr.ecr.<region-name>.amazonaws.com/<repository-name>@<image-digest>
	uris = append(uris, fmt.Sprintf("%s@%s", aws.StringValue(r.RepositoryUri), aws.StringValue(i.Detail.ImageDigest)))

	// 012345678910.dkr.ecr.<region-name>.amazonaws.com/<repository-name>:tag
	for _, tag := range i.Detail.ImageTags {
		uris = append(uris, fmt.Sprintf("%s:%s", aws.StringValue(r.RepositoryUri), aws.StringValue(tag)))
	}

	return uris
}

// BatchDeleteImages ... 指定したrepositoryのimageを削除する。
func (c *Client) BatchDeleteImages(r Repository, imageCountMoreThan int, ecs ecs.Client) (*ecr.BatchDeleteImageOutput, error) {
	input, err := c.BatchDeleteImageInput(*r.Detail, imageCountMoreThan, ecs)
	if err != nil {
		return nil, err
	}

	//imageの存在判定
	if len(input.ImageIds) == 0 {
		return nil, nil
	}

	result, err := c.ecr.BatchDeleteImage(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// BatchDeleteImageInput ... DeleteするImageを絞り込む。
func (c *Client) BatchDeleteImageInput(r ecr.Repository, imageCountMoreThan int, ecs ecs.Client) (*ecr.BatchDeleteImageInput, error) {
	images, err := c.BatchGetImages(r)
	if err != nil {
		return nil, err
	}

	runningTasks, err := ecs.ListAllRunningTasks()
	if err != nil {
		return nil, err
	}

	var imageIds []*ecr.ImageIdentifier

	for i, image := range images {
		if i < imageCountMoreThan {
			continue
		}

		// 現在実行中のタスクがある場合は削除しない
		if image.IsImageUsedRunningTasks(runningTasks, r) {
			continue
		} else {
			imageIds = append(imageIds, &ecr.ImageIdentifier{ImageDigest: image.Detail.ImageDigest})
		}
	}

	input := &ecr.BatchDeleteImageInput{
		ImageIds:       imageIds,
		RepositoryName: r.RepositoryName,
	}

	return input, nil
}

// IsImageUsedRunningTasks ... 今動いてるタスクでイメージが使われてないか
func (i *Image) IsImageUsedRunningTasks(tasks []ecs.Task, r ecr.Repository) bool {
	uris := i.Uris(r)

	for _, task := range tasks {
		for _, uri := range uris {
			if task.Image == uri {
				return true
			}
		}
	}

	return false
}

// BatchGetImages ... イメージの詳細を取得
func (c *Client) BatchGetImages(r ecr.Repository) ([]*Image, error) {
	input := &ecr.DescribeImagesInput{
		RepositoryName: r.RepositoryName,
	}

	result, err := c.ecr.DescribeImages(input)
	if err != nil {
		return nil, err
	}

	var images []*Image
	for _, image := range result.ImageDetails {
		images = append(images, &Image{Detail: image})
	}
	sortedImages := sortImages(images)

	return sortedImages, nil
}

// sortImages ... ImagePushedAtの降順になるようにソートする
func sortImages(images []*Image) []*Image {
	sort.SliceStable(images, func(i, j int) bool {
		return images[i].Detail.ImagePushedAt.Before(*images[j].Detail.ImagePushedAt)
	})

	return images
}
