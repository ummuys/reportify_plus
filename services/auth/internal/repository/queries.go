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
        INSERT INTO identity.users (user_id, username, password, is_protected)
        VALUES ($1, $2, $3, $4);
    `

	createUserRolesQuery = `
        INSERT INTO identity.user_roles (user_id, role_id)
        SELECT
            u.user_id,
            r.role_id
        FROM identity.users u
        JOIN identity.roles r ON r.name = $1
        WHERE u.user_id = $2;
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

	// P0-09 2 июля
	getUserQuery = `
        SELECT u.user_id, u.username, r.name
        FROM identity.users AS u
        JOIN identity.user_roles AS ur ON u.user_id = ur.user_id
        JOIN identity.roles AS r ON r.role_id = ur.role_id
        WHERE u.user_id = $1;
    `
	// ------

    // P0-08 6 июля
    getIsProtectesAndRoleQuery = `
        SELECT u.is_protected, r.name
        FROM identity.users u
        JOIN identity.user_roles ur ON u.user_id = ur.user_id
        JOIN identity.roles r ON r.role_id = ur.role_id
        WHERE u.user_id = $1;
    `

    countAdminsQuery = `
        SELECT COUNT(*)
        FROM identity.users u
        JOIN identity.user_roles ur ON u.user_id = ur.user_id
        JOIN identity.roles r ON r.role_id = ur.role_id
        WHERE r.name = 'admin';
    `
    // -----------

	deleteUserQuery = `DELETE FROM identity.users WHERE user_id = $1`

	ListUsersQuery = `
    SELECT 
        u.user_id,
        u.username,
        r.name AS role
    FROM identity.users AS u
    JOIN identity.user_roles AS ur ON ur.user_id = u.user_id
    JOIN identity.roles AS r ON r.role_id = ur.role_id
    `
)
