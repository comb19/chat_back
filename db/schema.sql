CREATE TABLE "todos" (
    "id" serial NOT NULL,
    "title" varchar(255) NULL,
    "description" text NULL,
    "completed" boolean NULL,
    PRIMARY KEY ("id")
);