package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/quicksight"
)

func TestAccAWSQuickSightDataSource_disappears(t *testing.T) {
	var dataSource quicksight.DataSource
	resourceName := "aws_quicksight_datasource.default"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckQuickSightGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSQuickSightDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQuickSightDataSourceExists(resourceName, &dataSource),
					testAccCheckQuickSightDataSourceDisappears(&dataSource),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckQuickSightDataSourceExists(resourceName string, dataSource *quicksight.DataSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := testAccProvider.Meta().(*AWSClient).quicksightconn

		input := &quicksight.DescribeDataSourceInput{
			AwsAccountId: aws.String(awsAccountID),
			DataSourceId: aws.String(dataSourceID),
		}

		output, err := conn.DescribeDataSource(input)

		if err != nil {
			return err
		}

		if output == nil || output.DataSource == nil {
			return fmt.Errorf("QuickSight DataSource (%s) not found", rs.Primary.ID)
		}

		*dataSource = *output.DataSource

		return nil
	}
}

func testAccCheckQuickSightDataSourceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).quicksightconn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_quicksight_datasource" {
			continue
		}

		awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = conn.DescribeDataSource(&quicksight.DescribeDataSourceInput{
			AwsAccountId: aws.String(awsAccountID),
			DataSourceId: aws.String(dataSourceID),
		})
		if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("QuickSight DataSource '%s' was not deleted properly", rs.Primary.ID)
	}

	return nil
}

func testAccCheckQuickSightDataSourceDisappears(v *quicksight.DataSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).quicksightconn

		arn, err := arn.Parse(aws.StringValue(v.Arn))
		if err != nil {
			return err
		}

		dataSourceID, err := arn.Parse(aws.StringValue(v.DataSourceId))
		if err != nil {
			return err
		}

		parts := strings.SplitN(arn.Resource, "/", 3)

		input := &quicksight.DeleteDataSourceInput{
			AwsAccountId: aws.String(arn.AccountID),
			DataSourceId: aws.String(dataSourceID),
		}

		if _, err := conn.DeleteDataSource(input); err != nil {
			return err
		}

		return nil
	}
}

func testAccAWSQuickSightDataSourceConfig(rName, rType, rID, rUserName, rPassword, rAwsAccountID string) string {
	return fmt.Sprintf(`
						resource "aws_quicksight_datasource" "default" {
							data_source_id = %[1]q
							data_source_name = %[2]q
							data_source_type = %[3]q
							user_name = %[4]q
							password = %[5]q
							aws_account_id = %[6]q
						}`, rName, rType, rID, rUserName, rPassword, rAwsAccountId)
}
