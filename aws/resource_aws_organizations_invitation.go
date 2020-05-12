package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsOrganizationsInvitation() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsOrganizationsInvitationCreate,
		Read:   resourceAwsOrganizationsInvitationRead,
		Delete: resourceAwsOrganizationsInvitationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsOrganizationsInvitationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	params := &organizations.InviteAccountToOrganizationInput{
		Target: &organizations.HandshakeParty{
			Id:   aws.String(d.Get("account_id").(string)),
			Type: aws.String("ACCOUNT"),
		},
	}

	resp, err := conn.InviteAccountToOrganization(params)

	if err != nil {
		return fmt.Errorf("error inviting account to organization: %s", err)
	}

	d.SetId(*resp.Handshake.Id)
	d.Set("arn", resp.Handshake.Arn)

	return resourceAwsOrganizationsInvitationRead(d, meta)
}

func resourceAwsOrganizationsInvitationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	params := &organizations.DescribeHandshakeInput{
		HandshakeId: aws.String(d.Id()),
	}
	resp, err := conn.DescribeHandshake(params)

	if err != nil {
		return fmt.Errorf("error describing handshake (%s): %s", d.Id(), err)
	}

	handshake := resp.Handshake
	if handshake == nil {
		log.Printf("[WARN] Handshake does not exist, removing from state: %s", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("arn", handshake.Arn)

	return nil
}

func resourceAwsOrganizationsInvitationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	input := &organizations.CancelHandshakeInput{
		HandshakeId: aws.String(d.Id()),
	}
	_, err := conn.CancelHandshake(input)
	if err != nil {
		return err
	}
	return nil
}
