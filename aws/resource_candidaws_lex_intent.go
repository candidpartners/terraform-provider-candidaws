package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceAwsLexIntent() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsLexIntentCreate,
		Read:   resourceAwsLexIntentRead,
		Update: resourceAwsLexIntentUpdate,
		Delete: resourceAwsLexIntentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed:true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed:true,
			},
			"publish": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAwsLexIntentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	params := &lexmodelbuildingservice.PutIntentInput{
		Name:              aws.String(d.Get("name").(string)),
		Description:              aws.String(d.Get("description").(string)),
		CreateVersion: aws.Bool(d.Get("publish").(bool)),
	}
	resp, err := conn.PutIntent(params)
	if err != nil {
		return fmt.Errorf("error putting Lex slot type: %s", err)
	}

	d.SetId(aws.StringValue(resp.Name))
	d.Set("version", resp.Version)
	d.Set("checksum", resp.Checksum)
	return resourceAwsLexIntentRead(d, meta)
}

func resourceAwsLexIntentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	version := d.Get("version").(string)
	slotType, err := getLexIntent(d.Id(), version, conn)
	if err != nil {
		return fmt.Errorf("error getting Lex slot type %q: %s", d.Id(), err)
	}
	if slotType == nil {
		log.Printf("[WARN] LexModelBuildingService Intent %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", slotType.Description)
	d.Set("checksum", slotType.Checksum)
	d.Set("version", slotType.Version)

	return nil
}

func resourceAwsLexIntentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	name := d.Id()

	params := &lexmodelbuildingservice.PutIntentInput{
		Name:     aws.String(name),
		Checksum: aws.String(d.Get("checksum").(string)),
		Description: aws.String(d.Get("description").(string)),
		CreateVersion: aws.Bool(d.Get("publish").(bool)),
	}

	resp, err := conn.PutIntent(params)
	if err != nil {
		return err
	}

	d.Set("version", resp.Version)
	d.Set("checksum", resp.Checksum)
	return resourceAwsLexIntentRead(d, meta)
}

func resourceAwsLexIntentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	name := d.Id()

	input := &lexmodelbuildingservice.DeleteIntentInput{
		Name: aws.String(name),
	}
	_, err := conn.DeleteIntent(input)
	if err != nil {
		if isAWSErr(err, lexmodelbuildingservice.ErrCodeNotFoundException, "") {
			return nil
		}
		return err
	}

	return nil
}

func getLexIntent(name, version string, conn *lexmodelbuildingservice.LexModelBuildingService) (*lexmodelbuildingservice.GetIntentOutput, error) {
	input := &lexmodelbuildingservice.GetIntentInput{
		Name:    aws.String(name),
		Version: aws.String(version),
	}
	return conn.GetIntent(input)
}
