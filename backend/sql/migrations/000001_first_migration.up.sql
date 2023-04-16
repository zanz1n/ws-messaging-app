CREATE TYPE "UserRole" AS ENUM ('ADMIN', 'USER');

CREATE TABLE "user" (
    "id" VARCHAR(18) NOT NULL,
    "createdAt" BIGINT NOT NULL,
    "updatedAt" BIGINT NOT NULL,
    "role" "UserRole" NOT NULL DEFAULT 'USER',
    "username" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "message" (
    "id" VARCHAR(18) NOT NULL,
    "createdAt" BIGINT NOT NULL,
    "updatedAt" BIGINT NOT NULL,
    "content" TEXT,
    "imageUrl" TEXT,
    "userId" VARCHAR(36) NOT NULL,

    CONSTRAINT "message_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "user_username_key" ON "user"("username");

CREATE INDEX "user_id_username_idx" ON "user"("id", "username");

CREATE INDEX "message_id_createdAt_userId_idx" ON "message"("id", "createdAt", "userId");

ALTER TABLE "message" ADD CONSTRAINT "message_userId_fkey" FOREIGN KEY ("userId") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
