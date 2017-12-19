--
layout: "aws"
page_title: "AWS: aws_cognito_group"
side_bar_current: "docs-aws-resource-cognito-group"
description: |-
  Provides a Cognito Group resource.
---

# aws_cognito_group

Provides a Cognito Group resource.

## Example Usage

```hcl
resource "aws_cognito_group" "main" {
  group_name = "example-group"
  user_pool_id = "${aws_cognito_user_pool.example.id}"
}

resource "aws_cognito_user_pool" "example" {
  name = "example-pool"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required) The group name string.
* `user_pool_id` - (Required) The user pool ID.
* `description` - (Optional) The group description.
* `precedence` - (Optional) The group precedence.
* `role_arn` - (Optional) The role ARN for the group.
