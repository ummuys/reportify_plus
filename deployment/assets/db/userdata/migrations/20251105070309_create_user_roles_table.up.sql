CREATE TABLE IF NOT EXISTS identity.user_roles (
    user_id UUID PRIMARY KEY REFERENCES identity.users(user_id) ON DELETE CASCADE,
    role_id INT REFERENCES identity.roles(role_id) ON DELETE CASCADE
);

COMMENT ON TABLE identity.user_roles IS 'Роли пользователей';
COMMENT ON COLUMN identity.user_roles.user_id IS 'Ссылка на пользователя в identity.users';
COMMENT ON COLUMN identity.user_roles.role_id IS 'Ссылка на роли в identity.roles';