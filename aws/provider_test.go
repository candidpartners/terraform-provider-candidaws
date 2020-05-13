package aws

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviderFactories func(providers *[]*schema.Provider) map[string]terraform.ResourceProviderFactory
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviderFactories = func(providers *[]*schema.Provider) map[string]terraform.ResourceProviderFactory {
		return map[string]terraform.ResourceProviderFactory{
			"aws": func() (terraform.ResourceProvider, error) {
				p := Provider()
				*providers = append(*providers, p.(*schema.Provider))
				return p, nil
			},
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("AWS_PROFILE") == "" && os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Fatal("AWS_ACCESS_KEY_ID or AWS_PROFILE must be set for acceptance tests")
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}

	region := testAccGetRegion()
	log.Printf("[INFO] Test: Using %s as test region", region)
	os.Setenv("AWS_DEFAULT_REGION", region)

	err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}

func testAccGetRegion() string {
	v := os.Getenv("AWS_DEFAULT_REGION")
	if v == "" {
		return "us-west-2"
	}
	return v
}

func TestAccAWSProvider_Endpoints(t *testing.T) {
	var providers []*schema.Provider
	var endpoints strings.Builder

	// Initialize each endpoint configuration with matching name and value
	for _, endpointServiceName := range endpointServiceNames {
		// Skip deprecated endpoint configurations as they will override expected values
		if endpointServiceName == "kinesis_analytics" || endpointServiceName == "r53" {
			continue
		}

		endpoints.WriteString(fmt.Sprintf("%s = \"http://%s\"\n", endpointServiceName, endpointServiceName))
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigEndpoints(endpoints.String()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderEndpoints(&providers),
				),
			},
		},
	})
}

func TestAccAWSProvider_Endpoints_Deprecated(t *testing.T) {
	var providers []*schema.Provider
	var endpointsDeprecated strings.Builder

	// Initialize each deprecated endpoint configuration with matching name and value
	for _, endpointServiceName := range endpointServiceNames {
		// Only configure deprecated endpoint configurations
		if endpointServiceName != "kinesis_analytics" && endpointServiceName != "r53" {
			continue
		}

		endpointsDeprecated.WriteString(fmt.Sprintf("%s = \"http://%s\"\n", endpointServiceName, endpointServiceName))
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigEndpoints(endpointsDeprecated.String()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderEndpointsDeprecated(&providers)),
			},
		},
	})
}

func TestAccAWSProvider_IgnoreTagPrefixes_None(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigIgnoreTagPrefixes0(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderIgnoreTagPrefixes(&providers, []string{}),
				),
			},
		},
	})
}

func TestAccAWSProvider_IgnoreTagPrefixes_One(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigIgnoreTagPrefixes1("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderIgnoreTagPrefixes(&providers, []string{"test"}),
				),
			},
		},
	})
}

func TestAccAWSProvider_IgnoreTagPrefixes_Multiple(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigIgnoreTagPrefixes2("test1", "test2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderIgnoreTagPrefixes(&providers, []string{"test1", "test2"}),
				),
			},
		},
	})
}

func TestAccAWSProvider_IgnoreTags_None(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigIgnoreTags0(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderIgnoreTags(&providers, []string{}),
				),
			},
		},
	})
}

func TestAccAWSProvider_IgnoreTags_One(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigIgnoreTags1("test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderIgnoreTags(&providers, []string{"test"}),
				),
			},
		},
	})
}

func TestAccAWSProvider_IgnoreTags_Multiple(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigIgnoreTags2("test1", "test2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderIgnoreTags(&providers, []string{"test1", "test2"}),
				),
			},
		},
	})
}

func TestAccAWSProvider_Region_AwsChina(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigRegion("cn-northwest-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderDnsSuffix(&providers, "amazonaws.com.cn"),
					testAccCheckAWSProviderPartition(&providers, "aws-cn"),
				),
				PlanOnly: true,
			},
		},
	})
}

func TestAccAWSProvider_Region_AwsCommercial(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigRegion("us-west-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderDnsSuffix(&providers, "amazonaws.com"),
					testAccCheckAWSProviderPartition(&providers, "aws"),
				),
				PlanOnly: true,
			},
		},
	})
}

func TestAccAWSProvider_Region_AwsGovCloudUs(t *testing.T) {
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSProviderConfigRegion("us-gov-west-1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSProviderDnsSuffix(&providers, "amazonaws.com"),
					testAccCheckAWSProviderPartition(&providers, "aws-us-gov"),
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckAWSProviderDnsSuffix(providers *[]*schema.Provider, expectedDnsSuffix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if providers == nil {
			return fmt.Errorf("no providers initialized")
		}

		for _, provider := range *providers {
			if provider == nil || provider.Meta() == nil || provider.Meta().(*AWSClient) == nil {
				continue
			}

			providerDnsSuffix := provider.Meta().(*AWSClient).dnsSuffix

			if providerDnsSuffix != expectedDnsSuffix {
				return fmt.Errorf("expected DNS Suffix (%s), got: %s", expectedDnsSuffix, providerDnsSuffix)
			}
		}

		return nil
	}
}

func testAccCheckAWSProviderEndpoints(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if providers == nil {
			return fmt.Errorf("no providers initialized")
		}

		// Match AWSClient struct field names to endpoint configuration names
		endpointFieldNameF := func(endpoint string) func(string) bool {
			return func(name string) bool {
				switch endpoint {
				case "applicationautoscaling":
					endpoint = "appautoscaling"
				case "budgets":
					endpoint = "budget"
				case "cloudformation":
					endpoint = "cf"
				case "cloudhsm":
					endpoint = "cloudhsmv2"
				case "cognitoidentity":
					endpoint = "cognito"
				case "configservice":
					endpoint = "config"
				case "cur":
					endpoint = "costandusagereport"
				case "directconnect":
					endpoint = "dx"
				case "lexmodels":
					endpoint = "lexmodel"
				case "route53":
					endpoint = "r53"
				case "sdb":
					endpoint = "simpledb"
				case "serverlessrepo":
					endpoint = "serverlessapplicationrepository"
				case "servicecatalog":
					endpoint = "sc"
				case "servicediscovery":
					endpoint = "sd"
				case "stepfunctions":
					endpoint = "sfn"
				}

				return name == fmt.Sprintf("%sconn", endpoint)
			}
		}

		for _, provider := range *providers {
			if provider == nil || provider.Meta() == nil || provider.Meta().(*AWSClient) == nil {
				continue
			}

			providerClient := provider.Meta().(*AWSClient)

			for _, endpointServiceName := range endpointServiceNames {
				// Skip deprecated endpoint configurations as they will override expected values
				if endpointServiceName == "kinesis_analytics" || endpointServiceName == "r53" {
					continue
				}

				providerClientField := reflect.Indirect(reflect.ValueOf(providerClient)).FieldByNameFunc(endpointFieldNameF(endpointServiceName))

				if !providerClientField.IsValid() {
					return fmt.Errorf("unable to match AWSClient struct field name for endpoint name: %s", endpointServiceName)
				}

				actualEndpoint := reflect.Indirect(reflect.Indirect(providerClientField).FieldByName("Config").FieldByName("Endpoint")).String()
				expectedEndpoint := fmt.Sprintf("http://%s", endpointServiceName)

				if actualEndpoint != expectedEndpoint {
					return fmt.Errorf("expected endpoint (%s) value (%s), got: %s", endpointServiceName, expectedEndpoint, actualEndpoint)
				}
			}
		}

		return nil
	}
}

func testAccCheckAWSProviderEndpointsDeprecated(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if providers == nil {
			return fmt.Errorf("no providers initialized")
		}

		// Match AWSClient struct field names to endpoint configuration names
		endpointFieldNameF := func(endpoint string) func(string) bool {
			return func(name string) bool {
				switch endpoint {
				case "kinesis_analytics":
					endpoint = "kinesisanalytics"
				}

				return name == fmt.Sprintf("%sconn", endpoint)
			}
		}

		for _, provider := range *providers {
			if provider == nil || provider.Meta() == nil || provider.Meta().(*AWSClient) == nil {
				continue
			}

			providerClient := provider.Meta().(*AWSClient)

			for _, endpointServiceName := range endpointServiceNames {
				// Only check deprecated endpoint configurations
				if endpointServiceName != "kinesis_analytics" && endpointServiceName != "r53" {
					continue
				}

				providerClientField := reflect.Indirect(reflect.ValueOf(providerClient)).FieldByNameFunc(endpointFieldNameF(endpointServiceName))

				if !providerClientField.IsValid() {
					return fmt.Errorf("unable to match AWSClient struct field name for endpoint name: %s", endpointServiceName)
				}

				actualEndpoint := reflect.Indirect(reflect.Indirect(providerClientField).FieldByName("Config").FieldByName("Endpoint")).String()
				expectedEndpoint := fmt.Sprintf("http://%s", endpointServiceName)

				if actualEndpoint != expectedEndpoint {
					return fmt.Errorf("expected endpoint (%s) value (%s), got: %s", endpointServiceName, expectedEndpoint, actualEndpoint)
				}
			}
		}

		return nil
	}
}

func testAccCheckAWSProviderIgnoreTagPrefixes(providers *[]*schema.Provider, expectedIgnoreTagPrefixes []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if providers == nil {
			return fmt.Errorf("no providers initialized")
		}

		for _, provider := range *providers {
			if provider == nil || provider.Meta() == nil || provider.Meta().(*AWSClient) == nil {
				continue
			}

			providerClient := provider.Meta().(*AWSClient)

			actualIgnoreTagPrefixes := providerClient.ignoreTagPrefixes.Keys()

			if len(actualIgnoreTagPrefixes) != len(expectedIgnoreTagPrefixes) {
				return fmt.Errorf("expected ignore_tag_prefixes (%d) length, got: %d", len(expectedIgnoreTagPrefixes), len(actualIgnoreTagPrefixes))
			}

			for _, expectedElement := range expectedIgnoreTagPrefixes {
				var found bool

				for _, actualElement := range actualIgnoreTagPrefixes {
					if actualElement == expectedElement {
						found = true
						break
					}
				}

				if !found {
					return fmt.Errorf("expected ignore_tag_prefixes element, but was missing: %s", expectedElement)
				}
			}

			for _, actualElement := range actualIgnoreTagPrefixes {
				var found bool

				for _, expectedElement := range expectedIgnoreTagPrefixes {
					if actualElement == expectedElement {
						found = true
						break
					}
				}

				if !found {
					return fmt.Errorf("unexpected ignore_tag_prefixes element: %s", actualElement)
				}
			}
		}

		return nil
	}
}

func testAccCheckAWSProviderIgnoreTags(providers *[]*schema.Provider, expectedIgnoreTags []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if providers == nil {
			return fmt.Errorf("no providers initialized")
		}

		for _, provider := range *providers {
			if provider == nil || provider.Meta() == nil || provider.Meta().(*AWSClient) == nil {
				continue
			}

			providerClient := provider.Meta().(*AWSClient)

			actualIgnoreTags := providerClient.ignoreTags.Keys()

			if len(actualIgnoreTags) != len(expectedIgnoreTags) {
				return fmt.Errorf("expected ignore_tags (%d) length, got: %d", len(expectedIgnoreTags), len(actualIgnoreTags))
			}

			for _, expectedElement := range expectedIgnoreTags {
				var found bool

				for _, actualElement := range actualIgnoreTags {
					if actualElement == expectedElement {
						found = true
						break
					}
				}

				if !found {
					return fmt.Errorf("expected ignore_tags element, but was missing: %s", expectedElement)
				}
			}

			for _, actualElement := range actualIgnoreTags {
				var found bool

				for _, expectedElement := range expectedIgnoreTags {
					if actualElement == expectedElement {
						found = true
						break
					}
				}

				if !found {
					return fmt.Errorf("unexpected ignore_tags element: %s", actualElement)
				}
			}
		}

		return nil
	}
}

func testAccCheckAWSProviderPartition(providers *[]*schema.Provider, expectedPartition string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if providers == nil {
			return fmt.Errorf("no providers initialized")
		}

		for _, provider := range *providers {
			if provider == nil || provider.Meta() == nil || provider.Meta().(*AWSClient) == nil {
				continue
			}

			providerPartition := provider.Meta().(*AWSClient).partition

			if providerPartition != expectedPartition {
				return fmt.Errorf("expected DNS Suffix (%s), got: %s", expectedPartition, providerPartition)
			}
		}

		return nil
	}
}

func testAccAWSProviderConfigEndpoints(endpoints string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    %[1]s
  }
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`, endpoints)
}

func testAccAWSProviderConfigIgnoreTagPrefixes0() string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`)
}

func testAccAWSProviderConfigIgnoreTagPrefixes1(tagPrefix1 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  ignore_tag_prefixes         = [%[1]q]
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`, tagPrefix1)
}

func testAccAWSProviderConfigIgnoreTagPrefixes2(tagPrefix1, tagPrefix2 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  ignore_tag_prefixes         = [%[1]q, %[2]q]
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`, tagPrefix1, tagPrefix2)
}

func testAccAWSProviderConfigIgnoreTags0() string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`)
}

func testAccAWSProviderConfigIgnoreTags1(tag1 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  ignore_tags                 = [%[1]q]
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`, tag1)
}

func testAccAWSProviderConfigIgnoreTags2(tag1, tag2 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  ignore_tags                 = [%[1]q, %[2]q]
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`, tag1, tag2)
}

func testAccAWSProviderConfigRegion(region string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  region                      = %[1]q
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:aws:s3:::test"
}
`, region)
}
