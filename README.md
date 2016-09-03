# Terraform shell provider

This is a terraform provider that lets you wrap shell based tools to 
[Terraform](https://terraform.io/) resources in a simple way.

## Installing

[Copied from the Terraform documentation](https://www.terraform.io/docs/plugins/basics.html):
> To install a plugin, put the binary somewhere on your filesystem, then configure Terraform to be able to find it. The configuration where plugins are defined is ~/.terraformrc for Unix-like systems and %APPDATA%/terraform.rc for Windows.

The binary should be renamed to terraform-provider-shell

You should update your .terraformrc and refer to the binary:

```hcl
providers {
  libvirt = "/path/to/terraform-provider-shell"
}
```

## Using the provider

Here is an example that will setup the following:


Now you can see the plan, apply it, and then destroy the infrastructure:

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

## Running

1.  create the example file main.tf in your working directory
2.  terraform plan
3.  terraform apply

## Running acceptance tests

```console
make test
```

## Known Problems

* Whenever command is changed the resource will be rebuilt.
* The provider won't support `Update` operation.
* a

## Author

* Toni Ylenius

The structure is inspired from the [Softlayer](https://github.com/finn-no/terraform-provider-softlayer) and [libvirt](https://github.com/dmacvicar/terraform-provider-libvirt) Terraform provider sources.

## License

* MIT, See LICENSE file
