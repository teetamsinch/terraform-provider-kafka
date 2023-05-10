package kafka

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mdhwk/terraform-provider-kafka/kafka/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bootstrap_servers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Comma separated list of brokers. Format <broker>:<port>.",
			},
			"aws_iam": awsIamSchema(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"kafka_acl": resourceACL(),
			//"kafka_acls": resourceACLs(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: configure,
	}
}

func awsIamSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Amazon Resource Name (ARN) of an IAM Role to assume prior to making API calls.",
					//ValidateFunc: verify.ValidARN,
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "An identifier for the assumed role session.",
					//ValidateFunc: validAssumeRoleSessionName,
				},
			},
		},
	}
}

func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := client.Config{
		BootstrapServers: stringValueSlice(d.Get("bootstrap_servers").([]interface{})),
		IAM:              parseAwsIAM(d),
	}

	var diags diag.Diagnostics

	c, err := client.NewClient(&config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}

func parseAwsIAM(d *schema.ResourceData) *client.IAM {
	if v, ok := d.GetOkExists("aws_iam"); ok {
		m := terraformMap(v)
		return &client.IAM{
			RoleArn:     m["role_arn"].(string),
			SessionName: m["session_name"].(string),
		}
	}
	return nil
}

func terraformMap(in interface{}) map[string]interface{} {
	if in != nil && len(in.([]interface{})) > 0 && in.([]interface{})[0] != nil {
		return in.([]interface{})[0].(map[string]interface{})
	}
	return nil
}

func stringValueSlice(in []interface{}) []string {
	s := make([]string, len(in))
	for i, v := range in {
		s[i] = v.(string)
	}
	return s
}

func stringValueMap(in map[string]interface{}) map[string]string {
	m := make(map[string]string, len(in))
	for k, v := range in {
		m[k] = v.(string)
	}
	return m
}