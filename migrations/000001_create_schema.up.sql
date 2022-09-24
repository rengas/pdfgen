-- profile schema
CREATE table profile(
    id uuid PRIMARY KEY,
    email varchar(256),
    firebase_id varchar(32),
    provider varchar(64),
    CONSTRAINT unique_firebase UNIQUE (firebase_id),
    CONSTRAINT unique_email UNIQUE (email, provider)
);

-- design schema
CREATE table design(
    id uuid PRIMARY KEY,
    name varchar(256),
    fields json DEFAULT NULL,
    profile_id uuid REFERENCES profile(id),
    template TEXT NOT NULL
);