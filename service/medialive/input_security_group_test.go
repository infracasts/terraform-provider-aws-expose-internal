package medialive_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/medialive"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/infracasts/terraform-provider-aws-expose-internal/acctest"
	"github.com/infracasts/terraform-provider-aws-expose-internal/conns"
	"github.com/infracasts/terraform-provider-aws-expose-internal/create"
	tfmedialive "github.com/infracasts/terraform-provider-aws-expose-internal/service/medialive"
	"github.com/infracasts/terraform-provider-aws-expose-internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccMediaLiveInputSecurityGroup_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var inputSecurityGroup medialive.DescribeInputSecurityGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_input_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccInputSecurityGroupsPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInputSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInputSecurityGroupConfig_basic(rName, "10.0.0.8/32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputSecurityGroupExists(resourceName, &inputSecurityGroup),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "whitelist_rules.*", map[string]string{
						"cidr": "10.0.0.8/32",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMediaLiveInputSecurityGroup_updateCIDR(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var inputSecurityGroup medialive.DescribeInputSecurityGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_input_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccInputSecurityGroupsPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInputSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInputSecurityGroupConfig_basic(rName, "10.0.0.8/32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputSecurityGroupExists(resourceName, &inputSecurityGroup),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "whitelist_rules.*", map[string]string{
						"cidr": "10.0.0.8/32",
					}),
				),
			},
			{
				Config: testAccInputSecurityGroupConfig_basic(rName, "10.2.0.0/16"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputSecurityGroupExists(resourceName, &inputSecurityGroup),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "whitelist_rules.*", map[string]string{
						"cidr": "10.2.0.0/16",
					}),
				),
			},
		},
	})
}

func TestAccMediaLiveInputSecurityGroup_updateTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var inputSecurityGroup medialive.DescribeInputSecurityGroupOutput
	resourceName := "aws_medialive_input_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccInputSecurityGroupsPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInputSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInputSecurityGroupConfig_tags1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputSecurityGroupExists(resourceName, &inputSecurityGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccInputSecurityGroupConfig_tags2("key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputSecurityGroupExists(resourceName, &inputSecurityGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccInputSecurityGroupConfig_tags1("key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputSecurityGroupExists(resourceName, &inputSecurityGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccMediaLiveInputSecurityGroup_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var inputSecurityGroup medialive.DescribeInputSecurityGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_input_security_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccInputSecurityGroupsPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInputSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInputSecurityGroupConfig_basic(rName, "10.0.0.8/32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputSecurityGroupExists(resourceName, &inputSecurityGroup),
					acctest.CheckResourceDisappears(acctest.Provider, tfmedialive.ResourceInputSecurityGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckInputSecurityGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).MediaLiveConn
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_medialive_input_security_group" {
			continue
		}

		_, err := tfmedialive.FindInputSecurityGroupByID(ctx, conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return create.Error(names.MediaLive, create.ErrActionCheckingDestroyed, tfmedialive.ResNameInputSecurityGroup, rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckInputSecurityGroupExists(name string, inputSecurityGroup *medialive.DescribeInputSecurityGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.MediaLive, create.ErrActionCheckingExistence, tfmedialive.ResNameInputSecurityGroup, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.MediaLive, create.ErrActionCheckingExistence, tfmedialive.ResNameInputSecurityGroup, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).MediaLiveConn
		ctx := context.Background()
		resp, err := tfmedialive.FindInputSecurityGroupByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return create.Error(names.MediaLive, create.ErrActionCheckingExistence, tfmedialive.ResNameInputSecurityGroup, rs.Primary.ID, err)
		}

		*inputSecurityGroup = *resp

		return nil
	}
}

func testAccInputSecurityGroupsPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).MediaLiveConn
	ctx := context.Background()

	input := &medialive.ListInputSecurityGroupsInput{}
	_, err := conn.ListInputSecurityGroups(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccInputSecurityGroupConfig_basic(rName, cidr string) string {
	return fmt.Sprintf(`
resource "aws_medialive_input_security_group" "test" {
  whitelist_rules {
    cidr = %[2]q
  }

  tags = {
    Name = %[1]q
  }
}
`, rName, cidr)
}

func testAccInputSecurityGroupConfig_tags1(key1, value1 string) string {
	return acctest.ConfigCompose(
		fmt.Sprintf(`
resource "aws_medialive_input_security_group" "test" {
  whitelist_rules {
    cidr = "10.2.0.0/16"
  }

  tags = {
    %[1]q = %[2]q
  }
}
`, key1, value1))
}

func testAccInputSecurityGroupConfig_tags2(key1, value1, key2, value2 string) string {
	return acctest.ConfigCompose(
		fmt.Sprintf(`
resource "aws_medialive_input_security_group" "test" {
  whitelist_rules {
    cidr = "10.2.0.0/16"
  }

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, key1, value1, key2, value2))
}
