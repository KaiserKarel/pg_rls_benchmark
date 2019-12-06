CREATE TABLE users_groups (
	user_id INTEGER REFERENCES users (id),
	group_id INTEGER REFERENCES groups (id),
	PRIMARY KEY (user_id, group_id)
);
