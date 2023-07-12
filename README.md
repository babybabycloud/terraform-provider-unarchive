# terraform-provider-unarchive

## Description
The unarchive provider is helpful to extract files from archvie files. It now only support ZIP, TAR and TAR in GZip.

## Installation
Install with Terraform 0.13+.

```hcl
terraform {
  required_providers {
    unarchive = {
      source = "babybabycloud/unarchive"
      version = "1.0.0"
    }
  }
}

provider "unarchive" {
  # Configuration options
}
```

## Usage

Please refer to [Docs](https://registry.terraform.io/providers/babybabycloud/unarchive/latest/docs) for more details.
