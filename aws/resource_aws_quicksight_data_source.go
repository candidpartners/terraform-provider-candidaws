package aws

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/quicksight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAwsQuickSightDataSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsQuickSightDataSourceCreate,
		Read:   resourceAwsQuickSightDataSourceRead,
		Update: resourceAwsQuickSightDataSourceUpdate,
		Delete: resourceAwsQuickSightDataSourceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"aws_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"principal": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_source_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_parameters": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mysql_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"athena_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workgroup": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"amazon_elasticsearch_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"aurora_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"aurora_postgre_sql_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"aws_iot_analytics_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_set_name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"jira_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"site_base_url": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"maria_db_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"postgre_sql_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"presto_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"catalog": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"rds_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"redshift_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"s3_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mainifest_file_location": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:     schema.TypeString,
													Required: true,
												},
												"key": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"service_now_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"site_base_url": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"snowflake_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"warehouse": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"spark_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"sql_server_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.teradata_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"teradata_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.twitter_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:     schema.TypeString,
										Required: true,
									},
									"host": {
										Type:     schema.TypeString,
										Required: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"twitter_parameters": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"data_source_parameters.0.mysql_parameters", "data_source_parameters.0.athena_parameters", "data_source_parameters.0.amazon_elasticsearch_parameters", "data_source_parameters.0.aurora_parameters", "data_source_parameters.0.aurora_postgre_sql_parameters", "data_source_parameters.0.aws_iot_analytics_parameters", "data_source_parameters.0.jira_parameters", "data_source_parameters.0.maria_db_parameters", "data_source_parameters.0.postgre_sql_parameters", "data_source_parameters.0.presto_parameters", "data_source_parameters.0.rds_parameters", "data_source_parameters.0.redshift_parameters", "data_source_parameters.0.s3_parameters", "data_source_parameters.0.service_now_parameters", "data_source_parameters.0.snowflake_parameters", "data_source_parameters.0.spark_parameters", "data_source_parameters.0.sql_server_parameters", "data_source_parameters.0.teradata_parameters"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"query": {
										Type:     schema.TypeString,
										Required: true,
									},
									"max_rows": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			// add update support
			"data_source_type": {
				Type:    schema.TypeString,
				Default: quicksight.DataSourceTypeAthena,
				ValidateFunc: validation.StringInSlice([]string{
					quicksight.DataSourceTypeAthena,
					quicksight.DataSourceTypeAurora,
					quicksight.DataSourceTypeAuroraPostgresql,
					quicksight.DataSourceTypeMariadb,
					quicksight.DataSourceTypeMysql,
					quicksight.DataSourceTypePostgresql,
					quicksight.DataSourceTypePresto,
					quicksight.DataSourceTypeRedshift,
					quicksight.DataSourceTypeS3,
					quicksight.DataSourceTypeServicenow,
					quicksight.DataSourceTypeSnowflake,
					quicksight.DataSourceTypeSpark,
					quicksight.DataSourceTypeSqlserver,
					quicksight.DataSourceTypeTeradata,
					quicksight.DataSourceTypeTwitter,
					quicksight.DataSourceTypeJira,
					quicksight.DataSourceTypeAwsIotAnalytics,
				}, false),
				Optional: true,
				ForceNew: true,
			},
		},
	}
}
func expandPermissions(values []interface{}) []*quicksight.ResourcePermission {
	valueSlice := []*quicksight.ResourcePermission{}
	for _, element := range values {
		e := element.(map[string]interface{})
		m := &quicksight.ResourcePermission{
			Principal: aws.String(e["principal"].(string)),
		}
		if a, ok := e["actions"]; ok {
			m.Actions = expandStringList(a.([]interface{}))
		}
		valueSlice = append(valueSlice, m)
	}
	return valueSlice
}
func expandMysqlParameters(values []interface{}) *quicksight.MySqlParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.MySqlParameters{
		Database: aws.String(string(val["database"].(string))),
		Host:     aws.String(string(val["host"].(string))),
		Port:     aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandAmazonElasticsearchParameters(values []interface{}) *quicksight.AmazonElasticsearchParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.AmazonElasticsearchParameters{
		Domain: aws.String(string(val["domain"].(string))),
	}
	return res
}
func expandAthenaParameters(values []interface{}) *quicksight.AthenaParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.AthenaParameters{
		WorkGroup: aws.String(string(val["workgroup"].(string))),
	}
	return res
}
func expandAuroraParameters(values []interface{}) *quicksight.AuroraParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.AuroraParameters{
		Database: aws.String(string(val["database"].(string))),
		Host:     aws.String(string(val["host"].(string))),
		Port:     aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandAuroraPostgreSQLParameters(values []interface{}) *quicksight.AuroraPostgreSqlParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.AuroraPostgreSqlParameters{
		Database: aws.String(string(val["database"].(string))),
		Host:     aws.String(string(val["host"].(string))),
		Port:     aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandAwsIotAnalyticsParameters(values []interface{}) *quicksight.AwsIotAnalyticsParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.AwsIotAnalyticsParameters{
		DataSetName: aws.String(string(val["data_set_name"].(string))),
	}
	return res
}
func expandJiraParameters(values []interface{}) *quicksight.JiraParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.JiraParameters{
		SiteBaseUrl: aws.String(string(val["site_base_url"].(string))),
	}
	return res
}
func expandMariaDbParameters(values []interface{}) *quicksight.MariaDbParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.MariaDbParameters{
		Database: aws.String(string(val["database"].(string))),
		Host:     aws.String(string(val["host"].(string))),
		Port:     aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandPostgreSQLParameters(values []interface{}) *quicksight.PostgreSqlParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.PostgreSqlParameters{
		Database: aws.String(string(val["database"].(string))),
		Host:     aws.String(string(val["host"].(string))),
		Port:     aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandPrestoParameters(values []interface{}) *quicksight.PrestoParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.PrestoParameters{
		Catalog: aws.String(string(val["catalog"].(string))),
		Host:    aws.String(string(val["host"].(string))),
		Port:    aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandRdsParameters(values []interface{}) *quicksight.RdsParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.RdsParameters{
		Database:   aws.String(string(val["database"].(string))),
		InstanceId: aws.String(string(val["instance_id"].(string))),
	}
	return res
}
func expandRedshiftParameters(values []interface{}) *quicksight.RedshiftParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.RedshiftParameters{
		ClusterId: aws.String(string(val["cluster_id"].(string))),
		Database:  aws.String(string(val["database"].(string))),
		Host:      aws.String(string(val["host"].(string))),
		Port:      aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandManifestFileLocation(values []interface{}) *quicksight.ManifestFileLocation {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.ManifestFileLocation{
		Bucket: aws.String(string(val["bucket"].(string))),
		Key:    aws.String(string(val["key"].(string))),
	}
	return res
}
func expandS3Parameters(values []interface{}) *quicksight.S3Parameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.S3Parameters{
		ManifestFileLocation: expandManifestFileLocation(val["mainifest_file_location"].([]interface{})),
	}
	return res
}
func expandServiceNowParameters(values []interface{}) *quicksight.ServiceNowParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.ServiceNowParameters{
		SiteBaseUrl: aws.String(string(val["site_base_url"].(string))),
	}
	return res
}
func expandSnowflakeParameters(values []interface{}) *quicksight.SnowflakeParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.SnowflakeParameters{
		Database:  aws.String(string(val["database"].(string))),
		Host:      aws.String(string(val["host"].(string))),
		Warehouse: aws.String(string(val["warehouse"].(string))),
	}
	return res
}
func expandSparkParameters(values []interface{}) *quicksight.SparkParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.SparkParameters{
		Host: aws.String(string(val["host"].(string))),
		Port: aws.Int64(int64(val["warehouse"].(int))),
	}
	return res
}
func expandSQLServerParameters(values []interface{}) *quicksight.SqlServerParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.SqlServerParameters{
		Database: aws.String(string(val["database"].(string))),
		Host:     aws.String(string(val["host"].(string))),
		Port:     aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandTeradataParameters(values []interface{}) *quicksight.TeradataParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.TeradataParameters{
		Database: aws.String(string(val["database"].(string))),
		Host:     aws.String(string(val["host"].(string))),
		Port:     aws.Int64(int64(val["port"].(int))),
	}
	return res
}
func expandTwitterParameters(values []interface{}) *quicksight.TwitterParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.TwitterParameters{
		Query:   aws.String(string(val["query"].(string))),
		MaxRows: aws.Int64(int64(val["max_rows"].(int))),
	}
	return res
}

func expandDataSourceParameters(values []interface{}) *quicksight.DataSourceParameters {
	if len(values) == 0 {
		return nil
	}

	val := values[0].(map[string]interface{})
	res := &quicksight.DataSourceParameters{
		AthenaParameters:              expandAthenaParameters(val["athena_parameters"].([]interface{})),
		MySqlParameters:               expandMysqlParameters(val["mysql_parameters"].([]interface{})),
		AmazonElasticsearchParameters: expandAmazonElasticsearchParameters(val["amazon_elasticsearch_parameters"].([]interface{})),
		AuroraParameters:              expandAuroraParameters(val["aurora_parameters"].([]interface{})),
		AuroraPostgreSqlParameters:    expandAuroraPostgreSQLParameters(val["aurora_postgre_sql_parameters"].([]interface{})),
		AwsIotAnalyticsParameters:     expandAwsIotAnalyticsParameters(val["aws_iot_analytics_parameters"].([]interface{})),
		JiraParameters:                expandJiraParameters(val["jira_parameters"].([]interface{})),
		MariaDbParameters:             expandMariaDbParameters(val["maria_db_parameters"].([]interface{})),
		PostgreSqlParameters:          expandPostgreSQLParameters(val["postgre_sql_parameters"].([]interface{})),
		PrestoParameters:              expandPrestoParameters(val["presto_parameters"].([]interface{})),
		RdsParameters:                 expandRdsParameters(val["rds_parameters"].([]interface{})),
		RedshiftParameters:            expandRedshiftParameters(val["redshift_parameters"].([]interface{})),
		S3Parameters:                  expandS3Parameters(val["s3_parameters"].([]interface{})),
		ServiceNowParameters:          expandServiceNowParameters(val["service_now_parameters"].([]interface{})),
		SnowflakeParameters:           expandSnowflakeParameters(val["snowflake_parameters"].([]interface{})),
		SparkParameters:               expandSparkParameters(val["spark_parameters"].([]interface{})),
		SqlServerParameters:           expandSQLServerParameters(val["sql_server_parameters"].([]interface{})),
		TeradataParameters:            expandTeradataParameters(val["teradata_parameters"].([]interface{})),
		TwitterParameters:             expandTwitterParameters(val["twitter_parameters"].([]interface{})),
	}
	return res
}
func resourceAwsQuickSightDataSourceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID := meta.(*AWSClient).accountid
	dataSourceID := d.Get("data_source_id").(string)
	dataSourceName := d.Get("data_source_name").(string)
	dataSourceType := d.Get("data_source_type").(string)
	if v, ok := d.GetOk("aws_account_id"); ok {
		awsAccountID = v.(string)
	}

	createOpts := &quicksight.CreateDataSourceInput{
		AwsAccountId:         aws.String(awsAccountID),
		DataSourceParameters: expandDataSourceParameters(d.Get("data_source_parameters").([]interface{})),
		Permissions:          expandPermissions(d.Get("permissions").([]interface{})),
		DataSourceId:         aws.String(dataSourceID),
		Name:                 aws.String(dataSourceName),
		Type:                 aws.String(dataSourceType),
	}

	resp, err := conn.CreateDataSource(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating QuickSight DataSource: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", awsAccountID, aws.StringValue(resp.DataSourceId)))

	return resourceAwsQuickSightDataSourceRead(d, meta)
}

func resourceAwsQuickSightDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(d.Id())
	if err != nil {
		return err
	}

	descOpts := &quicksight.DescribeDataSourceInput{
		AwsAccountId: aws.String(awsAccountID),
		DataSourceId: aws.String(dataSourceID),
	}

	resp, err := conn.DescribeDataSource(descOpts)
	if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] QuickSight DataSource %s is already gone", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error describing QuickSight DataSource (%s): %s", d.Id(), err)
	}

	d.Set("arn", resp.DataSource.Arn)
	d.Set("aws_account_id", awsAccountID)
	d.Set("data_source_id", resp.DataSource.DataSourceId)
	d.Set("data_source_name", resp.DataSource.Name)
	d.Set("data_source_type", resp.DataSource.Type)

	return nil
}

func resourceAwsQuickSightDataSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(d.Id())
	if err != nil {
		return err
	}
	updateOpts := &quicksight.UpdateDataSourceInput{
		AwsAccountId:         aws.String(awsAccountID),
		DataSourceId:         aws.String(dataSourceID),
		DataSourceParameters: expandDataSourceParameters(d.Get("data_source_parameters").([]interface{})),
	}
	if v, ok := d.GetOk("data_source_name"); ok {
		updateOpts.Name = aws.String(v.(string))
	}

	_, err = conn.UpdateDataSource(updateOpts)
	if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] QuickSight DataSource %s is already gone", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error updating QuickSight DataSource %s: %s", d.Id(), err)
	}
	return resourceAwsQuickSightDataSourceRead(d, meta)
}

func resourceAwsQuickSightDataSourceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).quicksightconn

	awsAccountID, dataSourceID, err := resourceAwsQuickSightDataSourceParseID(d.Id())
	if err != nil {
		return err
	}

	deleteOpts := &quicksight.DeleteDataSourceInput{
		AwsAccountId: aws.String(awsAccountID),
		DataSourceId: aws.String(dataSourceID),
	}

	if _, err := conn.DeleteDataSource(deleteOpts); err != nil {
		if isAWSErr(err, quicksight.ErrCodeResourceNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("Error deleting QuickSight DataSource %s: %s", d.Id(), err)
	}

	return nil
}

func resourceAwsQuickSightDataSourceParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected AWS_ACCOUNT_ID/DATA_SOURCE_ID", id)
	}
	return parts[0], parts[1], nil
}
