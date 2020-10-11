package rdsdataservice

import (
	"log"

	"github.com/awsiv/terraform-provider-rdsdataservice/rdsdataservice/internal/keyvaluetags"
	"github.com/awsiv/terraform-provider-rdsdataservice/rdsdataservice/internal/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	// TODO: Move the validation to this, requires conditional schemas
	// TODO: Move the configuration to this, requires validation

	// The actual provider
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["secret_key"],
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["profile"],
			},

			"assume_role": assumeRoleSchema(),

			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["shared_credentials_file"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["token"],
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				Description:  descriptions["region"],
				InputDefault: "us-east-1",
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     25,
				Description: descriptions["max_retries"],
			},

			"allowed_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"forbidden_account_ids"},
				Set:           schema.HashString,
			},

			"forbidden_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"allowed_account_ids"},
				Set:           schema.HashString,
			},

			"endpoints": endpointsSchema(),

			"ignore_tag_prefixes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Resource tag key prefixes to ignore across all resources.",
			},

			"ignore_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Resource tag keys to ignore across all resources.",
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},

			"skip_credentials_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_credentials_validation"],
			},

			"skip_get_ec2_platforms": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_get_ec2_platforms"],
			},

			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_region_validation"],
			},

			"skip_requesting_account_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_requesting_account_id"],
			},

			"skip_metadata_api_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_metadata_api_check"],
			},

			"s3_force_path_style": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["s3_force_path_style"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rdsdataservice_postgres_database": resourceAwsRdsdataservicePostgresDatabase(),
			"rdsdataservice_postgres_role":     resourceAwsRdsdataservicePostgresRole(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string
var endpointServiceNames []string

func init() {
	descriptions = map[string]string{
		"region": "The region where AWS operations will take place. Examples\n" +
			"are us-east-1, us-west-2, etc.",

		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",

		"profile": "The profile for API operations. If not set, the default profile\n" +
			"created with `aws configure` will be used.",

		"shared_credentials_file": "The path to the shared credentials file. If not set\n" +
			"this defaults to ~/.aws/credentials.",

		"token": "session token. A session token is only required if you are\n" +
			"using temporary security credentials.",

		"max_retries": "The maximum number of times an AWS API request is\n" +
			"being executed. If the API request still fails, an error is\n" +
			"thrown.",

		"endpoint": "Use this to override the default service endpoint URL",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",

		"skip_credentials_validation": "Skip the credentials validation via STS API. " +
			"Used for AWS API implementations that do not have STS available/implemented.",

		"skip_get_ec2_platforms": "Skip getting the supported EC2 platforms. " +
			"Used by users that don't have ec2:DescribeAccountAttributes permissions.",

		"skip_region_validation": "Skip static validation of region name. " +
			"Used by users of alternative AWS-like APIs or users w/ access to regions that are not public (yet).",

		"skip_requesting_account_id": "Skip requesting the account ID. " +
			"Used for AWS API implementations that do not have IAM/STS API and/or metadata API.",

		"skip_medatadata_api_check": "Skip the AWS Metadata API check. " +
			"Used for AWS API implementations that do not have a metadata api endpoint.",

		"s3_force_path_style": "Set this to true to force the request to use path-style addressing,\n" +
			"i.e., http://s3.amazonaws.com/BUCKET/KEY. By default, the S3 client will\n" +
			"use virtual hosted bucket addressing when possible\n" +
			"(http://BUCKET.s3.amazonaws.com/KEY). Specific to the Amazon S3 service.",
	}

	endpointServiceNames = []string{
		"acm",
		"acmpca",
		"amplify",
		"apigateway",
		"applicationautoscaling",
		"applicationinsights",
		"appmesh",
		"appstream",
		"appsync",
		"athena",
		"autoscaling",
		"autoscalingplans",
		"backup",
		"batch",
		"budgets",
		"cloud9",
		"cloudformation",
		"cloudfront",
		"cloudhsm",
		"cloudsearch",
		"cloudtrail",
		"cloudwatch",
		"cloudwatchevents",
		"cloudwatchlogs",
		"codebuild",
		"codecommit",
		"codedeploy",
		"codepipeline",
		"cognitoidentity",
		"cognitoidp",
		"configservice",
		"cur",
		"datapipeline",
		"datasync",
		"dax",
		"devicefarm",
		"directconnect",
		"dlm",
		"dms",
		"docdb",
		"ds",
		"dynamodb",
		"ec2",
		"ecr",
		"ecs",
		"efs",
		"eks",
		"elasticache",
		"elasticbeanstalk",
		"elastictranscoder",
		"elb",
		"emr",
		"es",
		"firehose",
		"fms",
		"forecast",
		"fsx",
		"gamelift",
		"glacier",
		"globalaccelerator",
		"glue",
		"greengrass",
		"guardduty",
		"iam",
		"inspector",
		"iot",
		"iotanalytics",
		"iotevents",
		"kafka",
		"kinesis_analytics",
		"kinesis",
		"kinesisanalytics",
		"kinesisvideo",
		"kms",
		"lakeformation",
		"lambda",
		"lexmodels",
		"licensemanager",
		"lightsail",
		"macie",
		"managedblockchain",
		"mediaconnect",
		"mediaconvert",
		"medialive",
		"mediapackage",
		"mediastore",
		"mediastoredata",
		"mq",
		"neptune",
		"opsworks",
		"organizations",
		"personalize",
		"pinpoint",
		"pricing",
		"qldb",
		"quicksight",
		"r53",
		"ram",
		"rds",
		"redshift",
		"resourcegroups",
		"route53",
		"route53resolver",
		"s3",
		"s3control",
		"sagemaker",
		"sdb",
		"secretsmanager",
		"securityhub",
		"serverlessrepo",
		"servicecatalog",
		"servicediscovery",
		"servicequotas",
		"ses",
		"shield",
		"sns",
		"sqs",
		"ssm",
		"stepfunctions",
		"storagegateway",
		"sts",
		"swf",
		"transfer",
		"waf",
		"wafregional",
		"worklink",
		"workspaces",
		"xray",
	}
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		AccessKey:               d.Get("access_key").(string),
		SecretKey:               d.Get("secret_key").(string),
		Profile:                 d.Get("profile").(string),
		Token:                   d.Get("token").(string),
		Region:                  d.Get("region").(string),
		CredsFilename:           d.Get("shared_credentials_file").(string),
		Endpoints:               make(map[string]string),
		MaxRetries:              d.Get("max_retries").(int),
		IgnoreTagsConfig:        expandProviderIgnoreTags(d.Get("ignore_tags").([]interface{})),
		Insecure:                d.Get("insecure").(bool),
		SkipCredsValidation:     d.Get("skip_credentials_validation").(bool),
		SkipGetEC2Platforms:     d.Get("skip_get_ec2_platforms").(bool),
		SkipRegionValidation:    d.Get("skip_region_validation").(bool),
		SkipRequestingAccountId: d.Get("skip_requesting_account_id").(bool),
		SkipMetadataApiCheck:    d.Get("skip_metadata_api_check").(bool),
		S3ForcePathStyle:        d.Get("s3_force_path_style").(bool),
		terraformVersion:        terraformVersion,
	}

	if l, ok := d.Get("assume_role").([]interface{}); ok && len(l) > 0 && l[0] != nil {
		m := l[0].(map[string]interface{})

		if v, ok := m["duration_seconds"].(int); ok && v != 0 {
			config.AssumeRoleDurationSeconds = v
		}

		if v, ok := m["external_id"].(string); ok && v != "" {
			config.AssumeRoleExternalID = v
		}

		if v, ok := m["policy"].(string); ok && v != "" {
			config.AssumeRolePolicy = v
		}

		if policyARNSet, ok := m["policy_arns"].(*schema.Set); ok && policyARNSet.Len() > 0 {
			for _, policyARNRaw := range policyARNSet.List() {
				policyARN, ok := policyARNRaw.(string)

				if !ok {
					continue
				}

				config.AssumeRolePolicyARNs = append(config.AssumeRolePolicyARNs, policyARN)
			}
		}

		if v, ok := m["role_arn"].(string); ok && v != "" {
			config.AssumeRoleARN = v
		}

		if v, ok := m["session_name"].(string); ok && v != "" {
			config.AssumeRoleSessionName = v
		}

		if tagMapRaw, ok := m["tags"].(map[string]interface{}); ok && len(tagMapRaw) > 0 {
			config.AssumeRoleTags = make(map[string]string)

			for k, vRaw := range tagMapRaw {
				v, ok := vRaw.(string)

				if !ok {
					continue
				}

				config.AssumeRoleTags[k] = v
			}
		}

		if transitiveTagKeySet, ok := m["transitive_tag_keys"].(*schema.Set); ok && transitiveTagKeySet.Len() > 0 {
			for _, transitiveTagKeyRaw := range transitiveTagKeySet.List() {
				transitiveTagKey, ok := transitiveTagKeyRaw.(string)

				if !ok {
					continue
				}

				config.AssumeRoleTransitiveTagKeys = append(config.AssumeRoleTransitiveTagKeys, transitiveTagKey)
			}
		}

		log.Printf("[INFO] assume_role configuration set: (ARN: %q, SessionID: %q, ExternalID: %q)", config.AssumeRoleARN, config.AssumeRoleSessionName, config.AssumeRoleExternalID)
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		for _, endpointServiceName := range endpointServiceNames {
			config.Endpoints[endpointServiceName] = endpoints[endpointServiceName].(string)
		}
	}

	if v, ok := d.GetOk("allowed_account_ids"); ok {
		for _, accountIDRaw := range v.(*schema.Set).List() {
			config.AllowedAccountIds = append(config.AllowedAccountIds, accountIDRaw.(string))
		}
	}

	if v, ok := d.GetOk("forbidden_account_ids"); ok {
		for _, accountIDRaw := range v.(*schema.Set).List() {
			config.ForbiddenAccountIds = append(config.ForbiddenAccountIds, accountIDRaw.(string))
		}
	}

	return config.Client()
}

// This is a global MutexKV for use within this plugin.
var awsMutexKV = mutexkv.NewMutexKV()

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_role_arn"],
				},

				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
				},

				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_external_id"],
				},

				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
			},
		},
	}
}

func endpointsSchema() *schema.Schema {
	endpointsAttributes := make(map[string]*schema.Schema)

	for _, endpointServiceName := range endpointServiceNames {
		endpointsAttributes[endpointServiceName] = &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: descriptions["endpoint"],
		}
	}

	// Since the endpoints attribute is a TypeSet we cannot use ConflictsWith
	endpointsAttributes["kinesis_analytics"].Deprecated = "use `endpoints` configuration block `kinesisanalytics` argument instead"
	endpointsAttributes["r53"].Deprecated = "use `endpoints` configuration block `route53` argument instead"

	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: endpointsAttributes,
		},
	}
}

func expandProviderIgnoreTags(l []interface{}) *keyvaluetags.IgnoreConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	ignoreConfig := &keyvaluetags.IgnoreConfig{}
	m := l[0].(map[string]interface{})

	if v, ok := m["keys"].(*schema.Set); ok {
		ignoreConfig.Keys = keyvaluetags.New(v.List())
	}

	if v, ok := m["key_prefixes"].(*schema.Set); ok {
		ignoreConfig.KeyPrefixes = keyvaluetags.New(v.List())
	}

	return ignoreConfig
}
