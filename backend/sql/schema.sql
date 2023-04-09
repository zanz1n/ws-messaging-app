CREATE TYPE "UserRole" AS ENUM ('ADMIN', 'USER');

CREATE TABLE "user" (
    "id" VARCHAR(36) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "role" "UserRole" NOT NULL DEFAULT 'USER',
    "username" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "message" (
    "id" VARCHAR(36) NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "content" TEXT,
    "imageUrl" TEXT,
    "userId" VARCHAR(36) NOT NULL,

    CONSTRAINT "message_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "user_username_key" ON "user"("username");

CREATE INDEX "user_id_username_idx" ON "user"("id", "username");

CREATE INDEX "message_id_userId_createdAt_idx" ON "message"("id", "userId", "createdAt");

ALTER TABLE "message" ADD CONSTRAINT "message_userId_fkey" FOREIGN KEY ("userId") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
