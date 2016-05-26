variable "aws_region" {}

module "iam" {
  source = "../modules/iam"
}

module "api_gateway" {
  source = "../modules/api_gateway"
  api_gateway_aws_account_id = "${var.aws_account_id}"
  api_gateway_aws_region = "${var.aws_region}"
  api_gateway_invoke_discfg_lambda_role_arn = "${module.iam.gateway_invoke_discfg_lambda_role_arn}"
  api_gateway_stage = "${var.api_stage}"
  api_gateway_api_name = "${var.api_name}"
}

# Note, eventually this can be set by Discfg commands via the API.
# For now, this will be the simplest way to get it going.
module "dynamodb" {
  source = "../modules/dynamodb"
  dynamodb_discfg_table = "${var.discfg_table}"
}

output "lambda_function_role_id" {
  value = "${module.iam.lambda_function_role_id}"
}
