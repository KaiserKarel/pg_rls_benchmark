CREATE TYPE permission_level AS ENUM (
    'read', 'alter', 'owner', 'admin'
    );

CREATE TABLE object_user_permissions (
	user_id INTEGER REFERENCES users (id),
	object_id INTEGER REFERENCES objects (id),
	level permission_level,

    PRIMARY KEY (user_id, object_id)
);

CREATE TABLE object_group_permissions (
    group_id INTEGER REFERENCES groups (id),
    object_id INTEGER REFERENCES objects (id),
    level permission_level,

    PRIMARY KEY (group_id, object_id)
);

