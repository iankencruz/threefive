# backend/atlas.hcl


env "local" {

  // Reference an input variable.
  url = "postgres://user:pass@localhost:5432/threefive?sslmode=disable"

  // The path to your schema files
  src = "file://backend/sql/migrations/"


  dev = "docker://postgres/17?search_path=public"

}
