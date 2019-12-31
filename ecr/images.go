package ecr

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func (c *client) BatchDeleteImages(repositoryName string) error {
	imageIds, err := c.deleteTargetImageIds(repositoryName, 50)
	if err != nil {
		return err
	}

	//imageの存在判定
	if len(imageIds) == 0 {
		return nil
	}

	//全てを消す
	err = c.batchDeleteImage(imageIds, repositoryName)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) describeImages(repositoryName string) ([]*ecr.ImageDetail, error) {
	input := &ecr.DescribeImagesInput{
		RepositoryName: &repositoryName,
	}

	result, err := c.ecr.DescribeImages(input)
	if err != nil {
		return nil, err
	}

	return result.ImageDetails, nil
}

func (c *client) deleteTargetImageIds(repositoryName string, imageCountMoreThan int) ([]*ecr.ImageIdentifier, error) {
	images, err := c.describeImages(*aws.String(repositoryName))
	if err != nil {
		return nil, err
	}

	var imageIds []*ecr.ImageIdentifier

	for i, image := range images {
		if i <= imageCountMoreThan {
			continue
		}

		imageIds = append(imageIds, &ecr.ImageIdentifier{
			ImageDigest: image.ImageDigest,
		})
	}

	return imageIds, nil
}

func (c *client) batchDeleteImage(imageIds []*ecr.ImageIdentifier, repositoryName string) error {
	input := &ecr.BatchDeleteImageInput{
		ImageIds: imageIds,
		RepositoryName: &repositoryName,
	}

	result, err := c.ecr.BatchDeleteImage(input)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
