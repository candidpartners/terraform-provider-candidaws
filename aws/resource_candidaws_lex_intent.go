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
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sample_utterances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"fulfillment_activity": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code_hook": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message_version": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"slots": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Managed by Terraform",
						},
						"obfuscation_setting": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "NONE",
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"response_card": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sample_utterances": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"slot_constraint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"slot_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"slot_type_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_elicitation_prompt": {
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
					},
				},
			},
			"confirmation_prompt": {
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
			"conclusion_statement": {
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
			"rejection_statement": {
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
		Name:                aws.String(d.Get("name").(string)),
		Description:         aws.String(d.Get("description").(string)),
		CreateVersion:       aws.Bool(d.Get("publish").(bool)),
		SampleUtterances:    expandSampleUtterances(d.Get("sample_utterances").([]interface{})),
		FulfillmentActivity: expandFulfillmentActivity(d.Get("fulfillment_activity").([]interface{})),
		ConfirmationPrompt:  expandPrompt(d.Get("confirmation_prompt").([]interface{})),
		RejectionStatement:  expandStatement(d.Get("rejection_statement").([]interface{})),
		ConclusionStatement:  expandStatement(d.Get("conclusion_statement").([]interface{})),
		Slots:               expandSlots(d.Get("slots").([]interface{})),
	}
	resp, err := conn.PutIntent(params)
	if err != nil {
		return fmt.Errorf("error putting Lex intent: %s", err)
	}

	d.SetId(aws.StringValue(resp.Name))
	d.Set("version", resp.Version)
	d.Set("checksum", resp.Checksum)
	return resourceAwsLexIntentRead(d, meta)
}

func resourceAwsLexIntentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	version := d.Get("version").(string)
	intent, err := getLexIntent(d.Id(), version, conn)
	if err != nil {
		return fmt.Errorf("error getting Lex intent %q: %s", d.Id(), err)
	}
	if intent == nil {
		log.Printf("[WARN] LexModelBuildingService Intent %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", intent.Description)
	d.Set("checksum", intent.Checksum)
	d.Set("version", intent.Version)

	if err := d.Set("sample_utterances", flattenSampleUtterances(intent.SampleUtterances)); err != nil {
		return fmt.Errorf("error setting sample_utterances: %s", err)
	}

	if err := d.Set("fulfillment_activity", flattenFulfillmentActivity(intent.FulfillmentActivity)); err != nil {
		return fmt.Errorf("error setting sample_utterances: %s", err)
	}

	if err := d.Set("confirmation_prompt", flattenPrompt(intent.ConfirmationPrompt)); err != nil {
		return fmt.Errorf("error setting confirmation_prompt: %s", err)
	}

	if err := d.Set("rejection_statement", flattenStatement(intent.RejectionStatement)); err != nil {
		return fmt.Errorf("error setting rejection_statement: %s", err)
	}

	if err := d.Set("conclusion_statement", flattenStatement(intent.ConclusionStatement)); err != nil {
		return fmt.Errorf("error setting conclusion_statement: %s", err)
	}

	if err := d.Set("slots", flattenSlots(intent.Slots)); err != nil {
		return fmt.Errorf("error setting slots: %s", err)
	}

	return nil
}

func resourceAwsLexIntentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	name := d.Id()

	params := &lexmodelbuildingservice.PutIntentInput{
		Name:                aws.String(name),
		Checksum:            aws.String(d.Get("checksum").(string)),
		Description:         aws.String(d.Get("description").(string)),
		CreateVersion:       aws.Bool(d.Get("publish").(bool)),
		SampleUtterances:    expandSampleUtterances(d.Get("sample_utterances").([]interface{})),
		FulfillmentActivity: expandFulfillmentActivity(d.Get("fulfillment_activity").([]interface{})),
		ConfirmationPrompt:  expandPrompt(d.Get("confirmation_prompt").([]interface{})),
		RejectionStatement:  expandStatement(d.Get("rejection_statement").([]interface{})),
		ConclusionStatement:  expandStatement(d.Get("conclusion_statement").([]interface{})),
		Slots:               expandSlots(d.Get("slots").([]interface{})),
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

func expandSampleUtterances(values []interface{}) []*string {
	valueSlice := []*string{}
	for _, element := range values {
		e := element.(string)

		valueSlice = append(valueSlice, &e)
	}

	return valueSlice
}

func expandFulfillmentActivity(values []interface{}) *lexmodelbuildingservice.FulfillmentActivity {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	return &lexmodelbuildingservice.FulfillmentActivity{
		Type:     aws.String(val["type"].(string)),
		CodeHook: expandCodeHook(val["code_hook"].([]interface{})),
	}
}

func expandCodeHook(values []interface{}) *lexmodelbuildingservice.CodeHook {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	return &lexmodelbuildingservice.CodeHook{
		MessageVersion: aws.String(val["message_version"].(string)),
		Uri:            aws.String(val["uri"].(string)),
	}
}

func expandPrompt(values []interface{}) *lexmodelbuildingservice.Prompt {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &lexmodelbuildingservice.Prompt{
		MaxAttempts: aws.Int64(int64(val["max_attempts"].(int))),
		Messages:    expandMessages(val["messages"].([]interface{})),
	}
	if responseCard, ok := val["response_card"]; ok {
		res.ResponseCard = aws.String(responseCard.(string))
	}
	return res
}

func expandMessages(values []interface{}) []*lexmodelbuildingservice.Message {
	valueSlice := []*lexmodelbuildingservice.Message{}
	for _, element := range values {
		e := element.(map[string]interface{})

		m := &lexmodelbuildingservice.Message{
			Content:     aws.String(e["content"].(string)),
			ContentType: aws.String(e["content_type"].(string)),
		}

		if groupNumber, ok := e["group_number"]; ok {
			m.GroupNumber = aws.Int64(int64(groupNumber.(int)))
		}

		valueSlice = append(valueSlice, m)
	}

	return valueSlice
}

func expandStatement(values []interface{}) *lexmodelbuildingservice.Statement {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &lexmodelbuildingservice.Statement{
		Messages: expandMessages(val["messages"].([]interface{})),
	}
	if responseCard, ok := val["response_card"]; ok {
		res.ResponseCard = aws.String(responseCard.(string))
	}
	return res
}

func expandSlots(values []interface{}) []*lexmodelbuildingservice.Slot {
	valueSlice := []*lexmodelbuildingservice.Slot{}
	for _, element := range values {
		e := element.(map[string]interface{})

		m := &lexmodelbuildingservice.Slot{
			Name:           aws.String(e["name"].(string)),
			SlotConstraint: aws.String(e["slot_constraint"].(string)),
		}

		if description, ok := e["description"]; ok {
			m.Description = aws.String(description.(string))
		}
		if obs, ok := e["obfuscation_setting"]; ok && obs.(string) != "" {
			m.ObfuscationSetting = aws.String(obs.(string))
		}
		if pri, ok := e["priority"]; ok {
			m.Priority = aws.Int64(int64(pri.(int)))
		}
		if rc, ok := e["response_card"]; ok && rc.(string) != "" {
			m.ResponseCard = aws.String(rc.(string))
		}
		if ut, ok := e["sample_utterances"]; ok {
			m.SampleUtterances = expandStringList(ut.([]interface{}))
		}
		if st, ok := e["slot_type"]; ok && st.(string) != "" {
			m.SlotType = aws.String(st.(string))
		}
		if st, ok := e["slot_type_version"]; ok && st.(string) != "" {
			m.SlotTypeVersion = aws.String(st.(string))
		}
		if p, ok := e["value_elicitation_prompt"]; ok {
			m.ValueElicitationPrompt = expandPrompt(p.([]interface{}))
		}

		valueSlice = append(valueSlice, m)
	}

	return valueSlice
}

func flattenSampleUtterances(cs []*string) []string {
	valuesSlice := make([]string, len(cs))
	if len(cs) > 0 {
		for i, v := range cs {
			valuesSlice[i] = *v
		}
	}

	return valuesSlice
}

func flattenFulfillmentActivity(cs *lexmodelbuildingservice.FulfillmentActivity) []map[string]interface{} {
	if cs == nil {
		return nil
	}

	val := make(map[string]interface{})

	val["type"] = aws.StringValue(cs.Type)

	if cs.CodeHook != nil {
		val["code_hook"] = flattenCodeHook(cs.CodeHook)
	}

	return []map[string]interface{}{val}
}

func flattenCodeHook(cs *lexmodelbuildingservice.CodeHook) []map[string]interface{} {
	if cs == nil {
		return nil
	}

	val := make(map[string]interface{})

	val["message_version"] = aws.StringValue(cs.MessageVersion)
	val["uri"] = aws.StringValue(cs.Uri)

	return []map[string]interface{}{val}
}

func flattenPrompt(cs *lexmodelbuildingservice.Prompt) []map[string]interface{} {
	if cs == nil {
		return nil
	}

	val := make(map[string]interface{})

	val["max_attempts"] = aws.Int64Value(cs.MaxAttempts)
	val["messages"] = flattenMessages(cs.Messages)

	if cs.ResponseCard != nil {
		val["response_card"] = aws.StringValue(cs.ResponseCard)
	}

	return []map[string]interface{}{val}
}

func flattenMessages(cs []*lexmodelbuildingservice.Message) []interface{} {
	valuesSlice := make([]interface{}, len(cs))
	if len(cs) > 0 {
		for i, v := range cs {
			m := make(map[string]interface{})
			m["content"] = aws.StringValue(v.Content)
			m["content_type"] = aws.StringValue(v.ContentType)
			if v.GroupNumber != nil {
				m["group_number"] = aws.Int64Value(v.GroupNumber)
			}
			valuesSlice[i] = m
		}
	}

	return valuesSlice
}

func flattenStatement(cs *lexmodelbuildingservice.Statement) []map[string]interface{} {
	if cs == nil {
		return nil
	}

	val := make(map[string]interface{})

	val["messages"] = flattenMessages(cs.Messages)

	if cs.ResponseCard != nil {
		val["response_card"] = aws.StringValue(cs.ResponseCard)
	}

	return []map[string]interface{}{val}
}

func flattenSlots(cs []*lexmodelbuildingservice.Slot) []interface{} {
	valuesSlice := make([]interface{}, len(cs))
	if len(cs) > 0 {
		for i, v := range cs {
			m := make(map[string]interface{})
			m["name"] = aws.StringValue(v.Name)
			m["slot_constraint"] = aws.StringValue(v.SlotConstraint)
			if v.Description != nil {
				m["description"] = aws.StringValue(v.Description)
			}
			if v.ObfuscationSetting != nil {
				m["obfuscation_setting"] = aws.StringValue(v.ObfuscationSetting)
			}
			if v.Priority != nil {
				m["priority"] = aws.Int64Value(v.Priority)
			}
			if v.ResponseCard != nil {
				m["response_card"] = aws.StringValue(v.ResponseCard)
			}
			if v.SampleUtterances != nil {
				m["sample_utterances"] = flattenStringList(v.SampleUtterances)
			}
			if v.SlotType != nil {
				m["slot_type"] = aws.StringValue(v.SlotType)
			}
			if v.SlotTypeVersion != nil {
				m["slot_type_version"] = aws.StringValue(v.SlotTypeVersion)
			}
			if v.ValueElicitationPrompt != nil {
				m["value_elicitation_prompt"] = flattenPrompt(v.ValueElicitationPrompt)
			}
			valuesSlice[i] = m
		}
	}

	return valuesSlice
}

func getLexIntent(name, version string, conn *lexmodelbuildingservice.LexModelBuildingService) (*lexmodelbuildingservice.GetIntentOutput, error) {
	input := &lexmodelbuildingservice.GetIntentInput{
		Name:    aws.String(name),
		Version: aws.String(version),
	}
	return conn.GetIntent(input)
}
