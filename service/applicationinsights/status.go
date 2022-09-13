package applicationinsights

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationinsights"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-expose-internal/tfresource"
)

func statusApplication(conn *applicationinsights.ApplicationInsights, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindApplicationByName(conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.LifeCycle), nil
	}
}
