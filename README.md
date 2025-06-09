# Forum Go

A simple forum application built with Go, HTML, and CSS.

## Features

- User authentication (signup, login, logout)
- Discussion threads with categories
- Messages within discussions
- Thread status management (open, closed, archived)
- User profiles

## Prerequisites

- Go 1.16 or higher
- MySQL 5.7 or higher

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/forum-go.git
cd forum-go
```

2. Set up the database:
```bash
# Create a new MySQL database
mysql -u root -p
CREATE DATABASE forum_go;
USE forum_go;
exit;

# Run migrations
mysql -u root -p forum_go < src/migrations/01_initial.sql
mysql -u root -p forum_go < src/migrations/02_discussions.sql
mysql -u root -p forum_go < src/migrations/03_sessions.sql
```

3. Configure environment variables:
```bash
# Create a .env file in the src directory
cd src
touch .env

# Add the following content to the .env file
DB_USER=your_mysql_username
DB_PWD=your_mysql_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=forum_go
```

4. Build and run the application:
```bash
cd src
go build -o forum-go
./forum-go
```

5. Access the application at http://localhost:8080

## Default Admin Account

- Username: admin
- Password: admin123

## Project Structure

- `src/` - Source code
  - `config/` - Configuration files
  - `controllers/` - HTTP request handlers
  - `migrations/` - SQL migration files
  - `models/` - Data models
  - `public/` - Static assets (CSS, JS, images)
  - `repositories/` - Database access layer
  - `routes/` - HTTP routes
  - `services/` - Business logic
  - `templates/` - HTML templates
  - `utils/` - Utility functions
  - `main.go` - Application entry point

## License

This project is licensed under the MIT License - see the LICENSE file for details. 