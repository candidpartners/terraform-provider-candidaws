package aws

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/quicksight"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsQuickSightNamespace() *schema.Resource {
	return &schema.Resource{
		// NOTE: It is possible for a namespace to get stuck in "CREATING" status if an account has
		//		 not completed QuickSight signup. 
		Create: resourceAwsQuickSightNamespaceCreate,
		Read:   resourceAwsQuickSightNamespaceRead,

		// NOTE: AWS QuickSight Namespace does not have a dedicated edit/update endpoint.
		//Update: resourceAwsQuickSightNamespaceUpdate,

		// NOTE: Deleting an AWS QuickSight Namespace will also delete users and groups
		//		 associated with that namespace.
		//		 ref: https://docs.aws.amazon.com/sdk-for-go/api/service/quicksight/#QuickSight.DeleteNamespace
		Delete: resourceAwsQuickSightNamespaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"aws_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"identity_store": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					quicksight.IdentityTypeQuicksight,
				}, false),
			},

			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			//"tags": tagsSchemaForceNew(), // TODO use this helper later in place of inline below
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
func resourceAwsQuickSightNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID := meta.(*AWSClient).accountid
	namespace := d.Get("namespace").(string)

	if v, ok := d.GetOk("aws_account_id"); ok {
		awsAccountID = v.(string)
	}

	createOpts := &quicksight.CreateNamespaceInput{
		AwsAccountId:  aws.String(awsAccountID),
		Namespace:     aws.String(namespace),
		IdentityStore: aws.String(d.Get("identity_store").(string)),
	}

	if attr, ok := d.GetOk("tags"); ok {
        createOpts.Tags = keyvaluetags.New(attr.(map[string]interface{})).IgnoreAws().QuicksightTags()
    }

	_, err := conn.CreateNamespace(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating QuickSight Namespace: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", awsAccountID, namespace))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"CREATED"},
		Refresh:    stateRefresh(conn, awsAccountID, namespace),
		Timeout:    15 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
		return fmt.Errorf("Error waiting for QuickSight Namespace (%s) to become deleted: %s", d.Id(), err)
    }
	return resourceAwsQuickSightNamespaceRead(d, meta)
}

func resourceAwsQuickSightNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, namespace, err := resourceAwsQuickSightNamespaceParseID(d.Id())
	if err != nil {
		return err
	}

	descOpts := &quicksight.DescribeNamespaceInput{
		AwsAccountId:  aws.String(awsAccountID),
		Namespace:     aws.String(namespace),
	}

	resp, err := conn.DescribeNamespace(descOpts)
	if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] QuickSight Namespace %s is not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error describing QuickSight Namespace (%s): %s", d.Id(), err)
	}

	d.Set("namespace", resp.Namespace.Name)
	d.Set("aws_account_id", awsAccountID)

	return nil
}

func resourceAwsQuickSightNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, namespace, err := resourceAwsQuickSightNamespaceParseID(d.Id())
	if err != nil {
		return err
	}

	deleteOpts := &quicksight.DeleteNamespaceInput{
		AwsAccountId: aws.String(awsAccountID),
		Namespace: aws.String(namespace),
	}

	if _, err := conn.DeleteNamespace(deleteOpts); err != nil {
		if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("Error deleting QuickSight Namespace %s: %s", d.Id(), err)
	}

    stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "CREATING"},
		Target:     []string{"DELETED"},
		Refresh:    stateRefresh(conn, awsAccountID, namespace),
		Timeout:    15 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
		return fmt.Errorf("Error waiting for QuickSight Namespace (%s) to become deleted: %s", d.Id(), err)
    }

	return nil
}

func resourceAwsQuickSightNamespaceParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected AWS_ACCOUNT_ID/NAMESPACE", id)
	}
	return parts[0], parts[1], nil
}

func stateRefresh(conn *quicksight.QuickSight, awsAccountID, namespace string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		emptyResp := &quicksight.DescribeNamespaceOutput{}

		resp, err := conn.DescribeNamespace(&quicksight.DescribeNamespaceInput {
			AwsAccountId: aws.String(awsAccountID),
			Namespace: aws.String(namespace),
		})
		if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
			return emptyResp, "DELETED", nil
		}
		creationStatus := *resp.Namespace.CreationStatus
		if err != nil {
			return nil, creationStatus, err
		}

		return resp, creationStatus, nil
	}
}