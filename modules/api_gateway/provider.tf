provider "aws" {
  region = "eu-north-1"
}

terraform {
  required_providers {
    aws = {
      version = "~> 5.82.2"
    }
  }
}
