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

// ACL Network ACLs all contain explicit deny-all rules that cannot be
// destroyed or changed by users. This rules are numbered very high to be a
// catch-all.
// See http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_ACLs.html#default-network-acl
const (
	awsDefaultAclRuleNumberIpv4 = 32767
	awsDefaultAclRuleNumberIpv6 = 32768
)

func resourceAwsDefaultNetworkAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsDefaultNetworkAclCreate,
		// We reuse aws_network_acl's read method, the operations are the same
		Read:   resourceAwsNetworkAclRead,
		Delete: resourceAwsDefaultNetworkAclDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"default_network_acl_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// We want explicit management of Subnets here, so we do not allow them to be
			// computed. Instead, an empty config will enforce just that; removal of the
			// any Subnets that have been assigned to the Default Network ACL. Because we
			// can't actually remove them, this will be a continual plan until the
			// Subnets are themselves destroyed or reassigned to a different Network
			// ACL
			"subnet_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				ForceNew: true,
			},
			// We want explicit management of Rules here, so we do not allow them to be
			// computed. Instead, an empty config will enforce just that; removal of the
			// rules
			"ingress": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_no": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv6_cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"icmp_code": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				Set: resourceAwsNetworkAclEntryHash,
			},
			"egress": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_no": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv6_cidr_block": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"icmp_code": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				Set: resourceAwsNetworkAclEntryHash,
			},

			"tags": tagsSchema2(),

			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAwsDefaultNetworkAclCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn
	d.SetId(d.Get("default_network_acl_id").(string))

	// revoke all default and pre-existing rules on the default network acl.
	// In the UPDATE method, we'll apply only the rules in the configuration.
	log.Printf("[DEBUG] Revoking default ingress and egress rules for Default Network ACL for %s", d.Id())
	err1 := revokeAllNetworkACLEntries(d.Id(), meta)
	if err1 != nil {
		return err1
	}

	log.Printf("[INFO] Deleting Network Acl: %s", d.Id())
	input := &ec2.DeleteNetworkAclInput{
		NetworkAclId: aws.String(d.Id()),
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteNetworkAcl(input)
		if err != nil {
			if isAWSErr(err, "InvalidNetworkAclID.NotFound", "") {
				return nil
			}
			if isAWSErr(err, "DependencyViolation", "") {
				err = cleanUpDependencyViolations(d, conn)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return resource.RetryableError(fmt.Errorf("Dependencies found and cleaned up, retrying"))
			}

			return resource.NonRetryableError(err)

		}
		log.Printf("[Info] Deleted network ACL %s successfully", d.Id())
		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.DeleteNetworkAcl(input)
		if err != nil && isAWSErr(err, "InvalidNetworkAclID.NotFound", "") {
			return nil
		}
		err = cleanUpDependencyViolations(d, conn)
		if err != nil {
			// This seems excessive but is probably the best way to make sure it's actually deleted
			_, err = conn.DeleteNetworkAcl(input)
			if err != nil && isAWSErr(err, "InvalidNetworkAclID.NotFound", "") {
				return nil
			}
		}
	}
	if err != nil {
		return fmt.Errorf("Error destroying Network ACL (%s): %s", d.Id(), err)
	}

	return resourceAwsDefaultNetworkAclRead(d, meta)
}

func resourceAwsDefaultNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Default rt has been deleted")
	return nil
}
func resourceAwsDefaultNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy Default Network ACL. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

// revokeAllNetworkACLEntries revoke all ingress and egress rules that the Default
// Network ACL currently has
func revokeAllNetworkACLEntries(netaclId string, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	resp, err := conn.DescribeNetworkAcls(&ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []*string{aws.String(netaclId)},
	})

	if err != nil {
		log.Printf("[DEBUG] Error looking up Network ACL: %s", err)
		return err
	}

	if resp == nil {
		return fmt.Errorf("Error looking up Default Network ACL Entries: No results")
	}

	networkAcl := resp.NetworkAcls[0]
	for _, e := range networkAcl.Entries {
		// Skip the default rules added by AWS. They can be neither
		// configured or deleted by users. See http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_ACLs.html#default-network-acl
		if *e.RuleNumber == awsDefaultAclRuleNumberIpv4 ||
			*e.RuleNumber == awsDefaultAclRuleNumberIpv6 {
			continue
		}

		// track if this is an egress or ingress rule, for logging purposes
		rt := "ingress"
		if *e.Egress {
			rt = "egress"
		}

		log.Printf("[DEBUG] Destroying Network ACL (%s) Entry number (%d)", rt, int(*e.RuleNumber))
		_, err := conn.DeleteNetworkAclEntry(&ec2.DeleteNetworkAclEntryInput{
			NetworkAclId: aws.String(netaclId),
			RuleNumber:   e.RuleNumber,
			Egress:       e.Egress,
		})
		if err != nil {
			return fmt.Errorf("Error deleting entry (%s): %s", e, err)
		}
	}

	return nil
}
