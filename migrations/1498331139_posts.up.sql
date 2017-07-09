DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
CREATE EXTENSION pgcrypto;

------------------------------------------------------------
------------------------------------------------------------

CREATE TYPE user_role AS ENUM ('ADMIN', 'MEMBER');

CREATE TABLE users (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  uname          text NOT NULL,
  role           user_role NOT NULL DEFAULT 'MEMBER',
  digest         text NOT NULL,
  email          text NOT NULL,
  gpg_key        text NOT NULL,
  last_online_at timestamptz NOT NULL DEFAULT NOW(),
  created_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE posts (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id),
  title          text NOT NULL,
  slug           text NOT NULL,
  sub_title      text NOT NULL,
  short          text NOT NULL,
  post_content   text NOT NULL,
  digest         text NOT NULL,
  published      boolean DEFAULT FALSE,
  updated_at     timestamptz NOT NULL DEFAULT NOW(),
  created_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE images (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        uuid NOT NULL REFERENCES users(id),
  url            text NOT NULL,
  medium         text NOT NULL,
  small          text NOT NULL,
  caption        text NOT NULL,
  updated_at     timestamptz NOT NULL DEFAULT NOW(),
  created_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE tags (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name          text NOT NULL,
  slug          text NOT NULL
);

CREATE TABLE posts_tags (
  post_id       uuid NOT NULL REFERENCES posts(id),
  tag_id        uuid NOT NULL REFERENCES tags(id)
);
------------------------------------------------------------
------------------------------------------------------------

CREATE TABLE sessions (
  id            uuid PRIMARY KEY,
  user_id       uuid NOT NULL REFERENCES users(id),
  ip_address    inet NOT NULL,
  user_agent    text NULL,
  logged_out_at timestamptz NULL,
  expired_at    timestamptz NOT NULL DEFAULT NOW() + INTERVAL '2 weeks',
  created_at    timestamptz NOT NULL DEFAULT NOW()
);

-- Speed up user_id FK joins
CREATE INDEX sessions__user_id ON sessions (user_id);

CREATE VIEW active_sessions AS
  SELECT *
  FROM sessions
  WHERE expired_at > NOW()
    AND logged_out_at IS NULL
;

------------------------------------------------------------
------------------------------------------------------------

CREATE OR REPLACE FUNCTION ip_root(ip_address inet) RETURNS inet AS
$$
  DECLARE
    masklen int;
  BEGIN
    masklen := CASE family(ip_address) WHEN 4 THEN 24 ELSE 48 END;
    RETURN host(network(set_masklen(ip_address, masklen)));
  END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE TABLE ratelimits (
  id             bigserial        PRIMARY KEY,
  ip_address     inet             NOT NULL,
  created_at     timestamptz      NOT NULL DEFAULT NOW()
);

CREATE INDEX ratelimits__ip_root ON ratelimits (ip_root(ip_address));