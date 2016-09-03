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
			"arguments": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The input arguments for commands",
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
	config := meta.(*Config)

	wd := config.WorkingDirectory
	command, err := interpolateCommand(config.CreateCommand, config.CreateParameters, argumentsAsStrings(d))
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Creating generic resource: %s", command)
	_, err = runCommand(command, wd)
	if err != nil {
		return err
	}

	d.SetId(hash(command))
	log.Printf("[INFO] Created generic resource: %s", d.Id())

	return resourceGenericShellRead(d, meta)
}

func resourceGenericShellRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	wd := config.WorkingDirectory
	command, err := interpolateCommand(config.ReadCommand, config.ReadParameters, argumentsAsStrings(d))
	if err != nil {
		return err
	}
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
	config := meta.(*Config)

	wd := config.WorkingDirectory
	command, err := interpolateCommand(config.DeleteCommand, config.DeleteParameters, argumentsAsStrings(d))
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Deleting generic resource: %s", command)
	_, err = runCommand(command, wd)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func argumentsAsStrings(d *schema.ResourceData) map[string]string {
	args := make(map[string]string)
	for key, val := range d.Get("arguments").(map[string]interface{}) {
		args[key] = val.(string)
	}
	return args
}

func interpolateCommand(command string, parameters []interface{}, arguments map[string]string) (string, error) {
	if len(parameters) == 0 {
		return command, nil
	}

	inputArgs := make([]interface{}, len(parameters))
	for i, p := range parameters {
		if v, ok := arguments[p.(string)]; ok {
			inputArgs[i] = v
		} else {
			return "", fmt.Errorf("Error interpolating command '%s', parameter '%s' missing.", command, p)
		}
	}
	log.Printf("[DEBUG] Interpolating, command '%s' and args: '%v'", command, inputArgs)
	newCommand := fmt.Sprintf(command, inputArgs...)

	pos := strings.Index(newCommand, "%!")
	if pos != -1 {
		return "", fmt.Errorf("Error interpolating command '%s' using args '%v'", newCommand, inputArgs)
	}

	return newCommand, nil
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
