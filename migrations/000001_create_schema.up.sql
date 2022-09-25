-- profile schema
CREATE table profile(
    id uuid PRIMARY KEY,
    email varchar(256) NOT NULL,
    firebase_id varchar(32)NOT NULL,
    provider varchar(64)NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    updated_at timestamp without time zone default (now() at time zone 'utc'),
    deleted_at timestamp without time zone default NULL,
    CONSTRAINT unique_firebase UNIQUE (firebase_id),
    CONSTRAINT unique_email UNIQUE (email, provider)
);

-- design schema
CREATE table design(
    id uuid PRIMARY KEY,
    name varchar(256) NOT NULL,
    fields json DEFAULT NULL,
    profile_id uuid REFERENCES profile(id),
    template TEXT NOT NULL,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    updated_at timestamp without time zone default (now() at time zone 'utc'),
    deleted_at timestamp without time zone default NULL
);