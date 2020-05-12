package aws

import (
	"testing"
)

func TestValidateArn(t *testing.T) {
	v := ""
	_, errors := validateArn(v, "arn")
	if len(errors) != 0 {
		t.Fatalf("%q should not be validated as an ARN: %q", v, errors)
	}

	validNames := []string{
		"arn:aws:elasticbeanstalk:us-east-1:123456789012:environment/My App/MyEnvironment", // Beanstalk
		"arn:aws:iam::123456789012:user/David",                                             // IAM User
		"arn:aws:iam::aws:policy/CloudWatchReadOnlyAccess",                                 // Managed IAM policy
		"arn:aws:rds:eu-west-1:123456789012:db:mysql-db",                                   // RDS
		"arn:aws:s3:::my_corporate_bucket/exampleobject.png",                               // S3 object
		"arn:aws:events:us-east-1:319201112229:rule/rule_name",                             // CloudWatch Rule
		"arn:aws:lambda:eu-west-1:319201112229:function:myCustomFunction",                  // Lambda function
		"arn:aws:lambda:eu-west-1:319201112229:function:myCustomFunction:Qualifier",        // Lambda func qualifier
		"arn:aws-cn:ec2:cn-north-1:123456789012:instance/i-12345678",                       // China EC2 ARN
		"arn:aws-cn:s3:::bucket/object",                                                    // China S3 ARN
		"arn:aws-iso:ec2:us-iso-east-1:123456789012:instance/i-12345678",                   // C2S EC2 ARN
		"arn:aws-iso:s3:::bucket/object",                                                   // C2S S3 ARN
		"arn:aws-iso-b:ec2:us-isob-east-1:123456789012:instance/i-12345678",                // SC2S EC2 ARN
		"arn:aws-iso-b:s3:::bucket/object",                                                 // SC2S S3 ARN
		"arn:aws-us-gov:ec2:us-gov-west-1:123456789012:instance/i-12345678",                // GovCloud EC2 ARN
		"arn:aws-us-gov:s3:::bucket/object",                                                // GovCloud S3 ARN
	}
	for _, v := range validNames {
		_, errors := validateArn(v, "arn")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid ARN: %q", v, errors)
		}
	}

	invalidNames := []string{
		"arn",
		"123456789012",
		"arn:aws",
		"arn:aws:logs",
		"arn:aws:logs:region:*:*",
	}
	for _, v := range invalidNames {
		_, errors := validateArn(v, "arn")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid ARN", v)
		}
	}
}

func TestValidateIAMPolicyJsonString(t *testing.T) {
	type testCases struct {
		Value    string
		ErrCount int
	}

	invalidCases := []testCases{
		{
			Value:    `{0:"1"}`,
			ErrCount: 1,
		},
		{
			Value:    `{'abc':1}`,
			ErrCount: 1,
		},
		{
			Value:    `{"def":}`,
			ErrCount: 1,
		},
		{
			Value:    `{"xyz":[}}`,
			ErrCount: 1,
		},
		{
			Value:    ``,
			ErrCount: 1,
		},
		{
			Value:    `    {"xyz": "foo"}`,
			ErrCount: 1,
		},
	}

	for _, tc := range invalidCases {
		_, errors := validateIAMPolicyJson(tc.Value, "json")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error.", tc.Value)
		}
	}

	validCases := []testCases{
		{
			Value:    `{}`,
			ErrCount: 0,
		},
		{
			Value:    `{"abc":["1","2"]}`,
			ErrCount: 0,
		},
	}

	for _, tc := range validCases {
		_, errors := validateIAMPolicyJson(tc.Value, "json")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q not to trigger a validation error.", tc.Value)
		}
	}
}
