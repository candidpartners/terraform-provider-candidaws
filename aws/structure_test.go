package aws

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestExpandStringList(t *testing.T) {
	expanded := []interface{}{"us-east-1a", "us-east-1b"}
	stringList := expandStringList(expanded)
	expected := []*string{
		aws.String("us-east-1a"),
		aws.String("us-east-1b"),
	}

	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func TestExpandStringListEmptyItems(t *testing.T) {
	expanded := []interface{}{"foo", "bar", "", "baz"}
	stringList := expandStringList(expanded)
	expected := []*string{
		aws.String("foo"),
		aws.String("bar"),
		aws.String("baz"),
	}

	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}

func TestCheckYamlString(t *testing.T) {
	var err error
	var actual string

	validYaml := `---
abc:
  def: 123
  xyz:
    -
      a: "ホリネズミ"
      b: "1"
`

	actual, err = checkYamlString(validYaml)
	if err != nil {
		t.Fatalf("Expected not to throw an error while parsing YAML, but got: %s", err)
	}

	// We expect the same YAML string back
	if actual != validYaml {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, validYaml)
	}

	invalidYaml := `abc: [`

	actual, err = checkYamlString(invalidYaml)
	if err == nil {
		t.Fatalf("Expected to throw an error while parsing YAML, but got: %s", err)
	}

	// We expect the invalid YAML to be shown back to us again.
	if actual != invalidYaml {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, invalidYaml)
	}
}

func TestNormalizeCloudFormationTemplate(t *testing.T) {
	var err error
	var actual string

	validNormalizedJson := `{"abc":"1"}`
	actual, err = normalizeCloudFormationTemplate(validNormalizedJson)
	if err != nil {
		t.Fatalf("Expected not to throw an error while parsing template, but got: %s", err)
	}
	if actual != validNormalizedJson {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, validNormalizedJson)
	}

	validNormalizedYaml := `abc: 1
`
	actual, err = normalizeCloudFormationTemplate(validNormalizedYaml)
	if err != nil {
		t.Fatalf("Expected not to throw an error while parsing template, but got: %s", err)
	}
	if actual != validNormalizedYaml {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, validNormalizedYaml)
	}
}
