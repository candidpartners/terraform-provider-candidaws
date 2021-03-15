package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsInternetGatewayDetach() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsInternetGatewayDetachCreate,
		Read:   resourceAwsInternetGatewayDetachRead,
		Delete: resourceAwsInternetGatewayDetachDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}
func resourceAwsInternetGatewayDetachCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn
	awsVpcID := d.Get("vpc_id").(string)
	awsInternetGatewayID := d.Get("internet_gateway_id").(string)
	createOpts := &ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String(awsInternetGatewayID),
		VpcId:             aws.String(awsVpcID),
	}
	resp, err := conn.DetachInternetGateway(createOpts)
	fmt.Println("resp", resp)
	if err != nil {
		return fmt.Errorf("error detaching igw: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", awsVpcID, awsInternetGatewayID))
	return resourceAwsInternetGatewayDetachRead(d, meta)
}
func resourceAwsInternetGatewayDetachRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Default IGW has been detached")
	return nil
}
func resourceAwsInternetGatewayDetachDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Default IGW has been detached")
	return nil
}
