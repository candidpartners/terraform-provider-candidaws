lz_account_id      = "271482790176"
network_account_id = "724330398347"
is_nonprod_account = true
full_lz_code       = "lz-services"
account_name       = "dev"
vpce_subnetname    =  "private"
tgw_subnets_ids    =  "private"
routeTables          = "base"
interface_service_names = [
    "ssm",
    "ec2messages",
    "ec2",
    "ssmmessages",
    "sts",
    "logs",
    "secretsmanager",
    "codecommit",
    "codecommit-fips",
    "git-codecommit",
    "git-codecommit-fips",
    "codebuild-fips",
    "codebuild",
    "monitoring",
    "kms",
    "ecr.api",
    "ecr.dkr",
    "ecs",
    "ecs-agent",
    "ecs-telemetry"
  ]
gateway_service_names = [
  "dynamodb",
  "s3"
]
tags = {
  "Base"                       = "01"
  "LZ"                         = "services"
  "Line of Business"           = "System IT Cloud"
  "2nd Level support"          = "P1-OPS Cloud"
  "Information Classification" = "INTERNAL FR"
  "Application System CI Name" = "CFS LZ INFRASTRUCTURE - DEV"
  "CI Environment"             = "Dev"
}