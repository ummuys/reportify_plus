package repository

const LoginQuery = `
SELECT 
    u.user_id,
    u.password,
    r.name AS role
FROM identity.users AS u
JOIN identity.user_roles AS ur ON ur.user_id = u.user_id
JOIN identity.roles AS r ON r.role_id = ur.role_id
WHERE u.username = $1;
`

const CreateUserQuery = `
WITH new_user AS (
    INSERT INTO identity.users (username, password)
    VALUES ($1, $2)
    RETURNING user_id
),
ins_role AS (
    INSERT INTO identity.user_roles (user_id, role_id)
    SELECT
        new_user.user_id,
        r.role_id
    FROM new_user
    JOIN identity.roles r ON r.name = $3
)
SELECT user_id FROM new_user;
`
