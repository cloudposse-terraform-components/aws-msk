output "cluster_name" {
  value       = try(module.kafka.cluster_name, null)
  description = "The cluster name of the MSK cluster"
}

output "cluster_arn" {
  value       = try(module.kafka.cluster_arn, null)
  description = "Amazon Resource Name (ARN) of the MSK cluster"
}

output "storage_mode" {
  value       = try(module.kafka.storage_mode, null)
  description = "Storage mode for supported storage tiers"
}

output "bootstrap_brokers" {
  value       = try(module.kafka.bootstrap_brokers, null)
  description = "Comma separated list of one or more hostname:port pairs of Kafka brokers suitable to bootstrap connectivity to the Kafka cluster"
}

output "bootstrap_brokers_tls" {
  value       = try(module.kafka.bootstrap_brokers_tls, null)
  description = "Comma separated list of one or more DNS names (or IP addresses) and TLS port pairs for access to the Kafka cluster using TLS"
}

output "bootstrap_brokers_public_tls" {
  value       = try(module.kafka.bootstrap_brokers_public_tls, null)
  description = "Comma separated list of one or more DNS names (or IP addresses) and TLS port pairs for public access to the Kafka cluster using TLS"
}

output "bootstrap_brokers_sasl_scram" {
  value       = try(module.kafka.bootstrap_brokers_sasl_scram, null)
  description = "Comma separated list of one or more DNS names (or IP addresses) and SASL SCRAM port pairs for access to the Kafka cluster using SASL/SCRAM"
}

output "bootstrap_brokers_public_sasl_scram" {
  value       = try(module.kafka.bootstrap_brokers_public_sasl_scram, null)
  description = "Comma separated list of one or more DNS names (or IP addresses) and SASL SCRAM port pairs for public access to the Kafka cluster using SASL/SCRAM"
}

output "bootstrap_brokers_sasl_iam" {
  value       = try(module.kafka.bootstrap_brokers_sasl_iam, null)
  description = "Comma separated list of one or more DNS names (or IP addresses) and SASL IAM port pairs for access to the Kafka cluster using SASL/IAM"
}

output "bootstrap_brokers_public_sasl_iam" {
  value       = try(module.kafka.bootstrap_brokers_public_sasl_iam, null)
  description = "Comma separated list of one or more DNS names (or IP addresses) and SASL IAM port pairs for public access to the Kafka cluster using SASL/IAM"
}

output "zookeeper_connect_string" {
  value       = try(module.kafka.zookeeper_connect_string, null)
  description = "Comma separated list of one or more hostname:port pairs to connect to the Apache Zookeeper cluster"
}

output "zookeeper_connect_string_tls" {
  value       = try(module.kafka.zookeeper_connect_string_tls, null)
  description = "Comma separated list of one or more hostname:port pairs to connect to the Apache Zookeeper cluster via TLS"
}

output "broker_endpoints" {
  value       = try(module.kafka.broker_endpoints, null)
  description = "List of broker endpoints"
}

output "current_version" {
  value       = try(module.kafka.current_version, null)
  description = "Current version of the MSK Cluster"
}

output "config_arn" {
  value       = try(module.kafka.config_arn, null)
  description = "Amazon Resource Name (ARN) of the MSK configuration"
}

output "latest_revision" {
  value       = try(module.kafka.latest_revision, null)
  description = "Latest revision of the MSK configuration"
}

output "hostnames" {
  value       = try(module.kafka.hostnames, null)
  description = "List of MSK Cluster broker DNS hostnames"
}

output "security_group_id" {
  value       = try(module.kafka.security_group_id, null)
  description = "The ID of the created security group"
}

output "security_group_arn" {
  value       = try(module.kafka.security_group_arn, null)
  description = "The ARN of the created security group"
}

output "security_group_name" {
  value       = try(module.kafka.security_group_name, null)
  description = "The name of the created security group"
}
