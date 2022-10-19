package achievement

var Schema = `
CREATE TABLE achievement (
    id uuid primary key,
    name text,
    createdOn text
);
`
