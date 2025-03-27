INSERT INTO user_roles(name, slug) VALUES('Admin', 'admin');

INSERT INTO user_permissions(name, slug) VALUES('Create User', 'create-user');
INSERT INTO user_permissions(name, slug) VALUES('Read User', 'read-user');
INSERT INTO user_permissions(name, slug) VALUES('Update User', 'update-user');
INSERT INTO user_permissions(name, slug) VALUES('Delete User', 'create-user');

-- Associate permissions with role
INSERT INTO user_roles_permissions(role_id, permission_id) VALUES(1, 1);
INSERT INTO user_roles_permissions(role_id, permission_id) VALUES(1, 2);
INSERT INTO user_roles_permissions(role_id, permission_id) VALUES(1, 3);
INSERT INTO user_roles_permissions(role_id, permission_id) VALUES(1, 4);