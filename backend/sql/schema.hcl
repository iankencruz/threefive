table "galleries" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
  }
  primary_key {
    columns = [column.id]
  }
}
table "users" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
  }
  column "first_name" {
    null = false
    type = character_varying(255)
  }
  column "last_name" {
    null = false
    type = character_varying(255)
  }
  column "email" {
    null = false
    type = character_varying(255)
  }
  column "password_hash" {
    null = false
    type = character_varying(255)
  }
  column "role" {
    null    = false
    type    = character_varying(50)
    default = "user"
  }
  column "is_active" {
    null    = false
    type    = boolean
    default = true
  }
  column "avatar_url" {
    null = true
    type = text
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_email" {
    columns = [column.email]
  }
  index "idx_users_is_active" {
    columns = [column.is_active]
  }
  index "idx_users_role" {
    columns = [column.role]
  }
  unique "users_email_key" {
    columns = [column.email]
  }
}
schema "public" {
  comment = "standard public schema"
}
