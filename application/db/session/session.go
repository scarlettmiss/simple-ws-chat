package session

var Schema = `
CREATE TABLE sessions (
    id uuid primary key,
    capacity numeric,
    min_rating numeric,
    max_rating numeric,
    session_constraint text,
    owner uuid,
    created_at timestamp
);
`
