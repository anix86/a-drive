# A-Drive Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-08-17

### Added
- Initial release of A-Drive full-stack file storage system
- Modern React frontend with TypeScript
- Go backend with Gin web framework
- User authentication and registration system with JWT tokens
- Role-based access control (user/admin roles)
- File upload, download, and management functionality
- Hierarchical folder structure with nested folder support
- User profile management and password change functionality
- Admin panel for user and system management
- Search functionality for files and folders
- RESTful API with comprehensive endpoint coverage
- SQLite database with GORM ORM integration
- Docker containerization support
- CORS configuration for frontend integration
- Comprehensive technical and functional documentation

### Frontend Features
- React 19.1.1 with TypeScript for type safety
- TanStack React Query for data fetching and caching
- Tailwind CSS for responsive, utility-first styling
- React Router for client-side routing with protected routes
- Drag and drop file uploads with React DnD and React Dropzone
- Real-time upload progress tracking
- File explorer with grid and list views
- Modal system for file/folder operations
- Search interface with filtering capabilities
- Admin panel interface for user management
- Mobile-responsive design
- Authentication flow with automatic token management

### Backend Features
- Go backend with Gin web framework
- JWT-based authentication middleware
- bcrypt password hashing for security
- File storage in user-specific directories
- Soft delete functionality for data recovery
- Configuration management through environment variables
- Health check endpoint for monitoring
- Bulk operations support for admin functions

### Database
- SQLite database with automatic migrations
- User, File, and Folder models with proper relationships
- Foreign key constraints and indexing
- Soft delete support with DeletedAt timestamps
- Default admin user creation on first run

### Security
- Password hashing with bcrypt and automatic salt generation
- JWT token validation on protected routes
- Role-based middleware for admin functions
- Input validation and sanitization
- CORS configuration for secure frontend communication
- Path traversal attack prevention
- Client-side token management with automatic refresh

### API Endpoints
- Authentication: `/api/auth/register`, `/api/auth/login`
- User Profile: `/api/profile/*`
- File Operations: `/api/files/*`
- Folder Operations: `/api/folders/*`
- Admin Functions: `/api/admin/*`
- Search: `/api/search`
- Health Check: `/health`

### Documentation
- Comprehensive technical documentation for full-stack architecture
- Frontend-specific technical documentation for React implementation
- Functional documentation covering backend API and features
- Frontend functional documentation covering UI and user workflows
- Documentation update guide for maintaining currency
- Unified documentation index with role-based navigation

---

**Note**: This changelog will be automatically updated with each significant change to the codebase.