module "network_init_primary" {
  source = "s3::https://s3-us-gov-west-1.amazonaws.com/cfs-base01-foundation-library-release-bucket/networking-infra-gov-v2/networking-infra-gov-v2.0.0.1068.zip"
  tgw = {
    destroy                            = false
    transit_gateway_id                 = var.routeTables == "mediation" ? local.mediation_transit_gateway_id_primary : local.transit_gateway_id_primary
    tgw_cloud_route_table_id           = local.cloud_route_table_primary
    tgw_post_inspection_route_table_id = local.post_inspection_route_table_primary
  }
  tags                    = var.tags
  account_name            = var.account_name
  full_lz_code            = var.full_lz_code
  gateway_service_names   = var.gateway_service_names
  interface_service_names = var.interface_service_names
  vpce_subnetname         = var.vpce_subnetname
  tgw_subnets_ids         = var.tgw_subnets_ids
  providers = {
    aws.network_account = aws.network_account_primary
    aws.account         = aws.account_primary
  }
}

module "network_init_secondary" {
  source = "s3::https://s3-us-gov-west-1.amazonaws.com/cfs-base01-foundation-library-release-bucket/networking-infra-gov-v2/networking-infra-gov-v2.0.0.1068.zip"
  tgw = {
    destroy                            = false
    transit_gateway_id                 = var.routeTables == "mediation" ? local.mediation_transit_gateway_id_secondary : local.transit_gateway_id_secondary
    tgw_cloud_route_table_id           = local.cloud_route_table_secondary
    tgw_post_inspection_route_table_id = local.post_inspection_route_table_secondary
  }
  tags                    = var.tags
  account_name            = var.account_name
  full_lz_code            = var.full_lz_code
  gateway_service_names   = var.gateway_service_names
  interface_service_names = var.interface_service_names
  vpce_subnetname         = var.vpce_subnetname
  tgw_subnets_ids         = var.tgw_subnets_ids
  providers = {
    aws.network_account = aws.network_account_secondary
    aws.account         = aws.account_secondary
  }
}
