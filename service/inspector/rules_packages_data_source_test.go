package inspector_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/inspector"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/infracasts/terraform-provider-aws-public/acctest"
)

func TestAccInspectorRulesPackagesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, inspector.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRulesPackagesDataSourceConfig_basic,
				Check:  resource.TestCheckResourceAttrSet("data.aws_inspector_rules_packages.test", "arns.#"),
			},
		},
	})
}

const testAccRulesPackagesDataSourceConfig_basic = `
data "aws_inspector_rules_packages" "test" {}
`
