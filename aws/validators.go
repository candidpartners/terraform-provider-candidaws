package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"regexp"
)

const (
	awsAccountIDRegexpPattern = `^(aws|\d{12})$`
	awsPartitionRegexpPattern = `^aws(-[a-z]+)*$`
	awsRegionRegexpPattern    = `^[a-z]{2}(-[a-z]+)+-\d$`
)

var awsAccountIDRegexp = regexp.MustCompile(awsAccountIDRegexpPattern)
var awsPartitionRegexp = regexp.MustCompile(awsPartitionRegexpPattern)
var awsRegionRegexp = regexp.MustCompile(awsRegionRegexpPattern)

func validateTransferServerID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	// https://docs.aws.amazon.com/transfer/latest/userguide/API_CreateUser.html
	pattern := `^s-([0-9a-f]{17})$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q isn't a valid transfer server id (only lowercase alphanumeric characters are allowed): %q",
			k, value))
	}

	return
}

func validateTransferUserName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	// https://docs.aws.amazon.com/transfer/latest/userguide/API_CreateUser.html
	if !regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_-]{2,31}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("Invalid %q: must be between 3 and 32 alphanumeric or special characters hyphen and underscore. However, %q cannot begin with a hyphen", k, k))
	}
	return
}

func validateArn(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return
	}

	parsedARN, err := arn.Parse(value)

	if err != nil {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: %s", k, value, err))
		return
	}

	if parsedARN.Partition == "" {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: missing partition value", k, value))
	} else if !awsPartitionRegexp.MatchString(parsedARN.Partition) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid partition value (expecting to match regular expression: %s)", k, value, awsPartitionRegexpPattern))
	}

	if parsedARN.Region != "" && !awsRegionRegexp.MatchString(parsedARN.Region) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid region value (expecting to match regular expression: %s)", k, value, awsRegionRegexpPattern))
	}

	if parsedARN.AccountID != "" && !awsAccountIDRegexp.MatchString(parsedARN.AccountID) {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: invalid account ID value (expecting to match regular expression: %s)", k, value, awsAccountIDRegexpPattern))
	}

	if parsedARN.Resource == "" {
		errors = append(errors, fmt.Errorf("%q (%s) is an invalid ARN: missing resource value", k, value))
	}

	return
}

func validateIAMPolicyJson(v interface{}, k string) (ws []string, errors []error) {
	// IAM Policy documents need to be valid JSON, and pass legacy parsing
	value := v.(string)
	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON policy", k))
		return
	}
	if value[:1] != "{" {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON policy", k))
		return
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func validateIamRolePolicyName(v interface{}, k string) (ws []string, errors []error) {
	// https://github.com/boto/botocore/blob/2485f5c/botocore/data/iam/2010-05-08/service-2.json#L8291-L8296
	value := v.(string)
	if len(value) > 128 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 128 characters", k))
	}
	if !regexp.MustCompile(`^[\w+=,.@-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(`%q must match [\w+=,.@-]`, k))
	}
	return
}

func validateIamRolePolicyNamePrefix(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 100 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 100 characters", k))
	}
	if !regexp.MustCompile(`^[\w+=,.@-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(`%q must match [\w+=,.@-]`, k))
	}
	return
}
