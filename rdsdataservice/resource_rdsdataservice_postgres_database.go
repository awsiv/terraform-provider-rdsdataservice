package rdsdataservice

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsRdsdataservicePostgresDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsRdsdataservicePostgresDatabaseCreate,
		Read:   resourceAwsRdsdataservicePostgresDatabaseRead,
		Update: resourceAwsRdsdataservicePostgresDatabaseUpdate,
		Delete: resourceAwsRdsdataservicePostgresDatabaseDelete,
		Exists: resourceAwsRdsdataservicePostgresDatabaseExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name.",
			},
			"resource_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DB ARN.",
			},
			"secret_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DBA Secret ARN.",
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "postgres",
				Description: "The ROLE which owns the database.",
			},
		},
	}
}

func resourceAwsRdsdataservicePostgresDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("CREATE DATABASE %s OWNER %s;",
		d.Get("name").(string),
		d.Get("owner").(string))

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Create Postgres Database: %#v", createOpts)

	_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error creating Postgres Database: %#v", err)
	}

	d.SetId(d.Get("name").(string))
	log.Printf("[INFO] Postgres Database ID: %s", d.Id())

	return err
}

func resourceAwsRdsdataservicePostgresDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("DROP DATABASE %s;",
		d.Get("name").(string))

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Drop Postgres Database: %#v", createOpts)

	_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error dropping Postgres Database: %#v", err)
	}

	d.SetId("")
	return err
}

func resourceAwsRdsdataservicePostgresDatabaseExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("SELECT datname FROM pg_database WHERE datname='%s';",
		d.Get("name").(string))

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Check Postgres Database exists: %#v", createOpts)

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return false, fmt.Errorf("Error checking Postgres Database exists: %#v", err)
	}

	if len(output.Records) != 1 {
		return false, nil
	}

	return true, nil
}

func resourceAwsRdsdataservicePostgresDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("SELECT d.datname, pg_catalog.pg_get_userbyid(d.datdba) from pg_database d WHERE datname='%s';",
		d.Get("name").(string))

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Read Postgres Database: %#v", createOpts)

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error reading Postgres Database: %#v", err)
	}

	if len(output.Records) != 1 {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read Postgres Database details: %#v", output)

	/*
		sqlFmt := `SELECT %s` +
			`FROM pg_catalog.pg_database AS d, pg_catalog.pg_tablespace AS ts ` +
			`WHERE d.datname = '%s' AND d.dattablespace = ts.oid`
		sql = fmt.Sprintf(sqlFmt, d.Get("name").(string))

		createOpts = rdsdataservice.ExecuteStatementInput{
			ResourceArn: aws.String(d.Get("resource_arn").(string)),
			SecretArn:   aws.String(d.Get("secret_arn").(string)),
			Sql:         aws.String(sql),
		}
		log.Printf("[DEBUG] Read Postgres Database details: %#v", createOpts)

		if err != nil {
			return fmt.Errorf("Error reading Postgres Database: %#v", err)
		}

		if len(output.Records) != 1 {
			d.SetId("")
			return nil
		}

		log.Printf("[DEBUG] Read Postgres Database details: %#v", output.Records)
	*/
	d.Set("name", output.Records[0][0].StringValue)
	d.Set("owner", output.Records[0][1].StringValue)

	return err
}

func resourceAwsRdsdataservicePostgresDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	if d.HasChange("name") {
		oraw, nraw := d.GetChange("name")
		o := oraw.(string)
		n := nraw.(string)
		if n == "" {
			return fmt.Errorf("Error setting database name to an empty string")
		}

		sql := fmt.Sprintf("ALTER DATABASE %s RENAME TO %s", o, n)

		createOpts := rdsdataservice.ExecuteStatementInput{
			ResourceArn: aws.String(d.Get("resource_arn").(string)),
			SecretArn:   aws.String(d.Get("secret_arn").(string)),
			Sql:         aws.String(sql),
		}

		log.Printf("[DEBUG] Update Postgres Database name: %#v", createOpts)

		_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

		if err != nil {
			return fmt.Errorf("Error updating Postgres Database name: %#v", err)
		}
		d.SetId(n)
	}

	if d.HasChange("owner") {
		oraw, nraw := d.GetChange("owner")
		o := oraw.(string)
		n := nraw.(string)
		if n == "" {
			return fmt.Errorf("Error setting database owner to an empty string")
		}

		sql := fmt.Sprintf("ALTER DATABASE %s OWNER TO %s", o, n)

		createOpts := rdsdataservice.ExecuteStatementInput{
			ResourceArn: aws.String(d.Get("resource_arn").(string)),
			SecretArn:   aws.String(d.Get("secret_arn").(string)),
			Sql:         aws.String(sql),
		}

		log.Printf("[DEBUG] Update Postgres Database owner: %#v", createOpts)

		_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

		if err != nil {
			return fmt.Errorf("Error updating Postgres Database owner: %#v", err)
		}
	}

	return nil
}

func rdsDataserviceExecuteStatement(d *schema.ResourceData, sql string, meta interface{}) (*rdsdataservice.ExecuteStatementOutput, error) {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Update Postgres Database name: %#v", createOpts)

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return nil, fmt.Errorf("Error updating Postgres Database name: %#v", err)
	}

	return output, nil
}
