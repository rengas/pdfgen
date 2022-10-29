CREATE TYPE role AS ENUM ('normal', 'admin');

--api
CREATE TABLE IF NOT EXISTS users (
     id uuid PRIMARY KEY,
     email VARCHAR(256) NOT NULL UNIQUE,
     password_hash VARCHAR(256) NOT NULL,
     first_name VARCHAR(256) DEFAULT NULL,
     last_name VARCHAR(256) DEFAULT NULL,
     role role NOT NULL,
     created_at timestamp without time zone default (now() at time zone 'utc'),
     updated_at timestamp without time zone default (now() at time zone 'utc'),
     deleted_At timestamp without time zone default NULL
);

-- design schema
CREATE table design(
    id uuid PRIMARY KEY,
    name varchar(256) NOT NULL,
    fields json DEFAULT NULL,
    user_id uuid REFERENCES users(id),
    template TEXT NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    updated_at timestamp without time zone default (now() at time zone 'utc'),
    deleted_at timestamp without time zone default NULL
);