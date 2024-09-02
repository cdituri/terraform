
data "http" "gitlab_schema" {
  url = "https://gitlab.com/gitlab-org/gitlab/-/raw/master/app/assets/javascripts/editor/schema/ci.json"
}

data "local_file" "valid_test_yaml" {
  filename = "invalid_test.yml"
}

locals {
  valid_test_yaml_content = file("${data.local_file.valid_test_yaml.filename}")
  valid_test_yaml_json = jsonencode(yamldecode(local.valid_test_yaml_content))
  json_schema = data.http.gitlab_schema.response_body
}

resource "terraform_data" "yaml2json" {
  input = local.valid_test_yaml_json

  lifecycle {
    precondition {
      condition = jsonschema(local.json_schema, local.valid_test_yaml_json)
      error_message = "The JSON instance does not match the schema."
    }
  }
}

output "validated_yaml" {
    value = {
        #schema = local.json_schema
        content = {
          json = local.valid_test_yaml_json
          yaml = local.valid_test_yaml_content
        }
        isValid = jsonschema(local.json_schema, local.valid_test_yaml_json)
    }
}

