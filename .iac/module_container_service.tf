module "sandman-service" {
  source = "https://nexus.pkg.dev.vm.co.mz/repository/zip-public/tf-modules/aws-container-svc/2.0.0-RELEASE.zip"

  expose = {
    reachability="internal-private"
    hostname=var.hostname
  }

  tasks_deployment = {
    rolling_strategy = contains(["dev", "test", "demo"], var.environment)? "AGRESSIVE" : "CAUTIOUS"
    auto_rollback = true
  }




  depends_on = [
    aws_iam_role.sandman_iam_role
  ]

  name  = var.service_name


  # This is used to query environment variables in Gitlab
  environment = var.environment

  type = "backend"

  traffic_control = {
    from_composite_services = {
      allow = contains(["dev", "test", "demo"], var.environment)? "private_only" : "any"
    }
    from_container_services = {
      allow_any =  true
      allow_selected = []
    }
  }

  #
  iam_role_name = local.iam_role_name



  # Optional: The maximum time the envoy-proxy should wait for a response
  response_timeout_milis = 80000



  compute = {
    memory = var.memory
    cpu = var.cpu
  }


  containers = {

    # optional
    envoy = {
      log_level = contains(["dev", "test", "demo"], var.environment)? "info" : "off"
    }

    main = {

      #      health_check = {
      #      success_status=200
      #      target = {
      #        port=8080
      #        path="/health/liveness"
      #        }
      #      }

      image = var.app_docker_image
      #Allports  being exposed by container
      expose_ports = [local.port.container]
      # The  port to  Listen for requests in the mesh
      service_port = local.port.service_port


      variables = {

        clear_text = {

          "SQS_QUEUE_URL"="${aws_sqs_queue.compse_sqs_sandman.url}"

        }
      }

    }


  }

  logs_retention_days =  var.logs_retention_days

  auto_scaling = {

    enabled                = true

    replicas = {
      max                    = var.container_max_replicas
      min                    = var.container_min_replicas
    }

    targets  =   {
      cpu_utilization    = 80
      memory_utilization = 80
    }

    # Optional
    timing = {
      scale_in_cooldown  = 60
      scale_out_cooldown = 60
    }

  }

  observability = {
    alarm_evaluation_periods = {
      compute_alarms = 4
      task_count_alarms =2
      request_alarms=5
    }
    thresholds = {
      client_errors_percent  = 20
      server_errors_percent  = 20
      incoming_requests_per_min   = 5
    }
    create_alarms       = contains(["dev", "test", "demo"], var.environment)? false : true
    create_dashboard    = contains(["dev", "test", "demo"], var.environment)? false : true
  }



}
