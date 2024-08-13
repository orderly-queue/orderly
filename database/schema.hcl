schema "public" {}

table "users" {
  schema = schema.public

  column "id" {
    type = uuid
    null = false
  }

  column "name" {
    type = text
    null = false
  }

  column "email" {
    type = text
    null = false
  }

  column "password" {
    type = text
    null = false
  }

  column "admin" {
    type    = bool
    default = false
  }

  column "created_at" {
    type    = bigint
    null    = false
    default = "EXTRACT(epoch FROM NOW())"
  }

  column "updated_at" {
    type    = bigint
    null    = false
    default = "EXTRACT(epoch FROM NOW())"
  }

  column "deleted_at" {
    type = bigint
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_users_email" {
    columns = [column.email]
    unique  = true
  }
}
