package ec2

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/infracasts/terraform-provider-aws-expose-internal/conns"
	tftags "github.com/infracasts/terraform-provider-aws-expose-internal/tags"
	"github.com/infracasts/terraform-provider-aws-expose-internal/tfresource"
)

func DataSourceTransitGatewayVPNAttachment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTransitGatewayVPNAttachmentRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"filter": CustomFiltersSchema(),
			"tags":   tftags.TagsSchemaComputed(),
			"transit_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceTransitGatewayVPNAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &ec2.DescribeTransitGatewayAttachmentsInput{
		Filters: BuildAttributeFilterList(map[string]string{
			"resource-type": ec2.TransitGatewayAttachmentResourceTypeVpn,
		}),
	}

	input.Filters = append(input.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)

	if v, ok := d.GetOk("tags"); ok {
		input.Filters = append(input.Filters, BuildTagFilterList(
			Tags(tftags.New(v.(map[string]interface{}))),
		)...)
	}

	if v, ok := d.GetOk("vpn_connection_id"); ok {
		input.Filters = append(input.Filters, BuildAttributeFilterList(map[string]string{
			"resource-id": v.(string),
		})...)
	}

	if v, ok := d.GetOk("transit_gateway_id"); ok {
		input.Filters = append(input.Filters, BuildAttributeFilterList(map[string]string{
			"transit-gateway-id": v.(string),
		})...)
	}

	transitGatewayAttachment, err := FindTransitGatewayAttachment(conn, input)

	if err != nil {
		return tfresource.SingularDataSourceFindError("EC2 Transit Gateway VPN Attachment", err)
	}

	d.SetId(aws.StringValue(transitGatewayAttachment.TransitGatewayAttachmentId))
	d.Set("transit_gateway_id", transitGatewayAttachment.TransitGatewayId)
	d.Set("vpn_connection_id", transitGatewayAttachment.ResourceId)

	if err := d.Set("tags", KeyValueTags(transitGatewayAttachment.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	return nil
}
