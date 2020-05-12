package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsOrganizationsInvitationAcceptance() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsOrganizationsInvitationAcceptanceCreate,
		Read:   resourceAwsOrganizationsInvitationAcceptanceRead,
		Delete: resourceAwsOrganizationsInvitationAcceptanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"invitation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAwsOrganizationsInvitationAcceptanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	params := &organizations.AcceptHandshakeInput{
		HandshakeId: aws.String(d.Get("invitation_id").(string)),
	}

	resp, err := conn.AcceptHandshake(params)

	if err != nil {
		return fmt.Errorf("error accepting inviting to organization: %s", err)
	}

	d.SetId(*resp.Handshake.Id)

	return resourceAwsOrganizationsInvitationAcceptanceRead(d, meta)
}

func resourceAwsOrganizationsInvitationAcceptanceRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceAwsOrganizationsInvitationAcceptanceDelete(_ *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).organizationsconn

	input := &organizations.LeaveOrganizationInput{}
	_, err := conn.LeaveOrganization(input)
	if err != nil {
		return err
	}
	return nil
}
