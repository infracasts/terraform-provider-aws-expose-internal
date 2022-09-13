package redshiftdata

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-public/tfresource"
)

func statusStatement(conn *redshiftdataapiservice.RedshiftDataAPIService, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindStatementByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
