package services

import (
	"forum-go/models"
	"forum-go/repositories"
)

// AdminStats represents admin dashboard statistics
type AdminStats struct {
	TotalUsers          int
	TotalDiscussions    int
	TotalMessages       int
	ActiveUsers         int
	BannedUsers         int
	OpenDiscussions     int
	ClosedDiscussions   int
	ArchivedDiscussions int
}

// GetAdminStatsService retrieves statistics for the admin dashboard
func GetAdminStatsService() (AdminStats, error) {
	var stats AdminStats
	var err error

	// Get total users
	stats.TotalUsers, err = repositories.GetTotalUsersCount()
	if err != nil {
		return stats, err
	}

	// Get banned users count
	stats.BannedUsers, err = repositories.GetBannedUsersCount()
	if err != nil {
		return stats, err
	}

	// Calculate active users
	stats.ActiveUsers = stats.TotalUsers - stats.BannedUsers

	// Get discussion counts by status
	stats.OpenDiscussions, err = repositories.GetDiscussionCountByStatus("open")
	if err != nil {
		return stats, err
	}

	stats.ClosedDiscussions, err = repositories.GetDiscussionCountByStatus("closed")
	if err != nil {
		return stats, err
	}

	stats.ArchivedDiscussions, err = repositories.GetDiscussionCountByStatus("archived")
	if err != nil {
		return stats, err
	}

	// Calculate total discussions
	stats.TotalDiscussions = stats.OpenDiscussions + stats.ClosedDiscussions + stats.ArchivedDiscussions

	// Get total messages count
	stats.TotalMessages, err = repositories.GetTotalMessagesCount()
	if err != nil {
		return stats, err
	}

	return stats, nil
}

// GetRecentMessagesService retrieves recent messages for admin review
func GetRecentMessagesService(limit int) ([]models.Message, error) {
	messages, err := repositories.GetRecentMessages(limit)
	if err != nil {
		return nil, err
	}

	// If no messages, return empty slice instead of nil
	if messages == nil {
		return []models.Message{}, nil
	}

	// Get author information for each message
	result := *messages
	for i := range result {
		author, err := repositories.ReadUserById(result[i].AuthorId)
		if err == nil && author != nil {
			result[i].Author = author
		}
	}

	return result, nil
}

// GetRecentUsersService retrieves recently registered users
func GetRecentUsersService(limit int) ([]models.User, error) {
	return repositories.GetRecentUsers(limit)
}

// GetUsersForAdminService retrieves users with filtering for admin management
func GetUsersForAdminService(searchQuery string, statusFilter string, limit int, offset int) ([]models.User, int, error) {
	// Get total count first
	totalCount, err := repositories.GetUsersCountForAdmin(searchQuery, statusFilter)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated users
	users, err := repositories.GetUsersForAdmin(searchQuery, statusFilter, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// BanUserService bans a user account
func BanUserService(userId int) error {
	_, err := repositories.BanUser(userId, true)
	return err
}

// UnbanUserService unbans a user account
func UnbanUserService(userId int) error {
	_, err := repositories.BanUser(userId, false)
	return err
}
