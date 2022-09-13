// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package cloudformation

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	tftags "github.com/infracasts/terraform-provider-aws-expose-internal/tags"
)

// []*SERVICE.Tag handling

// Tags returns cloudformation service tags.
func Tags(tags tftags.KeyValueTags) []*cloudformation.Tag {
	result := make([]*cloudformation.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &cloudformation.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from cloudformation service tags.
func KeyValueTags(tags []*cloudformation.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}
