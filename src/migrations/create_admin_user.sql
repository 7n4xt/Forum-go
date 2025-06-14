-- Script to create an admin user account
-- Username: admin
-- Password: admin123
-- This script creates a hashed password using SHA512 (same as the Go application)

-- First, check if the admin user already exists and delete if necessary
DELETE FROM users WHERE username = 'admin' OR email = 'admin@forum.local';

-- Insert the admin user with SHA512 hashed password
-- Password 'admin123' -> SHA512 hash: '7fcf4ba391c48784edde599889d6e3f1e47a27db36ecc050cc92f259bfac38afad2c68a1ae804d77075e8fb722503f3eca2b2c1006ee6f6c7b7628cb45fffd1d'
INSERT INTO users (
    username, 
    email, 
    password_hash, 
    role, 
    created_at,
    is_banned,
    message_count,
    discussion_count
) VALUES (
    'admin',
    'admin@forum.local',
    '7fcf4ba391c48784edde599889d6e3f1e47a27db36ecc050cc92f259bfac38afad2c68a1ae804d77075e8fb722503f3eca2b2c1006ee6f6c7b7628cb45fffd1d',
    'admin',
    NOW(),
    FALSE,
    0,
    0
);

-- Verify the admin user was created successfully
SELECT 
    user_id,
    username,
    email,
    role,
    created_at,
    is_banned
FROM users 
WHERE username = 'admin';

-- Display confirmation message
SELECT 'Admin user created successfully!' AS message,
       'Username: admin' AS username,
       'Password: admin123' AS password,
       'Role: admin' AS role;
