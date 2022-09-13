package apigateway_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-public/acctest"
)

func init() {
	acctest.RegisterServiceErrorCheckFunc(apigateway.EndpointsID, testAccErrorCheckSkip)
}

// skips tests that have error messages indicating unsupported features
func testAccErrorCheckSkip(t *testing.T) resource.ErrorCheckFunc {
	return acctest.ErrorCheckSkipMessagesContaining(t,
		"no matching Route53Zone found",
	)
}
