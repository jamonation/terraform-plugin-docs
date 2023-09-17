terraform {
  required_providers {
    docs = {
      source = "chainguard.dev/edu/docs"
    }
  }
}

provider "docs" {
  name = "testy"
}

data "docs_readme" "all" {}
