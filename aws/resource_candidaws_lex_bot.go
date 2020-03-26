package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceAwsLexBot() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsLexBotCreate,
		Read:   resourceAwsLexBotRead,
		Update: resourceAwsLexBotUpdate,
		Delete: resourceAwsLexBotDelete,
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
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en-US",
			},
			"child_directed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"publish": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAwsLexBotCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	params := &lexmodelbuildingservice.PutBotInput{
		Name:                aws.String(d.Get("name").(string)),
		Description:         aws.String(d.Get("description").(string)),
		CreateVersion:       aws.Bool(d.Get("publish").(bool)),
		ChildDirected:       aws.Bool(d.Get("child_directed").(bool)),
		Locale:       aws.String(d.Get("locale").(string)),
	}
	resp, err := conn.PutBot(params)
	if err != nil {
		return fmt.Errorf("error putting Lex bot: %s", err)
	}

	d.SetId(aws.StringValue(resp.Name))
	d.Set("version", resp.Version)
	d.Set("checksum", resp.Checksum)
	return resourceAwsLexBotRead(d, meta)
}

func resourceAwsLexBotRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	version := d.Get("version").(string)
	bot, err := getLexBot(d.Id(), version, conn)
	if err != nil {
		return fmt.Errorf("error getting Lex bot %q: %s", d.Id(), err)
	}
	if bot == nil {
		log.Printf("[WARN] LexModelBuildingService Bot %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", bot.Description)
	d.Set("checksum", bot.Checksum)
	d.Set("version", bot.Version)

	return nil
}

func resourceAwsLexBotUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	name := d.Id()

	params := &lexmodelbuildingservice.PutBotInput{
		Name:                aws.String(name),
		Checksum:            aws.String(d.Get("checksum").(string)),
		Description:         aws.String(d.Get("description").(string)),
		CreateVersion:       aws.Bool(d.Get("publish").(bool)),
		ChildDirected:       aws.Bool(d.Get("child_directed").(bool)),
		Locale:       aws.String(d.Get("locale").(string)),
	}

	resp, err := conn.PutBot(params)
	if err != nil {
		return err
	}

	d.Set("version", resp.Version)
	d.Set("checksum", resp.Checksum)
	return resourceAwsLexBotRead(d, meta)
}

func resourceAwsLexBotDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	name := d.Id()

	input := &lexmodelbuildingservice.DeleteBotInput{
		Name: aws.String(name),
	}
	_, err := conn.DeleteBot(input)
	if err != nil {
		if isAWSErr(err, lexmodelbuildingservice.ErrCodeNotFoundException, "") {
			return nil
		}
		return err
	}

	return nil
}

func getLexBot(name, version string, conn *lexmodelbuildingservice.LexModelBuildingService) (*lexmodelbuildingservice.GetBotOutput, error) {
	input := &lexmodelbuildingservice.GetBotInput{
		Name:    aws.String(name),
		VersionOrAlias: aws.String(version),
	}
	return conn.GetBot(input)
}
