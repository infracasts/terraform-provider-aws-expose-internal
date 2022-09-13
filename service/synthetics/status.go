package synthetics

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/synthetics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-expose-internal/tfresource"
)

func statusCanaryState(conn *synthetics.Synthetics, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindCanaryByName(conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status.State), nil
	}
}
