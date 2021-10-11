locals {
  region                                        = data.aws_region.current.name
  account_id                                    = data.aws_caller_identity.current.account_id
  partition                                     = data.aws_partition.current.partition
  account_alias                                 = join("-", compact([var.resource_prefix, var.resource_slug, "base${var.base_code}", "gov", var.lz_code, var.account_name]))
  project_name                                  = join("-", compact([var.resource_prefix, var.resource_slug, var.lz_code, var.account_name, "lz-default-app"]))
  prod_dns_components                           = data.terraform_remote_state.prod_network.outputs.dns_components
  dns_vpc_id_primary                            = local.prod_dns_components.dns_vpc_id_primary.id
  dns_vpc_id_secondary                          = local.prod_dns_components.dns_vpc_id_secondary.id
  awscfs_frb_pvt_forward_inbound_rule_primary   = local.prod_dns_components.awscfs_frb_pvt_forward_inbound_rule_primary.id
  awscfs_frb_pvt_forward_inbound_rule_secondary = local.prod_dns_components.awscfs_frb_pvt_forward_inbound_rule_secondary.id
  frb_org_forward_rule_primary                  = local.prod_dns_components.frb_org_forward_rule_primary.id
  frb_org_forward_rule_secondary                = local.prod_dns_components.frb_org_forward_rule_secondary.id
  frb_pvt_forward_rule_primary                  = local.prod_dns_components.frb_pvt_forward_rule_primary.id
  frb_pvt_forward_rule_secondary                = local.prod_dns_components.frb_pvt_forward_rule_secondary.id
}