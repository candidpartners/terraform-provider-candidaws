module lz_init_cicd_services {
  source             = "s3::https://s3-us-gov-west-1.amazonaws.com/cfs-base01-foundation-library-release-bucket/cicd-multi-stage/cicd-multi-stage.0.0.1079.zip"
  access_logs_bucket = "cfs-base01-prod-svc-west-log-bucket"
  artifacts_bucket   = data.aws_s3_bucket.artifacts_bucket.id
  project_name       = join("-", compact([var.resource_prefix, var.resource_slug, "landing-zone-init-services"]))
  context            = "${local.base_name}.landing-zone-init-services"
  vpc_id             = var.vpc_id
  subnet_ids         = var.subnet_ids
  security_group_ids = var.security_group_ids
  with_okta          = true
  okta_secret_name   = var.okta_secret_name
  state_bucket_id    = "cfs-base01-terraform-state-bucket"
  prod_promo_accounts = { 
  }
  nonprod_promo_accounts = { 
    "dev" = {
      name               = "1-dev",
      env                = "dev",
      accountid          = "271482790176",
      compliance_context = "cfs-base01-gov.landing-zone-init-services",
      enable_plan        = false,
      enable_approval    =  true 
    }
    "proving-ground" = {
      name               = "2-proving-ground",
      env                = "proving-ground",
      accountid          = "271059040458",
      compliance_context = "cfs-base01-gov.landing-zone-init-services",
      enable_plan        = false,
      enable_approval    =  true 
    }
    "test" = {
      name               = "3-test",
      env                = "test",
      accountid          = "387795238509",
      compliance_context = "cfs-base01-gov.landing-zone-init-services",
      enable_plan        = false,
      enable_approval    =  true 
    }
  }
  prod_deploy_env          = []
  nonprod_deploy_env       = ["dev","proving-ground","test"]
  proxy_config_secret_name = var.proxy_config_secret_name
  root_ca_secret_name      = var.root_ca_secret_name
  deploy_role_arn          = "arn:aws-us-gov:iam::701303807346:role/cfs-deploy-role"
  logs_kms_key_arn         = data.aws_kms_key.logs_key.arn
  s3_kms_key_arn           = data.aws_kms_key.s3_key.arn
  lambda_kms_key_arn       = data.aws_kms_key.lambda_key.arn
  tags                     = local.tags
  is_enabled_okta_lz       = true
  lz_code                  = "services"
  providers = {
    aws = aws
    okta = okta
  }
}

module lz_org_init_cicd_services {
  source             = "s3::https://s3-us-gov-west-1.amazonaws.com/cfs-base01-foundation-library-release-bucket/cicd-multi-stage/cicd-multi-stage.0.0.1079.zip"
  access_logs_bucket = "cfs-base01-prod-svc-west-log-bucket"
  artifacts_bucket   = data.aws_s3_bucket.artifacts_bucket.id
  project_name       = join("-", compact([var.resource_prefix, var.resource_slug, "landing-zone-org-services"]))
  context            = "${local.base_name}.landing-zone-init-services"
  vpc_id             = var.vpc_id
  subnet_ids         = var.subnet_ids
  security_group_ids = var.security_group_ids
  with_okta          = true
  okta_secret_name   = var.okta_secret_name
  state_bucket_id    = "cfs-base01-terraform-state-bucket"
  prod_promo_accounts = { 
  }
  nonprod_promo_accounts = { 
    "dev" = {
      name               = "1-dev",
      env                = "dev",
      accountid          = "271482790176",
      compliance_context = "cfs-base01-gov.landing-zone-init-services",
      enable_plan        = false,
      enable_approval    =  true 
    }
    "proving-ground" = {
      name               = "2-proving-ground",
      env                = "proving-ground",
      accountid          = "271059040458",
      compliance_context = "cfs-base01-gov.landing-zone-init-services",
      enable_plan        = false,
      enable_approval    =  true 
    }
    "test" = {
      name               = "3-test",
      env                = "test",
      accountid          = "387795238509",
      compliance_context = "cfs-base01-gov.landing-zone-init-services",
      enable_plan        = false,
      enable_approval    =  true 
    }
  }
  prod_deploy_env          = []
  nonprod_deploy_env       = ["dev","proving-ground","test"]
  proxy_config_secret_name = var.proxy_config_secret_name
  root_ca_secret_name      = var.root_ca_secret_name
  deploy_role_arn          = "arn:aws-us-gov:iam::701303807346:role/cfs-deploy-role"
  logs_kms_key_arn         = data.aws_kms_key.logs_key.arn
  s3_kms_key_arn           = data.aws_kms_key.s3_key.arn
  lambda_kms_key_arn       = data.aws_kms_key.lambda_key.arn
  tags                     = local.tags
  lz_code                  = "services"
  providers = {
    aws  = aws
    okta = okta
  }
}