package aws

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/transfer"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsTransferServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsTransferServerCreate,
		Read:   resourceAwsTransferServerRead,
		Update: resourceAwsTransferServerUpdate,
		Delete: resourceAwsTransferServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  transfer.EndpointTypePublic,
				ValidateFunc: validation.StringInSlice([]string{
					transfer.EndpointTypePublic,
					transfer.EndpointTypeVpc,
					transfer.EndpointTypeVpcEndpoint,
				}, false),
			},

			"endpoint_details": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_endpoint_id": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"endpoint_details.0.address_allocation_ids", "endpoint_details.0.subnet_ids", "endpoint_details.0.vpc_id"},
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								validNamePattern := "^vpce-[0-9a-f]{17}$"
								validName, nameMatchErr := regexp.MatchString(validNamePattern, value)
								if !validName || nameMatchErr != nil {
									errors = append(errors, fmt.Errorf(
										"%q must match regex '%v'", k, validNamePattern))
								}
								return
							},
							DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
								if n == "" && d.Get("endpoint_type").(string) == transfer.EndpointTypeVpc {
									return true
								}
								return false
							},
						},
						"address_allocation_ids": {
							Type:          schema.TypeSet,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							ConflictsWith: []string{"endpoint_details.0.vpc_endpoint_id"},
						},
						"subnet_ids": {
							Type:          schema.TypeSet,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							ConflictsWith: []string{"endpoint_details.0.vpc_endpoint_id"},
						},
						"vpc_id": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  validation.NoZeroValues,
							ConflictsWith: []string{"endpoint_details.0.vpc_endpoint_id"},
						},
					},
				},
			},

			"vpce_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"host_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(0, 4096),
			},

			"host_key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"invocation_role": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateArn,
			},

			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"identity_provider_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  transfer.IdentityProviderTypeServiceManaged,
				ValidateFunc: validation.StringInSlice([]string{
					transfer.IdentityProviderTypeServiceManaged,
					transfer.IdentityProviderTypeApiGateway,
				}, false),
			},

			"logging_role": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateArn,
			},

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsTransferServerCreate(d *schema.ResourceData, meta interface{}) error {
	updateAfterCreate := false
	conn := meta.(*AWSClient).transferconn
	tags := keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws().TransferTags()
	createOpts := &transfer.CreateServerInput{}

	if len(tags) != 0 {
		createOpts.Tags = tags
	}

	identityProviderDetails := &transfer.IdentityProviderDetails{}
	if attr, ok := d.GetOk("invocation_role"); ok {
		identityProviderDetails.InvocationRole = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("url"); ok {
		identityProviderDetails.Url = aws.String(attr.(string))
	}

	if identityProviderDetails.Url != nil || identityProviderDetails.InvocationRole != nil {
		createOpts.IdentityProviderDetails = identityProviderDetails
	}

	if attr, ok := d.GetOk("identity_provider_type"); ok {
		createOpts.IdentityProviderType = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("logging_role"); ok {
		createOpts.LoggingRole = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("endpoint_type"); ok {
		createOpts.EndpointType = aws.String(attr.(string))
	}

	if attr, ok := d.GetOk("endpoint_details"); ok {
		createOpts.EndpointDetails = expandTransferServerEndpointDetails(attr.([]interface{}))

		if createOpts.EndpointDetails.AddressAllocationIds != nil {
			createOpts.EndpointDetails.AddressAllocationIds = nil
			updateAfterCreate = true
		}
	}

	if attr, ok := d.GetOk("host_key"); ok {
		createOpts.HostKey = aws.String(attr.(string))
	}

	log.Printf("[DEBUG] Create Transfer Server Option: %#v", createOpts)

	resp, err := conn.CreateServer(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating Transfer Server: %s", err)
	}

	d.SetId(*resp.ServerId)

	if updateAfterCreate {
		updateOpts := &transfer.UpdateServerInput{
			ServerId:        aws.String(d.Id()),
			EndpointDetails: expandTransferServerEndpointDetails(d.Get("endpoint_details").([]interface{})),
		}

		if err := doTransferServerUpdate(d, updateOpts, conn, meta.(*AWSClient).ec2conn, true); err != nil {
			return err
		}
	}

	return resourceAwsTransferServerRead(d, meta)
}

func resourceAwsTransferServerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).transferconn

	descOpts := &transfer.DescribeServerInput{
		ServerId: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Describe Transfer Server Option: %#v", descOpts)

	resp, err := conn.DescribeServer(descOpts)
	if err != nil {
		if isAWSErr(err, transfer.ErrCodeResourceNotFoundException, "") {
			log.Printf("[WARN] Transfer Server (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	endpoint := meta.(*AWSClient).RegionalHostname(fmt.Sprintf("%s.server.transfer", d.Id()))

	d.Set("arn", resp.Server.Arn)
	d.Set("endpoint", endpoint)
	d.Set("invocation_role", "")
	d.Set("url", "")
	if resp.Server.IdentityProviderDetails != nil {
		d.Set("invocation_role", aws.StringValue(resp.Server.IdentityProviderDetails.InvocationRole))
		d.Set("url", aws.StringValue(resp.Server.IdentityProviderDetails.Url))
	}
	d.Set("endpoint_type", resp.Server.EndpointType)
	d.Set("endpoint_details", flattenTransferServerEndpointDetails(resp.Server.EndpointDetails))
	d.Set("identity_provider_type", resp.Server.IdentityProviderType)
	d.Set("logging_role", resp.Server.LoggingRole)
	d.Set("host_key_fingerprint", resp.Server.HostKeyFingerprint)

	if resp.Server.EndpointDetails.VpcEndpointId != nil {
		out, err := describeTransferServerVPCEndpoint(meta.(*AWSClient).ec2conn, resp.Server.EndpointDetails.VpcEndpointId)
		if err != nil {
			return err
		}

		d.Set("vpce_security_group_ids", flattenVpcEndpointSecurityGroupIds(out.Groups))
	}

	if err := d.Set("tags", keyvaluetags.TransferKeyValueTags(resp.Server.Tags).IgnoreAws().Map()); err != nil {
		return fmt.Errorf("Error setting tags: %s", err)
	}
	return nil
}

func resourceAwsTransferServerUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).transferconn
	updateFlag := false
	stopFlag := false
	updateOpts := &transfer.UpdateServerInput{
		ServerId: aws.String(d.Id()),
	}

	if d.HasChange("logging_role") {
		updateFlag = true
		updateOpts.LoggingRole = aws.String(d.Get("logging_role").(string))
	}

	if d.HasChange("invocation_role") || d.HasChange("url") {
		identityProviderDetails := &transfer.IdentityProviderDetails{}
		updateFlag = true
		if attr, ok := d.GetOk("invocation_role"); ok {
			identityProviderDetails.InvocationRole = aws.String(attr.(string))
		}

		if attr, ok := d.GetOk("url"); ok {
			identityProviderDetails.Url = aws.String(attr.(string))
		}
		updateOpts.IdentityProviderDetails = identityProviderDetails
	}

	if d.HasChange("endpoint_type") {
		updateFlag = true
		if attr, ok := d.GetOk("endpoint_type"); ok {
			updateOpts.EndpointType = aws.String(attr.(string))
		}
	}

	if d.HasChange("endpoint_details") {
		updateFlag = true
		if attr, ok := d.GetOk("endpoint_details"); ok {
			updateOpts.EndpointDetails = expandTransferServerEndpointDetails(attr.([]interface{}))
		}

		if d.HasChange("endpoint_details.0.address_allocation_ids") {
			stopFlag = true
		}
	}

	if d.HasChange("host_key") {
		updateFlag = true
		if attr, ok := d.GetOk("host_key"); ok {
			updateOpts.HostKey = aws.String(attr.(string))
		}
	}

	if updateFlag {
		if err := doTransferServerUpdate(d, updateOpts, conn, meta.(*AWSClient).ec2conn, stopFlag); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if err := keyvaluetags.TransferUpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	return resourceAwsTransferServerRead(d, meta)
}

func resourceAwsTransferServerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).transferconn

	if d.Get("force_destroy").(bool) {
		log.Printf("[DEBUG] Transfer Server (%s) attempting to forceDestroy", d.Id())
		if err := deleteTransferUsers(conn, d.Id(), nil); err != nil {
			return err
		}
	}

	delOpts := &transfer.DeleteServerInput{
		ServerId: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Delete Transfer Server Option: %#v", delOpts)

	_, err := conn.DeleteServer(delOpts)
	if err != nil {
		if isAWSErr(err, transfer.ErrCodeResourceNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("error deleting Transfer Server (%s): %s", d.Id(), err)
	}

	if err := waitForTransferServerDeletion(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for Transfer Server (%s): %s", d.Id(), err)
	}

	return nil
}

func waitForTransferServerDeletion(conn *transfer.Transfer, serverID string) error {
	params := &transfer.DescribeServerInput{
		ServerId: aws.String(serverID),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := conn.DescribeServer(params)

		if isAWSErr(err, transfer.ErrCodeResourceNotFoundException, "") {
			return nil
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Transfer Server (%s) still exists", serverID))
	})
	if isResourceTimeoutError(err) {
		_, err = conn.DescribeServer(params)
		if isAWSErr(err, transfer.ErrCodeResourceNotFoundException, "") {
			return nil
		}
		if err == nil {
			return fmt.Errorf("Transfer server (%s) still exists", serverID)
		}
	}
	if err != nil {
		return fmt.Errorf("Error waiting for transfer server deletion: %s", err)
	}
	return nil
}

func deleteTransferUsers(conn *transfer.Transfer, serverID string, nextToken *string) error {
	listOpts := &transfer.ListUsersInput{
		ServerId:  aws.String(serverID),
		NextToken: nextToken,
	}

	log.Printf("[DEBUG] List Transfer User Option: %#v", listOpts)

	resp, err := conn.ListUsers(listOpts)
	if err != nil {
		return err
	}

	for _, user := range resp.Users {

		delOpts := &transfer.DeleteUserInput{
			ServerId: aws.String(serverID),
			UserName: user.UserName,
		}

		log.Printf("[DEBUG] Delete Transfer User Option: %#v", delOpts)

		_, err = conn.DeleteUser(delOpts)
		if err != nil {
			if isAWSErr(err, transfer.ErrCodeResourceNotFoundException, "") {
				continue
			}
			return fmt.Errorf("error deleting Transfer User (%s) for Server(%s): %s", *user.UserName, serverID, err)
		}
	}

	if resp.NextToken != nil {
		return deleteTransferUsers(conn, serverID, resp.NextToken)
	}

	return nil
}

func expandTransferServerEndpointDetails(l []interface{}) *transfer.EndpointDetails {
	if len(l) == 0 || l[0] == nil {
		return nil
	}
	e := l[0].(map[string]interface{})

	out := &transfer.EndpointDetails{}

	if v, ok := e["vpc_endpoint_id"]; ok && v != "" {
		out.VpcEndpointId = aws.String(v.(string))
	}

	if v, ok := e["address_allocation_ids"]; ok {
		out.AddressAllocationIds = expandStringSet(v.(*schema.Set))
	}

	if v, ok := e["subnet_ids"]; ok {
		out.SubnetIds = expandStringSet(v.(*schema.Set))
	}

	if v, ok := e["vpc_id"]; ok {
		out.VpcId = aws.String(v.(string))
	}

	return out
}

func flattenTransferServerEndpointDetails(endpointDetails *transfer.EndpointDetails) []interface{} {
	if endpointDetails == nil {
		return []interface{}{}
	}

	e := make(map[string]interface{})
	if endpointDetails.VpcEndpointId != nil {
		e["vpc_endpoint_id"] = *endpointDetails.VpcEndpointId
	}
	if endpointDetails.AddressAllocationIds != nil {
		e["address_allocation_ids"] = endpointDetails.AddressAllocationIds
	}
	if endpointDetails.SubnetIds != nil {
		e["subnet_ids"] = endpointDetails.SubnetIds
	}
	if endpointDetails.VpcId != nil {
		e["vpc_id"] = *endpointDetails.VpcId
	}

	return []interface{}{e}
}

func doTransferServerUpdate(d *schema.ResourceData, updateOpts *transfer.UpdateServerInput, transferConn *transfer.Transfer, ec2Conn *ec2.EC2, stopFlag bool) error {
	if stopFlag {
		if err := waitForTransferServerVPCEndpointState(transferConn, ec2Conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
			return fmt.Errorf("error waiting for Transfer Server VPC Endpoint (%s) to start: %s", d.Id(), err)
		}

		if err := stopAndWaitForTransferServer(d.Id(), transferConn, d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	_, err := transferConn.UpdateServer(updateOpts)
	if err != nil {
		if isAWSErr(err, transfer.ErrCodeResourceNotFoundException, "") {
			log.Printf("[WARN] Transfer Server (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error updating Transfer Server (%s): %s", d.Id(), err)
	}

	if stopFlag {
		if err := startAndWaitForTransferServer(d.Id(), transferConn, d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	if err := updateTransferServerVPCEndpointSecurityGroup(d, transferConn, ec2Conn); err != nil {
		return err
	}

	return nil
}

func stopAndWaitForTransferServer(serverId string, conn *transfer.Transfer, timeout time.Duration) error {
	stopReq := &transfer.StopServerInput{
		ServerId: aws.String(serverId),
	}
	if _, err := conn.StopServer(stopReq); err != nil {
		return fmt.Errorf("error stopping Transfer Server (%s): %s", serverId, err)
	}

	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{transfer.StateStarting, transfer.StateOnline, transfer.StateStopping},
		Target:  []string{transfer.StateOffline},
		Refresh: refreshTransferServerStatus(conn, serverId),
		Timeout: timeout,
		Delay:   10 * time.Second,
	}

	if _, err := stateChangeConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for Transfer Server (%s) to stop: %s", serverId, err)
	}

	return nil
}

func startAndWaitForTransferServer(serverId string, conn *transfer.Transfer, timeout time.Duration) error {
	stopReq := &transfer.StartServerInput{
		ServerId: aws.String(serverId),
	}

	if _, err := conn.StartServer(stopReq); err != nil {
		return fmt.Errorf("error starting Transfer Server (%s): %s", serverId, err)
	}

	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{transfer.StateStarting, transfer.StateOffline, transfer.StateStopping},
		Target:  []string{transfer.StateOnline},
		Refresh: refreshTransferServerStatus(conn, serverId),
		Timeout: timeout,
		Delay:   10 * time.Second,
	}

	if _, err := stateChangeConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for Transfer Server (%s) to start: %s", serverId, err)
	}

	return nil
}

func refreshTransferServerStatus(conn *transfer.Transfer, serverId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		server, err := describeTransferServer(conn, serverId)

		if server == nil {
			return 42, "destroyed", nil
		}

		return server, aws.StringValue(server.State), err
	}
}

func describeTransferServer(conn *transfer.Transfer, serverId string) (*transfer.DescribedServer, error) {
	params := &transfer.DescribeServerInput{
		ServerId: aws.String(serverId),
	}

	resp, err := conn.DescribeServer(params)

	return resp.Server, err
}

func waitForTransferServerVPCEndpointState(transferConn *transfer.Transfer, ec2Conn *ec2.EC2, serverId string, timeout time.Duration) error {
	server, err := describeTransferServer(transferConn, serverId)

	if err != nil {
		return err
	}

	if err := vpcEndpointWaitUntilAvailable(ec2Conn, *server.EndpointDetails.VpcEndpointId, timeout); err != nil {
		return err
	}

	return nil
}

func describeTransferServerVPCEndpoint(conn *ec2.EC2, vpceId *string) (*ec2.VpcEndpoint, error) {
	params := &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []*string{vpceId},
	}

	resp, err := conn.DescribeVpcEndpoints(params)

	return resp.VpcEndpoints[0], err
}

func updateTransferServerVPCEndpointSecurityGroup(d *schema.ResourceData, transferConn *transfer.Transfer, ec2conn *ec2.EC2) error {
	server, err := describeTransferServer(transferConn, d.Id())

	if err != nil {
		return fmt.Errorf("error describing Transfer Server VPC Endpoint: %s", err)
	}

	req := &ec2.ModifyVpcEndpointInput{
		VpcEndpointId: aws.String(*server.EndpointDetails.VpcEndpointId),
	}

	if d.IsNewResource() {
		out, err := describeTransferServerVPCEndpoint(ec2conn, server.EndpointDetails.VpcEndpointId)
		if err != nil {
			return err
		}

		req.RemoveSecurityGroupIds = append(req.RemoveSecurityGroupIds, out.Groups[0].GroupId)

		setVpcEndpointCreateList(d, "vpce_security_group_ids", &req.AddSecurityGroupIds)
	} else {
		setVpcEndpointUpdateLists(d, "vpce_security_group_ids", &req.AddSecurityGroupIds, &req.RemoveSecurityGroupIds)
	}

	if _, err := ec2conn.ModifyVpcEndpoint(req); err != nil {
		return fmt.Errorf("error updating VPC Endpoint: %s", err)
	}

	return nil
}
