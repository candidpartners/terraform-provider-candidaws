variable "resource_prefix" {
  type        = string
  description = "The prefix used for all resources created to identify"
  default     = "cfs"
}
variable "resource_slug" {
  type        = string
  description = "The suffix used for all resources created to allow side by side deployments.  The slug can also be used to differentiate paths as needed."
  default     = ""
}
variable "base_code" {
  type    = string
  default = "01"
}
variable "base_prod_services_account_id" {
  type    = string
  default = "701303807346"
}
variable "gov_payer_account_id" {
  type    = string
  default = "483019498278"
}
variable "root_ca_secret_name" {
  type = string
}

variable "lz_account_id" {
  type = string
}

variable "lz_code" {
  type = string
}
variable "tags" {
  type = map(string)
}
variable "logs_bucket_arn" {
  type = string
}
variable "logs_bucket_kms_arn" {
  type = string
}
variable "proxy_components" {
  type = object({
    http_proxy  = string
    https_proxy = string
    no_proxy    = string
  })
}
variable "is_prod_services" {
  type = bool
}
variable "vpc_cidr_block" {}
variable "account_name" {
  type = string
}
variable "okta_role_configs" {}
variable "primary_services_account_id" {
  type        = string
  description = "primary services account id"
  default     = "701303807346"
}
variable "okta_secret_name" {
  type    = string
  default = "/cfs/base01/okta_configs"
}
variable "dns" {}
variable "s3_access_logs_bucket_name" {}
variable "tfstate_bucket_name" {
  type = string
}
variable "artifacts_bucket_name" {
  type = string
}
variable "compliance_policy_bucket_name" {
  type = string
}
variable "enable_account_wide_key" {
  type = bool
}
variable "account_short_name" {
  type = string
}
variable "cidr_information_primary" {
  type = map(object({
    subnets = map(object({
      attach_to_tgw = bool,
      cidr_blocks   = map(string)
    }))
  }))
  default = {}
}
variable "cidr_information_secondary" {
  type = map(object({
    subnets = map(object({
      attach_to_tgw = bool,
      cidr_blocks   = map(string)
    }))
  }))
  default = {}
}
variable "private_subnet_name" {
  type        = string
  description = "name of the private subnets to associate to the pipeline"
}
variable "is_network_enabled" {
  type = bool
}
