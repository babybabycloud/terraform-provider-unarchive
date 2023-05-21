terraform {
  required_providers {
    unarchive = {
      source = "hashicorp.com/edu/unarchive"
    }
  }
}

provider "unarchive" {}

data "unarchive_file" "zip" {
  file_name = "jedis-mock-1.0.8.jar"
  includes = ["server"]
  excludes = ["Redis"]
  output = "zip"
  type = ".zip"
}

output "zip_file_names" {
  value = data.unarchive_file.zip.file_names
}

data "unarchive_file" "targz" {
  file_name = "Python-3.11.3.tgz"
  includes = ["Objects"]
  excludes = ["txt"]
  output = "targz"
  type = ".tar.gz"
  flat = true
}

output "targz_file_names" {
  value = data.unarchive_file.targz.file_names
}
