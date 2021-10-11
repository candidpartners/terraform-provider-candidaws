lz_account_id = "271482790176"
lz_code       = "services"
tags = {
  "Base"                       = "01"
  "LZ"                         = "services"
  "Line of Business"           = "System IT Cloud"
  "2nd Level support"          = "P1-OPS Cloud"
  "Information Classification" = "INTERNAL FR"
  "Application System CI Name" = "CFS LZ INFRASTRUCTURE - DEV"
  "CI Environment"             = "Dev"
}
root_ca_secret_name = "fr-enterprise-root-ca1-v2"
logs_bucket_arn     = "arn:aws-us-gov:s3:::cfs-01-logging-ingest-bucket"
# need to read from output of base init
logs_bucket_kms_arn = "arn:aws-us-gov:kms:us-gov-west-1:723318254883:key/431f2b41-d1c6-4ee8-b8ce-a135bf5ae66e" 
proxy_components = {
  http_proxy  = "http://p3proxy.frb.org:8080"
  https_proxy = "http://p3proxy.frb.org:8080"
  no_proxy    = "compliance.base.awscfs.frb.pvt,.base.awscfs.frb.pvt,.awscfs.frb.pvt"
}
# service renamed to svc and nonprod to np
s3_access_logs_bucket_name = {
  primary = "cfs-base01-services-dev-west-log-bucket" 
  secondary = "cfs-base01-services-dev-east-log-bucket"
}
artifacts_bucket_name         = "cfs-base01-services-dev-artifacts-bucket"
compliance_policy_bucket_name = "cfs-base01-services-dev-mce-policies-bucket"
is_prod_services              = false
private_subnet_name = "private"
is_network_enabled = true
vpc_cidr_block = { 
  primary = "100.84.52.0/24"
  secondary = "100.100.52.0/24"
}
cidr_information_primary = { 
  primary = { 
    subnets = { 
      private = {
        attach_to_tgw = true
        cidr_blocks = {
          "us-gov-west-1a" = "100.84.52.0/26",
          "us-gov-west-1b" = "100.84.52.64/26",
          "us-gov-west-1c" = "100.84.52.128/26",
          
        }
      }
    }
  }
}
cidr_information_secondary = { 
  secondary = { 
    subnets =  { 
      private = {
        attach_to_tgw = true
        cidr_blocks = {
          "us-gov-east-1a" = "100.100.52.0/26",
          "us-gov-east-1b" = "100.100.52.64/26",
          "us-gov-east-1c" = "100.100.52.128/26",
          
        }
      }
    }
  } 
}
account_name = "dev"
okta_role_configs = {
  "adt_engineer_enabled"                             = true,
  "adt_secrets_admin_enabled"                        = true,
  "compliance_engineer_enabled"                      = true,
  "adt_engineer_pg_enabled"                          = false,
  "foundation_engineer_iam_pass_codebuild_role_name" = "cfs-landing-zone-codebuild-role",
  "adt_engineer_iam_pass_codebuild_role_name"        = "cfs-landing-zone-codebuild-role",
  "adfs_adt_engineer_iam_pass_codebuild_role_name"   = "cfs-landing-zone-codebuild-role",
  "foundation_engineer_accepted_actions" = [
    "codepipeline:StartPipelineExecution",
    "codepipeline:RetryStageExecution",
    "codepipeline:StopPipelineExecution",
    "codebuild:RetryBuild",
    "codebuild:RetryBuildBatch",
    "codebuild:StopBuild",
    "codebuild:StopBuildBatch"
  ],
  "adt_engineer_accepted_actions" = [
    "codepipeline:StartPipelineExecution",
    "codepipeline:RetryStageExecution",
    "codepipeline:StopPipelineExecution",
    "codebuild:RetryBuild",
    "codebuild:RetryBuildBatch",
    "codebuild:StopBuild",
    "codebuild:StopBuildBatch"
  ],
  "adfs_adt_engineer_accepted_actions" = [
    "codepipeline:StartPipelineExecution",
    "codepipeline:RetryStageExecution",
    "codepipeline:StopPipelineExecution",
    "codebuild:RetryBuild",
    "codebuild:RetryBuildBatch",
    "codebuild:StopBuild",
    "codebuild:StopBuildBatch"
  ],
  "compliance_engineer_accepted_actions" = [
    "codecommit:AssociateApprovalRuleTemplateWithRepository",
    "codecommit:BatchAssociateApprovalRuleTemplateWithRepositories",
    "codecommit:BatchDisassociateApprovalRuleTemplateFromRepositories",
    "codecommit:BatchGetCommits",
    "codecommit:BatchGetPullRequests",
    "codecommit:BatchDescribeMergeConflicts",
    "codecommit:CreateBranch",
    "codecommit:CreateCommit",
    "codecommit:CreatePullRequest",
    "codecommit:DescribeMergeConflicts",
    "codecommit:DescribePullRequestEvents",
    "codecommit:DisassociateApprovalRuleTemplateFromRepository",
    "codecommit:EvaluatePullRequestApprovalRules",
    "codecommit:GetRepository",
    "codecommit:GetReferences",
    "codecommit:GetPullRequestApprovalStates",
    "codecommit:GetPullRequest",
    "codecommit:GetMergeOptions",
    "codecommit:GetMergeConflicts",
    "codecommit:GetMergeCommit",
    "codecommit:GetFolder",
    "codecommit:GetFile",
    "codecommit:GetDifferences",
    "codecommit:GetCommitHistory",
    "codecommit:GetCommit",
    "codecommit:GetCommentsForPullRequest",
    "codecommit:GetCommentsForComparedCommit",
    "codecommit:GetCommentReactions",
    "codecommit:GetComment",
    "codecommit:GetBranch",
    "codecommit:GetBlob",
    "codecommit:GetApprovalRuleTemplate",
    "codecommit:ListRepositories",
    "codecommit:ListPullRequests",
    "codecommit:ListBranches",
    "codecommit:MergeBranchesByFastForward",
    "codecommit:MergeBranchesBySquash",
    "codecommit:MergeBranchesByThreeWay",
    "codecommit:MergePullRequestByFastForward",
    "codecommit:MergePullRequestBySquash",
    "codecommit:MergePullRequestByThreeWay",
    "codecommit:PutCommentReaction",
    "codecommit:PutFile",
    "codecommit:PostCommentForComparedCommit",
    "codecommit:PostCommentForPullRequest",
    "codecommit:PostCommentReply",
    "codecommit:UntagResource",
    "codecommit:UpdateComment",
    "codecommit:UpdatePullRequestDescription",
    "codecommit:UpdatePullRequestStatus",
    "codecommit:UpdatePullRequestTitle",
    "codecommit:GitPull",
    "codecommit:GitPush"
  ],
  "compliance_engineer_deny_actions" = [
    "codecommit:GitPush",
    "codecommit:PutFile",
    "codecommit:DeleteBranch"
  ],
}
dns = {
  hosted_zone_name = "dev.services.awscfs.frb.pvt"
  enable           = true
}
tfstate_bucket_name = "cfs-base01-lz-services-dev-tf-state-bucket"
enable_account_wide_key = false
account_short_name      = "dev"