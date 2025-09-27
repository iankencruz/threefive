CREATE TABLE IF NOT EXISTS "users" (
  "id" bigint,
  "name" varchar NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "repos" (
  "id" bigint,
  "name" varchar NOT NULL,
  "owner_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "owner_id" FOREIGN KEY ("owner_id") REFERENCES "users" ("id")
);


CREATE TABLE IF NOT EXISTS "commits" (
  "id" bigint,
  "message" varchar NOT NULL,
  "repo_id" bigint NOT NULL,
  "author_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "repo_id" FOREIGN KEY ("repo_id") REFERENCES "repos" ("id"),
  CONSTRAINT "author_id" FOREIGN KEY ("author_id") REFERENCES "users" ("id")
);
