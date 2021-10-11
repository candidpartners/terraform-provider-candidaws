data "terraform_remote_state" "prod_network" {
  backend = "s3"
  config = {
    region = "us-gov-west-1"
    bucket = "cfs-base01-terraform-state-bucket"
    key    = "workloads/network-definition/network-def-terraform-1-rhel-no-alpha.tfstate"
  }
}

data "terraform_remote_state" "nonprod_network" {
  backend = "s3"
  config = {
    region = "us-gov-west-1"
    bucket = "cfs-base01-terraform-state-bucket"
    key    = "workloads/non-prod-network-definition/terraform.tfstate"
  }
}

##################################################
### Mediation Network Info
##################################################
data "terraform_remote_state" "mediation_prod_network" {
  backend = "s3"
  config = {
    region = "us-gov-west-1"
    bucket = "cfs-base01-terraform-state-bucket"
    key = "workloads/perimeter-frfs-prod-network-definition/terraform.tfstate"
  }
}

data "terraform_remote_state" "mediation_nonprod_network" {
  backend = "s3"
  config = {
    region = "us-gov-west-1"
    bucket = "cfs-base01-terraform-state-bucket"
    key = "workloads/perimeter-frfs-np-network-definition/terraform.tfstate"
  }
}


##################################################
### FRFS Network Info
##################################################
data "terraform_remote_state" "frfs_prod_network" {
  backend = "s3"
  config = {
    region = "us-gov-west-1"
    bucket = "cfs-base01-terraform-state-bucket"
    key = "workloads/base-prod-frfs-network-definition/terraform.tfstate"
  }
}

data "terraform_remote_state" "frfs_nonprod_network" {
  backend = "s3"
  config = {
    region = "us-gov-west-1"
    bucket = "cfs-base01-terraform-state-bucket"
    key = "workloads/base-nonprod-frfs-network-definition/terraform.tfstate"
  }
}