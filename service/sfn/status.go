package sfn

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-expose-internal/tfresource"
)

func statusStateMachine(conn *sfn.SFN, stateMachineArn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindStateMachineByARN(conn, stateMachineArn)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
