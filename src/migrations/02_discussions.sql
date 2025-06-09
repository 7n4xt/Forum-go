-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    category_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(7) DEFAULT '#1877f2',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create discussions table
CREATE TABLE IF NOT EXISTS discussions (
    discussion_id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status ENUM('open', 'closed', 'archived') DEFAULT 'open',
    visibility ENUM('public', 'private') DEFAULT 'public',
    author_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    view_count INT DEFAULT 0,
    message_count INT DEFAULT 0,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Create discussion_categories table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS discussion_categories (
    discussion_id INT,
    category_id INT,
    PRIMARY KEY (discussion_id, category_id),
    FOREIGN KEY (discussion_id) REFERENCES discussions(discussion_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);

-- Add discussion_id column to messages table if it doesn't exist
ALTER TABLE messages ADD COLUMN IF NOT EXISTS discussion_id INT;
ALTER TABLE messages ADD CONSTRAINT fk_message_discussion FOREIGN KEY (discussion_id) REFERENCES discussions(discussion_id) ON DELETE CASCADE;

-- Insert some default categories
INSERT INTO categories (name, description, color) VALUES
('General', 'General discussions about any topic', '#1877f2'),
('Help', 'Get help with problems or questions', '#31a24c'),
('Announcements', 'Official announcements and news', '#f7b928'),
('Tech', 'Discussions about technology and programming', '#8c44db'),
('Off-topic', 'Discussions that don\'t fit in other categories', '#65676b')
ON DUPLICATE KEY UPDATE name = name; 