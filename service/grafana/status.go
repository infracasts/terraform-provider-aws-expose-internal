package grafana

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-public/tfresource"
)

func statusWorkspaceStatus(conn *managedgrafana.ManagedGrafana, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindWorkspaceByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusWorkspaceSAMLConfiguration(conn *managedgrafana.ManagedGrafana, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindSamlConfigurationByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
