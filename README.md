# terraform-provider-postmark
This provider is built to handel terraform for a well-known mailing system called [postmark](https://postmarkapp.com/)

# Requirements
Terraform 1.3.x
Go 1.19.1 (to build the provider plugin)

# Installation

```
    $ git clone git@github.com:shebang-labs/terraform-provider-postmark.git
    $ make installLinux (for linux) or make installMac (for Mac)
```


# Usage
``` 
provider "postmark" {
  account_token = "xxxxxxxxxxxxxxxxxxxx"
}
```
# Run Tests
Run the following:
```
POSTMARK_ACCOUNT_TOKEN=xxxxxxxxxxxxxxxxxxxxxx TF_ACC=1 go test
```

Note: you need to replace the xxxxxxxxx with your account token.
# Data Sources Example

The next line will create a datasource that have all servers data from your postmark portal

```
data "postmark_servers" "servers" {}
```

# Resources

## Server

You can edit/delete for edit you can edit those two fields

```

resource "postmark_server" "s1" {
    name = "Server 1"
    color = "blue"
}
```

## Stream

You can edit/delete for edit you can edit description and name only

```
resource "postmark_stream" "st1" {
  stream_id = "transactional-dev-1"
  name = "Stream 1"
  description = "This is my first transactional stream"
  message_stream_type = "Transactional"
  server_token = "xxxxxxxxxxxxxxxxxxxx"
}
```

## Domain

You can edit/delete for edit you can only edit the return_path_domain
```
resource "postmark_domain" "d1" {
  name = "exmaple.com"
  return_path_domain = "test.example.com"
}
```