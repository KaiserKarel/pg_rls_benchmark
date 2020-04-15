CREATE FUNCTION user_permissions(userid INTEGER, level permission_level) RETURNS TABLE (objectid INTEGER) AS $$
    BEGIN
        RETURN QUERY (SELECT object_id FROM object_user_permissions
                    WHERE user_id = userid);
    END;
    $$ LANGUAGE 'plpgsql' STABLE;

CREATE FUNCTION  group_permissions(groupid INTEGER, level permission_level) RETURNS TABLE (objectid INTEGER) AS $$
    BEGIN
        RETURN QUERY (SELECT object_id FROM object_group_permissions
                      WHERE group_id = groupid);
    END;
    $$ LANGUAGE 'plpgsql' STABLE;

CREATE FUNCTION  user_groups_permissions(userid INTEGER, level permission_level) RETURNS TABLE (objectid INTEGER) AS $$
    BEGIN
        RETURN QUERY (
            WITH groupids AS (SELECT group_id FROM users_groups WHERE user_id = userid)
            SELECT group_permissions(groupids.group_id, level) FROM groupids);
    END;
    $$ LANGUAGE 'plpgsql' STABLE;


CREATE POLICY user_has_read_permission ON objects
    FOR SELECT
    USING (
        id IN (SELECT user_permissions(current_setting('jwt.userid')::INTEGER, 'read'))
    );

CREATE POLICY user_group_has_read_permission ON objects
    FOR SELECT
    USING (
        id IN (SELECT user_groups_permissions(current_setting('jwt.userid')::INTEGER, 'read'))
    );

CREATE POLICY user_has_insert_permission ON objects
    FOR UPDATE
    USING (
        (SELECT level FROM users WHERE id = current_setting('jwt.userid')::INTEGER) >= 'regular'
    );

CREATE POLICY user_has_update_permission ON objects
    FOR UPDATE
    USING (
        id IN (SELECT user_permissions(current_setting('jwt.userid')::INTEGER, 'alter'))
    );

CREATE POLICY user_group_has_update_permission ON objects
    FOR UPDATE
    USING (
        id IN (SELECT user_groups_permissions(current_setting('jwt.userid')::INTEGER, 'alter'))
    );

CREATE POLICY user_has_delete_permission ON objects
    FOR DELETE
    USING (
        id IN (SELECT user_permissions(current_setting('jwt.userid')::INTEGER, 'owner'))
    );

CREATE POLICY user_group_has_alter_permission ON objects
    FOR DELETE
    USING (
        id IN (SELECT user_groups_permissions(current_setting('jwt.userid')::INTEGER, 'owner'))
    );

ALTER TABLE objects ENABLE ROW LEVEL SECURITY;
