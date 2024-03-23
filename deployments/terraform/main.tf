terraform {
  required_version = ">= 1.7.5"
  required_providers {
    postgresql = {
      source = "cyrilgdn/postgresql"
      version = "1.22.0"
    }
  }
}

variable "postgres_db" {
    default = "my_db"
}

variable "postgres_user" {
    default = "my_user"
}

variable "postgres_password" {
    default = "my_password"
}

variable "postgres_host" {
    default = "127.0.0.1"
}

# Define PostgreSQL Provider
provider "postgresql" {
  host     = var.postgres_host
  username = var.postgres_user
  password = var.postgres_password
  database = var.postgres_db
  sslmode = "disable"
}

# Define Schemas
resource "postgresql_schema" "user_account" {
  name = "user_account"
}