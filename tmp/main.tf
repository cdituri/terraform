data "http" "taskfile_schema" {
    url = "https://taskfile.dev/schema.json"
}

data "http" "zloeber_chatgpt_taskfile" {
    url = "https://raw.githubusercontent.com/zloeber/taskfiles/main/tasks/Taskfile.chatgpt.yml"
}

locals {
  zloeber_yaml2json = jsonencode({
    invalid = "now it is invalid",
    original = yamldecode(data.http.zloeber_chatgpt_taskfile.response_body),
  })

  zloeber_taskfile_isValid = jsonschema(data.http.taskfile_schema.response_body, local.zloeber_yaml2json)
}

resource "terraform_data" "yaml2json" {
  input = local.zloeber_yaml2json

  lifecycle {
    precondition {
        condition = local.zloeber_taskfile_isValid
        error_message = "The JSON instance does not match the schema."
    }
  }
}

output "taskfile" {
    value = {
        schema = data.http.taskfile_schema.response_body
        content = {
            json = local.zloeber_yaml2json
            yaml = data.http.zloeber_chatgpt_taskfile.response_body
        }
        isValid = local.zloeber_taskfile_isValid
    }
}