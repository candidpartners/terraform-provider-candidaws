package aws

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^(r-[0-9a-z]{4,32})|(ou-[0-9a-z]{4,32}-[a-z0-9]{8,32})$"), "see https://docs.aws.amazon.com/organizations/latest/APIReference/API_MoveAccount.html#organizations-MoveAccount-request-DestinationParentId"),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAwsOrganizationsAccountEmail,
			},
			"iam_user_access_to_billing": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{organizations.IAMUserAccessToBillingAllow, organizations.IAMUserAccessToBillingDeny}, true),
			},
			"role_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAwsOrganizationsAccountRoleName,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAwsGovcloudAccountCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	params := &organizations.CreateGovCloudAccountInput{
		AccountName: aws.String(d.Get("account_name").(string)),
		Email:       aws.String(d.Get("email").(string)),
	}
	if role, ok := d.GetOk("role_name"); ok {
		params.RoleName = aws.String(role.(string))
	}
	if iam_user, ok := d.GetOk("iam_user_access_to_billing"); ok {
		params.IamUserAccessToBilling = aws.String(iam_user.(string))
	}

	log.Printf("[DEBUG] Creating AWS GovCloud Accounnt: %s", params)

	var resp *organizations.CreateGovCloudAccountOutput
	err := resource.Retry(4*time.Minute, func() *resource.RetryError {
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

	requestID := *resp.CreateAccountStatus.Id

	// Wait for the account to become available
	log.Printf("[DEBUG] Waiting for account request (%s) to succeed", requestID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{organizations.CreateAccountStateInProgress},
		Target:       []string{organizations.CreateAccountStateSucceeded},
		Refresh:      resourceAwsOrganizationsAccountStateRefreshFunc(conn, requestID),
		PollInterval: 10 * time.Second,
		Timeout:      5 * time.Minute,
	}
	stateResp, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"error waiting for account request (%s) to become available: %s",
			requestID, stateErr)
	}

	// Store the ID
	accountID := stateResp.(*organizations.CreateAccountStatus).GovCloudAccountId
	d.SetId(*accountID)

	if v, ok := d.GetOk("parent_id"); ok {
		newParentID := v.(string)

		existingParentID, err := resourceAwsOrganizationsAccountGetParentId(conn, d.Id())

		if err != nil {
			return fmt.Errorf("error getting AWS Organizations Account (%s) parent: %s", d.Id(), err)
		}

		if newParentID != existingParentID {
			input := &organizations.MoveAccountInput{
				AccountId:           accountID,
				SourceParentId:      aws.String(existingParentID),
				DestinationParentId: aws.String(newParentID),
			}

			if _, err := conn.MoveAccount(input); err != nil {
				return fmt.Errorf("error moving AWS Organizations Account (%s): %s", d.Id(), err)
			}
		}
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		if err := keyvaluetags.OrganizationsUpdateTags(conn, d.Id(), nil, v); err != nil {
			return fmt.Errorf("error adding AWS Organizations GovCloud Account (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceAwsGovcloudAccountRead(d, meta)
}

func resourceAwsGovcloudAccountRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn
	describeOpts := &organizations.DescribeAccountInput{
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

	parentID, err := resourceAwsOrganizationsAccountGetParentId(conn, d.Id())
	if err != nil {
		return fmt.Errorf("error getting AWS Organizations Account (%s): %s", d.Id(), err)
	}

	d.Set("arn", account.Arn)
	d.Set("joined_method", account.JoinedMethod)
	d.Set("joined_timestamp", aws.TimeValue(account.JoinedTimestamp).Format(time.RFC3339))
	d.Set("account_name", account.Name)
	d.Set("email", account.Email)
	d.Set("parent_id", parentID)
	d.Set("status", account.Status)

	tags, err := keyvaluetags.OrganizationsListTags(conn, d.Id())

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
			DestinationParentId: aws.String(n.(string)),
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

	return resourceAwsGovcloudAccountRead(d, meta)
}

func resourceAwsGovcloudAccountDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	input := &organizations.RemoveAccountFromOrganizationInput{
		AccountId: aws.String(d.Id()),
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
