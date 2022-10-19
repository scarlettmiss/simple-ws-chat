package user

var Schema = `
CREATE TABLE users (
    id uuid primary key,
    username text,
    password text,
    created_at timestamp,
  	updated_at timestamp,
  	last_online timestamp,
  	online bool,
    skill_points numeric,
	points numeric,
	games_played numeric,
	reputation numeric,
	deleted bool
);
`
