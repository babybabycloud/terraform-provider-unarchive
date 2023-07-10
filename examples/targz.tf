data "unarchive_file" "targz" {
  file_name = "Python-3.11.3.tgz"
  output = "targz"
  type = ".tar.gz"
  filters = [
    {
      "includes": ["operations/lists"],
      "excludes": ["P"]
    }
  ]
}

output "targz_output" {
    value = data.unarchive_file.targz.file_names
}
