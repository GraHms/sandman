terraform {
  backend "s3" {}

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }

      external = {
      source = "hashicorp/external"
      version = "2.2.2"
    }

  }
}



provider "aws" {
  region = var.region

  assume_role {
    role_arn     = var.aws_iam_role_arn
    session_name = "GL-${var.pipeline_id}"
    external_id  = var.aws_external_id
  }

  default_tags {
    tags = local.resource_tags
  }
}
