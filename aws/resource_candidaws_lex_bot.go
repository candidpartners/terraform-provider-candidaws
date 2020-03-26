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
			"voice_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clarification_prompt": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"messages": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_number": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
								},
							},
						},
						"response_card": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  1,
						},
						"max_attempts": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"idle_session_ttl_in_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
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
			"intents": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"intent_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"intent_version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"abort_statement": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"messages": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_number": {
										Type:     schema.TypeInt,
										Optional: true,
										Default: 1,
									},
								},
							},
						},
						"response_card": {
							Type:     schema.TypeString,
							Optional: true,
							Default: 1,
						},
					},
				},
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
		AbortStatement:  expandStatement(d.Get("abort_statement").([]interface{})),
		Intents:  expandIntents(d.Get("intents").([]interface{})),
		ClarificationPrompt:  expandPrompt(d.Get("clarification_prompt").([]interface{})),
	}
	if vid, ok := d.GetOk("voice_id"); ok {
		params.VoiceId = aws.String(vid.(string))
	}
	if ttl, ok := d.GetOk("idle_session_ttl_in_seconds"); ok {
		params.IdleSessionTTLInSeconds = aws.Int64(int64(ttl.(int)))
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
	d.Set("voice_id", aws.StringValue(bot.VoiceId))
	d.Set("idle_session_ttl_in_seconds", aws.Int64Value(bot.IdleSessionTTLInSeconds))

	if err := d.Set("abort_statement", flattenStatement(bot.AbortStatement)); err != nil {
		return fmt.Errorf("error setting abort_statement: %s", err)
	}

	if err := d.Set("intents", flattenIntents(bot.Intents)); err != nil {
		return fmt.Errorf("error setting intents: %s", err)
	}

	if err := d.Set("clarification_prompt", flattenPrompt(bot.ClarificationPrompt)); err != nil {
		return fmt.Errorf("error setting clarification_prompt: %s", err)
	}

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
		AbortStatement:  expandStatement(d.Get("abort_statement").([]interface{})),
		Intents:  expandIntents(d.Get("intents").([]interface{})),
		ClarificationPrompt:  expandPrompt(d.Get("clarification_prompt").([]interface{})),
	}

	if d.HasChange("voice_id") {
		params.VoiceId = aws.String(d.Get("voice_id").(string))
	}
	if ttl, ok := d.GetOk("idle_session_ttl_in_seconds"); ok {
		params.IdleSessionTTLInSeconds = aws.Int64(int64(ttl.(int)))
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

func expandIntents(values []interface{}) []*lexmodelbuildingservice.Intent {
	valueSlice := []*lexmodelbuildingservice.Intent{}
	for _, element := range values {
		e := element.(map[string]interface{})

		m := &lexmodelbuildingservice.Intent{
			IntentName:           aws.String(e["intent_name"].(string)),
			IntentVersion:           aws.String(e["intent_version"].(string)),
		}

		valueSlice = append(valueSlice, m)
	}

	return valueSlice
}

func flattenIntents(cs []*lexmodelbuildingservice.Intent) []interface{} {
	valuesSlice := make([]interface{}, len(cs))
	if len(cs) > 0 {
		for i, v := range cs {
			m := make(map[string]interface{})
			m["intent_name"] = aws.StringValue(v.IntentName)
			m["intent_version"] = aws.StringValue(v.IntentVersion)
			valuesSlice[i] = m
		}
	}

	return valuesSlice
}

func getLexBot(name, version string, conn *lexmodelbuildingservice.LexModelBuildingService) (*lexmodelbuildingservice.GetBotOutput, error) {
	input := &lexmodelbuildingservice.GetBotInput{
		Name:    aws.String(name),
		VersionOrAlias: aws.String(version),
	}
	return conn.GetBot(input)
}
