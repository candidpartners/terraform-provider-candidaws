locals {
  prod_network_components    = data.terraform_remote_state.prod_network.outputs.networking_components
  nonprod_network_components = data.terraform_remote_state.nonprod_network.outputs.networking_components
  
  transit_gateway_id_primary = var.is_nonprod_account ? local.nonprod_network_components.transit_gateway_id_primary : local.prod_network_components.transit_gateway_id_primary
  tgw_cloud_route_table_id_primary = var.is_nonprod_account ? local.nonprod_network_components.tgw_cloud_route_table_id_primary : local.prod_network_components.tgw_cloud_route_table_id_primary
  tgw_post_inspection_route_table_id_primary = var.is_nonprod_account ? local.nonprod_network_components.tgw_post_inspection_route_table_id_primary : local.prod_network_components.tgw_post_inspection_route_table_id_primary
  
  transit_gateway_id_secondary = var.is_nonprod_account ? local.nonprod_network_components.transit_gateway_id_secondary : local.prod_network_components.transit_gateway_id_secondary
  tgw_cloud_route_table_id_secondary = var.is_nonprod_account ? local.nonprod_network_components.tgw_cloud_route_table_id_secondary : local.prod_network_components.tgw_cloud_route_table_id_secondary
  tgw_post_inspection_route_table_id_secondary = var.is_nonprod_account ? local.nonprod_network_components.tgw_post_inspection_route_table_id_secondary : local.prod_network_components.tgw_post_inspection_route_table_id_secondary
  
  
##################################################
### Mediation Network
##################################################
  mediation_prod_network_components    = data.terraform_remote_state.mediation_prod_network.outputs.network_definition
  mediation_nonprod_network_components = data.terraform_remote_state.mediation_nonprod_network.outputs.network_definition
 
  mediation_transit_gateway_id_primary = "dummy value" #  var.is_nonprod_account ? local.mediation_nonprod_network_components.tgw_primary : local.mediation_prod_network_components.tgw_primary
  # These aren't available in the state yet, however they are present and just commented out for now so adding them here.
  mediation_tgw_cloud_route_table_id_primary = "dummy value" # var.is_nonprod_account ? local.mediation_nonprod_network_components.tgw_cloud_route_table_id_primary : local.mediation_prod_network_components.tgw_cloud_route_table_id_primary
  mediation_tgw_post_inspection_route_table_id_primary = "dummy value" # var.is_nonprod_account ? local.mediation_nonprod_network_components.tgw_post_inspection_route_table_id_primary : local.mediation_prod_network_components.tgw_post_inspection_route_table_id_primary

  mediation_transit_gateway_id_secondary = "dummy value" #  var.is_nonprod_account ? local.mediation_nonprod_network_components.tgw_secondary : local.mediation_prod_network_components.tgw_secondary
  # These aren't available in the state yet, however they are present and just commented out for now so adding them here.
  mediation_tgw_cloud_route_table_id_secondary = "dummy value" # var.is_nonprod_account ? local.mediation_nonprod_network_components.tgw_cloud_route_table_id_secondary : local.mediation_prod_network_components.tgw_cloud_route_table_id_secondary
  mediation_tgw_post_inspection_route_table_id_secondary = "dummy value" # var.is_nonprod_account ? local.mediation_nonprod_network_components.tgw_post_inspection_route_table_id_secondary : local.mediation_prod_network_components.tgw_post_inspection_route_table_id_secondary
  
    
##################################################
### FRFS Network
##################################################
  frfs_prod_network_components    = data.terraform_remote_state.frfs_prod_network.outputs.networking_components
  frfs_nonprod_network_components = data.terraform_remote_state.frfs_nonprod_network.outputs.networking_components
  
  frfs_tgw_cloud_route_table_id_primary = var.is_nonprod_account ? local.frfs_nonprod_network_components.tgw_cloud_route_table_id_primary : local.frfs_prod_network_components.tgw_cloud_route_table_id_primary
  frfs_tgw_post_inspection_route_table_id_primary = var.is_nonprod_account ? local.frfs_nonprod_network_components.tgw_post_inspection_route_table_id_primary : local.frfs_prod_network_components.tgw_post_inspection_route_table_id_primary
  
 
  frfs_tgw_cloud_route_table_id_secondary = var.is_nonprod_account ? local.frfs_nonprod_network_components.tgw_cloud_route_table_id_secondary : local.frfs_prod_network_components.tgw_cloud_route_table_id_secondary
  frfs_tgw_post_inspection_route_table_id_secondary = var.is_nonprod_account ? local.frfs_nonprod_network_components.tgw_post_inspection_route_table_id_secondary : local.frfs_prod_network_components.tgw_post_inspection_route_table_id_secondary


##################################################
### Route Table Logic
##################################################
  cloud_route_table_primary = var.routeTables == "mediation" ? local.mediation_tgw_cloud_route_table_id_primary : var.routeTables == "frfs" ? local.frfs_tgw_cloud_route_table_id_primary : local.tgw_cloud_route_table_id_primary
  post_inspection_route_table_primary = var.routeTables == "mediation" ? local.mediation_tgw_post_inspection_route_table_id_primary : var.routeTables == "frfs" ? local.frfs_tgw_post_inspection_route_table_id_primary : local.tgw_post_inspection_route_table_id_primary

  cloud_route_table_secondary = var.routeTables == "mediation" ? local.mediation_tgw_cloud_route_table_id_secondary : var.routeTables == "frfs" ? local.frfs_tgw_cloud_route_table_id_secondary : local.tgw_cloud_route_table_id_secondary
  post_inspection_route_table_secondary = var.routeTables == "mediation" ? local.mediation_tgw_post_inspection_route_table_id_secondary : var.routeTables == "frfs" ? local.frfs_tgw_post_inspection_route_table_id_secondary : local.tgw_post_inspection_route_table_id_secondary
}
