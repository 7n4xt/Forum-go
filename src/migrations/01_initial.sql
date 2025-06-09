-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('user', 'admin') DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    is_banned BOOLEAN DEFAULT FALSE,
    banned_at TIMESTAMP NULL,
    profile_picture VARCHAR(255) DEFAULT '',
    biography TEXT,
    message_count INT DEFAULT 0,
    discussion_count INT DEFAULT 0
);

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    message_id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    discussion_id INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    has_image BOOLEAN DEFAULT FALSE,
    image_path VARCHAR(255) DEFAULT '',
    like_count INT DEFAULT 0,
    dislike_count INT DEFAULT 0,
    popularity_score INT DEFAULT 0,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Insert admin user
INSERT INTO users (username, email, password_hash, role, created_at)
VALUES (
    'admin',
    'admin@example.com',
    -- SHA512 hash of 'admin123'
    'c7ad44cbad762a5da0a452f9e854fdc1e0e7a52a38015f23f3eab1d80b931dd472634dfac71cd34ebc35d16ab7fb8a90c81f975113d6c7538dc69dd8de9077ec',
    'admin',
    NOW()
) ON DUPLICATE KEY UPDATE username = username; 