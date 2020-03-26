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
	intent, err := getLexIntent(d.Id(), version, conn)
	if err != nil {
		return fmt.Errorf("error getting Lex slot type %q: %s", d.Id(), err)
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
		MessageVersion:     aws.String(val["message_version"].(string)),
		Uri:     aws.String(val["uri"].(string)),
	}
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

func getLexIntent(name, version string, conn *lexmodelbuildingservice.LexModelBuildingService) (*lexmodelbuildingservice.GetIntentOutput, error) {
	input := &lexmodelbuildingservice.GetIntentInput{
		Name:    aws.String(name),
		Version: aws.String(version),
	}
	return conn.GetIntent(input)
}
