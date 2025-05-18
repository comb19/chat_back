CREATE TABLE "todos" (
    "id" serial NOT NULL,
    "title" varchar(255) NULL,
    "description" text NULL,
    "completed" boolean NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "users" (
    "id" varchar(32) PRIMARY KEY,
    "user_name" varchar(255) NOT NULL
);

CREATE TABLE "guilds" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar(255) NOT NULL,
    "description" text NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "channels" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar(255) NOT NULL,
    "description" text NULL,
    "private" boolean NOT NULL DEFAULT false,
    "guild_id" UUID,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE
);

CREATE TABLE "messages" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "content" text NOT NULL,
    "user_id" varchar(32),
    "channel_id" UUID NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE
);

CREATE TABLE "user_channels" (
    user_id varchar(32) NOT NULL,
    channel_id UUID NOT NULL,

    PRIMARY KEY (user_id, channel_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE
);

CREATE TABLE "user_guilds" (
    user_id varchar(32) NOT NULL,
    guild_id UUID NOT NULL,

    PRIMARY KEY (user_id, guild_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE
);

CREATE TABLE "guild_invitations" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    "owner_id" varchar(32) NOT NULL,
    "guild_id" UUID NOT NULL,
    "expiration" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE, 
    FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_timestamp
BEFORE UPDATE ON channels
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();