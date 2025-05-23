CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(128) NOT NULL, -- SHA512 = 128 caractères
    role ENUM('user', 'admin') DEFAULT 'user',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME,
    is_banned BOOLEAN DEFAULT FALSE,
    banned_at DATETIME NULL,
    profile_picture VARCHAR(255),
    biography TEXT,
    message_count INT DEFAULT 0,
    discussion_count INT DEFAULT 0
);

-- Table CATEGORIES
CREATE TABLE categories (
    category_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(7) DEFAULT '#007bff', -- Code couleur hexadécimal
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table DISCUSSIONS
CREATE TABLE discussions (
    discussion_id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status ENUM('open', 'closed', 'archived') DEFAULT 'open',
    visibility ENUM('public', 'private') DEFAULT 'public',
    author_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    view_count INT DEFAULT 0,
    message_count INT DEFAULT 0,
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Table MESSAGES
CREATE TABLE messages (
    message_id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    discussion_id INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    has_image BOOLEAN DEFAULT FALSE,
    image_path VARCHAR(255),
    like_count INT DEFAULT 0,
    dislike_count INT DEFAULT 0,
    popularity_score INT DEFAULT 0, -- like_count - dislike_count
    FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (discussion_id) REFERENCES discussions(discussion_id) ON DELETE CASCADE
);

-- Table MESSAGE_REACTIONS
CREATE TABLE message_reactions (
    reaction_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    message_id INT NOT NULL,
    reaction_type ENUM('like', 'dislike') NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_message_reaction (user_id, message_id)
);

-- Table IMAGES
CREATE TABLE images (
    image_id INT AUTO_INCREMENT PRIMARY KEY,
    original_name VARCHAR(255) NOT NULL,
    stored_path VARCHAR(255) NOT NULL,
    file_size INT,
    mime_type VARCHAR(50),
    message_id INT NOT NULL,
    uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (message_id) REFERENCES messages(message_id) ON DELETE CASCADE
);

-- Table FRIENDSHIPS
CREATE TABLE friendships (
    friendship_id INT AUTO_INCREMENT PRIMARY KEY,
    requester_id INT NOT NULL,
    addressee_id INT NOT NULL,
    status ENUM('pending', 'accepted', 'rejected', 'blocked') DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (requester_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (addressee_id) REFERENCES users(user_id) ON DELETE CASCADE,
    UNIQUE KEY unique_friendship (requester_id, addressee_id),
    CHECK (requester_id != addressee_id)
);

-- Table USER_SESSIONS
CREATE TABLE user_sessions (
    session_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Table ADMIN_ACTIONS
CREATE TABLE admin_actions (
    action_id INT AUTO_INCREMENT PRIMARY KEY,
    admin_id INT NOT NULL,
    action_type ENUM('delete_discussion', 'delete_message', 'ban_user', 'change_discussion_status', 'unban_user') NOT NULL,
    target_user_id INT NULL,
    target_discussion_id INT NULL,
    target_message_id INT NULL,
    reason TEXT,
    performed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (target_user_id) REFERENCES users(user_id) ON DELETE SET NULL,
    FOREIGN KEY (target_discussion_id) REFERENCES discussions(discussion_id) ON DELETE SET NULL,
    FOREIGN KEY (target_message_id) REFERENCES messages(message_id) ON DELETE SET NULL
);

-- ========================================
-- TABLES DE LIAISON (Relations N:N)
-- ========================================

-- Table DISCUSSION_CATEGORIES (Liaison Discussion-Catégorie)
CREATE TABLE discussion_categories (
    discussion_id INT,
    category_id INT,
    PRIMARY KEY (discussion_id, category_id),
    FOREIGN KEY (discussion_id) REFERENCES discussions(discussion_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);

-- Table DISCUSSION_ACCESS (Accès aux discussions privées)
CREATE TABLE discussion_access (
    access_id INT AUTO_INCREMENT PRIMARY KEY,
    discussion_id INT NOT NULL,
    user_id INT NOT NULL,
    granted_by INT NOT NULL,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (discussion_id) REFERENCES discussions(discussion_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(user_id) ON DELETE CASCADE,
    UNIQUE KEY unique_access (discussion_id, user_id)
);

-- ========================================
-- INDEX POUR OPTIMISATION
-- ========================================

-- Index pour les recherches fréquentes
CREATE INDEX idx_discussions_author ON discussions(author_id);
CREATE INDEX idx_discussions_status ON discussions(status);
CREATE INDEX idx_discussions_visibility ON discussions(visibility);
CREATE INDEX idx_discussions_created_at ON discussions(created_at DESC);

CREATE INDEX idx_messages_discussion ON messages(discussion_id);
CREATE INDEX idx_messages_author ON messages(author_id);
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);
CREATE INDEX idx_messages_popularity ON messages(popularity_score DESC);

CREATE INDEX idx_message_reactions_message ON message_reactions(message_id);
CREATE INDEX idx_message_reactions_user ON message_reactions(user_id);

CREATE INDEX idx_user_sessions_token ON user_sessions(token);
CREATE INDEX idx_user_sessions_expires ON user_sessions(expires_at);

CREATE INDEX idx_friendships_requester ON friendships(requester_id);
CREATE INDEX idx_friendships_addressee ON friendships(addressee_id);
CREATE INDEX idx_friendships_status ON friendships(status);

-- ========================================
-- TRIGGERS POUR AUTOMATISATION
-- ========================================

-- Trigger pour mettre à jour le score de popularité
DELIMITER $$
CREATE TRIGGER update_message_popularity
AFTER INSERT ON message_reactions
FOR EACH ROW
BEGIN
    UPDATE messages 
    SET 
        like_count = (SELECT COUNT(*) FROM message_reactions WHERE message_id = NEW.message_id AND reaction_type = 'like'),
        dislike_count = (SELECT COUNT(*) FROM message_reactions WHERE message_id = NEW.message_id AND reaction_type = 'dislike'),
        popularity_score = (SELECT COUNT(*) FROM message_reactions WHERE message_id = NEW.message_id AND reaction_type = 'like') - 
                          (SELECT COUNT(*) FROM message_reactions WHERE message_id = NEW.message_id AND reaction_type = 'dislike')
    WHERE message_id = NEW.message_id;
END$$

CREATE TRIGGER update_message_popularity_delete
AFTER DELETE ON message_reactions
FOR EACH ROW
BEGIN
    UPDATE messages 
    SET 
        like_count = (SELECT COUNT(*) FROM message_reactions WHERE message_id = OLD.message_id AND reaction_type = 'like'),
        dislike_count = (SELECT COUNT(*) FROM message_reactions WHERE message_id = OLD.message_id AND reaction_type = 'dislike'),
        popularity_score = (SELECT COUNT(*) FROM message_reactions WHERE message_id = OLD.message_id AND reaction_type = 'like') - 
                          (SELECT COUNT(*) FROM message_reactions WHERE message_id = OLD.message_id AND reaction_type = 'dislike')
    WHERE message_id = OLD.message_id;
END$$

-- Trigger pour mettre à jour le compteur de messages des discussions
CREATE TRIGGER update_discussion_message_count_insert
AFTER INSERT ON messages
FOR EACH ROW
BEGIN
    UPDATE discussions 
    SET message_count = message_count + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE discussion_id = NEW.discussion_id;
    
    UPDATE users 
    SET message_count = message_count + 1
    WHERE user_id = NEW.author_id;
END$$

CREATE TRIGGER update_discussion_message_count_delete
AFTER DELETE ON messages
FOR EACH ROW
BEGIN
    UPDATE discussions 
    SET message_count = message_count - 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE discussion_id = OLD.discussion_id;
    
    UPDATE users 
    SET message_count = message_count - 1
    WHERE user_id = OLD.author_id;
END$$

-- Trigger pour mettre à jour le compteur de discussions des utilisateurs
CREATE TRIGGER update_user_discussion_count_insert
AFTER INSERT ON discussions
FOR EACH ROW
BEGIN
    UPDATE users 
    SET discussion_count = discussion_count + 1
    WHERE user_id = NEW.author_id;
END$$

CREATE TRIGGER update_user_discussion_count_delete
AFTER DELETE ON discussions
FOR EACH ROW
BEGIN
    UPDATE users 
    SET discussion_count = discussion_count - 1
    WHERE user_id = OLD.author_id;
END$$

DELIMITER ;

-- ========================================
-- DONNÉES DE TEST
-- ========================================

-- Insertion des catégories par défaut
INSERT INTO categories (name, description, color) VALUES
('Général', 'Discussions générales', '#007bff'),
('Aide', 'Questions et aide', '#28a745'),
('Annonces', 'Annonces importantes', '#dc3545'),
('Développement', 'Discussions sur le développement', '#6f42c1'),
('Off-topic', 'Discussions hors sujet', '#6c757d');

-- Insertion d'un utilisateur admin par défaut
INSERT INTO users (username, email, password_hash, role, biography) VALUES
('admin', 'admin@forum.com', SHA2('AdminPassword123!', 512), 'admin', 'Administrateur du forum');

-- Insertion d'utilisateurs de test
INSERT INTO users (username, email, password_hash, biography) VALUES
('johndoe', 'john@example.com', SHA2('MyPassword123!', 512), 'Développeur passionné'),
('janedoe', 'jane@example.com', SHA2('SecurePass123!', 512), 'Designer UI/UX'),
('testuser', 'test@example.com', SHA2('TestPassword123!', 512), 'Utilisateur de test');

-- Insertion de discussions de test
INSERT INTO discussions (title, description, author_id) VALUES
('Bienvenue sur le forum', 'Discussion de bienvenue pour tous les nouveaux membres', 1),
('Questions fréquentes', 'Retrouvez ici les réponses aux questions les plus courantes', 1),
('Suggestions d''amélioration', 'Proposez vos idées pour améliorer le forum', 2);

-- Liaison discussions-catégories
INSERT INTO discussion_categories (discussion_id, category_id) VALUES
(1, 1), -- Bienvenue -> Général
(1, 3), -- Bienvenue -> Annonces
(2, 2), -- FAQ -> Aide
(3, 1); -- Suggestions -> Général

-- Messages de test
INSERT INTO messages (content, author_id, discussion_id) VALUES
('Bienvenue à tous sur notre nouveau forum ! N''hésitez pas à vous présenter.', 1, 1),
('Merci pour l''accueil ! Hâte de participer aux discussions.', 2, 1),
('Comment puis-je modifier mon profil ?', 3, 2),
('Il serait bien d''avoir un mode sombre pour le forum.', 2, 3);

-- Réactions de test
INSERT INTO message_reactions (user_id, message_id, reaction_type) VALUES
(2, 1, 'like'),
(3, 1, 'like'),
(1, 4, 'like'),
(3, 4, 'like');

-- ========================================
-- VUES UTILES
-- ========================================

-- Vue pour les discussions avec informations complètes
CREATE VIEW v_discussions_complete AS
SELECT 
    d.discussion_id,
    d.title,
    d.description,
    d.status,
    d.visibility,
    d.created_at,
    d.updated_at,
    d.view_count,
    d.message_count,
    u.username as author_username,
    u.profile_picture as author_picture,
    GROUP_CONCAT(c.name) as categories
FROM discussions d
JOIN users u ON d.author_id = u.user_id
LEFT JOIN discussion_categories dc ON d.discussion_id = dc.discussion_id
LEFT JOIN categories c ON dc.category_id = c.category_id
WHERE d.status != 'archived'
GROUP BY d.discussion_id;

-- Vue pour les messages avec informations complètes
CREATE VIEW v_messages_complete AS
SELECT 
    m.message_id,
    m.content,
    m.created_at,
    m.updated_at,
    m.has_image,
    m.image_path,
    m.like_count,
    m.dislike_count,
    m.popularity_score,
    u.username as author_username,
    u.profile_picture as author_picture,
    d.title as discussion_title
FROM messages m
JOIN users u ON m.author_id = u.user_id
JOIN discussions d ON m.discussion_id = d.discussion_id;

-- Vue pour les statistiques utilisateurs
CREATE VIEW v_user_stats AS
SELECT 
    u.user_id,
    u.username,
    u.email,
    u.role,
    u.created_at,
    u.last_login,
    u.is_banned,
    u.message_count,
    u.discussion_count,
    COUNT(DISTINCT f1.friendship_id) + COUNT(DISTINCT f2.friendship_id) as friend_count
FROM users u
LEFT JOIN friendships f1 ON u.user_id = f1.requester_id AND f1.status = 'accepted'
LEFT JOIN friendships f2 ON u.user_id = f2.addressee_id AND f2.status = 'accepted'
GROUP BY u.user_id;

-- ========================================
-- PROCÉDURES STOCKÉES UTILES
-- ========================================

DELIMITER $$

-- Procédure pour bannir un utilisateur
CREATE PROCEDURE BanUser(
    IN p_admin_id INT,
    IN p_user_id INT,
    IN p_reason TEXT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- Bannir l'utilisateur
    UPDATE users 
    SET is_banned = TRUE, banned_at = CURRENT_TIMESTAMP 
    WHERE user_id = p_user_id;
    
    -- Enregistrer l'action admin
    INSERT INTO admin_actions (admin_id, action_type, target_user_id, reason)
    VALUES (p_admin_id, 'ban_user', p_user_id, p_reason);
    
    -- Supprimer les sessions actives
    DELETE FROM user_sessions WHERE user_id = p_user_id;
    
    COMMIT;
END$$

-- Procédure pour supprimer une discussion et ses dépendances
CREATE PROCEDURE DeleteDiscussion(
    IN p_admin_id INT,
    IN p_discussion_id INT,
    IN p_reason TEXT
)
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- Enregistrer l'action admin
    INSERT INTO admin_actions (admin_id, action_type, target_discussion_id, reason)
    VALUES (p_admin_id, 'delete_discussion', p_discussion_id, p_reason);
    
    -- Supprimer la discussion (CASCADE supprimera les messages et dépendances)
    DELETE FROM discussions WHERE discussion_id = p_discussion_id;
    
    COMMIT;
END$$

DELIMITER ;

-- Message de confirmation
SELECT 'Base de données du forum créée avec succès!' as status;