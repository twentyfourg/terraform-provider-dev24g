24G Terraform Bitbucket Provider
==================

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

In order to authenticate with Bitbucket, an App Password must be created (Personal Settings > App Passwords). Give it all permissions as it is still scoped to the permissions of your account.

The following environment variables are required:
```sh
export BITBUCKET_USERNAME=<bitbucket username>
export BITBUCKET_PASSWORD=<bitbucket app password>
```

Building and Installing The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/twentyfourg/terraform-provider-dev24g`

```sh
$ export GOPATH=$HOME/go
$ mkdir -p $GOPATH/src/github.com/twentyfourg; cd $GOPATH/src/github.com/twentyfourg
$ git clone git@github.com:twentyfourg/terraform-provider-dev24g.git
```

Enter the provider directory and build the provider. `go install` puts the new binary into `$HOME/go/bin`.

```sh
$ cd $GOPATH/src/src/github.com/twentyfourg/terraform-provider-dev24g
$ go install
```

Copy installed provider into the terraform plugins directory

```sh
mkdir -p $HOME/.terraform.d/plugins/
cp $GOPATH/bin/terraform-provider-dev24g $HOME/.terraform.d/plugins
```

Using the provider
----------------------

```hcl
# Configure the Provider
provider "dev24g" {
  workspace = "24g"
}

data "dev24g_bitbucket_repository" "api" {
  name      = "796-4-1-vxp-api"
}

resource "dev24g_bitbucket_deployment" "evan" {
  name       = "evan"
  stage      = "Test"
  repository = "${data.dev24g_bitbucket_repository.api.workspace}/${data.dev24g_bitbucket_repository.api.name}"
}

resource "dev24g_bitbucket_deployment_variable" "foobar" {
  key        = "foo"
  value      = "bar"
  secured    = false
  deployment = dev24g_bitbucket_deployment.evan.id
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-bitbucket
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```