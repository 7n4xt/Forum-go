# Forum Go - Advanced Discussion Platform

## 📋 Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Project Structure](#project-structure)
- [API Routes](#api-routes)
- [Database Schema](#database-schema)
- [Team Composition](#team-composition)
- [Project Management](#project-management)

## 🌟 Overview

Forum Go is a modern, feature-rich discussion platform built with Go, MySQL, and vanilla HTML/CSS/JavaScript. It provides a comprehensive forum solution with advanced features like pagination, search, categorization, and administrative capabilities.

## ✨ Features

### Core Features (FT-1 to FT-7)
- **User Authentication**: Registration, login, logout with secure session management
- **Discussion Management**: Create, view, edit, and delete discussions
- **Message System**: Post, edit, delete messages within discussions
- **Category System**: Organize discussions by categories/tags
- **Reaction System**: Like/dislike messages with popularity scoring
- **User Roles**: Regular users and administrators
- **Discussion Status**: Open, closed, and archived states

### Advanced Features (FT-8 to FT-12)
- **Message Sorting**: Chronological order or popularity-based sorting
- **Pagination System**: Configurable pagination (10, 20, 30, or all items per page)
- **Category Filtering**: Browse discussions by specific categories/tags
- **Search Functionality**: Search discussions by title, description, or tags
- **Admin Dashboard**: Comprehensive platform management interface
  - User management (ban/unban users)
  - Discussion moderation (change status, delete discussions)
  - Message moderation (delete inappropriate messages)
  - Platform statistics and analytics

### Technical Features
- **Responsive Design**: Stack Overflow-inspired modern UI
- **Database Integration**: MySQL with proper foreign key relationships
- **Security**: Password hashing, SQL injection prevention, session management
- **Template System**: Go template engine with reusable components
- **RESTful Architecture**: Clean separation of concerns (MVC pattern)

## 🚀 Installation

### Prerequisites
- Go 1.19 or higher
- MySQL 8.0 or higher
- Git

### Setup Instructions

1. **Clone the repository**
   ```bash
   git clone https://github.com/7n4xt/Forum-go.git
   cd Forum-go
   ```

2. **Install dependencies**
   ```bash
   cd src
   go mod download
   ```

3. **Database Setup**
   ```bash
   # Create MySQL database
   mysql -u root -p
   CREATE DATABASE forum;
   
   # Import schema
   mysql -u root -p forum < migrations/migration.sql
   ```

4. **Environment Configuration**
   ```bash
   # Create .env file in src/ directory
   cp .env.example .env
   
   # Edit .env with your database credentials
   DB_NAME = "forum"
   DB_USER = "root"
   DB_PWD = "your_password"
   DB_HOST = "127.0.0.1"
   DB_PORT = 3306
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

6. **Access the application**
   - Open your browser and navigate to `http://localhost:8080`
   - Default admin credentials: username: `admin`, password: `admin123`

## 📁 Project Structure

```
Forum-go/
├── README.md
├── TASK_DISTRIBUTION.md
├── todo.md
└── src/
    ├── main.go                          # Application entry point
    ├── go.mod                           # Go module dependencies
    ├── go.sum                           # Dependency checksums
    ├── .env                             # Environment variables
    │
    ├── config/                          # Configuration files
    │   ├── db.config.go                 # Database configuration
    │   └── env.config.go                # Environment configuration
    │
    ├── models/                          # Data models
    │   └── models.go                    # User, Discussion, Message, Category models
    │
    ├── repositories/                    # Data access layer
    │   ├── users.repositorie.go         # User CRUD operations
    │   ├── discussion_repository.go     # Discussion CRUD operations
    │   ├── message_repository.go        # Message CRUD operations
    │   ├── reaction_repository.go       # Reaction system operations
    │   └── session_repository.go        # Session management
    │
    ├── services/                        # Business logic layer
    │   ├── auth.services.go             # Authentication services
    │   ├── discussion.services.go       # Discussion business logic
    │   ├── message.services.go          # Message business logic
    │   ├── admin.services.go            # Admin dashboard services
    │   ├── reaction.services.go         # Reaction system logic
    │   └── main.service.go              # Core services
    │
    ├── controllers/                     # HTTP request handlers
    │   ├── auth.controller.go           # Authentication endpoints
    │   ├── auth.controllers.go          # Additional auth handlers
    │   ├── discussion.controllers.go    # Discussion endpoints
    │   ├── message.controllers.go       # Message endpoints
    │   ├── reaction.controllers.go      # Reaction endpoints
    │   └── admin.controllers.go         # Admin dashboard endpoints
    │
    ├── routes/                          # Route definitions
    │   ├── auth.routes.go               # Authentication routes
    │   ├── discussion.routes.go         # Discussion routes
    │   ├── message.route.go             # Message routes
    │   └── admin.routes.go              # Admin routes
    │
    ├── templates/                       # HTML templates
    │   ├── config.templates.go          # Template configuration
    │   ├── index.html                   # Home page
    │   ├── login.html                   # Login page
    │   ├── register.html                # Registration page
    │   ├── discussions.template.html    # Discussions listing
    │   ├── discussion.template.html     # Single discussion view
    │   ├── new_discussion.template.html # Create discussion form
    │   ├── admin_dashboard.html         # Admin dashboard
    │   ├── admin_discussions.html       # Admin discussions management
    │   └── admin_users.html             # Admin users management
    │
    ├── public/                          # Static assets
    │   ├── css/
    │   │   ├── styles.css               # Main stylesheet
    │   │   └── icons.css                # Icon definitions
    │   └── images/
    │       └── image.png                # Static images
    │
    ├── utils/                           # Utility functions
    │   ├── auth.utils.go                # Authentication utilities
    │   └── hash.utils.go                # Password hashing utilities
    │
    └── migrations/                      # Database migrations
        └── migration.sql                # Initial database schema
```

## 🛣️ API Routes

### View Routes (Template Rendering)

| Method | Route | Description | Auth Required |
|--------|-------|-------------|---------------|
| GET | `/` | Home page | No |
| GET | `/login` | Login page | No |
| GET | `/register` | Registration page | No |
| GET | `/discussions` | Discussions listing with search/filter | No |
| GET | `/discussions/{id}` | Single discussion view | No |
| GET | `/discussions/new` | Create discussion form | Yes |
| GET | `/admin` | Admin dashboard | Admin |
| GET | `/admin/discussions` | Admin discussions management | Admin |
| GET | `/admin/users` | Admin users management | Admin |

### Data Processing Routes (API Endpoints)

#### Authentication
| Method | Route | Description | Auth Required |
|--------|-------|-------------|---------------|
| POST | `/register` | User registration | No |
| POST | `/login` | User login | No |
| POST | `/logout` | User logout | Yes |

#### Discussions
| Method | Route | Description | Auth Required |
|--------|-------|-------------|---------------|
| POST | `/discussions/create` | Create new discussion | Yes |
| POST | `/discussions/update` | Update discussion content | Yes (Owner/Admin) |
| POST | `/discussions/status` | Update discussion status | Yes (Owner/Admin) |
| POST | `/discussions/delete` | Delete discussion | Yes (Owner/Admin) |
| POST | `/discussions/delete/{id}` | Delete specific discussion | Yes (Owner/Admin) |

#### Messages
| Method | Route | Description | Auth Required |
|--------|-------|-------------|---------------|
| GET | `/messages` | Get all messages | No |
| GET | `/messages/view/{id}` | Get specific message | No |
| POST | `/messages/create` | Create new message | Yes |
| POST | `/messages/create/{discussionId}` | Create message in discussion | Yes |
| POST | `/messages/update/{id}` | Update message | Yes (Owner/Admin) |
| POST | `/messages/delete/{id}` | Delete message | Yes (Owner/Admin) |

#### Reactions
| Method | Route | Description | Auth Required |
|--------|-------|-------------|---------------|
| POST | `/messages/react/{id}/like` | Like a message | Yes |
| POST | `/messages/react/{id}/dislike` | Dislike a message | Yes |

#### Admin Routes
| Method | Route | Description | Auth Required |
|--------|-------|-------------|---------------|
| POST | `/admin/users/ban` | Ban/unban user | Admin |
| POST | `/admin/discussions/status` | Update discussion status | Admin |
| POST | `/admin/discussions/delete` | Delete discussion | Admin |
| POST | `/admin/discussions/delete/{id}` | Delete specific discussion | Admin |
| POST | `/admin/messages/delete` | Delete message | Admin |
| POST | `/admin/messages/delete/{id}` | Delete specific message | Admin |

### Query Parameters

#### Discussions Listing (`/discussions`)
- `search`: Search term for title/description/tags
- `category`: Filter by category ID
- `status`: Filter by status (open/closed/archived)
- `author`: Filter by author ID
- `limit`: Items per page (10/20/30/1000)
- `page`: Page number

#### Discussion View (`/discussions/{id}`)
- `sort`: Sort messages (newest/chronological/popularity)
- `limit`: Messages per page (10/20/30/1000)
- `page`: Page number

## 🗄️ Database Schema

### Tables Overview
- **users**: User accounts and profiles
- **categories**: Discussion categories/tags
- **discussions**: Main discussion threads
- **messages**: Individual messages within discussions
- **discussion_categories**: Many-to-many relationship between discussions and categories
- **message_reactions**: User reactions (likes/dislikes) to messages
- **sessions**: User session management

### Key Relationships
- Users can create multiple discussions and messages
- Discussions belong to categories (many-to-many)
- Messages belong to discussions (one-to-many)
- Users can react to messages (many-to-many)
- Sessions track user authentication

## 👥 Team Composition

**Team Lead & Full-Stack Developer**: Malek
- **Role**: Project architecture, backend development, frontend integration, database design
- **Responsibilities**: 
  - Complete implementation of all features (FT-1 through FT-12)
  - Database schema design and optimization
  - Security implementation
  - UI/UX design and implementation
  - Admin dashboard development
  - Testing and deployment

## 📊 Project Management

### Project Decomposition and Phases

#### **Phase 1: Foundation (FT-1 to FT-3)**
- **Duration**: Week 1-2
- **Scope**: Core infrastructure and basic functionality
- **Deliverables**:
  - Database schema design and implementation
  - User authentication system
  - Basic discussion and message CRUD operations
  - Template system setup

#### **Phase 2: Core Features (FT-4 to FT-7)**
- **Duration**: Week 3-4
- **Scope**: Essential forum features
- **Deliverables**:
  - Category system implementation
  - Reaction system (likes/dislikes)
  - User role management
  - Discussion status management
  - Basic UI implementation

#### **Phase 3: Advanced Features (FT-8 to FT-11)**
- **Duration**: Week 5-6
- **Scope**: Enhanced user experience
- **Deliverables**:
  - Message sorting and pagination
  - Search functionality
  - Category filtering
  - Responsive UI design
  - Performance optimization

#### **Phase 4: Administration (FT-12)**
- **Duration**: Week 7
- **Scope**: Platform management
- **Deliverables**:
  - Complete admin dashboard
  - User management system
  - Content moderation tools
  - Analytics and statistics

#### **Phase 5: Polish and Documentation**
- **Duration**: Week 8
- **Scope**: Final touches and documentation
- **Deliverables**:
  - Bug fixes and optimization
  - Complete documentation
  - README and project reports
  - Final testing

### Task Distribution Strategy

As a solo developer, the strategy focused on:

1. **Feature-Based Development**: Each feature was developed as a complete vertical slice
2. **MVC Architecture**: Consistent separation of concerns across all features
3. **Iterative Approach**: Build basic functionality first, then enhance with advanced features
4. **Database-First Design**: Strong foundation with proper schema design
5. **Security-First Mindset**: Security considerations integrated from the beginning

### Time Management and Priorities

#### **Priority Classification**:
1. **Critical (P0)**: Authentication, basic CRUD, database security
2. **High (P1)**: Search, pagination, admin dashboard
3. **Medium (P2)**: Advanced UI features, analytics
4. **Low (P3)**: UI polish, advanced admin features

#### **Time Allocation**:
- **Backend Development**: 60% (Database, services, controllers)
- **Frontend Development**: 25% (Templates, CSS, JavaScript)
- **Testing & Debugging**: 10%
- **Documentation**: 5%

### Documentation Strategy

1. **Code Documentation**:
   - Comprehensive comments in French and English
   - Function and method documentation
   - API endpoint documentation

2. **Architecture Documentation**:
   - Database schema documentation
   - Route documentation
   - Project structure explanation

3. **User Documentation**:
   - Installation guide
   - Feature overview
   - Admin user guide

4. **Project Management Documentation**:
   - Task distribution tracking
   - Progress reports
   - Technical decision log

### Technical Decisions and Challenges

#### **Key Technical Decisions**:
1. **Go Standard Library**: Minimized external dependencies for better control
2. **Template Engine**: Used Go's built-in template system for server-side rendering
3. **Database Design**: Normalized schema with proper foreign key relationships
4. **Security**: Implemented SHA512 password hashing and session-based authentication
5. **UI Framework**: Vanilla CSS with Stack Overflow-inspired design

#### **Major Challenges Solved**:
1. **Nullable Database Fields**: Implemented proper handling of NULL values in Go structs
2. **Template Function Integration**: Added mathematical functions for pagination
3. **Admin Permission System**: Role-based access control throughout the application
4. **Search Performance**: Optimized search queries with proper indexing
5. **Responsive Design**: Mobile-friendly interface without external frameworks

### Success Metrics

- ✅ All 12 features (FT-1 through FT-12) successfully implemented
- ✅ Complete admin dashboard with full moderation capabilities
- ✅ Responsive design compatible with all screen sizes
- ✅ Secure authentication and session management
- ✅ Comprehensive search and filtering system
- ✅ Performance optimized for concurrent users
- ✅ Complete documentation and installation guide

---

## 🚀 Getting Started

1. Follow the [Installation](#installation) instructions
2. Create your admin account or use the default credentials
3. Start creating discussions and engaging with the community
4. Access the admin dashboard at `/admin` for platform management

For any issues or questions, please refer to the documentation or create an issue in the repository.
