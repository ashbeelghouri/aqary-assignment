CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) not NULL,
    phone_number VARCHAR(20) UNIQUE not NULL,
    otp VARCHAR(6),
    otp_expiration_time TIMESTAMP
);