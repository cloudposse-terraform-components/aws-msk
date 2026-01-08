locals {
  ssm_path_prefix = length(var.ssm_cluster_name_override) > 0 ? format("/%s/%s", var.ssm_path_prefix, var.ssm_cluster_name_override) : format("/%s/%s", var.ssm_path_prefix, module.this.name)

  all_kafka_parameters = [
    {
      name        = format("%s/broker_endpoints", local.ssm_path_prefix)
      value       = join(",", module.kafka.broker_endpoints)
      description = "List of broker endpoints"
      type        = "SecureString"
    },
    {
      name        = format("%s/bootstrap_brokers", local.ssm_path_prefix)
      value       = module.kafka.bootstrap_brokers
      description = "Bootstrap brokers"
      type        = "SecureString"
    },
    {
      name        = format("%s/bootstrap_brokers_tls", local.ssm_path_prefix)
      value       = module.kafka.bootstrap_brokers_tls
      description = "TLS bootstrap brokers"
      type        = "SecureString"
    },
    {
      name        = format("%s/bootstrap_brokers_public_tls", local.ssm_path_prefix)
      value       = module.kafka.bootstrap_brokers_public_tls
      description = "Public TLS bootstrap brokers"
      type        = "SecureString"
    },
    {
      name        = format("%s/bootstrap_brokers_sasl_scram", local.ssm_path_prefix)
      value       = module.kafka.bootstrap_brokers_sasl_scram
      description = "SASL SCRAM bootstrap brokers"
      type        = "SecureString"
    },
    {
      name        = format("%s/bootstrap_brokers_public_sasl_scram", local.ssm_path_prefix)
      value       = module.kafka.bootstrap_brokers_public_sasl_scram
      description = "Public SASL SCRAM bootstrap brokers"
      type        = "SecureString"
    },
    {
      name        = format("%s/bootstrap_brokers_sasl_iam", local.ssm_path_prefix)
      value       = module.kafka.bootstrap_brokers_sasl_iam
      description = "SASL IAM bootstrap brokers"
      type        = "SecureString"
    },
    {
      name        = format("%s/bootstrap_brokers_public_sasl_iam", local.ssm_path_prefix)
      value       = module.kafka.bootstrap_brokers_public_sasl_iam
      description = "Public SASL IAM bootstrap brokers"
      type        = "SecureString"
    },
    {
      name        = format("%s/zookeeper_connect_string", local.ssm_path_prefix)
      value       = module.kafka.zookeeper_connect_string
      description = "Zookeeper connect string"
      type        = "SecureString"
    },
    {
      name        = format("%s/zookeeper_connect_string_tls", local.ssm_path_prefix)
      value       = module.kafka.zookeeper_connect_string_tls
      description = "TLS Zookeeper connect string"
      type        = "SecureString"
    }
  ]
  kafka_parameters = [
    for p in local.all_kafka_parameters : p
    if p.value != "" && p.value != null
  ]
}

module "parameter_store_write" {
  source  = "cloudposse/ssm-parameter-store/aws"
  version = "0.13.0"

  count = var.ssm_parameters_enabled && local.enabled ? 1 : 0

  parameter_write = var.ssm_parameters_enabled && local.enabled ? local.kafka_parameters : []

  context = module.this.context
}

