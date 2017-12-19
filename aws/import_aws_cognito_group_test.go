package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAWSCognitoGroup_importBasic(t *testing.T) {
	//resourceName := "aws_cognito_group.main"
	groupName := fmt.Sprintf("tf-acc-test-group-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	poolName := fmt.Sprintf("tf-acc-test-pool_%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCognitoGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSCognitoGroupConfig_basic(groupName, poolName),
			},
			//resource.TestStep{
			//	ResourceName:      resourceName,
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}
