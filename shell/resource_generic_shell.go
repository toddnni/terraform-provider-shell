package shell

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/armon/circbuf"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGenericShell() *schema.Resource {
	return &schema.Resource{
		Create: resourceGenericShellCreate,
		Read:   resourceGenericShellRead,
		Delete: resourceGenericShellDelete,

		// desc: will always recreate the resource if something is changed
		// will output variables but we don't define them here
		// eg. if contains access_ipv4

		Schema: map[string]*schema.Schema{
			"create_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Command to create a resource",
				ForceNew:    true,
			},
			"working_directory": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A working directory where to run the commands",
				ForceNew:    true,
				Default:     ".",
			},
			"read_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Command to read status of a resource",
				ForceNew:    true,
			},
			"delete_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Command to delete a resource",
				ForceNew:    true,
			},
			"output": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Output from the read command",
			},
		},
	}
}

func resourceGenericShellCreate(d *schema.ResourceData, meta interface{}) error {

	command := d.Get("create_command").(string)
	wd := d.Get("working_directory").(string)
	log.Printf("[DEBUG] Creating generic resource: %s", command)
	_, err := runCommand(command, wd)
	if err != nil {
		return err
	}

	d.SetId(hash(command))
	log.Printf("[INFO] Created generic resource: %s", d.Id())

	return resourceGenericShellRead(d, meta)
}

func resourceGenericShellRead(d *schema.ResourceData, meta interface{}) error {

	command := d.Get("read_command").(string)
	wd := d.Get("working_directory").(string)
	log.Printf("[DEBUG] Reading generic resource: %s", command)
	output, err := runCommand(command, wd)
	if err != nil {
		log.Printf("[INFO] Read command returned error, marking resource deleted: %s", output)
		d.SetId("")
		return nil
	}

	outputs := make(map[string]string)
	split := strings.Split(output, "\n")
	for _, varline := range split {
		log.Printf("[DEBUG] Generic resource read line: %s", varline)

		if varline == "" {
			continue
		}

		pos := strings.Index(varline, "=")
		if pos == -1 {
			log.Printf("[INFO] Generic, ignoring line without equal sign: \"%s\"", varline)
			continue
		}

		key := varline[:pos]
		value := varline[pos+1:]
		log.Printf("[DEBUG] Generic: \"%s\" = \"%s\"", key, value)
		outputs[key] = value
	}
	d.Set("output", outputs)

	return nil
}

func resourceGenericShellDelete(d *schema.ResourceData, meta interface{}) error {

	command := d.Get("delete_command").(string)
	wd := d.Get("working_directory").(string)
	log.Printf("[DEBUG] Deleting generic resource: %s", command)
	_, err := runCommand(command, wd)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

const (
	// maxBufSize limits how much output we collect from a local
	// invocation. This is to prevent TF memory usage from growing
	// to an enormous amount due to a faulty process.
	maxBufSize = 8 * 1024
)

func runCommand(command string, working_dir string) (string, error) {
	// Execute the command using a shell
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	// Setup the command
	command = fmt.Sprintf("cd %s && %s", working_dir, command)
	cmd := exec.Command(shell, flag, command)
	output, _ := circbuf.NewBuffer(maxBufSize)
	cmd.Stderr = io.Writer(output)
	cmd.Stdout = io.Writer(output)

	// Output what we're about to run
	log.Printf("[DEBUG] generic shell resource going to execute: %s %s \"%s\"", shell, flag, command)

	// Run the command to completion
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("Error running command '%s': '%v'. Output: %s",
			command, err, output.Bytes())
	}

	log.Printf("[DEBUG] generic shell resource command output was: \"%s\"", output)

	return output.String(), nil
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
