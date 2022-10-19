package message

var Schema = `
CREATE TABLE messages (
    id uuid primary key,
    user_id uuid,
    message text,
    created_at timestamp
);
`
