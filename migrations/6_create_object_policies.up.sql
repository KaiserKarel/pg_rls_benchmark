CREATE POLICY has_read_permission ON objects
    FOR SELECT
    USING (
        (SELECT level FROM object_user_permissions WHERE user_id = current_setting('jwt.userid')::INTEGER AND object_id = id) >= 'read'
        OR
        (SELECT level FROM object_group_permissions WHERE group_id =
            (SELECT group_id FROM users_groups WHERE user_id = current_setting('jwt.userid')::INTEGER AND object_id = id)) >= 'read'
        );

CREATE POLICY has_insert_permission ON objects
    FOR UPDATE
    USING (
        (SELECT level FROM users WHERE id = current_setting('jwt.userid')::INTEGER) >= 'regular'
    );

CREATE POLICY has_alter_permission ON objects
    FOR UPDATE
    USING (
        (SELECT level FROM object_user_permissions WHERE user_id = current_setting('jwt.userid')::INTEGER AND object_id = id) >= 'alter'
        OR
        (SELECT level FROM object_group_permissions WHERE group_id IN
            (SELECT group_id FROM users_groups WHERE user_id = current_setting('jwt.userid')::INTEGER AND object_id = id)) >= 'alter'
    );

CREATE POLICY has_owner_permission ON objects
    FOR ALL
    USING (
        (SELECT level FROM object_user_permissions WHERE user_id = current_setting('jwt.userid')::INTEGER AND object_id = id) >= 'owner'
        OR
        (SELECT level FROM object_group_permissions WHERE group_id IN
            (SELECT group_id FROM users_groups WHERE user_id = current_setting('jwt.userid')::INTEGER AND object_id = id)) >= 'owner'
    );

ALTER TABLE objects ENABLE ROW LEVEL SECURITY ;
