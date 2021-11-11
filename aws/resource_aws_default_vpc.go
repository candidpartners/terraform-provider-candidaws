package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsDefaultVpc() *schema.Resource {
	// reuse aws_vpc schema, and methods for READ, UPDATE
	dvpc := resourceAwsVpc()
	dvpc.Create = resourceAwsDefaultVpcCreate
	dvpc.Delete = resourceAwsDefaultVpcDelete
	dvpc.Read = resourceAwsDefaultVpcRead

	// cidr_block is a computed value for Default VPCs
	dvpc.Schema["cidr_block"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	// instance_tenancy is a computed value for Default VPCs
	dvpc.Schema["instance_tenancy"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	// assign_generated_ipv6_cidr_block is a computed value for Default VPCs
	dvpc.Schema["assign_generated_ipv6_cidr_block"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}

	return dvpc
}

func resourceAwsDefaultVpcCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn
	req := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("isDefault"),
				Values: aws.StringSlice([]string{"true"}),
			},
		},
	}

	resp, err := conn.DescribeVpcs(req)
	if err != nil {
		d.SetId("vpc-removed")
		return nil
	}

	if resp.Vpcs == nil || len(resp.Vpcs) == 0 {
		d.SetId("vpc-removed")
		return nil
	}
	if resp.Vpcs != nil || len(resp.Vpcs) > 0 {
		d.SetId(aws.StringValue(resp.Vpcs[0].VpcId))
		vpcID := d.Id()
		deleteVpcOpts := &ec2.DeleteVpcInput{
			VpcId: &vpcID,
		}
		log.Printf("[INFO] Deleting VPC: %s", d.Id())

		err2 := resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err2 := conn.DeleteVpc(deleteVpcOpts)
			if err2 == nil {
				return nil
			}

			if isAWSErr(err2, "InvalidVpcID.NotFound", "") {
				return nil
			}
			if isAWSErr(err2, "DependencyViolation", "") {
				return resource.RetryableError(err2)
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting VPC: %s", err2))
		})
		if isResourceTimeoutError(err2) {
			_, err2 = conn.DeleteVpc(deleteVpcOpts)
			if isAWSErr(err2, "InvalidVpcID.NotFound", "") {
				return nil
			}
		}

		if err != nil {
			return fmt.Errorf("Error deleting VPC: %s", err)
		}
	}

	return resourceAwsDefaultVpcRead(d, meta)
}

func resourceAwsDefaultVpcRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Default vpc has been deleted")
	return nil
}

func resourceAwsDefaultVpcDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy Default VPC. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
