-- profile schema
CREATE table profile(
    id uuid PRIMARY KEY,
    email varchar(256)
);

-- design schema
CREATE table design(
    id uuid PRIMARY KEY,
    name varchar(256),
    fields json NOT NULL,
    profile_id uuid REFERENCES profile(id),
    template TEXT NOT NULL
);