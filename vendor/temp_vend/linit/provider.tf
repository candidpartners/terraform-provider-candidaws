provider "prismacloud" {
  version = "1.0.4"
}
provider "okta" {
  version = "3.10.1"
}
provider aws {
  alias  = "base_services_prod_primary"
  region = "us-gov-west-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.base_prod_services_account_id}:role/cfs-deploy-role"
  }
}
provider aws {
  alias  = "base_services_prod_secondary"
  region = "us-gov-east-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.base_prod_services_account_id}:role/cfs-deploy-role"
  }
}
provider aws {
  alias  = "lz_account_primary"
  region = "us-gov-west-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.lz_account_id}:role/cfs-landing-zone-deploy-role"
  }
}
provider aws {
  alias  = "lz_account_secondary"
  region = "us-gov-east-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.lz_account_id}:role/cfs-landing-zone-deploy-role"
  }
}
provider candidaws {
  alias  = "lz_account_primary"
  region = "us-gov-west-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.gov_payer_account_id}:role/cfs-landing-zone-init-role"
  }
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.lz_account_id}:role/awg-nit-iam-role-AwsOrganizations"
  }
}
provider candidaws {
  alias  = "lz_account_secondary"
  region = "us-gov-east-1"
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.gov_payer_account_id}:role/cfs-landing-zone-init-role"
  }
  assume_role {
    role_arn = "arn:aws-us-gov:iam::${var.lz_account_id}:role/awg-nit-iam-role-AwsOrganizations"
  }
}
