package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSCognitoGroup_basic(t *testing.T) {
	groupName := fmt.Sprintf("tf-acc-test-group-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	poolName := fmt.Sprintf("tf-acc-test-pool_%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCognitoGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCognitoGroupConfig_basic(groupName, poolName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSCognitoGroupExists("aws_cognito_group.main"),
					resource.TestCheckResourceAttr("aws_cognito_group.main", "group_name", groupName),
					resource.TestCheckResourceAttr("aws_cognito_user_pool.main", "name", poolName),
				),
			},
		},
	})
}

func testAccCheckAWSCognitoGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Cognito Group ID set")
		}

		conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

		_, err := conn.GetGroup(&cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(rs.Primary.ID),
			UserPoolId: aws.String(rs.Primary.Attributes["user_pool_id"]),
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAWSCognitoGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).cognitoidpconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cognito_group" {
			continue
		}

		_, err := conn.GetGroup(&cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(rs.Primary.ID),
			UserPoolId: aws.String(rs.Primary.Attributes["user_pool_id"]),
		})

		if err != nil {
			if isAWSErr(err, "ResourceNotFoundException", "") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccAWSCognitoGroupConfig_basic(groupName, poolName string) string {
	return fmt.Sprintf(`
resource "aws_cognito_group" "main" {
  group_name = "%s"
  user_pool_id = "${aws_cognito_user_pool.main.id}"
}
resource "aws_cognito_user_pool" "main" {
  name = "%s"
}
`, groupName, poolName)
}
