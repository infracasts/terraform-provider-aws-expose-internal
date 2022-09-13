// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package organizations

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
	tftags "github.com/infracasts/terraform-provider-aws-expose-internal/tags"
)

// ListTags lists organizations service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn organizationsiface.OrganizationsAPI, identifier string) (tftags.KeyValueTags, error) {
	return ListTagsWithContext(context.Background(), conn, identifier)
}

func ListTagsWithContext(ctx context.Context, conn organizationsiface.OrganizationsAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &organizations.ListTagsForResourceInput{
		ResourceId: aws.String(identifier),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns organizations service tags.
func Tags(tags tftags.KeyValueTags) []*organizations.Tag {
	result := make([]*organizations.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &organizations.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from organizations service tags.
func KeyValueTags(tags []*organizations.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates organizations service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn organizationsiface.OrganizationsAPI, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTagsWithContext(context.Background(), conn, identifier, oldTags, newTags)
}
func UpdateTagsWithContext(ctx context.Context, conn organizationsiface.OrganizationsAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &organizations.UntagResourceInput{
			ResourceId: aws.String(identifier),
			TagKeys:    aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &organizations.TagResourceInput{
			ResourceId: aws.String(identifier),
			Tags:       Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
