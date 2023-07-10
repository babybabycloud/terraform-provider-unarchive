data "unarchive_file" "targz" {
  file_name = "Python-3.11.3.tgz"
  output = "targz"
  type = ".tar.gz"
  filters = [
    {
      "includes": ["Objects"],
      "excludes": ["\\.h$"]
    },
    {
      "includes": ["Lib"],
      "excludes": ["\\.py$"],
    }
  ]
}

output "targz_output" {
    value = data.unarchive_file.targz.file_names
}
