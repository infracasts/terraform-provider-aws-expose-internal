package apigatewayv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/infracasts/terraform-provider-aws-public/flex"
)

func flattenCaseInsensitiveStringSet(list []*string) *schema.Set {
	return schema.NewSet(hashStringCaseInsensitive, flex.FlattenStringList(list))
}
