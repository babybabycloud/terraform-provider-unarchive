terraform {
  required_providers {
    unarchive = {
      source = "hashicorp.com/edu/unarchive"
    }
  }
}

provider "unarchive" {}

data "unarchive_zip_file" "example" {
  file_name = "master.zip"
}
