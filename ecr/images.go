package ecr

import (
	"fmt"
	"github.com/Taimee/ecr-lifecycle/ecs"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// Image is stored original ecr.Image
type Image struct {
	original *ecr.Image
}

// Image型にはURIがないので定義
// 012345678910.dkr.ecr.<region-name>.amazonaws.com/<repository-name>:latest
func (i *Image) uri(r *Repository) *string {
	uri := *r.original.RepositoryUri + ":" + *i.original.ImageId.ImageTag
	return &uri
}

// BatchDeleteImages ... 指定したrepositoryのimageを削除する。
func (c *Client) BatchDeleteImages(r *Repository, imageCountMoreThan *int) error {
	input, err := c.newRegisterBatchDeleteImageInput(r, imageCountMoreThan)
	if err != nil {
		return err
	}

	//imageの存在判定
	if len(input.ImageIds) == 0 {
		return nil
	}

	result, err := c.ecr.BatchDeleteImage(input)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

func (c *Client) newRegisterBatchDeleteImageInput(r *Repository, imageCountMoreThan *int) (*ecr.BatchDeleteImageInput, error) {
	images, err := c.batchGetImages(r)
	if err != nil {
		return nil, err
	}

	sess, err := ecs.RegisterECSNewSession("sandbox", "ap-northeast-1")
	if err != nil {
		return nil, err
	}
	ecsClient := ecs.NewClient(sess)
	runningTasks, err := ecsClient.ListAllRunningTasks()
	if err != nil {
		return nil, err
	}

	var imageIds []*ecr.ImageIdentifier

	for i, image := range images {
		if i < *imageCountMoreThan {
			continue
		}

		// 現在実行中のタスクがある場合は削除しない
		if image.isUsedRunningTasks(runningTasks, r) {
			continue
		} else {
			imageIds = append(imageIds, image.original.ImageId)
		}
	}

	input := &ecr.BatchDeleteImageInput{
		ImageIds:       imageIds,
		RepositoryName: r.original.RepositoryName,
	}

	return input, nil
}

func (i *Image) isUsedRunningTasks(tasks []*ecs.Task, r *Repository) bool {
	uri := i.uri(r)

	for _, task := range tasks {
		if *task.Image == *uri {
			return true
		}
	}

	return false
}

func (c *Client) batchGetImages(r *Repository) ([]*Image, error) {
	input, err := c.newRegisterBatchGetImageInput(r)
	if err != nil {
		return nil, err
	}

	result, err := c.ecr.BatchGetImage(input)
	if err != nil {
		return nil, err
	}

	var images []*Image
	for _, image := range result.Images {
		images = append(images, &Image{original: image})
	}

	return images, nil
}

func (c *Client) newRegisterBatchGetImageInput(r *Repository) (*ecr.BatchGetImageInput, error) {
	input := &ecr.DescribeImagesInput{
		RepositoryName: r.original.RepositoryName,
	}

	result, err := c.ecr.DescribeImages(input)
	if err != nil {
		return nil, err
	}

	var imageIds []*ecr.ImageIdentifier
	for _, imageDetail := range result.ImageDetails {
		imageIds = append(imageIds, &ecr.ImageIdentifier{ImageDigest: imageDetail.ImageDigest})
	}

	batchInput := &ecr.BatchGetImageInput{ImageIds: imageIds, RepositoryName: r.original.RepositoryName}

	return batchInput, nil
}
