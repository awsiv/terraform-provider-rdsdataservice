package rdsdataservice

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsRdsdataservicePostgresRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsRdsdataservicePostgresRoleCreate,
		Read:   resourceAwsRdsdataservicePostgresRoleRead,
		Update: resourceAwsRdsdataservicePostgresRoleUpdate,
		Delete: resourceAwsRdsdataservicePostgresRoleDelete,
		Exists: resourceAwsRdsdataservicePostgresRoleExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The PostgreSQL database name to connect to",
			},
			"login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determine whether a role is allowed to log in",
			},
			"inherit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Determine whether a role "inherits" the privileges of roles it is a member of`,
			},
			"create_database": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Define a role's ability to create databases",
			},
			"create_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determine whether this role will be permitted to create new roles",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Sets the role's password",
			},
			"roles": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				MinItems:    0,
				Description: "Role(s) to grant to this new role",
			},
			"superuser": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Determine whether the new role is a "superuser"`,
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

func resourceAwsRdsdataservicePostgresRoleCreate(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	password := ""
	if attr, ok := d.GetOk("password"); ok {
		password = fmt.Sprintf(" PASSWORD '%s' ", attr.(string))
	}
	superuser := ""
	if attr, ok := d.GetOk("superuser"); ok {
		if attr.(bool) {
			superuser = fmt.Sprintf(" SUPERUSER ")
		} else {
			superuser = fmt.Sprintf(" NOSUPERUSER ")
		}
	}
	createrole := ""
	if attr, ok := d.GetOk("create_role"); ok {
		if attr.(bool) {
			createrole = fmt.Sprintf(" CREATEROLE ")
		} else {
			createrole = fmt.Sprintf(" NOCREATEROLE ")
		}
	}

	createdatabase := ""
	if attr, ok := d.GetOk("create_database"); ok {
		if attr.(bool) {
			createrole = fmt.Sprintf(" CREATEDB ")
		} else {
			createrole = fmt.Sprintf(" NOCREATEDB ")
		}
	}

	inherit := ""
	if attr, ok := d.GetOk("inherit"); ok {
		if attr.(bool) {
			inherit = fmt.Sprintf(" INHERIT ")
		} else {
			inherit = fmt.Sprintf(" NOINHERIT ")
		}
	}

	login := ""
	if attr, ok := d.GetOk("login"); ok {
		if attr.(bool) {
			login = fmt.Sprintf(" LOGIN ")
		} else {
			login = fmt.Sprintf(" NOLOGIN ")
		}
	}
	sql := fmt.Sprintf("CREATE ROLE %s WITH %s %s %s %s %s;",
		password,
		superuser,
		createrole,
		createdatabase,
		inherit,
		login,
	)

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Create Postgres Role: %#v", createOpts)

	_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error creating Postgres Role: %#v", err)
	}

	// Grant roles

	for _, grantingRole := range d.Get("roles").(*schema.Set).List() {
		query := fmt.Sprintf(
			"GRANT %s TO %s", grantingRole.(string), d.Get("name").(string),
		)
		createOpts = rdsdataservice.ExecuteStatementInput{
			ResourceArn: aws.String(d.Get("resource_arn").(string)),
			SecretArn:   aws.String(d.Get("secret_arn").(string)),
			Sql:         aws.String(query),
		}
		log.Printf("[DEBUG] Create Postgres Grant Role: %#v", createOpts)

		_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

		if err != nil {
			return fmt.Errorf("Error granting Role: %s to %s: %#v", grantingRole, d.Get("name").(string), err)
		}
	}

	d.SetId(d.Get("name").(string))
	log.Printf("[INFO] Postgres Role ID: %s", d.Id())

	return err
}

func resourceAwsRdsdataservicePostgresRoleDelete(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("DROP ROLE %s",
		d.Get("name").(string))

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Drop Postgres Role: %#v", createOpts)

	_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return fmt.Errorf("Error dropping Postgres Role: %#v", err)
	}

	d.SetId("")
	return nil
}

func resourceAwsRdsdataservicePostgresRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf("SELECT rolname FROM pg_catalog.pg_roles WHERE rolname='%s'",
		d.Get("name").(string))

	createOpts := rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(d.Get("resource_arn").(string)),
		SecretArn:   aws.String(d.Get("secret_arn").(string)),
		Sql:         aws.String(sql),
	}

	log.Printf("[DEBUG] Check Postgres Role exists: %#v", createOpts)

	output, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

	if err != nil {
		return false, fmt.Errorf("Error checking Postgres Role exists: %#v", err)
	}

	if len(output.Records) != 1 {
		return false, nil
	}

	return true, nil
}

func resourceAwsRdsdataservicePostgresRoleRead(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	sql := fmt.Sprintf(`SELECT 
		rolname,
		rolsuper,
		rolinherit, 
		rolcreaterole,
		rolcreatedb,
		rolcanlogin 
		FROM pg_catalog.pg_roles WHERE rolname=%s`,
		d.Get("name").(string),
	)

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

	d.Set("name", output.Records[0][0].StringValue)
	d.Set("rolsuper", output.Records[0][1].StringValue)
	d.Set("rolinherit", output.Records[0][1].StringValue)
	d.Set("rolcreaterole", output.Records[0][1].StringValue)
	d.Set("rolcreatedb", output.Records[0][1].StringValue)
	d.Set("rolcanlogin", output.Records[0][1].StringValue)

	// TODO: password

	d.SetId(d.Get("name").(string))
	return err
}

func resourceAwsRdsdataservicePostgresRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	rdsdataserviceconn := meta.(*AWSClient).rdsdataserviceconn

	// TODO: run this in transaction

	if d.HasChange("name") {
		oraw, nraw := d.GetChange("name")
		o := oraw.(string)
		n := nraw.(string)
		if n == "" {
			return fmt.Errorf("Error setting role name to an empty string")
		}

		sql := fmt.Sprintf("ALTER ROLE %s RENAME TO %s", o, n)

		createOpts := rdsdataservice.ExecuteStatementInput{
			ResourceArn: aws.String(d.Get("resource_arn").(string)),
			SecretArn:   aws.String(d.Get("secret_arn").(string)),
			Sql:         aws.String(sql),
		}

		log.Printf("[DEBUG] Update Postgres Role name: %#v", createOpts)

		_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

		if err != nil {
			return fmt.Errorf("Error updating Postgres Role name: %#v", err)
		}
		d.SetId(n)
	}
	// TODO: Store secret arn for role in tfstate
	if d.HasChange("login") {
		login := d.Get("login").(bool)
		tok := "NOLOGIN"
		if login {
			tok = "LOGIN"
		}

		sql := fmt.Sprintf("ALTER ROLE %s WITH %s", d.Get("name").(string), tok)

		createOpts := rdsdataservice.ExecuteStatementInput{
			ResourceArn: aws.String(d.Get("resource_arn").(string)),
			SecretArn:   aws.String(d.Get("secret_arn").(string)),
			Sql:         aws.String(sql),
		}

		log.Printf("[DEBUG] Update Postgres Role login: %#v", createOpts)

		_, err := rdsdataserviceconn.ExecuteStatement(&createOpts)

		if err != nil {
			return fmt.Errorf("Error updating Postgres Role login: %#v", err)
		}
	}

	return nil
}
