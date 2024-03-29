package aws

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	homedir "github.com/mitchellh/go-homedir"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// TODO: Move the validation to this, requires conditional schemas
	// TODO: Move the configuration to this, requires validation

	// The actual provider
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["secret_key"],
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["profile"],
			},

			"assume_role": assumeRoleSchema(),

			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["shared_credentials_file"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["token"],
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				Description:  descriptions["region"],
				InputDefault: "us-east-1",
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     25,
				Description: descriptions["max_retries"],
			},

			"allowed_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"forbidden_account_ids"},
				Set:           schema.HashString,
			},

			"forbidden_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"allowed_account_ids"},
				Set:           schema.HashString,
			},

			"endpoints": endpointsSchema(),

			"ignore_tag_prefixes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Resource tag key prefixes to ignore across all resources.",
			},

			"ignore_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Resource tag keys to ignore across all resources.",
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},

			"skip_credentials_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_credentials_validation"],
			},

			"skip_get_ec2_platforms": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_get_ec2_platforms"],
			},

			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_region_validation"],
			},

			"skip_requesting_account_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_requesting_account_id"],
			},

			"skip_metadata_api_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_metadata_api_check"],
			},

			"s3_force_path_style": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["s3_force_path_style"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"aws_caller_identity":  dataSourceAwsCallerIdentity(),
			"aws_internet_gateway": dataSourceAwsInternetGateway(),
			"aws_vpc":              dataSourceAwsVpc(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"aws_transfer_server":                     resourceAwsTransferServer(),
			"aws_lex_slot_type":                       resourceAwsLexSlotType(),
			"aws_lex_intent":                          resourceAwsLexIntent(),
			"aws_lex_bot":                             resourceAwsLexBot(),
			"aws_organizations_gov_cloud_account":     resourceAwsOrganizationsGovCloudAccount(),
			"aws_organizations_invitation":            resourceAwsOrganizationsInvitation(),
			"aws_organizations_invitation_acceptance": resourceAwsOrganizationsInvitationAcceptance(),
			"aws_iam_role":                            resourceAwsIamRole(),
			"aws_iam_role_policy":                     resourceAwsIamRolePolicy(),
			"aws_iam_role_policy_attachment":          resourceAwsIamRolePolicyAttachment(),
			"aws_quicksight_data_source":              resourceAwsQuickSightDataSource(),
			"aws_quicksight_group_membership":         resourceAwsQuickSightGroupMembership(),
			"aws_quicksight_iam_policy_assignment":    resourceAwsQuickSightIAMPolicyAssignment(),
			"aws_quicksight_namespace":                resourceAwsQuickSightNamespace(),
			"aws_internet_gateway_detach":             resourceAwsInternetGatewayDetach(),
			"aws_internet_gateway_delete":             resourceAwsInternetGatewayDelete(),
			"aws_default_network_acl":                 resourceAwsDefaultNetworkAcl(),
			"aws_network_acl":                         resourceAwsNetworkAcl(),
			"aws_default_route_table":                 resourceAwsDefaultRouteTable(),
			"aws_route_table":                         resourceAwsRouteTable(),
			"aws_default_security_group":              resourceAwsDefaultSecurityGroup(),
			"aws_security_group":                      resourceAwsSecurityGroup(),
			"aws_security_group_rule":                 resourceAwsSecurityGroupRule(),
			"aws_subnet":                              resourceAwsSubnet(),
			"aws_default_subnet":                      resourceAwsDefaultSubnet(),
			"aws_network_interface":                   resourceAwsNetworkInterface(),
			"aws_default_vpc":                         resourceAwsDefaultVpc(),
			"aws_vpc":                                 resourceAwsVpc(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string
var endpointServiceNames []string

func init() {
	descriptions = map[string]string{
		"region": "The region where AWS operations will take place. Examples\n" +
			"are us-east-1, us-west-2, etc.",

		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",

		"profile": "The profile for API operations. If not set, the default profile\n" +
			"created with `aws configure` will be used.",

		"shared_credentials_file": "The path to the shared credentials file. If not set\n" +
			"this defaults to ~/.aws/credentials.",

		"token": "session token. A session token is only required if you are\n" +
			"using temporary security credentials.",

		"max_retries": "The maximum number of times an AWS API request is\n" +
			"being executed. If the API request still fails, an error is\n" +
			"thrown.",

		"endpoint": "Use this to override the default service endpoint URL",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",

		"skip_credentials_validation": "Skip the credentials validation via STS API. " +
			"Used for AWS API implementations that do not have STS available/implemented.",

		"skip_get_ec2_platforms": "Skip getting the supported EC2 platforms. " +
			"Used by users that don't have ec2:DescribeAccountAttributes permissions.",

		"skip_region_validation": "Skip static validation of region name. " +
			"Used by users of alternative AWS-like APIs or users w/ access to regions that are not public (yet).",

		"skip_requesting_account_id": "Skip requesting the account ID. " +
			"Used for AWS API implementations that do not have IAM/STS API and/or metadata API.",

		"skip_medatadata_api_check": "Skip the AWS Metadata API check. " +
			"Used for AWS API implementations that do not have a metadata api endpoint.",

		"s3_force_path_style": "Set this to true to force the request to use path-style addressing,\n" +
			"i.e., http://s3.amazonaws.com/BUCKET/KEY. By default, the S3 client will\n" +
			"use virtual hosted bucket addressing when possible\n" +
			"(http://BUCKET.s3.amazonaws.com/KEY). Specific to the Amazon S3 service.",

		"assume_role_role_arn": "The ARN of an IAM role to assume prior to making API calls.",

		"assume_role_session_name": "The session name to use when assuming the role. If omitted," +
			" no session name is passed to the AssumeRole call.",

		"assume_role_external_id": "The external ID to use when assuming the role. If omitted," +
			" no external ID is passed to the AssumeRole call.",

		"assume_role_policy": "The permissions applied when assuming a role. You cannot use," +
			" this policy to grant further permissions that are in excess to those of the, " +
			" role that is being assumed.",
	}

	endpointServiceNames = []string{
		"accessanalyzer",
		"acm",
		"acmpca",
		"amplify",
		"apigateway",
		"applicationautoscaling",
		"applicationinsights",
		"appmesh",
		"appstream",
		"appsync",
		"athena",
		"autoscaling",
		"autoscalingplans",
		"backup",
		"batch",
		"budgets",
		"cloud9",
		"cloudformation",
		"cloudfront",
		"cloudhsm",
		"cloudsearch",
		"cloudtrail",
		"cloudwatch",
		"cloudwatchevents",
		"cloudwatchlogs",
		"codebuild",
		"codecommit",
		"codedeploy",
		"codepipeline",
		"cognitoidentity",
		"cognitoidp",
		"configservice",
		"cur",
		"dataexchange",
		"datapipeline",
		"datasync",
		"dax",
		"devicefarm",
		"directconnect",
		"dlm",
		"dms",
		"docdb",
		"ds",
		"dynamodb",
		"ec2",
		"ecr",
		"ecs",
		"efs",
		"eks",
		"elasticache",
		"elasticbeanstalk",
		"elastictranscoder",
		"elb",
		"emr",
		"es",
		"firehose",
		"fms",
		"forecast",
		"fsx",
		"gamelift",
		"glacier",
		"globalaccelerator",
		"glue",
		"greengrass",
		"guardduty",
		"iam",
		"imagebuilder",
		"inspector",
		"iot",
		"iotanalytics",
		"iotevents",
		"kafka",
		"kinesis_analytics",
		"kinesis",
		"kinesisanalytics",
		"kinesisvideo",
		"kms",
		"lakeformation",
		"lambda",
		"lexmodels",
		"licensemanager",
		"lightsail",
		"macie",
		"managedblockchain",
		"marketplacecatalog",
		"mediaconnect",
		"mediaconvert",
		"medialive",
		"mediapackage",
		"mediastore",
		"mediastoredata",
		"mq",
		"neptune",
		"opsworks",
		"organizations",
		"personalize",
		"pinpoint",
		"pricing",
		"qldb",
		"quicksight",
		"r53",
		"ram",
		"rds",
		"redshift",
		"resourcegroups",
		"route53",
		"route53resolver",
		"s3",
		"s3control",
		"sagemaker",
		"sdb",
		"secretsmanager",
		"securityhub",
		"serverlessrepo",
		"servicecatalog",
		"servicediscovery",
		"servicequotas",
		"ses",
		"shield",
		"sns",
		"sqs",
		"ssm",
		"stepfunctions",
		"storagegateway",
		"sts",
		"swf",
		"transfer",
		"waf",
		"wafregional",
		"wafv2",
		"worklink",
		"workmail",
		"workspaces",
		"xray",
	}
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AccessKey:               d.Get("access_key").(string),
		SecretKey:               d.Get("secret_key").(string),
		Profile:                 d.Get("profile").(string),
		Token:                   d.Get("token").(string),
		Region:                  d.Get("region").(string),
		Endpoints:               make(map[string]string),
		MaxRetries:              d.Get("max_retries").(int),
		Insecure:                d.Get("insecure").(bool),
		SkipCredsValidation:     d.Get("skip_credentials_validation").(bool),
		SkipGetEC2Platforms:     d.Get("skip_get_ec2_platforms").(bool),
		SkipRegionValidation:    d.Get("skip_region_validation").(bool),
		SkipRequestingAccountId: d.Get("skip_requesting_account_id").(bool),
		SkipMetadataApiCheck:    d.Get("skip_metadata_api_check").(bool),
		S3ForcePathStyle:        d.Get("s3_force_path_style").(bool),
		terraformVersion:        terraformVersion,
	}

	// Set CredsFilename, expanding home directory
	credsPath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
	if err != nil {
		return nil, err
	}
	config.CredsFilename = credsPath

	assumeRoleList := d.Get("assume_role").([]interface{})

	if len(assumeRoleList) > 0 {

		var assumeRoleBlocks []AssumeRoleBlock

		for i := 0; i < len(assumeRoleList); i++ {
			assumeRole := assumeRoleList[i].(map[string]interface{})

			var newBlock AssumeRoleBlock
			newBlock.AssumeRoleARN = assumeRole["role_arn"].(string)
			newBlock.AssumeRoleSessionName = assumeRole["session_name"].(string)
			newBlock.AssumeRoleExternalID = assumeRole["external_id"].(string)

			if v := assumeRole["policy"].(string); v != "" {
				newBlock.AssumeRolePolicy = v
			}

			assumeRoleBlocks = append(assumeRoleBlocks, newBlock)
		}

		config.AssumeRoleBlocks = assumeRoleBlocks

		log.Printf("[INFO] assume_role configuration set: %q", config.AssumeRoleBlocks)
	} else {
		log.Printf("[INFO] No assume_role block read from configuration")
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		for _, endpointServiceName := range endpointServiceNames {
			config.Endpoints[endpointServiceName] = endpoints[endpointServiceName].(string)
		}
	}

	if v, ok := d.GetOk("ignore_tag_prefixes"); ok {
		for _, ignoreTagPrefixRaw := range v.(*schema.Set).List() {
			config.IgnoreTagPrefixes = append(config.IgnoreTagPrefixes, ignoreTagPrefixRaw.(string))
		}
	}

	if v, ok := d.GetOk("ignore_tags"); ok {
		for _, ignoreTagRaw := range v.(*schema.Set).List() {
			config.IgnoreTags = append(config.IgnoreTags, ignoreTagRaw.(string))
		}
	}

	if v, ok := d.GetOk("allowed_account_ids"); ok {
		for _, accountIDRaw := range v.(*schema.Set).List() {
			config.AllowedAccountIds = append(config.AllowedAccountIds, accountIDRaw.(string))
		}
	}

	if v, ok := d.GetOk("forbidden_account_ids"); ok {
		for _, accountIDRaw := range v.(*schema.Set).List() {
			config.ForbiddenAccountIds = append(config.ForbiddenAccountIds, accountIDRaw.(string))
		}
	}

	return config.Client()
}

// This is a global MutexKV for use within this plugin.
var awsMutexKV = mutexkv.NewMutexKV()

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_role_arn"],
				},

				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
				},

				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_external_id"],
				},

				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
			},
		},
	}
}

func endpointsSchema() *schema.Schema {
	endpointsAttributes := make(map[string]*schema.Schema)

	for _, endpointServiceName := range endpointServiceNames {
		endpointsAttributes[endpointServiceName] = &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: descriptions["endpoint"],
		}
	}

	// Since the endpoints attribute is a TypeSet we cannot use ConflictsWith
	endpointsAttributes["kinesis_analytics"].Deprecated = "use `endpoints` configuration block `kinesisanalytics` argument instead"
	endpointsAttributes["r53"].Deprecated = "use `endpoints` configuration block `route53` argument instead"

	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: endpointsAttributes,
		},
	}
}
