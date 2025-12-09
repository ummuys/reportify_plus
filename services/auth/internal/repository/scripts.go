package repository

const (
	loginQuery = `
SELECT 
    u.user_id,
    u.password,
    r.name AS role
FROM identity.users AS u
JOIN identity.user_roles AS ur ON ur.user_id = u.user_id
JOIN identity.roles AS r ON r.role_id = ur.role_id
WHERE u.username = $1;
`

	createUserQuery = `
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

	updateUsernameQuery = `
UPDATE identity.users SET username = $2 WHERE user_id = $1;
`

	updatePasswordQuery = `
UPDATE identity.users SET password = $2 WHERE user_id = $1;
`

	// #nosec G101 -- SQL query, not hardcoded password
	updateRoleQuery = `
    WITH selected_role AS (
        SELECT role_id FROM identity.roles WHERE name = $2
    )
    UPDATE identity.user_roles ur
    SET role_id = sr.role_id
    FROM selected_role sr
    WHERE ur.user_id = $1;
    `

	deleteUserQuery = `DELETE FROM identity.users WHERE user_id = $1`
)
