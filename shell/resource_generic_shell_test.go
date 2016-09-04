package shell

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGenericShellProvider_Basic(t *testing.T) {
	const testConfig = `
	provider "shell" {
		working_directory = "."
		create_command = "echo \"hi\" > test_file"
		read_command = "awk '{print \"out=\" $0}' test_file"
		delete_command = "rm test_file"
	}
	resource "shell_resource" "test" {
	}
`

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenericShellDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResource("shell_resource.test", "out", "hi"),
				),
			},
		},
	})
}

func TestAccGenericShellProvider_WeirdOutput(t *testing.T) {
	const testConfig = `
	provider "shell" {
		create_command = "echo \" can you = read this\" > test_file3"
		read_command = "awk '{print \"out=\" $0}' test_file3"
		delete_command = "rm test_file3"
	}
	resource "shell_resource" "test" {
	}
`

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenericShellDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResource("shell_resource.test", "out", " can you = read this"),
				),
			},
		},
	})
}

func TestAccGenericShellProvider_Parameters(t *testing.T) {
	const testConfig = `
	provider "shell" {
		create_command = "echo \"%s\" > %s"
		create_parameters = [ "output", "file" ]
		read_command = "awk '{print \"out=\" $0}' %s"
		read_parameters = [ "file" ]
		delete_command = "rm %s"
		delete_parameters = [ "file" ]
	}
	resource "shell_resource" "test" {
		arguments {
			output = "param value"
			file = "file4"
		}
	}
`

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenericShellDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResource("shell_resource.test", "out", "param value"),
				),
			},
		},
	})
}

func TestAccGenericShellProvider_Update(t *testing.T) {
	const testConfig1 = `
	provider "shell" {
		create_command = "echo \"%s\" > %s"
		create_parameters = [ "output", "file" ]
		read_command = "awk '{print \"out=\" $0}' %s"
		read_parameters = [ "file" ]
		delete_command = "rm %s"
		delete_parameters = [ "file" ]
	}
	resource "shell_resource" "test" {
		arguments {
			output = "hi"
			file = "testfileU1"
		}
	}
`
	const testConfig2 = `
	provider "shell" {
		create_command = "echo \"%s\" > %s"
		create_parameters = [ "output", "file" ]
		read_command = "awk '{print \"out=\" $0}' %s"
		read_parameters = [ "file" ]
		delete_command = "rm %s"
		delete_parameters = [ "file" ]
	}
	resource "shell_resource" "test" {
		arguments {
			output = "hi all"
			file = "testfileU2"
		}
	}
`

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenericShellDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testConfig1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResource("shell_resource.test", "out", "hi"),
				),
			},
			resource.TestStep{
				Config: testConfig2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResource("shell_resource.test", "out", "hi all"),
				),
			},
		},
	})
}

func testAccCheckResource(name string, outparam string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s, found: %s", name, s.RootModule().Resources)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if expected, got := value, rs.Primary.Attributes["output."+outparam]; got != expected {
			return fmt.Errorf("Wrong value in output '%s=%s', expected '%s'", outparam, got, expected)
		}
		return nil
	}
}

func testAccCheckGenericShellDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "shell_resource" {
			continue
		}

		splitted := strings.Split(rs.Primary.Attributes["create_command"], " ")
		file := splitted[len(splitted)-1]
		if _, err := os.Stat(file); err == nil {
			return fmt.Errorf("File '%s' exists after delete", file)
		}
	}
	return nil
}
