CREATE TYPE user_level AS ENUM (
    'anon', 'regular', 'enterprise', 'support', 'employee', 'sudo'
);

CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	username TEXT,
	level user_level
);
