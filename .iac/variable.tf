variable "aws_iam_role_arn" {
  type        = string
  description = "The IAM Role to be used to manage the resources"
}

variable "aws_external_id" {
  type        = string
  description = "The external ID used to authenticate on AWS"
}

variable "environment_tag" {
  type        = string
  description = "The environment tag to be set on resources"
}

variable "environment" {
  type        = string
  description = "The name of the environment to which resources are being deployed"
}

variable "pipeline_id" {
  type        = string
  description = "A reference of the change that created/modified the resource"
}


variable "region" {
  type        = string
  description = "The AWS Region that will hold the resources to be provisioned"
  default     = "af-south-1"
}

variable "git_repository_url" {
  type        = string
  description = "The URL of the Git repository keeping the Infrastructure code"
}

variable "gitlab_token" {
  type        = string
  description = "The Gitlab token to be used for IaC automation"
}

variable "team_slug" {
  type = string
  description = "Team slug"
}

variable "business_unit" {
  type        = string
  description = "The Business unit that owns the resources"
}

variable "squad_email" {
  type    = string
  description = "The email of the team that owns the infrastructure"
}

variable "api_gateway_timeout" {
  type = number
  description = "Timeout of the API-Gateway when contacting the backend"
  default = 1000
}


variable "service_name" {
  type = string
  description = "The name of deployed service"
}

variable "app_docker_image" {
  type = string
  description = "The url of the application image"

}

variable "hostname" {
  type = string
  description = "Domain name"

}
