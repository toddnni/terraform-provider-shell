# Terraform shell provider

This is a terraform provider that lets you wrap shell based tools to [Terraform](https://terraform.io/) resources in a simple way.

There is also [exec](https://github.com/gosuri/terraform-exec-provider) provider, but it only implements `Create` CRUD operation.

## Naming

The naming of this provider has been hard. The provider is about wrapping functionality by running shell scripts. Originally the name was `generic_shell_wrapper`, but currently the name is just `shell`. The naming in code is still inconsistent.

## Installing

[Copied from the Terraform documentation](https://www.terraform.io/docs/plugins/basics.html):
> To install a plugin, put the binary somewhere on your filesystem, then configure Terraform to be able to find it. The configuration where plugins are defined is ~/.terraformrc for Unix-like systems and %APPDATA%/terraform.rc for Windows.

Build it from source (instructions below) and move the binary `terraform-provider-shell` to `bin/` and it should work.

## Using the provider

First, an simple example that is used in tests too

```hcl
provider "shell" {
  create_command = "echo \"hi\" > test_file"
  read_command = "awk '{print \"out=\" $0}' test_file"
  delete_command = "rm test_file"
}
```

```console
$ terraform plan
$ terraform apply
$ terraform destroy
```

## Building from source

1.  [Install Go](https://golang.org/doc/install) on your machine
2.  [Set up Gopath](https://golang.org/doc/code.html)
3.  `git clone` this repository into `$GOPATH/src/github.com/toddnni/terraform-provider-shell`
4.  Get the dependencies. Run `go get`
6.  `make install`. You will now find the
    binary at `$GOPATH/bin/terraform-provider-shell`.

## Running acceptance tests

```console
make test
```

## Known Problems

* Whenever a command is changed the resource will be rebuilt.
* The provider won't support `Update` CRUD operation.
* The provider won't print output of the commands.
* The provider will error instead of removing the resource if the delete command fails. However, this is a safe default.

## Author

* Toni Ylenius

The structure is inspired from the [Softlayer](https://github.com/finn-no/terraform-provider-softlayer) and [libvirt](https://github.com/dmacvicar/terraform-provider-libvirt) Terraform provider sources.

Some code has been adapted from local-exec provisioner from terraform core.

## License

* MIT, See LICENSE file
