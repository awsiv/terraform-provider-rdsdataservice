package rdsdataservice

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lib/pq"
)

func dbExists(dbname string, d *schema.ResourceData, meta interface{}) (bool, error) {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("SELECT datname FROM pg_database WHERE datname='%s'", dbname)

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Check db exists: %#v", createOpts)

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return false, fmt.Errorf("Error checking db exists: %#v", err)
	}

	if len(output.Records) == 0 {
		return false, nil
	}

	return true, nil
}

func schemaExists(schemaname string, d *schema.ResourceData, meta interface{}) (bool, error) {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("SELECT 1 FROM pg_namespace WHERE nspname='%s'", schemaname)

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Check schema exists: %#v", createOpts)

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return false, fmt.Errorf("Error checking schema exists: %#v", err)
	}

	if len(output.Records) == 0 {
		return false, nil
	}

	return true, nil
}

func roleExists(rolename string, d *schema.ResourceData, meta interface{}) (bool, error) {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("SELECT 1 FROM pg_roles WHERE rolname='%s'", rolename)

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Check role exists: %#v", createOpts)

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return false, fmt.Errorf("Error checking role exists: %#v", err)
	}

	if len(output.Records) == 0 {
		return false, nil
	}

	return true, nil
}

func pgArrayToSet(arr pq.ByteaArray) *schema.Set {
	s := make([]interface{}, len(arr))
	for i, v := range arr {
		s[i] = string(v)
	}
	return schema.NewSet(schema.HashString, s)
}
