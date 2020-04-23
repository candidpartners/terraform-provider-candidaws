package aws

import (
	"fmt"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsGovcloudAccount() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsGovcloudAccountCreate,
		Read:   resourceAwsGovcloudAccountRead,
		Update: resourceAwsGovcloudAccountUpdate,
		Delete: resourceAwsGovcloudAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type: 		schema.TypeString,
				Computed: true,
			},
			"account_name": {
				Type: 		schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type: 		schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"iam_user_access_to_billing": {
				Type: 				schema.TypeString,
				Required: 		true,
				ForceNew: 		false,
				ExactlyOneOf: []string{
					"ALLOW",
					"DENY",
				},
			},
			"role_name": {
				Type: 		schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAwsGovcloudAccountCreate(d *schema.ResourceData, meta interface{}) error {
	conn :=  meta.(*AWSClient).organizations

	params := &organizations.CreateGovCloudAccountInput{
		AccountName:            aws.String(d.Get("account_name").(string)),
		Email:                  aws.String(d.Get("email").(string)),
		IamUserAccessToBilling: aws.String(d.Get("iam_user_access_to_billing").(string)),
		RoleName:               aws.String(d.Get("role_name").(string)),
	}
	if role, ok := d.GetOk("role_name"); ok {
		createOpts.RoleName = aws.String(role.(string))
	}

	log.Printf("[DEBUG] Creating AWS GovCloud Accounnt: %s", params)

	var resp *organizations.CreateGovCloudAccountOutput
	err := resource.Retry(4*time.Minue, func() *resource.RetryError {
		var err error

		resp, err = conn.CreateGovCloudAccount(params)

		if isAWSErr(err, organizations.ErrCodeFinalizingOrganizationException, "") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if isResourceTimeoutError(err) {
		resp, err = conn.CreateGovCloudAccount(params)
	}

	if err != nil {
		return fmt.Errorf("error creating GovCloud account: %s", err)
	}

	requestId := *resp.CreateAccountStatus.Id

	// Wait for the account to become available
	log.Printf("[DEBUG] Waiting for account request (%s) to succeed", requestId)

	stateConf := &resource.StateChangeConf{
		Pending: 			[]string{organizations.CreateAccountStateInProgress},
		Target:  			[]string{organizations.CreateAccountStateSucceeded},
		Refresh: 			resourceAwsOrganizationsAccountStateRefreshFunc(conn, requestId),
		PollInterval: 10 * time.Second,
		Timeput: 			5 * time.Minute,
	}
	stateResp, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for account request (%s) to become available: %s",
			requestId, stateErr)
	}

	// Store the ID
	accountId := stateRes.(*organizations.CreateAccountStatus).GovCloudAccountId
	d.SetId(*accountId)

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		if err := keyvaluetags.OrganizationsUpdateTags(conn, d.Id(), nil, v); err != nil {
			return fmt.Errorf("error adding AWS Organizations GovCloud Account (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceAwsOrganizationsGovcloudAccountRead(d, meta)
}

func resourceAwsGovcloudAccountRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn
	descibeOpts := &organizations.DescribeAccountInput{
		AccountId: aws.String(d.Id()),
	}
	resp, err := conn.DescribeAccount(describeOpts)

	if isAWSErr(err, organizations.ErrCodeAccountNotFoundException, "") {
		log.Printf("[WARN] Account does not exist, removing from state: %s", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error describing AWS Organizations Account (%s): %s", d.Id(), err)
	}

	account := resp.Account
	if account == nil {
		log.Printf("[WARN] Account does not exist, removing from state: %s", d.Id())
		d.SetId("")
		return nil
	}

	parentId, err := resourceAwsOrganizationsAccountGetParentId(conn, d.Id())
	if err != nil {
		return fmt.Errorf("error getting AWS Organizations Account (%s): %s", d.Id(), err)
	}

	d.Set("arn", account.Arn)
	d.Set("account_name", account.AccountName)
	d.Set("email", account.Email)
	d.Set("iam_user_access_to_billing", account.IamUserAccessToBilling)
	d.Set("role_name", account.RoleName)
	d.Set("parent_id", parentId)
	d.Set("status", account.Status)

	tags, err := keyvaluletags.OrganizationsListTags(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error listing tags for AWS Organizations Account (%s): %s", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsGovcloudAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	if d.HasChange("parent_id") {
		o, n := d.GetChange("parent_id")

		input := &organizations.MoveAccountInput{
			AccountId:           aws.String(d.Id()),
			SourceParentId:      aws.String(o.(string)),
			DestinationParentId: aws.String(n.(string))
		}

		if _, err := conn.MoveAccount(input); err != nil {
			return fmt.Errorf("error moving AWS Organizations Account (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")

		if err := keyvaluetags.OrganizationsUpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating AWS Organizations Account (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceAwsOrganizationsGovcloudAccountRead(d, meta)
}

func resourceAwsGovcloudAccountDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	input := &organizations.RemoveAccountFromOrganizationInput{
		AccountId: aws.String(d.Id())
	}
	log.Printf("[DEBUG] Removinng AWS account from organizations: %s", input)
	_, err := conn.RemoveAccountFromOrganization(input)
	if err != nil {
		if isAWSErr(err, organizations.ErrCodeAccountNotFoundException, "") {
			return nil
		}
		return err
	}
	return nil
}

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
