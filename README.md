# terraform-provider

Provisions and configures all the fundamental entities of [Bob](https://bob-cd.github.io/):
- [Resource Provider](https://bob-cd.github.io/pages/concepts/resource.html)
- [Artifact Store](https://bob-cd.github.io/pages/concepts/artifact.html)
- [Pipeline](https://bob-cd.github.io/pages/concepts/pipeline.html)

## Status

Experimental

## Requirements
- [Terraform](https://www.terraform.io/downloads) 1.0+
- A running instance of [Bob](https://bob-cd.github.io/pages/getting-started.html)
- [Go](https://go.dev/dl/) 1.16+
- A recent version of [Babashka](https://github.com/babashka/babashka#installation)

## Installing
- Run `bb install` to build and install it in the terraform plugins dir

## Usage
- The sample terraform code is located in the `examples` dir
- Run `terraform init` to intitialise terraform with the provider
- Run `terraform plan -out tfplan` to check the plan
- Run `terraform apply tfplan` to apply

## License
Copyright 2022 Rahul De under [MIT](https://opensource.org/licenses/MIT)
