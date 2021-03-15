package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsInternetGatewayDelete() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsInternetGatewayDeleteCreate,
		Read:   resourceAwsInternetGatewayDeleteRead,
		Delete: resourceAwsInternetGatewayDeleteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"internet_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
func resourceAwsInternetGatewayDeleteCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn
	awsInternetGatewayID := d.Get("internet_gateway_id").(string)
	createOpts := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(awsInternetGatewayID),
	}
	resp, err := conn.DeleteInternetGateway(createOpts)
	fmt.Println("resp", resp)
	if err != nil {
		return fmt.Errorf("error deleteing igw: %s", err)
	}
	d.SetId(fmt.Sprintf("%s", awsInternetGatewayID))
	return resourceAwsInternetGatewayDeleteRead(d, meta)
}
func resourceAwsInternetGatewayDeleteRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Default IGW has been detached")
	return nil
}
func resourceAwsInternetGatewayDeleteDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Default IGW has been detached")
	return nil
}
