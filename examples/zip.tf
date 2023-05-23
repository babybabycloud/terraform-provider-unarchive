data "unarchive_file" "zip" {
  file_name = "jedis-mock-1.0.8.jar"
  output = "zip"
  type = ".zip"
}

output "zip_output" {
    value = data.unarchive_file.zip.file_names
}
