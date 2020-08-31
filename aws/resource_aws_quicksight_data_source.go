package aws

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/quicksight"
)

func resourceAwsQuickSightDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsQuickSightDataSourceCreate,
		Read:   resourceAwsQuickSightDataSourceRead,
		Update: resourceAwsQuickSightDataSourceUpdate,
		Delete: resourceAwsQuickSightDataSourceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"aws_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_source_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_source_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAwsQuickSightDataSourceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID := meta.(*AWSClient).accountid
	userName := d.Get("user_name").(string)
	password := d.Get("password").(string)
	dataSourceID := d.Get("data_source_id").(string)
	dataSourceName := d.Get("data_source_name").(string)
	dataSourceType := d.Get("data_source_type").(string)

	if v, ok := d.GetOk("aws_account_id"); ok {
		awsAccountID = v.(string)
	}

	createOpts := &quicksight.CreateDataSourceInput{
		AwsAccountId: aws.String(awsAccountID),
		Credentials: &quicksight.DataSourceCredentials{
			CredentialPair: &quicksight.CredentialPair{
				Password: aws.String(password),
				Username: aws.String(userName),
			},
		},
		DataSourceId: aws.String(dataSourceID),
		Name:         aws.String(dataSourceName),
		Type:         aws.String(dataSourceType),
	}

	resp, err := conn.CreateDataSource(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating QuickSight DataSource: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", awsAccountID, aws.StringValue(resp.DataSourceId)))

	return resourceAwsQuickSightDataSourceRead(d, meta)
}

func resourceAwsQuickSightDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(d.Id())
	if err != nil {
		return err
	}

	descOpts := &quicksight.DescribeDataSourceInput{
		AwsAccountId: aws.String(awsAccountID),
		DataSourceId: aws.String(dataSourceID),
	}

	resp, err := conn.DescribeDataSource(descOpts)
	if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] QuickSight DataSource %s is already gone", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error describing QuickSight DataSource (%s): %s", d.Id(), err)
	}

	d.Set("arn", resp.DataSource.Arn)
	d.Set("aws_account_id", awsAccountID)
	d.Set("data_source_id", resp.DataSource.DataSourceId)
	d.Set("data_source_name", resp.DataSource.Name)
	d.Set("data_source_type", resp.DataSource.Type)

	return nil
}

func resourceAwsQuickSightDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(d.Id())
	if err != nil {
		return err
	}

	updateOpts := &quicksight.UpdateDataSourceInput{
		AwsAccountId: aws.String(awsAccountID),
		DataSourceId: aws.String(dataSourceID),
	}

	if v, ok := d.GetOk("data_source_name"); ok {
		updateOpts.Name = aws.String(v.(string))
	}

	_, err = conn.UpdateDataSource(updateOpts)
	if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] QuickSight DataSource %s is already gone", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error updating QuickSight DataSource %s: %s", d.Id(), err)
	}

	return resourceAwsQuickSightDataSourceRead(d, meta)
}

func resourceAwsQuickSightDataSourceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(d.Id())
	if err != nil {
		return err
	}

	deleteOpts := &quicksight.DeleteDataSourceInput{
		AwsAccountId: aws.String(awsAccountID),
		DataSourceId: aws.String(dataSourceID),
	}

	if _, err := conn.DeleteDataSource(deleteOpts); err != nil {
		if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("Error deleting QuickSight DataSource %s: %s", d.Id(), err)
	}

	return nil
}

func resourceAwsQuickSightDataSourceParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected AWS_ACCOUNT_ID/DATA_SOURCE_ID", id)
	}
	return parts[0], parts[1], nil
}
