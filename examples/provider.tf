terraform {
  required_providers {
    bob = {
      version = ">= 0.1.0"
      source  = "bob-cd/providers/bob"
    }
  }
}

provider "bob" {
  url     = "http://localhost:7777" # default
  timeout = 10000                   # in ms, default
}
