env "dev" {
  url = getenv(DATABASE_URL)
  src = "file://db/schema.sql"
}

