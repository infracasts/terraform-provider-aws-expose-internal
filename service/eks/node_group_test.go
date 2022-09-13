package eks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/infracasts/terraform-provider-aws-expose-internal/acctest"
	"github.com/infracasts/terraform-provider-aws-expose-internal/conns"
	tfeks "github.com/infracasts/terraform-provider-aws-expose-internal/service/eks"
	"github.com/infracasts/terraform-provider-aws-expose-internal/tfresource"
)

func init() {
	acctest.RegisterServiceErrorCheckFunc(eks.EndpointsID, testAccErrorCheckSkip)

}

func TestAccEKSNodeGroup_basic(t *testing.T) {
	var nodeGroup eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	eksClusterResourceName := "aws_eks_cluster.test"
	iamRoleResourceName := "aws_iam_role.node"
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_dataSourceName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup),
					resource.TestCheckResourceAttr(resourceName, "ami_type", eks.AMITypesAl2X8664),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "eks", regexp.MustCompile(fmt.Sprintf("nodegroup/%[1]s/%[1]s/.+", rName))),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_name", eksClusterResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "capacity_type", eks.CapacityTypesOnDemand),
					resource.TestCheckResourceAttr(resourceName, "disk_size", "20"),
					resource.TestCheckResourceAttr(resourceName, "instance_types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "labels.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "node_group_name_prefix", ""),
					resource.TestCheckResourceAttrPair(resourceName, "node_role_arn", iamRoleResourceName, "arn"),
					resource.TestMatchResourceAttr(resourceName, "release_version", regexp.MustCompile(`^\d+\.\d+\.\d+-\d{8}$`)),
					resource.TestCheckResourceAttr(resourceName, "remote_access.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.autoscaling_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", eks.NodegroupStatusActive),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "taint.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "update_config.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "version", eksClusterResourceName, "version"),
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

func TestAccEKSNodeGroup_Name_generated(t *testing.T) {
	var nodeGroup eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_nameGenerated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup),
					acctest.CheckResourceAttrNameGenerated(resourceName, "node_group_name"),
					resource.TestCheckResourceAttr(resourceName, "node_group_name_prefix", "terraform-"),
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

func TestAccEKSNodeGroup_namePrefix(t *testing.T) {
	var nodeGroup eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_namePrefix(rName, "tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "node_group_name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "node_group_name_prefix", "tf-acc-test-prefix-"),
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

func TestAccEKSNodeGroup_disappears(t *testing.T) {
	var nodeGroup eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_dataSourceName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup),
					acctest.CheckResourceDisappears(acctest.Provider, tfeks.ResourceNodeGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEKSNodeGroup_amiType(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_amiType(rName, eks.AMITypesAl2X8664Gpu),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "ami_type", eks.AMITypesAl2X8664Gpu),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_amiType(rName, eks.AMITypesAl2Arm64),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "ami_type", eks.AMITypesAl2Arm64),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_CapacityType_spot(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_capacityType(rName, eks.CapacityTypesSpot),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "capacity_type", eks.CapacityTypesSpot),
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

func TestAccEKSNodeGroup_diskSize(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_diskSize(rName, 21),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "disk_size", "21"),
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

func TestAccEKSNodeGroup_forceUpdateVersion(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_forceUpdateVersion(rName, "1.19"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "version", "1.19"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_update_version"},
			},
			{
				Config: testAccNodeGroupConfig_forceUpdateVersion(rName, "1.20"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "version", "1.20"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_InstanceTypes_multiple(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"
	instanceTypes := fmt.Sprintf("%q, %q, %q, %q", "t2.medium", "t3.medium", "t2.large", "t3.large")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_instanceTypesMultiple(rName, instanceTypes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "instance_types.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "instance_types.0", "t2.medium"),
					resource.TestCheckResourceAttr(resourceName, "instance_types.1", "t3.medium"),
					resource.TestCheckResourceAttr(resourceName, "instance_types.2", "t2.large"),
					resource.TestCheckResourceAttr(resourceName, "instance_types.3", "t3.large"),
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

func TestAccEKSNodeGroup_InstanceTypes_single(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_instanceTypesSingle(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "instance_types.#", "1"),
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

func TestAccEKSNodeGroup_labels(t *testing.T) {
	var nodeGroup1, nodeGroup2, nodeGroup3 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_labels1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "labels.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "labels.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_labels2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "labels.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "labels.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "labels.key2", "value2"),
				),
			},
			{
				Config: testAccNodeGroupConfig_labels1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup3),
					resource.TestCheckResourceAttr(resourceName, "labels.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "labels.key2", "value2"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_LaunchTemplate_id(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	launchTemplateResourceName1 := "aws_launch_template.test1"
	launchTemplateResourceName2 := "aws_launch_template.test2"
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_launchTemplateId1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "launch_template.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.id", launchTemplateResourceName1, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_launchTemplateId2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "launch_template.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.id", launchTemplateResourceName2, "id"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_LaunchTemplate_name(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	launchTemplateResourceName1 := "aws_launch_template.test1"
	launchTemplateResourceName2 := "aws_launch_template.test2"
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_launchTemplateName1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "launch_template.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.name", launchTemplateResourceName1, "name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_launchTemplateName2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "launch_template.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.name", launchTemplateResourceName2, "name"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_LaunchTemplate_version(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	launchTemplateResourceName := "aws_launch_template.test"
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_launchTemplateVersion1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "launch_template.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.version", launchTemplateResourceName, "default_version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_launchTemplateVersion2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "launch_template.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "launch_template.0.version", launchTemplateResourceName, "default_version"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_releaseVersion(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	ssmParameterDataSourceName := "data.aws_ssm_parameter.test"
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_releaseVersion(rName, "1.17"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttrPair(resourceName, "release_version", ssmParameterDataSourceName, "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_releaseVersion(rName, "1.18"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttrPair(resourceName, "release_version", ssmParameterDataSourceName, "value"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_RemoteAccess_ec2SSHKey(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_remoteAccessEC2SSHKey(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "remote_access.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "remote_access.0.ec2_ssh_key", rName),
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

func TestAccEKSNodeGroup_RemoteAccess_sourceSecurityGroupIDs(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_remoteAccessSourceSecurityIds1(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "remote_access.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "remote_access.0.source_security_group_ids.#", "1"),
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

func TestAccEKSNodeGroup_Scaling_desiredSize(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 2, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 1, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "1"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_Scaling_maxSize(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 1, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 1, 1, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "1"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_Scaling_minSize(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 2, 2, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 2, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "1"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_ScalingZeroDesiredSize_minSize(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 0, 1, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 1, 2, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "1"),
				),
			},
			{
				Config: testAccNodeGroupConfig_scalingSizes(rName, 0, 1, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.desired_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.max_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "scaling_config.0.min_size", "0"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_tags(t *testing.T) {
	var nodeGroup1, nodeGroup2, nodeGroup3 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccNodeGroupConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup3),
					testAccCheckNodeGroupNotRecreated(&nodeGroup2, &nodeGroup3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_taints(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_taints1(rName, "key1", "value1", "NO_SCHEDULE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "taint.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "taint.*", map[string]string{
						"key":    "key1",
						"value":  "value1",
						"effect": "NO_SCHEDULE",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_taints2(rName,
					"key1", "value1updated", "NO_EXECUTE",
					"key2", "value2", "NO_SCHEDULE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "taint.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "taint.*", map[string]string{
						"key":    "key1",
						"value":  "value1updated",
						"effect": "NO_EXECUTE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "taint.*", map[string]string{
						"key":    "key2",
						"value":  "value2",
						"effect": "NO_SCHEDULE",
					}),
				),
			},
			{
				Config: testAccNodeGroupConfig_taints1(rName, "key2", "value2", "NO_SCHEDULE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "taint.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "taint.*", map[string]string{
						"key":    "key2",
						"value":  "value2",
						"effect": "NO_SCHEDULE",
					}),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_update(t *testing.T) {
	var nodeGroup1 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "update_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "update_config.0.max_unavailable", "1"),
					resource.TestCheckResourceAttr(resourceName, "update_config.0.max_unavailable_percentage", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "update_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "update_config.0.max_unavailable", "0"),
					resource.TestCheckResourceAttr(resourceName, "update_config.0.max_unavailable_percentage", "40"),
				),
			},
		},
	})
}

func TestAccEKSNodeGroup_version(t *testing.T) {
	var nodeGroup1, nodeGroup2 eks.Nodegroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_node_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, eks.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeGroupConfig_version(rName, "1.19"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup1),
					resource.TestCheckResourceAttr(resourceName, "version", "1.19"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNodeGroupConfig_version(rName, "1.20"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodeGroupExists(resourceName, &nodeGroup2),
					testAccCheckNodeGroupNotRecreated(&nodeGroup1, &nodeGroup2),
					resource.TestCheckResourceAttr(resourceName, "version", "1.20"),
				),
			},
		},
	})
}

func testAccErrorCheckSkip(t *testing.T) resource.ErrorCheckFunc {
	return acctest.ErrorCheckSkipMessagesContaining(t,
		"InvalidParameterException: The following supplied instance types do not exist",
	)
}

func testAccCheckNodeGroupExists(resourceName string, nodeGroup *eks.Nodegroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EKS Node Group ID is set")
		}

		clusterName, nodeGroupName, err := tfeks.NodeGroupParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EKSConn

		output, err := tfeks.FindNodegroupByClusterNameAndNodegroupName(conn, clusterName, nodeGroupName)

		if err != nil {
			return err
		}

		*nodeGroup = *output

		return nil
	}
}

func testAccCheckNodeGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EKSConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_eks_node_group" {
			continue
		}

		clusterName, nodeGroupName, err := tfeks.NodeGroupParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		_, err = tfeks.FindNodegroupByClusterNameAndNodegroupName(conn, clusterName, nodeGroupName)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EKS Node Group %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckNodeGroupNotRecreated(i, j *eks.Nodegroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !aws.TimeValue(i.CreatedAt).Equal(aws.TimeValue(j.CreatedAt)) {
			return fmt.Errorf("EKS Node Group (%s) was recreated", aws.StringValue(j.NodegroupName))
		}

		return nil
	}
}

func testAccCheckNodeGroupRecreated(i, j *eks.Nodegroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.TimeValue(i.CreatedAt).Equal(aws.TimeValue(j.CreatedAt)) {
			return fmt.Errorf("EKS Node Group (%s) was not recreated", aws.StringValue(j.NodegroupName))
		}

		return nil
	}
}

func testAccNodeGroupBaseIAMAndVPCConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_partition" "current" {}

resource "aws_iam_role" "cluster" {
  name = "%[1]s-cluster"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = [
          "eks.${data.aws_partition.current.dns_suffix}",
          "eks-nodegroup.${data.aws_partition.current.dns_suffix}",
        ]
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "cluster-AmazonEKSClusterPolicy" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.cluster.name
}

resource "aws_iam_role" "node" {
  name = "%[1]s-node"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.${data.aws_partition.current.dns_suffix}"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "node-AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.node.name
}

resource "aws_iam_role_policy_attachment" "node-AmazonEKS_CNI_Policy" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.node.name
}

resource "aws_iam_role_policy_attachment" "node-AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.node.name
}

resource "aws_vpc" "test" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name                          = %[1]q
    "kubernetes.io/cluster/%[1]s" = "shared"
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.test.id
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_main_route_table_association" "test" {
  route_table_id = aws_route_table.test.id
  vpc_id         = aws_vpc.test.id
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = -1
    to_port     = 0
  }

  ingress {
    cidr_blocks = [aws_vpc.test.cidr_block]
    from_port   = 0
    protocol    = -1
    to_port     = 0
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  count = 2

  availability_zone       = data.aws_availability_zones.available.names[count.index]
  cidr_block              = cidrsubnet(aws_vpc.test.cidr_block, 8, count.index)
  map_public_ip_on_launch = true
  vpc_id                  = aws_vpc.test.id

  tags = {
    Name                          = %[1]q
    "kubernetes.io/cluster/%[1]s" = "shared"
  }
}
`, rName)
}

func testAccNodeGroupBaseConfig(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseIAMAndVPCConfig(rName),
		fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.cluster.arn

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [
    aws_iam_role_policy_attachment.cluster-AmazonEKSClusterPolicy,
    aws_main_route_table_association.test,
  ]
}
`, rName))
}

func testAccNodeGroupBaseVersionConfig(rName string, version string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseIAMAndVPCConfig(rName),
		fmt.Sprintf(`
resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.cluster.arn
  version  = %[2]q

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [
    aws_iam_role_policy_attachment.cluster-AmazonEKSClusterPolicy,
    aws_main_route_table_association.test,
  ]
}
`, rName, version))
}

func testAccNodeGroupConfig_dataSourceName(rName string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_nameGenerated(rName string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), `
resource "aws_eks_node_group" "test" {
  cluster_name  = aws_eks_cluster.test.name
  node_role_arn = aws_iam_role.node.arn
  subnet_ids    = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`)
}

func testAccNodeGroupConfig_namePrefix(rName, namePrefix string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name           = aws_eks_cluster.test.name
  node_group_name_prefix = %[1]q
  node_role_arn          = aws_iam_role.node.arn
  subnet_ids             = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    "aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy",
    "aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy",
    "aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly",
  ]
}
`, namePrefix))
}

func testAccNodeGroupConfig_amiType(rName, amiType string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  ami_type        = %[2]q
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, amiType))
}

func testAccNodeGroupConfig_capacityType(rName, capacityType string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  capacity_type   = %[2]q
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, capacityType))
}

func testAccNodeGroupConfig_diskSize(rName string, diskSize int) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  disk_size       = %[2]d
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, diskSize))
}

func testAccNodeGroupConfig_forceUpdateVersion(rName, version string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseVersionConfig(rName, version), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name         = aws_eks_cluster.test.name
  force_update_version = true
  node_group_name      = %[1]q
  node_role_arn        = aws_iam_role.node.arn
  subnet_ids           = aws_subnet.test[*].id
  version              = aws_eks_cluster.test.version

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_instanceTypesMultiple(rName, instanceTypes string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name = aws_eks_cluster.test.name
  # use a predetermined string instead of aws_ec2_instance_type_offerings data source for the instance_types (TypeList)
  # as the apply-time values of the data source's instance_types can change in order
  instance_types  = [%[1]s]
  node_group_name = %[2]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, instanceTypes, rName))
}

func testAccNodeGroupConfig_instanceTypesSingle(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
data "aws_ec2_instance_type_offering" "available" {
  filter {
    name   = "instance-type"
    values = ["t3.large", "t2.large"]
  }

  preferred_instance_types = ["t3.large", "t2.large"]
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  instance_types  = [data.aws_ec2_instance_type_offering.available.instance_type]
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_labels1(rName, labelKey1, labelValue1 string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  labels = {
    %[2]q = %[3]q
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, labelKey1, labelValue1))
}

func testAccNodeGroupConfig_labels2(rName, labelKey1, labelValue1, labelKey2, labelValue2 string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  labels = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, labelKey1, labelValue1, labelKey2, labelValue2))
}

func testAccNodeGroupConfig_launchTemplateId1(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
data "aws_ssm_parameter" "test" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.test.version}/amazon-linux-2/recommended/image_id"
}

resource "aws_launch_template" "test1" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-1"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_launch_template" "test2" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-2"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  launch_template {
    id      = aws_launch_template.test1.id
    version = aws_launch_template.test1.default_version
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_launchTemplateId2(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
data "aws_ssm_parameter" "test" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.test.version}/amazon-linux-2/recommended/image_id"
}

resource "aws_launch_template" "test1" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-1"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_launch_template" "test2" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-2"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  launch_template {
    id      = aws_launch_template.test2.id
    version = aws_launch_template.test2.default_version
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_launchTemplateName1(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
data "aws_ssm_parameter" "test" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.test.version}/amazon-linux-2/recommended/image_id"
}

resource "aws_launch_template" "test1" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-1"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_launch_template" "test2" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-2"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  launch_template {
    name    = aws_launch_template.test1.name
    version = aws_launch_template.test1.default_version
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_launchTemplateName2(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
data "aws_ssm_parameter" "test" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.test.version}/amazon-linux-2/recommended/image_id"
}

resource "aws_launch_template" "test1" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-1"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_launch_template" "test2" {
  image_id      = data.aws_ssm_parameter.test.value
  instance_type = "t3.medium"
  name          = "%[1]s-2"
  user_data     = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  launch_template {
    name    = aws_launch_template.test2.name
    version = aws_launch_template.test2.default_version
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_launchTemplateVersion1(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
data "aws_ssm_parameter" "test" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.test.version}/amazon-linux-2/recommended/image_id"
}

resource "aws_launch_template" "test" {
  image_id               = data.aws_ssm_parameter.test.value
  instance_type          = "t3.medium"
  name                   = %[1]q
  update_default_version = true
  user_data              = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  launch_template {
    name    = aws_launch_template.test.name
    version = aws_launch_template.test.default_version
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_launchTemplateVersion2(rName string) string {
	return acctest.ConfigCompose(
		testAccNodeGroupBaseConfig(rName),
		fmt.Sprintf(`
data "aws_ssm_parameter" "test" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.test.version}/amazon-linux-2/recommended/image_id"
}

resource "aws_launch_template" "test" {
  image_id               = data.aws_ssm_parameter.test.value
  instance_type          = "t3.large"
  name                   = %[1]q
  update_default_version = true
  user_data              = base64encode(templatefile("testdata/node-group-launch-template-user-data.sh.tmpl", { cluster_name = aws_eks_cluster.test.name }))
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  launch_template {
    name    = aws_launch_template.test.name
    version = aws_launch_template.test.default_version
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_releaseVersion(rName string, version string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseVersionConfig(rName, version), fmt.Sprintf(`
data "aws_ssm_parameter" "test" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.test.version}/amazon-linux-2/recommended/release_version"
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  release_version = data.aws_ssm_parameter.test.value
  subnet_ids      = aws_subnet.test[*].id
  version         = aws_eks_cluster.test.version

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_remoteAccessEC2SSHKey(rName, publicKey string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = %[2]q
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  remote_access {
    ec2_ssh_key = aws_key_pair.test.key_name
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, publicKey))
}

func testAccNodeGroupConfig_remoteAccessSourceSecurityIds1(rName, publicKey string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = %[2]q
}

resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  remote_access {
    ec2_ssh_key               = aws_key_pair.test.key_name
    source_security_group_ids = [aws_security_group.test.id]
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, publicKey))
}

func testAccNodeGroupConfig_scalingSizes(rName string, desiredSize, maxSize, minSize int) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = %[2]d
    max_size     = %[3]d
    min_size     = %[4]d
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, desiredSize, maxSize, minSize))
}

func testAccNodeGroupConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  tags = {
    %[2]q = %[3]q
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, tagKey1, tagValue1))
}

func testAccNodeGroupConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccNodeGroupConfig_taints1(rName, taintKey1, taintValue1, taintEffect1 string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  taint {
    key    = %[2]q
    value  = %[3]q
    effect = %[4]q
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, taintKey1, taintValue1, taintEffect1))
}

func testAccNodeGroupConfig_taints2(rName, taintKey1, taintValue1, taintEffect1, taintKey2, taintValue2, taintEffect2 string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  taint {
    key    = %[2]q
    value  = %[3]q
    effect = %[4]q
  }

  taint {
    key    = %[5]q
    value  = %[6]q
    effect = %[7]q
  }

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName, taintKey1, taintValue1, taintEffect1, taintKey2, taintValue2, taintEffect2))
}

func testAccNodeGroupConfig_update1(rName string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 3
    min_size     = 1
  }

  update_config {
    max_unavailable = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_update2(rName string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseConfig(rName), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id

  scaling_config {
    desired_size = 1
    max_size     = 3
    min_size     = 1
  }

  update_config {
    max_unavailable_percentage = 40
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}

func testAccNodeGroupConfig_version(rName, version string) string {
	return acctest.ConfigCompose(testAccNodeGroupBaseVersionConfig(rName, version), fmt.Sprintf(`
resource "aws_eks_node_group" "test" {
  cluster_name    = aws_eks_cluster.test.name
  node_group_name = %[1]q
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.test[*].id
  version         = aws_eks_cluster.test.version

  scaling_config {
    desired_size = 1
    max_size     = 1
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node-AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.node-AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.node-AmazonEC2ContainerRegistryReadOnly,
  ]
}
`, rName))
}
