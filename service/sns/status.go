package sns

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-public/tfresource"
)

func statusSubscriptionPendingConfirmation(conn *sns.SNS, arn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindSubscriptionAttributesByARN(conn, arn)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, output[SubscriptionAttributeNamePendingConfirmation], nil
	}
}
