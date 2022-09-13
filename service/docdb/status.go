package docdb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-public/tfresource"
)

func statusEventSubscription(ctx context.Context, conn *docdb.DocDB, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindEventSubscriptionByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
