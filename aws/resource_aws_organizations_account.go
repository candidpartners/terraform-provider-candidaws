package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"regexp"
)

// resourceAwsOrganizationsAccountStateRefreshFunc returns a resource.StateRefreshFunc
// that is used to watch a CreateAccount request
func resourceAwsOrganizationsAccountStateRefreshFunc(conn *organizations.Organizations, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := &organizations.DescribeCreateAccountStatusInput{
			CreateAccountRequestId: aws.String(id),
		}
		resp, err := conn.DescribeCreateAccountStatus(opts)
		if err != nil {
			if isAWSErr(err, organizations.ErrCodeCreateAccountStatusNotFoundException, "") {
				resp = nil
			} else {
				log.Printf("Error on OrganizationAccountStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our account yet. Return an empty state.
			return nil, "", nil
		}

		accountStatus := resp.CreateAccountStatus
		if *accountStatus.State == organizations.CreateAccountStateFailed {
			return nil, *accountStatus.State, fmt.Errorf(*accountStatus.FailureReason)
		}
		return accountStatus, *accountStatus.State, nil
	}
}

func validateAwsOrganizationsAccountEmail(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid email address", value))
	}

	if len(value) < 6 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be less than 6 characters", value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be greater than 64 characters", value))
	}

	return
}

func validateAwsOrganizationsAccountRoleName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w+=,.@-]{1,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must consist of uppercase letters, lowercase letters, digits with no spaces, and any of the following characters: =,.@-", value))
	}

	return
}

func resourceAwsOrganizationsAccountGetParentId(conn *organizations.Organizations, childId string) (string, error) {
	input := &organizations.ListParentsInput{
		ChildId: aws.String(childId),
	}
	var parents []*organizations.Parent

	err := conn.ListParentsPages(input, func(page *organizations.ListParentsOutput, lastPage bool) bool {
		parents = append(parents, page.Parents...)

		return !lastPage
	})

	if err != nil {
		return "", err
	}

	if len(parents) == 0 {
		return "", nil
	}

	// assume there is only a single parent
	// https://docs.aws.amazon.com/organizations/latest/APIReference/API_ListParents.html
	parent := parents[0]
	return aws.StringValue(parent.Id), nil
}
