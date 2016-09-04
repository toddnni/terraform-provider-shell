package shell

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"working_directory": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PWD", nil),
				Description: "The working directory where to run.",
			},
			"create_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Create command",
			},
			"create_parameters": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Names of the parameters for the create command",
			},
			"read_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Read command",
			},
			"read_parameters": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Names of the parameters for the read command",
			},
			"delete_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Delete command",
			},
			"delete_parameters": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Names of the parameters for the delete command",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"shell_resource": resourceGenericShell(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		WorkingDirectory: d.Get("working_directory").(string),
		CreateCommand:    d.Get("create_command").(string),
		CreateParameters: d.Get("create_parameters").([]interface{}),
		ReadCommand:      d.Get("read_command").(string),
		ReadParameters:   d.Get("read_parameters").([]interface{}),
		DeleteCommand:    d.Get("delete_command").(string),
		DeleteParameters: d.Get("delete_parameters").([]interface{}),
	}

	return &config, nil
}
