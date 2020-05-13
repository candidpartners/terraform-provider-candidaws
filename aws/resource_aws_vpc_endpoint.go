package aws

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func vpcEndpointStateRefresh(conn *ec2.EC2, vpceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Reading VPC Endpoint: %s", vpceId)
		resp, err := conn.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{
			VpcEndpointIds: aws.StringSlice([]string{vpceId}),
		})
		if err != nil {
			if isAWSErr(err, "InvalidVpcEndpointId.NotFound", "") {
				return "", "deleted", nil
			}

			return nil, "", err
		}

		n := len(resp.VpcEndpoints)
		switch n {
		case 0:
			return "", "deleted", nil

		case 1:
			vpce := resp.VpcEndpoints[0]
			state := aws.StringValue(vpce.State)
			// No use in retrying if the endpoint is in a failed state.
			if state == "failed" {
				return nil, state, errors.New("VPC Endpoint is in a failed state")
			}
			return vpce, state, nil

		default:
			return nil, "", fmt.Errorf("Found %d VPC Endpoints for %s, expected 1", n, vpceId)
		}
	}
}

func vpcEndpointWaitUntilAvailable(conn *ec2.EC2, vpceId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "pendingAcceptance"},
		Refresh:    vpcEndpointStateRefresh(conn, vpceId),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC Endpoint (%s) to become available: %s", vpceId, err)
	}

	return nil
}

func setVpcEndpointCreateList(d *schema.ResourceData, key string, c *[]*string) {
	if v, ok := d.GetOk(key); ok {
		list := v.(*schema.Set).List()
		if len(list) > 0 {
			*c = expandStringList(list)
		}
	}
}

func setVpcEndpointUpdateLists(d *schema.ResourceData, key string, a, r *[]*string) {
	if d.HasChange(key) {
		o, n := d.GetChange(key)
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		add := expandStringList(ns.Difference(os).List())
		if len(add) > 0 {
			*a = add
		}

		remove := expandStringList(os.Difference(ns).List())
		if len(remove) > 0 {
			*r = remove
		}
	}
}

func flattenVpcEndpointSecurityGroupIds(groups []*ec2.SecurityGroupIdentifier) *schema.Set {
	vSecurityGroupIds := []interface{}{}

	for _, group := range groups {
		vSecurityGroupIds = append(vSecurityGroupIds, aws.StringValue(group.GroupId))
	}

	return schema.NewSet(schema.HashString, vSecurityGroupIds)
}
