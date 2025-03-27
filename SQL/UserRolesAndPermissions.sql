INSERT INTO roles(name, slug) VALUES('Admin', 'admin');
INSERT INTO roles(name, slug) VALUES('User Admin', 'user-admin');

INSERT INTO user_permissions(name, slug) VALUES('Create User', 'create-user');
INSERT INTO user_permissions(name, slug) VALUES('Read User', 'read-user');
INSERT INTO user_permissions(name, slug) VALUES('Update User', 'update-user');
INSERT INTO user_permissions(name, slug) VALUES('Delete User', 'create-user');

-- Associate permissions with role

-- Admin
INSERT INTO roles_permissions(role_id, permission_id) VALUES(1, 1);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(1, 2);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(1, 3);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(1, 4);

-- User Admin
select * from roles_permissions;
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 1);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 2);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 3);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 4);