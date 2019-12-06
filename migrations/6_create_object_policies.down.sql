DROP POLICY has_read_permission ON objects;
DROP POLICY has_insert_permission ON objects;
DROP POLICY has_alter_permission ON objects;
DROP POLICY has_owner_permission ON objects;
ALTER TABLE objects DISABLE ROW LEVEL SECURITY ;

