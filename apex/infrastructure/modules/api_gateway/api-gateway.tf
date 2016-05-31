# Creates the Discfg API
resource "aws_api_gateway_rest_api" "DiscfgAPI" {
  name = "${var.api_gateway_api_name}"
  description = "A simple distributed configuration service"
}

# Notes
# ----------
# General API endpoint is: /{name}/keys/{key}
# This allows for getting and setting key values for a given discfg name

# ------------------------------------------------ RESOURCE PATHS ------------------------------------------------
# 
# Creates the /{name} API resource path (for config table name)
resource "aws_api_gateway_resource" "NameResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_rest_api.DiscfgAPI.root_resource_id}"
  path_part = "{name}"
}

# Creates the /keys API resource path
resource "aws_api_gateway_resource" "KeysResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_resource.NameResource.id}"
  path_part = "keys"
}

# /{name}/keys/{key} resource path
resource "aws_api_gateway_resource" "KeysKeyResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_resource.KeysResource.id}"
  path_part = "{key}"
}

# Creates the /cfg API resource path
resource "aws_api_gateway_resource" "CfgResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_rest_api.DiscfgAPI.root_resource_id}"
  path_part = "cfg"
}

# /cfg/create resource path
resource "aws_api_gateway_resource" "CfgCreateResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_resource.CfgResource.id}"
  path_part = "create"
}

# /cfg/create/{name} resource path
resource "aws_api_gateway_resource" "CfgCreateNameResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_resource.CfgCreateResource.id}"
  path_part = "{name}"
}


# --------------------------------------------------- METHODS ----------------------------------------------------
# ---- Method Execution
# ---- /{name}/keys/{key} POST
# 
# Creates the POST method under /keys/{key}
resource "aws_api_gateway_method" "KeysPOSTMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "POST"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "KeysPOSTIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysPOSTMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_set/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "KeysPOSTMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysPOSTMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "KeysPOSTIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysPOSTMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.KeysPOSTMethod200.status_code}"
}

# ---- Method Execution
# ---- /{name}/keys/{key} GET
# 
# Creates the GET method under /keys/{key}
resource "aws_api_gateway_method" "KeysGETMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "GET"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "KeysGETIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysGETMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_get/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "KeysGETMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysGETMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "KeysGETIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysGETMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.KeysGETMethod200.status_code}"
}

# ---- Method Execution
# ---- /cfg/create/{name} POST
# 
# Creates the POST method under /cfg/create/{name}
resource "aws_api_gateway_method" "CfgCreatePOSTMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgCreateNameResource.id}"
  http_method = "POST"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "CfgCreatePOSTIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgCreateNameResource.id}"
  http_method = "${aws_api_gateway_method.CfgCreatePOSTMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_create_cfg/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "CfgCreatePOSTMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgCreateNameResource.id}"
  http_method = "${aws_api_gateway_method.CfgCreatePOSTMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "CfgCreatePOSTIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgCreateNameResource.id}"
  http_method = "${aws_api_gateway_method.CfgCreatePOSTMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.CfgCreatePOSTMethod200.status_code}"
}

# -------------------------------------------------- DEPLOYMENT --------------------------------------------------
# Creates the API stage
resource "aws_api_gateway_deployment" "stage" {
 depends_on = ["aws_api_gateway_integration.KeysPOSTIntegration", "aws_api_gateway_integration.KeysGETIntegration", "aws_api_gateway_integration.CfgCreatePOSTIntegration"]

 rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
 stage_name = "${var.api_gateway_stage}"
}