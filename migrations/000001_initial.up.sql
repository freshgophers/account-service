BEGIN;

-- EXTENSIONS --

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- TABLES --

CREATE TABLE IF NOT EXISTS secrets
  (
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id              UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    secret          VARCHAR NOT NULL,
    phone           VARCHAR NOT NULL,
    attempts        INTEGER NOT NULL DEFAULT 0,
    status          VARCHAR NOT NULL
  );

-- DATA --

COMMIT;
