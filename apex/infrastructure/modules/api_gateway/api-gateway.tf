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

# /{name}/cfg resource path
resource "aws_api_gateway_resource" "CfgResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_resource.NameResource.id}"
  path_part = "cfg"
}


# --------------------------------------------------- METHODS ----------------------------------------------------
# ---- Method Execution
# ---- /{name}/keys/{key} PUT
# 
# Creates the PUT method under /keys/{key}
resource "aws_api_gateway_method" "KeysPUTMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "PUT"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "KeysPUTIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysPUTMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_set_key/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "KeysPUTMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysPUTMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "KeysPUTIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysPUTMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.KeysPUTMethod200.status_code}"
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
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_get_key/invocations"
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
# ---- /{name}/keys/{key} DELETE
# 
# Creates the DELETE method under /keys/{key}
resource "aws_api_gateway_method" "KeysDELETEMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "DELETE"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "KeysDELETEIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysDELETEMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_delete_key/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "KeysDELETEMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysDELETEMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "KeysDELETEIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysKeyResource.id}"
  http_method = "${aws_api_gateway_method.KeysDELETEMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.KeysDELETEMethod200.status_code}"
}

# ---- Method Execution
# ---- /cfg/{name} PUT
# 
# Creates the PUT method under /cfg/{name}
resource "aws_api_gateway_method" "CfgPUTMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "PUT"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "CfgPUTIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgPUTMethod.http_method}"
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
resource "aws_api_gateway_method_response" "CfgPUTMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgPUTMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "CfgPUTIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgPUTMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.CfgPUTMethod200.status_code}"
}

# ---- Method Execution
# ---- /cfg/{name} PATCH
# 
# Creates the PATCH method under /cfg/{name}
resource "aws_api_gateway_method" "CfgPATCHMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "PATCH"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "CfgPATCHIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgPATCHMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_update_cfg/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "CfgPATCHMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgPATCHMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "CfgPATCHIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgPATCHMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.CfgPATCHMethod200.status_code}"
}

# ---- Method Execution
# ---- /cfg/{name} DELETE
# 
# Creates the DELETE method under /cfg/{name}
resource "aws_api_gateway_method" "CfgDELETEMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "DELETE"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "CfgDELETEIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgDELETEMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_delete_cfg/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "CfgDELETEMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgDELETEMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "CfgDELETEIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgDELETEMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.CfgDELETEMethod200.status_code}"
}

# ---- Method Execution
# ---- /cfg/{name} OPTIONS
# 
# Creates the OPTIONS method under /cfg/{name}
resource "aws_api_gateway_method" "CfgOPTIONSMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "OPTIONS"
  authorization = "NONE"
}
# Configures the integration for the Resource Method (in other words, what gets triggered)
# Client -> Method Request -> Integration Request -> *Integration*
resource "aws_api_gateway_integration" "CfgOPTIONSIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgOPTIONSMethod.http_method}"
  type = "AWS"
  # Must be POST for invoking Lambda function
  integration_http_method = "POST"
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_info_cfg/invocations"
  request_templates = {
    "application/json" = "${file("${path.module}/api_gateway_body_mapping.template")}"
  }
}
# Integration -> Integration Response -> *Method Response* -> Client
resource "aws_api_gateway_method_response" "CfgOPTIONSMethod200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgOPTIONSMethod.http_method}"
  status_code = "200"
}
# Integration -> *Integration Response* -> Method Response -> Client
resource "aws_api_gateway_integration_response" "CfgOPTIONSIntegrationResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.CfgResource.id}"
  http_method = "${aws_api_gateway_method.CfgOPTIONSMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.CfgOPTIONSMethod200.status_code}"
}

# -------------------------------------------------- DEPLOYMENT --------------------------------------------------
# Creates the API stage
resource "aws_api_gateway_deployment" "stage" {
 depends_on = ["aws_api_gateway_integration.KeysPUTIntegration", "aws_api_gateway_integration.KeysGETIntegration", "aws_api_gateway_integration.KeysDELETEIntegration", "aws_api_gateway_integration.CfgPUTIntegration", "aws_api_gateway_integration.CfgPATCHIntegration", "aws_api_gateway_integration.CfgDELETEIntegration", "aws_api_gateway_integration.CfgOPTIONSIntegration"]

 rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
 stage_name = "${var.api_gateway_stage}"
}