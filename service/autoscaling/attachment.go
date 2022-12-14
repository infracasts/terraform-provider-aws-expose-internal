package autoscaling

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/infracasts/terraform-provider-aws-expose-internal/conns"
	"github.com/infracasts/terraform-provider-aws-expose-internal/tfresource"
)

func ResourceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAttachmentCreate,
		Read:   resourceAttachmentRead,
		Delete: resourceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"alb_target_group_arn": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Deprecated:   "Use lb_target_group_arn instead",
				ExactlyOneOf: []string{"alb_target_group_arn", "elb", "lb_target_group_arn"},
			},
			"autoscaling_group_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"elb": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ExactlyOneOf: []string{"alb_target_group_arn", "elb", "lb_target_group_arn"},
			},
			"lb_target_group_arn": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ExactlyOneOf: []string{"alb_target_group_arn", "elb", "lb_target_group_arn"},
			},
		},
	}
}

func resourceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AutoScalingConn
	asgName := d.Get("autoscaling_group_name").(string)

	if v, ok := d.GetOk("elb"); ok {
		lbName := v.(string)
		input := &autoscaling.AttachLoadBalancersInput{
			AutoScalingGroupName: aws.String(asgName),
			LoadBalancerNames:    aws.StringSlice([]string{lbName}),
		}

		if _, err := conn.AttachLoadBalancers(input); err != nil {
			return fmt.Errorf("attaching Auto Scaling Group (%s) load balancer (%s): %w", asgName, lbName, err)
		}
	} else {
		var targetGroupARN string
		if v, ok := d.GetOk("alb_target_group_arn"); ok {
			targetGroupARN = v.(string)
		} else if v, ok := d.GetOk("lb_target_group_arn"); ok {
			targetGroupARN = v.(string)
		}

		input := &autoscaling.AttachLoadBalancerTargetGroupsInput{
			AutoScalingGroupName: aws.String(asgName),
			TargetGroupARNs:      aws.StringSlice([]string{targetGroupARN}),
		}

		if _, err := conn.AttachLoadBalancerTargetGroups(input); err != nil {
			return fmt.Errorf("attaching Auto Scaling Group (%s) target group (%s): %w", asgName, targetGroupARN, err)
		}
	}

	//lintignore:R016 // Allow legacy unstable ID usage in managed resource
	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%s-", asgName)))

	return resourceAttachmentRead(d, meta)
}

func resourceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AutoScalingConn
	asgName := d.Get("autoscaling_group_name").(string)

	var err error

	if v, ok := d.GetOk("elb"); ok {
		lbName := v.(string)
		err = FindAttachmentByLoadBalancerName(conn, asgName, lbName)
	} else {
		var targetGroupARN string
		if v, ok := d.GetOk("alb_target_group_arn"); ok {
			targetGroupARN = v.(string)
		} else if v, ok := d.GetOk("lb_target_group_arn"); ok {
			targetGroupARN = v.(string)
		}
		err = FindAttachmentByTargetGroupARN(conn, asgName, targetGroupARN)
	}

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Auto Scaling Group Attachment %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading Auto Scaling Group Attachment (%s): %w", d.Id(), err)
	}

	return nil
}

func resourceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AutoScalingConn
	asgName := d.Get("autoscaling_group_name").(string)

	if v, ok := d.GetOk("elb"); ok {
		lbName := v.(string)
		input := &autoscaling.DetachLoadBalancersInput{
			AutoScalingGroupName: aws.String(asgName),
			LoadBalancerNames:    aws.StringSlice([]string{lbName}),
		}

		if _, err := conn.DetachLoadBalancers(input); err != nil {
			return fmt.Errorf("detaching Auto Scaling Group (%s) load balancer (%s): %w", asgName, lbName, err)
		}
	} else {
		var targetGroupARN string
		if v, ok := d.GetOk("alb_target_group_arn"); ok {
			targetGroupARN = v.(string)
		} else if v, ok := d.GetOk("lb_target_group_arn"); ok {
			targetGroupARN = v.(string)
		}

		input := &autoscaling.DetachLoadBalancerTargetGroupsInput{
			AutoScalingGroupName: aws.String(asgName),
			TargetGroupARNs:      aws.StringSlice([]string{targetGroupARN}),
		}

		if _, err := conn.DetachLoadBalancerTargetGroups(input); err != nil {
			return fmt.Errorf("detaching Auto Scaling Group (%s) target group (%s): %w", asgName, targetGroupARN, err)
		}
	}

	return nil
}

func FindAttachmentByLoadBalancerName(conn *autoscaling.AutoScaling, asgName, loadBalancerName string) error {
	asg, err := FindGroupByName(conn, asgName)

	if err != nil {
		return err
	}

	for _, v := range asg.LoadBalancerNames {
		if aws.StringValue(v) == loadBalancerName {
			return nil
		}
	}

	return &resource.NotFoundError{
		LastError: fmt.Errorf("Auto Scaling Group (%s) load balancer (%s) attachment not found", asgName, loadBalancerName),
	}
}

func FindAttachmentByTargetGroupARN(conn *autoscaling.AutoScaling, asgName, targetGroupARN string) error {
	asg, err := FindGroupByName(conn, asgName)

	if err != nil {
		return err
	}

	for _, v := range asg.TargetGroupARNs {
		if aws.StringValue(v) == targetGroupARN {
			return nil
		}
	}

	return &resource.NotFoundError{
		LastError: fmt.Errorf("Auto Scaling Group (%s) target group (%s) attachment not found", asgName, targetGroupARN),
	}
}
