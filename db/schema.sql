CREATE TABLE "todos" (
    "id" serial NOT NULL,
    "title" varchar(255) NULL,
    "description" text NULL,
    "completed" boolean NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "users" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar(255) NOT NULL,
)

CREATE TABLE "guilds" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar(255) NOT NULL,
    "description" text NULL,
)

CREATE TABLE "channels" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar(255) NOT NULL,
    "description" text NULL,
    "private" boolean NOT NULL DEFAULT false,
    "guild_id" UUID NOT NULL,
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE
)

CREATE TABLE "messages" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "content" text NOT NULL,
    "user_id" UUID NOT NULL,
    "channel_id" UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE
)

CREATE TABLE "user_channels" (
    user_id UUID NOT NULL,
    channel_id UUID NOT NULL,

    PRIMARY KEY (user_id, channel_id),
)