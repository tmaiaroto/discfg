# Creates the Discfg API
resource "aws_api_gateway_rest_api" "DiscfgAPI" {
  name = "${var.api_gateway_api_name}"
  description = "A simple distributed configuration service"
}

# Creates the /keys API resource path
resource "aws_api_gateway_resource" "KeysResource" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  parent_id = "${aws_api_gateway_rest_api.DiscfgAPI.root_resource_id}"
  path_part = "keys"
}

# Creates the POST method under /keys
resource "aws_api_gateway_method" "KeysMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysResource.id}"
  http_method = "POST"
  authorization = "NONE"
}

# Creates a method response; 200, 404, 500 etc.
resource "aws_api_gateway_method_response" "200" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysResource.id}"
  http_method = "${aws_api_gateway_method.KeysMethod.http_method}"
  status_code = "200"
  response_models = {
    "application/json" = "Empty"
  }
}

# Configures the integration for the Resource Method (what gets triggered for GET /keys)
resource "aws_api_gateway_integration" "KeysPOSTIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysResource.id}"
  http_method = "${aws_api_gateway_method.KeysMethod.http_method}"
  type = "AWS"
  integration_http_method = "POST" # Must be POST for invoking Lambda function
  credentials = "${var.api_gateway_invoke_discfg_lambda_role_arn}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.api_gateway_aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.api_gateway_aws_region}:${var.api_gateway_aws_account_id}:function:discfg_set/invocations"
}

# Configures the integration (Lambda) response that maps to the method response via the status_code
resource "aws_api_gateway_integration_response" "KeysResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
  resource_id = "${aws_api_gateway_resource.KeysResource.id}"
  http_method = "${aws_api_gateway_method.KeysMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.200.status_code}"
}

# Creates the API stage
resource "aws_api_gateway_deployment" "stage" {
 depends_on = ["aws_api_gateway_integration.KeysPOSTIntegration"]

 rest_api_id = "${aws_api_gateway_rest_api.DiscfgAPI.id}"
 stage_name = "${var.api_gateway_stage}"
}