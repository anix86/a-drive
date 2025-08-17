# A-Drive Technical Documentation

## Overview
A-Drive is a comprehensive file storage and management application featuring a modern React frontend and robust Go backend. The application provides secure user authentication, role-based access control, hierarchical file organization, and an intuitive web interface for file management.

## Architecture

### Full Stack Technology Overview
#### Backend Stack
- **Go 1.23.0**: Modern Go backend with high performance
- **Gin Framework**: HTTP web framework for REST API
- **SQLite + GORM**: Database with ORM for data persistence
- **JWT Authentication**: Secure token-based authentication
- **bcrypt**: Password hashing for security

#### Frontend Stack
- **React 19.1.1**: Modern React with latest features
- **TypeScript 4.9.5**: Type-safe development
- **TanStack React Query**: Data fetching and state management
- **Tailwind CSS**: Utility-first styling framework
- **React Router**: Client-side routing
- **Axios**: HTTP client for API communication

#### Development & Build Tools
- **React Scripts**: Build tooling and development server
- **PostCSS & Autoprefixer**: CSS processing
- **React DnD**: Drag and drop functionality
- **React Dropzone**: File upload interface

### System Architecture
```
┌─────────────────────────────────────────────────────────────────┐
│                     Frontend (React App)                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │   Dashboard     │  │  File Explorer  │  │  Admin Panel    │  │
│  │   - Auth UI     │  │  - Upload Zone  │  │  - User Mgmt    │  │
│  │   - Navigation  │  │  - File List    │  │  - System Stats │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
│              │                 │                 │              │
│              └─────────────────┼─────────────────┘              │
│                                │                                │
└────────────────────────────────┼────────────────────────────────┘
                                 │ HTTP/REST API (Axios)
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (Go/Gin)                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │   Middleware    │  │    Handlers     │  │     Routes      │  │
│  │   - Auth JWT    │  │   - Files API   │  │   - Protected   │  │
│  │   - CORS        │  │   - Folders API │  │   - Public      │  │
│  │   - Database    │  │   - Admin API   │  │   - Admin       │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
│              │                 │                 │              │
│              └─────────────────┼─────────────────┘              │
└────────────────────────────────┼────────────────────────────────┘
                                 │
                    ┌────────────┼────────────┐
                    │            ▼            │
              ┌─────────────────────────────┐ │
              │       SQLite Database       │ │
              │    (GORM + Migrations)      │ │
              │   - Users, Files, Folders   │ │
              └─────────────────────────────┘ │
                                             │
                                             ▼
                                  ┌─────────────────┐
                                  │  File System    │
                                  │    Storage      │
                                  │ /storage/files/ │
                                  └─────────────────┘
```

## Frontend Architecture

### Directory Structure
```
frontend/
├── public/                     # Static assets
│   ├── index.html             # Main HTML template
│   ├── manifest.json          # PWA manifest
│   └── favicon.ico            # Application icon
├── src/
│   ├── App.tsx                # Main application component
│   ├── index.tsx              # Application entry point
│   ├── components/            # Reusable UI components
│   │   ├── FileExplorer/      # File browser components
│   │   ├── Modals/            # Modal dialog components
│   │   ├── Upload/            # File upload components
│   │   ├── BulkOperations.tsx # Bulk file operations
│   │   ├── Profile.tsx        # User profile management
│   │   └── Search.tsx         # Search functionality
│   ├── pages/                 # Page-level components
│   │   ├── Dashboard.tsx      # Main user dashboard
│   │   ├── Login.tsx          # User authentication
│   │   ├── Register.tsx       # User registration
│   │   └── AdminPanel.tsx     # Administrative interface
│   ├── hooks/                 # Custom React hooks
│   │   ├── useAuth.tsx        # Authentication state
│   │   └── useFiles.tsx       # File operations state
│   ├── services/              # API layer
│   │   └── api.ts             # HTTP client and API functions
│   ├── types/                 # TypeScript definitions
│   │   └── index.ts           # Shared interfaces
│   └── styles/                # Global styles
├── package.json               # Dependencies and scripts
├── tsconfig.json              # TypeScript configuration
├── tailwind.config.js         # Tailwind CSS configuration
└── postcss.config.js          # PostCSS configuration
```

### Frontend Key Features
- **React Router**: Client-side routing with protected routes
- **TanStack React Query**: Server state management and caching
- **Tailwind CSS**: Utility-first styling with custom design system
- **TypeScript**: Full type safety across the application
- **Drag & Drop**: React DnD for file management operations
- **Real-time Updates**: Live UI updates with React Query
- **Responsive Design**: Mobile-first responsive layout

For detailed frontend technical information, see [Frontend Technical Documentation](FRONTEND_TECHNICAL_DOCUMENTATION.md).

## Backend Architecture

### Directory Structure
```
backend/
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration management
├── database/
│   └── database.go        # Database initialization & migrations
├── models/                # Data models
│   ├── user.go           # User model
│   ├── file.go           # File metadata model
│   └── folder.go         # Folder structure model
├── handlers/              # HTTP request handlers
│   ├── auth.go           # Authentication handlers
│   ├── files.go          # File operation handlers
│   ├── folders.go        # Folder operation handlers
│   ├── admin.go          # Admin functionality handlers
│   ├── profile.go        # User profile handlers
│   ├── search.go         # Search functionality handlers
│   └── bulk.go           # Bulk operation handlers
├── middleware/            # HTTP middleware
│   └── auth.go           # Authentication & authorization middleware
├── routes/                # Route definitions
│   ├── auth.go           # Authentication routes
│   ├── files.go          # File operation routes
│   ├── folders.go        # Folder operation routes
│   └── admin.go          # Admin routes
└── utils/                 # Utility functions
    ├── auth.go           # Password hashing utilities
    └── jwt.go            # JWT token utilities
```

### Core Dependencies

#### Main Dependencies
```go
require (
    github.com/gin-contrib/cors v1.7.6      // CORS middleware
    github.com/gin-gonic/gin v1.10.1        // HTTP web framework
    github.com/golang-jwt/jwt/v5 v5.3.0     // JWT authentication
    github.com/joho/godotenv v1.5.1         // Environment variable loading
    golang.org/x/crypto v0.41.0             // Cryptographic functions
    gorm.io/driver/sqlite v1.6.0            // SQLite database driver
    gorm.io/gorm v1.30.1                    // ORM library
)
```

### Data Models

#### User Model (`models/user.go`)
```go
type User struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    Username     string         `json:"username" gorm:"unique;not null"`
    Email        string         `json:"email" gorm:"unique;not null"`
    PasswordHash string         `json:"-" gorm:"not null"`
    Role         string         `json:"role" gorm:"default:user"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
    
    Folders []Folder `json:"folders,omitempty" gorm:"foreignKey:UserID"`
    Files   []File   `json:"files,omitempty" gorm:"foreignKey:UserID"`
}
```

#### File Model (`models/file.go`)
```go
type File struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    Name         string         `json:"name" gorm:"not null"`
    OriginalName string         `json:"original_name" gorm:"not null"`
    FolderID     *uint          `json:"folder_id"`
    UserID       uint           `json:"user_id" gorm:"not null"`
    FilePath     string         `json:"file_path" gorm:"not null"`
    Size         int64          `json:"size" gorm:"not null"`
    MimeType     string         `json:"mime_type"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
    
    User   User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Folder *Folder `json:"folder,omitempty" gorm:"foreignKey:FolderID"`
}
```

#### Folder Model (`models/folder.go`)
```go
type Folder struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null"`
    ParentID  *uint          `json:"parent_id"`
    UserID    uint           `json:"user_id" gorm:"not null"`
    IconType  string         `json:"icon_type" gorm:"default:folder"`
    Path      string         `json:"path" gorm:"not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    User       User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Parent     *Folder  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Subfolders []Folder `json:"subfolders,omitempty" gorm:"foreignKey:ParentID"`
    Files      []File   `json:"files,omitempty" gorm:"foreignKey:FolderID"`
}
```

### Configuration System

#### Configuration Structure (`config/config.go`)
```go
type Config struct {
    DatabasePath    string  // SQLite database file path
    JWTSecret      string  // JWT signing secret
    RootDirectory  string  // File storage root directory
    MaxFileSize    int64   // Maximum file upload size in bytes
    AllowedTypes   string  // Allowed file types (comma-separated or "*")
    Port          string  // Server port
}
```

#### Environment Variables
- `DATABASE_PATH`: SQLite database file location (default: "./storage/database.db")
- `JWT_SECRET`: Secret key for JWT token signing (default: "your-secret-key")
- `ROOT_DIRECTORY`: Root directory for file storage (default: "./storage/files")
- `MAX_FILE_SIZE`: Maximum file upload size in bytes (default: 104857600 = 100MB)
- `ALLOWED_FILE_TYPES`: Comma-separated list of allowed file extensions (default: "*")
- `PORT`: Server port (default: "8080")

### Authentication & Authorization

#### JWT Implementation
- Token-based authentication using JWT
- Tokens contain user ID and expiration
- Middleware validates tokens on protected routes

#### Role-Based Access Control
- Two roles: `user` and `admin`
- Admin users have additional privileges
- Role checking implemented in middleware

#### Middleware Stack (`middleware/auth.go`)
1. **DatabaseMiddleware**: Injects database connection into request context
2. **AuthMiddleware**: Validates JWT tokens and loads user information
3. **AdminMiddleware**: Ensures user has admin role

### API Route Structure

#### Public Routes (No Authentication)
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /health` - Health check endpoint

#### Protected Routes (Authentication Required)
- `GET /api/auth/me` - Get current user info
- `GET /api/profile` - Get user profile
- `PUT /api/profile` - Update user profile
- `POST /api/profile/change-password` - Change password
- File operations: `/api/files/*`
- Folder operations: `/api/folders/*`

#### Admin Routes (Admin Role Required)
- `/api/admin/*` - Admin-only functionality

### Database Management

#### Initialization (`database/database.go`)
1. Creates database directory if it doesn't exist
2. Establishes SQLite connection with GORM
3. Runs automatic migrations for all models
4. Creates default admin user if no users exist

#### Default Admin User
- Username: `admin`
- Email: `admin@example.com`
- Password: `admin123`
- Role: `admin`

### File Storage System

#### Storage Structure
```
storage/
├── database/
│   └── database.db        # SQLite database file
└── files/
    └── root/
        └── {user_id}/     # User-specific directories
            ├── file1.ext
            ├── file2.ext
            └── folder1/
                └── file3.ext
```

#### File Management
- Files are stored in user-specific directories
- Original filenames are preserved in database
- File paths are stored relative to user directory
- MIME types are detected and stored

### Security Features

#### Password Security
- Passwords are hashed using bcrypt
- Salt is automatically generated
- Hash verification for login

#### CORS Configuration
- Configured for frontend origins (localhost:3000, localhost:3001)
- Supports credentials for authentication
- Allows standard HTTP methods

#### Input Validation
- Request validation through Gin's binding
- File type and size restrictions
- Path traversal prevention

## Performance Considerations

### Database Optimization
- SQLite with WAL mode for better concurrency
- Proper indexing on foreign keys
- Soft deletes for data recovery

### File Operations
- Streaming for large file uploads
- Efficient directory traversal
- Metadata caching in database

## Deployment

### Build Process
```bash
cd backend
go mod tidy
go build -o a-drive-backend
```

### Docker Support
- Dockerfile provided for containerization
- Multi-stage build for optimized image size
- Volume mounting for persistent storage

## Frontend-Backend Integration

### API Communication
- **RESTful Design**: Consistent REST API endpoints
- **JSON Data Exchange**: Structured request/response format
- **Authentication Flow**: JWT token-based authentication
- **Error Handling**: Standardized error response format

### State Synchronization
- **React Query**: Automatic cache invalidation and updates
- **Optimistic Updates**: Immediate UI feedback with rollback
- **Real-time Sync**: Live updates between frontend and backend
- **Offline Handling**: Graceful degradation when offline

### Security Integration
- **Token Management**: Secure JWT storage and transmission
- **Request Interceptors**: Automatic token attachment
- **Response Handling**: Authentication error handling
- **CORS Configuration**: Secure cross-origin communication

## Related Documentation

### Frontend Details
- [Frontend Technical Documentation](FRONTEND_TECHNICAL_DOCUMENTATION.md) - React architecture and implementation
- [Frontend Functional Documentation](FRONTEND_FUNCTIONAL_DOCUMENTATION.md) - UI features and user workflows

### Backend Details
- [Functional Documentation](FUNCTIONAL_DOCUMENTATION.md) - API endpoints and backend features

### Project Management
- [Changelog](CHANGELOG.md) - Version history and updates
- [Documentation Update Guide](DOC_UPDATE_GUIDE.md) - Maintenance procedures

---

**Last Updated**: 2025-08-17
**Version**: 1.0.0