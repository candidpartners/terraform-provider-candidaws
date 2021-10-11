
data "aws_region" "current" {}
data "aws_partition" "current" {}
data "aws_caller_identity" "current" {}
data "terraform_remote_state" "prod_network" {
  backend = "s3"
  config = {
    region = "us-gov-west-1"
    bucket = "cfs-base01-terraform-state-bucket"
    key    = "workloads/network-definition/network-def-terraform-1-rhel-no-alpha.tfstate"
  }
}
