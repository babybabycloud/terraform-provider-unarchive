terraform {
  required_providers {
    unarchive = {
      source = "hashicorp.com/edu/unarchive"
    }
  }
}

provider "unarchive" {}

data "unarchive_file" "example" {
  file_name = "master.zip"
  flat = true
  output = "zip"
  type = ".zip"
}

output "name" {
  value = data.unarchive_file.example.file_names
}

output "file_name" {
  value = data.unarchive_file.example.file_name
}

data "unarchive_file" "tar" {
  file_name = "h.tar.gz"
  output = "targz"
  type = ".tar.gz"
}

output "tar_names" {
  value = data.unarchive_file.tar.file_names
}
