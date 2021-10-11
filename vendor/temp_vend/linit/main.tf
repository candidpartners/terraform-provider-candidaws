module "lz_init" {
  source = "s3::https://s3-us-gov-west-1.amazonaws.com/cfs-base01-foundation-library-release-bucket/landing-zone-init-item-gov-v2/landing-zone-init-item-gov-v2.0.0.1104.zip"
  providers = {
    aws.base_prod_services_account = aws.base_services_prod_primary
    aws.primary                    = aws.lz_account_primary
    aws.secondary                  = aws.lz_account_secondary
    candidaws.primary              = candidaws.lz_account_primary
    candidaws.secondary            = candidaws.lz_account_secondary
  }
  lz_code = var.lz_code
  base_prod_services_dns_components = {
    "primary" = {
      "shared_resolver_rules_awscfs_frb_pvt_forward_inbound_rule_id" = local.awscfs_frb_pvt_forward_inbound_rule_primary,
      "shared_resolver_rules_frb_org_forward_rule_id"                = local.frb_org_forward_rule_primary,
      "shared_resolver_rules_frb_pvt_forward_rule_id"                = local.frb_pvt_forward_rule_primary
    }
    "secondary" = {
      "shared_resolver_rules_awscfs_frb_pvt_forward_inbound_rule_id" = local.awscfs_frb_pvt_forward_inbound_rule_secondary,
      "shared_resolver_rules_frb_org_forward_rule_id"                = local.frb_org_forward_rule_secondary,
      "shared_resolver_rules_frb_pvt_forward_rule_id"                = local.frb_pvt_forward_rule_secondary
    }
  }
  tags                          = var.tags
  logs_bucket_arn               = var.logs_bucket_arn
  logs_bucket_kms_arn           = var.logs_bucket_kms_arn
  proxy_components              = var.proxy_components
  is_prod_services              = var.is_prod_services
  vpc_cidr_block                = var.vpc_cidr_block
  account_name                  = var.account_name
  account_short_name            = var.account_short_name
  okta_role_configs             = var.okta_role_configs
  root_ca_secret_name           = var.root_ca_secret_name
  s3_access_logs_bucket_name    = var.s3_access_logs_bucket_name
  artifacts_bucket_name         = var.artifacts_bucket_name
  compliance_policy_bucket_name = var.compliance_policy_bucket_name
  tfstate_bucket_name           = var.tfstate_bucket_name
  account_alias                 = local.account_alias
  enable_account_wide_key       = var.enable_account_wide_key
  private_subnet_name           = var.private_subnet_name
  cidr_information_primary      = var.cidr_information_primary
  cidr_information_secondary    = var.cidr_information_secondary
  is_network_enabled            = var.is_network_enabled
  base_okta_app_secret_arn      = "arn:${local.partition}:secretsmanager:${local.region}:${local.account_id}:secret:/cfs/base01/compliance/cfs-services-lz-def"
}
module "private_hosted_zone_global" {
  source = "s3::https://s3-us-gov-west-1.amazonaws.com/cfs-base01-foundation-library-release-bucket/r53-private-hosted-zones/r53-private-hosted-zones.0.0.856.zip"
  providers = {
    aws.primary                      = aws.lz_account_primary
    aws.secondary                    = aws.lz_account_secondary
    aws.base_services_prod_primary   = aws.base_services_prod_primary
    aws.base_services_prod_secondary = aws.base_services_prod_secondary
  }
  dns = {
    tags                 = var.tags
    enable               = var.dns.enable
    hosted_zone_name     = var.dns.hosted_zone_name
    vpc_id_primary       = module.lz_init.basic_network_infra.vpc_id_primary
    vpc_id_secondary     = module.lz_init.basic_network_infra.vpc_id_secondary
    dns_vpc_id_primary   = local.dns_vpc_id_primary
    dns_vpc_id_secondary = local.dns_vpc_id_secondary
  }
}