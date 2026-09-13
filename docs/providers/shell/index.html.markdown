---
layout: "script"
page_title: "Provider: Script"
sidebar_current: "docs-script-index"
description: |-
  The Script provider enables one to wrap shell based tools to Terraform resources in a simple way.
---

# Script Provider

The Script provider enables one to wrap shell based tools to Terraform resources in a simple way. The provider needs only commands that provide `Create`, `Read` and `Delete` CRUD operations. The commands will be run when Terraform manages resources and well tailored scripts and this provider will provide an easy way to integrate existing functionality to Terraform using shell scripts.

However, script based resource can't reach the same level of integration to Terraform as native providers, because of the following reasons
- Shell scripts are Operating System specific and the scripts will depend on tools provided by the machine running Terraform.
- Script based resources don't support `Update` operation. The only option is to recreate the resource on change.

The provider is an easy way to wrap existing functionality or learn how terraform providers work, but in many cases the right solution is to write a native provider.

Read more about Terraform provider operations and internals from module documentation TODOTOD.

The path to script working directory can to be provided in provider.
# This also in resource help

Use the navigation to the left to read about the available resources.

## Example Usage

```
# Configure the Script provider
provider "script" {
    working_directory = "${path.module}/scripts"
}

# Create a simple file using a wrapper
resource "script_wrapper" "file" {
    sdfsdfasdfsad
    # TODO
}
```

## Argument Reference

The following arguments are supported:

* `working_directory` - (Optional) The working directory where the commands for script resources will be run. Defaults to current working directory.
