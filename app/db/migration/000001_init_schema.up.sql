CREATE TABLE "quiz" (
  "id" varchar(255) PRIMARY KEY,
  "owner" varchar(255) NOT NULL,
  "content" text,
  "hash_content" varchar(255) NOT NULL,
  "answer" text,
  "hash_answer" varchar(255),
  "duration" int NOT NULL,
  "status" int NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "answer" (
  "id" varchar(255) PRIMARY KEY,
  "quiz_id" varchar(255),
  "owner" varchar(255) NOT NULL,
  "content" text NOT NULL,
  "hash_content" varchar(255) NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "answer" ADD FOREIGN KEY ("quiz_id") REFERENCES "quiz" ("id");