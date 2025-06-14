package repositories

import (
	"database/sql"
	"fmt"
	"forum-go/config"
	"forum-go/models"
	"strings"
)

func CreateUser(user models.User) (int, error) {
	// 1. Définition de la requête SQL d'insertion.
	query := "INSERT INTO `users`(`username`, `email`, `password_hash`, `role`, `profile_picture`, `biography`) VALUES (?,?,?,?,?,?);"

	// Handle nullable fields
	var profilePicture, biography interface{}
	if user.ProfilePicture != nil {
		profilePicture = *user.ProfilePicture
	} else {
		profilePicture = nil
	}
	if user.Biography != nil {
		biography = *user.Biography
	} else {
		biography = nil
	}

	// 2. Exécution de la requête avec les valeurs extraites de l'objet user.
	sqlResult, err := config.DbContext.Exec(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		profilePicture,
		biography,
	)

	if err != nil {
		return -1, fmt.Errorf("Erreur ajout utilisateur - Erreur : %s", err.Error())
	}

	// 3. Récupération de l'identifiant généré (LAST_INSERT_ID).
	id, err := sqlResult.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Erreur ajout utilisateur récupération ID - Erreur : %s", err.Error())
	}

	// 4. Conversion de l'ID en int.
	return int(id), nil
}

func ReadAllUsers() (*[]models.User, error) {
	// 1. Exécution de la requête SELECT
	sqlResult, err := config.DbContext.Query("SELECT * FROM `users`;")
	if err != nil {
		return nil, fmt.Errorf("Erreur récupération liste utilisateurs - Erreur : \n\t %s", err.Error())
	}
	defer sqlResult.Close()

	// 2. Parcours des lignes retournées
	var listUsers []models.User
	for sqlResult.Next() {
		var user models.User
		var lastLogin, bannedAt sql.NullTime
		var profilePicture, biography sql.NullString

		// On scanne chaque colonne dans les champs du struct User
		scanErr := sqlResult.Scan(
			&user.UserId,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.CreatedAt,
			&lastLogin,
			&user.IsBanned,
			&bannedAt,
			&profilePicture,
			&biography,
			&user.MessageCount,
			&user.DiscussionCount,
		)

		if scanErr != nil {
			return nil, fmt.Errorf("Erreur récupération liste utilisateurs - Erreur : \n\t %s", scanErr.Error())
		}

		// Conversion des champs NullTime en *time.Time
		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}
		if bannedAt.Valid {
			user.BannedAt = &bannedAt.Time
		}

		// Conversion des champs NullString en *string
		if profilePicture.Valid {
			user.ProfilePicture = &profilePicture.String
		}
		if biography.Valid {
			user.Biography = &biography.String
		}

		// Ajout à la liste des utilisateurs
		listUsers = append(listUsers, user)
	}

	if sqlResult.Err() != nil {
		return nil, fmt.Errorf("Erreur récupération liste utilisateurs - Erreur : \n\t %s", sqlResult.Err())
	}

	// 3. Retour du slice d'utilisateurs
	return &listUsers, nil
}

func ReadUserById(userId int) (*models.User, error) {
	// 1. Exécution de la requête préparée pour un seul résultat
	sqlResult := config.DbContext.QueryRow(
		"SELECT * FROM `users` WHERE user_id = ?;",
		userId,
	)

	// 2. Scan des colonnes dans un struct User
	var user models.User
	var lastLogin, bannedAt sql.NullTime
	var profilePicture, biography sql.NullString

	errScan := sqlResult.Scan(
		&user.UserId,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&lastLogin,
		&user.IsBanned,
		&bannedAt,
		&profilePicture,
		&biography,
		&user.MessageCount,
		&user.DiscussionCount,
	)

	if errScan != nil {
		// Si aucune ligne, on renvoie (nil, nil) pour indiquer "non trouvé"
		if errScan == sql.ErrNoRows {
			return nil, nil
		}
		// Autre erreur, retour d'une erreur
		return nil, fmt.Errorf("Erreur récupération utilisateur par ID - Erreur : \n\t %s", errScan.Error())
	}

	// Conversion des champs NullTime en *time.Time
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	if bannedAt.Valid {
		user.BannedAt = &bannedAt.Time
	}

	// Conversion des champs NullString en *string
	if profilePicture.Valid {
		user.ProfilePicture = &profilePicture.String
	}
	if biography.Valid {
		user.Biography = &biography.String
	}

	// 3. Retour utilisateur trouvé
	return &user, nil
}

func ReadUserByUsername(username string) (*models.User, error) {
	// 1. Exécution de la requête préparée pour un seul résultat
	sqlResult := config.DbContext.QueryRow(
		"SELECT * FROM `users` WHERE username = ?;",
		username,
	)

	// 2. Scan des colonnes dans un struct User
	var user models.User
	var lastLogin, bannedAt sql.NullTime
	var profilePicture, biography sql.NullString

	errScan := sqlResult.Scan(
		&user.UserId,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&lastLogin,
		&user.IsBanned,
		&bannedAt,
		&profilePicture,
		&biography,
		&user.MessageCount,
		&user.DiscussionCount,
	)

	if errScan != nil {
		// Si aucune ligne, on renvoie (nil, nil) pour indiquer "non trouvé"
		if errScan == sql.ErrNoRows {
			return nil, nil
		}
		// Autre erreur, retour d'une erreur
		return nil, fmt.Errorf("Erreur récupération utilisateur par nom d'utilisateur - Erreur : \n\t %s", errScan.Error())
	}

	// Conversion des champs NullTime en *time.Time
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	if bannedAt.Valid {
		user.BannedAt = &bannedAt.Time
	}

	// Conversion des champs NullString en *string
	if profilePicture.Valid {
		user.ProfilePicture = &profilePicture.String
	}
	if biography.Valid {
		user.Biography = &biography.String
	}

	// 3. Retour utilisateur trouvé
	return &user, nil
}

func ReadUserByEmail(email string) (*models.User, error) {
	// 1. Exécution de la requête préparée pour un seul résultat
	sqlResult := config.DbContext.QueryRow(
		"SELECT * FROM `users` WHERE email = ?;",
		email,
	)

	// 2. Scan des colonnes dans un struct User
	var user models.User
	var lastLogin, bannedAt sql.NullTime
	var profilePicture, biography sql.NullString

	errScan := sqlResult.Scan(
		&user.UserId,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&lastLogin,
		&user.IsBanned,
		&bannedAt,
		&profilePicture,
		&biography,
		&user.MessageCount,
		&user.DiscussionCount,
	)

	if errScan != nil {
		// Si aucune ligne, on renvoie (nil, nil) pour indiquer "non trouvé"
		if errScan == sql.ErrNoRows {
			return nil, nil
		}
		// Autre erreur, retour d'une erreur
		return nil, fmt.Errorf("Erreur récupération utilisateur par email - Erreur : \n\t %s", errScan.Error())
	}

	// Conversion des champs NullTime en *time.Time
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	if bannedAt.Valid {
		user.BannedAt = &bannedAt.Time
	}

	// Conversion des champs NullString en *string
	if profilePicture.Valid {
		user.ProfilePicture = &profilePicture.String
	}
	if biography.Valid {
		user.Biography = &biography.String
	}

	// 3. Retour utilisateur trouvé
	return &user, nil
}

func UpdateUser(user models.User) (bool, error) {
	// 1. Définition de la requête d'UPDATE sur la table `users`
	query := "UPDATE `users` SET `username` = ?, `email` = ?, `role` = ?, `is_banned` = ?, `profile_picture` = ?, `biography` = ? WHERE `user_id` = ?;"

	// Handle nullable fields
	var profilePicture, biography interface{}
	if user.ProfilePicture != nil {
		profilePicture = *user.ProfilePicture
	} else {
		profilePicture = nil
	}
	if user.Biography != nil {
		biography = *user.Biography
	} else {
		biography = nil
	}

	// 2. Exécution de la requête avec les champs de l'utilisateur
	sqlResult, err := config.DbContext.Exec(
		query,
		user.Username,
		user.Email,
		user.Role,
		user.IsBanned,
		profilePicture,
		biography,
		user.UserId,
	)

	if err != nil {
		return false, fmt.Errorf("Erreur modification utilisateur %d - Erreur : \n\t %s", user.UserId, err)
	}

	// 3. Vérification du nombre de lignes affectées
	nbrRows, err := sqlResult.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Erreur modification utilisateur - Erreur : \n\t %s", err)
	}

	if nbrRows == 0 {
		// Aucune ligne n'a été mise à jour
		return false, nil
	}

	// 4. Mise à jour réussie
	return true, nil
}

func UpdateUserPassword(userId int, newPasswordHash string) (bool, error) {
	query := "UPDATE `users` SET `password_hash` = ? WHERE `user_id` = ?;"

	sqlResult, err := config.DbContext.Exec(query, newPasswordHash, userId)
	if err != nil {
		return false, fmt.Errorf("Erreur modification mot de passe - Erreur : \n\t %s", err)
	}

	nbrRows, err := sqlResult.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Erreur modification mot de passe - Erreur : \n\t %s", err)
	}

	return nbrRows > 0, nil
}

func UpdateLastLogin(userId int) error {
	query := "UPDATE `users` SET `last_login` = CURRENT_TIMESTAMP WHERE `user_id` = ?;"

	_, err := config.DbContext.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("Erreur mise à jour dernière connexion - Erreur : \n\t %s", err)
	}

	return nil
}

func BanUser(userId int, banned bool) (bool, error) {
	var query string
	if banned {
		query = "UPDATE `users` SET `is_banned` = TRUE, `banned_at` = CURRENT_TIMESTAMP WHERE `user_id` = ?;"
	} else {
		query = "UPDATE `users` SET `is_banned` = FALSE, `banned_at` = NULL WHERE `user_id` = ?;"
	}

	sqlResult, err := config.DbContext.Exec(query, userId)
	if err != nil {
		return false, fmt.Errorf("Erreur modification statut bannissement - Erreur : \n\t %s", err)
	}

	nbrRows, err := sqlResult.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Erreur modification statut bannissement - Erreur : \n\t %s", err)
	}

	return nbrRows > 0, nil
}

func DeleteUser(userId int) (bool, error) {
	// 1. Préparation et exécution de la requête DELETE
	query := "DELETE FROM `users` WHERE `user_id` = ?;"
	result, err := config.DbContext.Exec(
		query,
		userId,
	)

	if err != nil {
		return false, fmt.Errorf("Erreur suppression utilisateur - Erreur : \n\t %s", err.Error())
	}

	// 2. Récupération du nombre de lignes affectées
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Erreur suppression utilisateur - Erreur : \n\t %s", err.Error())
	}

	if rowsDeleted == 0 {
		// Aucun enregistrement supprimé : ID non existant
		return false, nil
	}

	// 3. Suppression réussie
	return true, nil
}

func IncrementMessageCount(userId int) error {
	query := "UPDATE `users` SET `message_count` = `message_count` + 1 WHERE `user_id` = ?;"

	_, err := config.DbContext.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("Erreur incrémentation compteur de messages - Erreur : \n\t %s", err)
	}

	return nil
}

func IncrementDiscussionCount(userId int) error {
	query := "UPDATE `users` SET `discussion_count` = `discussion_count` + 1 WHERE `user_id` = ?;"

	_, err := config.DbContext.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("Erreur incrémentation compteur de discussions - Erreur : \n\t %s", err)
	}

	return nil
}

// GetUserByEmailorUsername searches for a user by their username or email address
// If the provided string contains '@', it will be treated as an email address
// Otherwise, it will be treated as a username
func GetUserByEmailorUsername(usernameOrEmail string) (*models.User, error) {
	// Check if the input is an email (contains @)
	if strings.Contains(usernameOrEmail, "@") {
		// Treat as email
		return ReadUserByEmail(usernameOrEmail)
	}

	// Treat as username
	return ReadUserByUsername(usernameOrEmail)
}

// ExisteUserByUsernameAndEmail checks if a user exists with the given username or email
// Returns:
// - nil, nil if no user exists with the given username or email
// - user, nil if a user exists with either the given username or email
// - nil, error if there was an error checking for existence
func ExisteUserByUsernameAndEmail(username, email string) (*models.User, error) {
	// First check by username
	user, err := ReadUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// If found by username, return the user
	if user != nil {
		return user, nil
	}

	// If not found by username, check by email
	user, err = ReadUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Return the result (either nil or user found by email)
	return user, nil
}

// GetTotalUsersCount returns the total number of users
func GetTotalUsersCount() (int, error) {
	var count int
	err := config.DbContext.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// GetBannedUsersCount returns the number of banned users
func GetBannedUsersCount() (int, error) {
	var count int
	err := config.DbContext.QueryRow("SELECT COUNT(*) FROM users WHERE is_banned = 1").Scan(&count)
	return count, err
}

// GetRecentUsers retrieves recently registered users
func GetRecentUsers(limit int) ([]models.User, error) {
	query := "SELECT user_id, username, email, password_hash, role, created_at, last_login, is_banned, banned_at, profile_picture, biography, message_count, discussion_count FROM users ORDER BY created_at DESC LIMIT ?"

	rows, err := config.DbContext.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var profilePicture, biography sql.NullString
		err := rows.Scan(
			&user.UserId, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
			&user.CreatedAt, &user.LastLogin, &user.IsBanned, &user.BannedAt,
			&profilePicture, &biography, &user.MessageCount, &user.DiscussionCount,
		)
		if err != nil {
			return nil, err
		}

		// Conversion des champs NullString en *string
		if profilePicture.Valid {
			user.ProfilePicture = &profilePicture.String
		}
		if biography.Valid {
			user.Biography = &biography.String
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUsersForAdmin retrieves users with search and filtering for admin panel
func GetUsersForAdmin(searchQuery string, statusFilter string, limit int, offset int) ([]models.User, error) {
	query := `SELECT user_id, username, email, password_hash, role, created_at, last_login, is_banned, banned_at, profile_picture, biography, message_count, discussion_count 
			  FROM users WHERE 1=1`

	var args []interface{}

	// Add search condition
	if searchQuery != "" {
		query += " AND (username LIKE ? OR email LIKE ?)"
		searchPattern := "%" + searchQuery + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Add status filter
	if statusFilter == "banned" {
		query += " AND is_banned = 1"
	} else if statusFilter == "active" {
		query += " AND is_banned = 0"
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := config.DbContext.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var profilePicture, biography sql.NullString
		err := rows.Scan(
			&user.UserId, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
			&user.CreatedAt, &user.LastLogin, &user.IsBanned, &user.BannedAt,
			&profilePicture, &biography, &user.MessageCount, &user.DiscussionCount,
		)
		if err != nil {
			return nil, err
		}

		// Conversion des champs NullString en *string
		if profilePicture.Valid {
			user.ProfilePicture = &profilePicture.String
		}
		if biography.Valid {
			user.Biography = &biography.String
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUsersCountForAdmin returns the count of users matching admin search criteria
func GetUsersCountForAdmin(searchQuery string, statusFilter string) (int, error) {
	query := "SELECT COUNT(*) FROM users WHERE 1=1"
	var args []interface{}

	// Add search condition
	if searchQuery != "" {
		query += " AND (username LIKE ? OR email LIKE ?)"
		searchPattern := "%" + searchQuery + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Add status filter
	if statusFilter == "banned" {
		query += " AND is_banned = 1"
	} else if statusFilter == "active" {
		query += " AND is_banned = 0"
	}

	var count int
	err := config.DbContext.QueryRow(query, args...).Scan(&count)
	return count, err
}

// ...existing code...
