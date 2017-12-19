package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsCognitoGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsCognitoGroupCreate,
		Read:   resourceAwsCognitoGroupRead,
		Update: resourceAwsCognitoGroupUpdate,
		Delete: resourceAwsCognitoGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateGroup.html
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCognitoGroupName,
			},
			"user_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateMaxLength(2048),
			},
			"precedence": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateCognitoGroupPrecedence,
			},
			"role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCognitoGroupArn,
			},
		},
	}
}

func resourceAwsCognitoGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	group := d.Get("group_name").(string)

	params := &cognitoidentityprovider.CreateGroupInput{
		GroupName:  aws.String(group),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		params.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("precedence"); ok {
		params.Precedence = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("role_arn"); ok {
		params.RoleArn = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Cognito Group: %s", params)

	_, err := conn.CreateGroup(params)
	if err != nil {
		return fmt.Errorf("Error creating Cognito Group: %s", err)
	}

	d.SetId(group)

	return resourceAwsCognitoGroupRead(d, meta)
}

func resourceAwsCognitoGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	group, err := conn.GetGroup(&cognitoidentityprovider.GetGroupInput{
		GroupName:  aws.String(d.Get("group_name").(string)),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	})

	log.Printf("[DEBUG] Reading Cognito Group: %s", d.Id())

	if err != nil {
		if isAWSErr(err, "ResourceNotFoundException", "") {
			log.Printf("[WARN] Cognito Group %q not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	g := group.Group

	d.Set("group_name", d.Id())
	d.Set("user_pool_id", g.UserPoolId)

	return nil
}

func resourceAwsCognitoGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn
	log.Print("[DEBUG] Updating Cognito Group")

	params := &cognitoidentityprovider.UpdateGroupInput{
		GroupName:  aws.String(d.Get("group_name").(string)),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	}

	if d.HasChange("group_name") {
		params.GroupName = aws.String(d.Get("group_name").(string))
	}

	if d.HasChange("description") {
		params.Description = aws.String(d.Get("description").(string))
	}

	if d.HasChange("precedence") {
		params.Precedence = aws.Int64(int64(d.Get("precedence").(int)))
	}

	if d.HasChange("role_arn") {
		params.RoleArn = aws.String(d.Get("role_arn").(string))
	}

	_, err := conn.UpdateGroup(params)
	if err != nil {
		return fmt.Errorf("Error updating Cognito Group: %s", err)
	}

	return resourceAwsCognitoGroupRead(d, meta)
}

func resourceAwsCognitoGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).cognitoidpconn

	log.Printf("[DEBUG] Deleting Cognito Group: %s", d.Id())

	_, err := conn.DeleteGroup(&cognitoidentityprovider.DeleteGroupInput{
		GroupName:  aws.String(d.Id()),
		UserPoolId: aws.String(d.Get("user_pool_id").(string)),
	})
	if err != nil {
		return err
	}

	return nil
}
