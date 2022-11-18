locals {

  project = {
    slug = "sandman-service"
  }
  naming = {
    prefix = "${var.team_slug}-${local.project.slug}-${var.environment}"
  }

  resource_tags = {

    Project         = "Sandman Orchestration Service"
    BusinessUnit    = var.business_unit
    ManagedBy       = var.squad_email
    PII             = "No"
    GitRepository   = var.git_repository_url
    Confidentiality = "C3"
    Environment     = "${var.environment_tag}"
    TaggingVersion  = "V2.4"
  }

  port = {
    container      = 8080
    service_port   = 8080
  }

  iam_role_name = "sandman-service-ecs-task-iam-role"
}