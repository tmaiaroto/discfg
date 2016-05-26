variable "aws_account_id" {}

variable "aws_region" {
	default = "us-east-1"
}

variable "stage" {
	default = "dev"
}

variable "discfg_table" {
	default = "discfg"
}

variable "api_name" {
	default = "Discfg"
}

variable "api_stage" {
	default = "dev"
}