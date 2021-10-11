
provider aws {
  alias  = "account_primary"
  region = "us-gov-west-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.lz_account_id}:role/cfs-landing-zone-deploy-role"
  }
}
provider aws {
  alias  = "account_secondary"
  region = "us-gov-east-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.lz_account_id}:role/cfs-landing-zone-deploy-role"
  }
}

provider aws {
  alias  = "network_account_primary"
  region = "us-gov-west-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.network_account_id}:role/cfs-deploy-role"
  }
}
provider aws {
  alias  = "network_account_secondary"
  region = "us-gov-east-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.network_account_id}:role/cfs-deploy-role"
  }
}
