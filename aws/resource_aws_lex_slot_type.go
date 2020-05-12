package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var validValueSelectionStrategies = []string{
	lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue,
	lexmodelbuildingservice.SlotValueSelectionStrategyTopResolution,
}

func resourceAwsLexSlotType() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsLexSlotTypeCreate,
		Read:   resourceAwsLexSlotTypeRead,
		Update: resourceAwsLexSlotTypeUpdate,
		Delete: resourceAwsLexSlotTypeDelete,
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
			"publish": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"value_selection_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue,
				ValidateFunc: validation.StringInSlice(validValueSelectionStrategies, false),
			},
			"enumeration_values": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"synonyms": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAwsLexSlotTypeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	params := &lexmodelbuildingservice.PutSlotTypeInput{
		Name:                   aws.String(d.Get("name").(string)),
		Description:            aws.String(d.Get("description").(string)),
		EnumerationValues:      expandEnumerationValues(d.Get("enumeration_values").([]interface{})),
		CreateVersion:          aws.Bool(d.Get("publish").(bool)),
		ValueSelectionStrategy: aws.String(d.Get("value_selection_strategy").(string)),
	}
	resp, err := conn.PutSlotType(params)
	if err != nil {
		return fmt.Errorf("error putting Lex slot type: %s", err)
	}

	d.SetId(aws.StringValue(resp.Name))
	d.Set("version", resp.Version)
	d.Set("checksum", resp.Checksum)
	return resourceAwsLexSlotTypeRead(d, meta)
}

func resourceAwsLexSlotTypeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	version := d.Get("version").(string)
	slotType, err := getLexSlotType(d.Id(), version, conn)
	if err != nil {
		return fmt.Errorf("error getting Lex slot type %q: %s", d.Id(), err)
	}
	if slotType == nil {
		log.Printf("[WARN] LexModelBuildingService Slot Type %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", slotType.Description)
	d.Set("checksum", slotType.Checksum)
	d.Set("version", slotType.Version)

	if err := d.Set("enumeration_values", flattenEnumerationValues(slotType.EnumerationValues)); err != nil {
		return fmt.Errorf("error setting enumeration_values: %s", err)
	}

	return nil
}

func resourceAwsLexSlotTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	name := d.Id()

	params := &lexmodelbuildingservice.PutSlotTypeInput{
		Name:                   aws.String(name),
		Checksum:               aws.String(d.Get("checksum").(string)),
		Description:            aws.String(d.Get("description").(string)),
		EnumerationValues:      expandEnumerationValues(d.Get("enumeration_values").([]interface{})),
		CreateVersion:          aws.Bool(d.Get("publish").(bool)),
		ValueSelectionStrategy: aws.String(d.Get("value_selection_strategy").(string)),
	}

	resp, err := conn.PutSlotType(params)
	if err != nil {
		return err
	}

	d.Set("version", resp.Version)
	d.Set("checksum", resp.Checksum)
	return resourceAwsLexSlotTypeRead(d, meta)
}

func resourceAwsLexSlotTypeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelbuildingserviceconn

	name := d.Id()

	input := &lexmodelbuildingservice.DeleteSlotTypeInput{
		Name: aws.String(name),
	}
	_, err := conn.DeleteSlotType(input)
	if err != nil {
		if isAWSErr(err, lexmodelbuildingservice.ErrCodeNotFoundException, "") {
			return nil
		}
		return err
	}

	return nil
}

func expandEnumerationValues(values []interface{}) []*lexmodelbuildingservice.EnumerationValue {
	valueSlice := []*lexmodelbuildingservice.EnumerationValue{}
	for _, element := range values {
		elementMap := element.(map[string]interface{})

		value := &lexmodelbuildingservice.EnumerationValue{
			Value: aws.String(elementMap["value"].(string)),
		}

		if synonyms, ok := elementMap["synonyms"]; ok {
			value.Synonyms = expandSynonyms(synonyms.([]interface{}))
		}

		valueSlice = append(valueSlice, value)
	}

	return valueSlice
}

func expandSynonyms(values []interface{}) []*string {
	valueSlice := []*string{}
	for _, element := range values {
		e := element.(string)
		valueSlice = append(valueSlice, &e)
	}

	return valueSlice
}

func flattenEnumerationValues(cs []*lexmodelbuildingservice.EnumerationValue) []map[string]interface{} {
	valuesSlice := make([]map[string]interface{}, len(cs))
	if len(cs) > 0 {
		for i, v := range cs {
			valuesSlice[i] = flattenEnumerationValue(v)
		}
	}

	return valuesSlice
}

func flattenEnumerationValue(c *lexmodelbuildingservice.EnumerationValue) map[string]interface{} {
	column := make(map[string]interface{})

	if c == nil {
		return column
	}

	if v := aws.StringValue(c.Value); v != "" {
		column["value"] = v
	}

	if v := aws.StringValueSlice(c.Synonyms); len(v) > 0 {
		column["synonyms"] = v
	}

	return column
}

func getLexSlotType(name, version string, conn *lexmodelbuildingservice.LexModelBuildingService) (*lexmodelbuildingservice.GetSlotTypeOutput, error) {
	input := &lexmodelbuildingservice.GetSlotTypeInput{
		Name:    aws.String(name),
		Version: aws.String(version),
	}
	return conn.GetSlotType(input)
}
