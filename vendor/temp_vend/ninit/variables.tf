variable "is_nonprod_account" {
  type = bool
}
variable "lz_account_id" {
  type = string

}
variable "network_account_id" {
  type = string
}
variable "gateway_service_names" {
  type = set(string)
}
variable "interface_service_names" {
  type = set(string)
}
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
variable "full_lz_code" {
  type        = string
  description = "Code the full name for the lz for which accounts are setup.. example lz-<alpha>"
  default     = ""
}
variable "base_code" {
  type        = string
  description = "Code for the base for which accounts are setup"
  default     = "01"
}
variable "account_name" {
  type = string
}
variable "tags" {
  type    = map(string)
  default = {}
}
variable "vpce_subnetname" {
  type = string
  description = "subnet to associate to the vpce"
}
variable "tgw_subnets_ids" {
  type = string
  description = "subnet to associate to the tgw attachment"
}
variable "routeTables" {
  type = string
  description = "name of network allocation type: base, mediation, or frfs"
}