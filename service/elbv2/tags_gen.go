// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package elbv2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	tftags "github.com/infracasts/terraform-provider-aws-public/tags"
)

// ListTags lists elbv2 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn elbv2iface.ELBV2API, identifier string) (tftags.KeyValueTags, error) {
	return ListTagsWithContext(context.Background(), conn, identifier)
}

func ListTagsWithContext(ctx context.Context, conn elbv2iface.ELBV2API, identifier string) (tftags.KeyValueTags, error) {
	input := &elbv2.DescribeTagsInput{
		ResourceArns: aws.StringSlice([]string{identifier}),
	}

	output, err := conn.DescribeTagsWithContext(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.TagDescriptions[0].Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns elbv2 service tags.
func Tags(tags tftags.KeyValueTags) []*elbv2.Tag {
	result := make([]*elbv2.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &elbv2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from elbv2 service tags.
func KeyValueTags(tags []*elbv2.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates elbv2 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn elbv2iface.ELBV2API, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTagsWithContext(context.Background(), conn, identifier, oldTags, newTags)
}
func UpdateTagsWithContext(ctx context.Context, conn elbv2iface.ELBV2API, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &elbv2.RemoveTagsInput{
			ResourceArns: aws.StringSlice([]string{identifier}),
			TagKeys:      aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.RemoveTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &elbv2.AddTagsInput{
			ResourceArns: aws.StringSlice([]string{identifier}),
			Tags:         Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.AddTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
