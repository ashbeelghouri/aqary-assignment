-- name: CreateUser :one
INSERT INTO users (name, phone_number)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone_number = $1
LIMIT 1;

-- name: UpdateUserOTP :one
UPDATE users
SET otp = $2,
    otp_expiration_time  = NOW() + INTERVAL '1 minute'
WHERE phone_number = $1
RETURNING *;