package rdsdataservice

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/lib/pq"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var objectTypes = map[string]string{
	"table":    "r",
	"sequence": "S",
}

func resourceAwsRdsdataservicePostgresGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsRdsdataservicePostgresGrantCreate,
		Read:   resourceAwsRdsdataservicePostgresGrantRead,
		// As create revokes and grants we can use it to update too
		Update: resourceAwsRdsdataservicePostgresGrantCreate,
		Delete: resourceAwsRdsdataservicePostgresGrantDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the role to grant privileges on",
			},
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The database to grant privileges on for this role",
			},
			"schema": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The database schema to grant privileges on for this role",
			},
			"object_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				/*ValidateFunc: validation.StringInSlice([]string{
					"table",
					"sequence",
				}, false),
				*/
				Description: "The PostgreSQL object type to grant the privileges on (one of: table, sequence)",
			},
			"privileges": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				MinItems:    1,
				Description: "The list of privileges to grant",
			},
			"secret_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The PostgreSQL database name to connect to",
			},
			"resource_arn": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "postgres",
				Description: "The PostgreSQL database name to connect to",
			},
		},
	}
}

func resourceAwsRdsdataservicePostgresGrantCreate(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn
	// TODO: Run this on transaction
	sql := fmt.Sprintf(
		"REVOKE ALL PRIVILEGES ON ALL %sS IN SCHEMA %s FROM %s",
		strings.ToUpper(d.Get("object_type").(string)),
		pq.QuoteIdentifier(d.Get("schema").(string)),
		pq.QuoteIdentifier(d.Get("role").(string)),
	)

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Create Postgres Grant: step 1: revoke: %#v", createOpts)

	_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error revoking Postgres grant: %#v", err)
	}

	// Grant roles
	privileges := []string{}
	for _, priv := range d.Get("privileges").(*schema.Set).List() {
		privileges = append(privileges, priv.(string))
	}

	sql = fmt.Sprintf(
		"GRANT %s ON ALL %s IN SCHEMA %s TO %s",
		strings.Join(privileges, ","),
		strings.ToUpper(d.Get("object_type").(string)),
		pq.QuoteIdentifier(d.Get("schema").(string)),
		pq.QuoteIdentifier(d.Get("role").(string)),
	)

	createOpts = rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}
	log.Printf("[DEBUG] Create Postgres Grant: step 2: grant: %#v", createOpts)

	_, err = rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error granting priviliges: %s to %s: %#v", strings.Join(privileges, ","), d.Get("name").(string), err)
	}

	d.SetId(generateGrantID(d))
	log.Printf("[INFO] Postgres Role ID: %s", d.Id())

	return err
}

func generateGrantID(d *schema.ResourceData) string {
	return strings.Join([]string{
		d.Get("role").(string), d.Get("database").(string),
		d.Get("schema").(string), d.Get("object_type").(string),
	}, "_")
}

func resourceAwsRdsdataservicePostgresGrantDelete(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf(
		"REVOKE ALL PRIVILEGES ON ALL %sS IN SCHEMA %s FROM %s",
		strings.ToUpper(d.Get("object_type").(string)),
		pq.QuoteIdentifier(d.Get("schema").(string)),
		pq.QuoteIdentifier(d.Get("role").(string)),
	)

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Drop Postgres Grant: %#v", createOpts)

	_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error dropping Postgres Grant: %#v", err)
	}

	d.SetId("")
	return nil
}

func resourceAwsRdsdataservicePostgresGrantRead(d *schema.ResourceData, meta interface{}) error {
	exists, err := checkRoleDBSchemaExists(d, meta)
	if err != nil {
		return err
	}
	if !exists {
		d.SetId("")
		return nil
	}
	return readRolePrivileges(d, meta)
}

func checkRoleDBSchemaExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	// Check the role exists
	role := d.Get("role").(string)
	exists, err := roleExists(role, d, meta)
	if err != nil {
		return false, err
	}
	if !exists {
		log.Printf("[DEBUG] role %s does not exists", role)
		return false, nil
	}

	// Check the database exists
	database := d.Get("database").(string)
	exists, err = dbExists(database, d, meta)
	if err != nil {
		return false, err
	}
	if !exists {
		log.Printf("[DEBUG] database %s does not exists", database)
		return false, nil
	}

	// Check the schema exists (the SQL connection needs to be on the right database)
	schema := d.Get("schema").(string)
	exists, err = schemaExists(schema, d, meta)
	if err != nil {
		return false, err
	}
	if !exists {
		log.Printf("[DEBUG] schema %s does not exists", schema)
		return false, nil
	}

	return true, nil
}

func readRolePrivileges(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	// This returns, for the specified role (rolname),
	// the list of all object of the specified type (relkind) in the specified schema (namespace)
	// with the list of the currently applied privileges (aggregation of privilege_type)
	//
	// Our goal is to check that every object has the same privileges as saved in the state.
	sql := `
SELECT pg_class.relname, array_remove(array_agg(privilege_type), NULL)
FROM pg_class
JOIN pg_namespace ON pg_namespace.oid = pg_class.relnamespace
LEFT JOIN (
    SELECT acls.* FROM (
        SELECT relname, relnamespace, relkind, (aclexplode(relacl)).* FROM pg_class c
    ) as acls
    JOIN pg_roles on grantee = pg_roles.oid
    WHERE rolname=$1
) privs
USING (relname, relnamespace, relkind)
WHERE nspname = $2 AND relkind = $3
GROUP BY pg_class.relname;
`

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error reading Postgres Database: %#v", err)
	}

	if len(output.Records) != 1 {
		d.SetId("")
		return nil
	}

	// TODO: Complete this :)

	return nil
}
