terraform {
  required_providers {
    typesense = {
      source = "omarkhd.net/terraform/typesense"
    }
  }
}

provider "typesense" {
}

resource "typesense_cluster" "example" {
  memory = "0.5_gb"
  vcpu = "2_vcpus_1_hr_burst_per_day"
  region = "oregon"
  name = "example"
}

resource "typesense_cluster_api_keys" "example" {
  cluster_id = typesense_cluster.example.id
}

output "typesense-admin-cluster-api-key" {
  value       = typesense_cluster_api_keys.example.admin_key
  description = "Admin key"
}

output "typesense-search-only-cluster-api-key" {
  value       = typesense_cluster_api_keys.example.search_only_key
  description = "Search Only key"
}
