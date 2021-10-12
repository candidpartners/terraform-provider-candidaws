package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsDefaultSubnet() *schema.Resource {
	// reuse aws_subnet schema, and methods for READ, UPDATE
	dsubnet := resourceAwsSubnet()
	dsubnet.Create = resourceAwsDefaultSubnetCreate
	dsubnet.Read = resourceAwsDefaultSubnetRead
	dsubnet.Delete = resourceAwsDefaultSubnetDelete

	// availability_zone is a required value for Default Subnets
	dsubnet.Schema["availability_zone"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	// availability_zone_id is a computed value for Default Subnets
	dsubnet.Schema["availability_zone_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	// vpc_id is a computed value for Default Subnets
	dsubnet.Schema["vpc_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	// cidr_block is a computed value for Default Subnets
	dsubnet.Schema["cidr_block"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	// ipv6_cidr_block is a computed value for Default Subnets
	dsubnet.Schema["ipv6_cidr_block"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	// map_public_ip_on_launch is a computed value for Default Subnets
	dsubnet.Schema["map_public_ip_on_launch"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	}
	// assign_ipv6_address_on_creation is a computed value for Default Subnets
	dsubnet.Schema["assign_ipv6_address_on_creation"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}

	return dsubnet
}

func resourceAwsDefaultSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	req := &ec2.DescribeSubnetsInput{}
	req.Filters = buildEC2AttributeFilterList(
		map[string]string{
			"availabilityZone": d.Get("availability_zone").(string),
			"defaultForAz":     "true",
		},
	)

	log.Printf("[DEBUG] Reading Default Subnet: %s", req)
	resp, err := conn.DescribeSubnets(req)
	if err != nil {
		return nil
	}
	if len(resp.Subnets) != 1 || resp.Subnets[0] == nil {
		return nil
	}
	d.SetId(aws.StringValue(resp.Subnets[0].SubnetId))

	log.Printf("[INFO] Deleting subnet: %s", d.Id())

	if err := deleteLingeringLambdaENIs(conn, "subnet-id", d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error deleting Lambda ENIs using subnet (%s): %s", d.Id(), err)
	}

	req2 := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(d.Id()),
	}

	wait := resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"destroyed"},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 1 * time.Second,
		Refresh: func() (interface{}, string, error) {
			_, err := conn.DeleteSubnet(req2)
			if err != nil {
				if apiErr, ok := err.(awserr.Error); ok {
					if apiErr.Code() == "DependencyViolation" {
						// There is some pending operation, so just retry
						// in a bit.
						return 42, "pending", nil
					}

					if apiErr.Code() == "InvalidSubnetID.NotFound" {
						return 42, "destroyed", nil
					}
				}

				return 42, "failure", err
			}

			return 42, "destroyed", nil
		},
	}

	if _, err := wait.WaitForState(); err != nil {
		return fmt.Errorf("Error deleting subnet: %s", err)
	}

	return resourceAwsDefaultSubnetRead(d, meta)
}
func resourceAwsDefaultSubnetRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Default subnet has been deleted")
	return nil
}
func resourceAwsDefaultSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy Default Subnet. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
