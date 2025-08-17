# A-Drive Functional Documentation

## Overview
A-Drive is a comprehensive, secure web-based file storage and management system featuring both a modern React frontend and robust Go backend. Users can upload, organize, and manage their files through an intuitive web interface with hierarchical folder structure, drag-and-drop functionality, and real-time updates. The application provides role-based access control with both regular user and administrative capabilities.

### Key System Features
- **Web-Based Interface**: Modern React application with responsive design
- **Secure Authentication**: JWT-based authentication with bcrypt password hashing
- **File Management**: Upload, download, rename, delete, and organize files
- **Folder Hierarchy**: Unlimited nested folder structure
- **Search Functionality**: Global search across files and folders
- **Administrative Tools**: User management and system administration
- **Real-time Updates**: Live UI updates and progress tracking
- **Drag & Drop**: Intuitive file upload and organization

## User Management

### User Registration
- **Endpoint**: `POST /api/auth/register`
- **Functionality**: Creates new user accounts
- **Requirements**:
  - Unique username
  - Valid email address
  - Password (minimum requirements apply)
- **Default Role**: `user`
- **Auto-created**: User home directory in file system

### User Authentication
- **Endpoint**: `POST /api/auth/login`
- **Functionality**: Authenticates users and provides JWT access tokens
- **Session Management**: Token-based authentication with configurable expiration
- **Security**: Bcrypt password hashing with salt

### User Profile Management
- **View Profile**: `GET /api/profile`
  - Returns user information (username, email, role, creation date)
- **Update Profile**: `PUT /api/profile`
  - Modify username and email
  - Validation for unique constraints
- **Change Password**: `POST /api/profile/change-password`
  - Requires current password verification
  - New password hashing and storage

## File Management

### File Upload
- **Endpoint**: `POST /api/files/upload`
- **Features**:
  - Single and multiple file uploads
  - Configurable file size limits (default: 100MB)
  - File type restrictions (configurable)
  - Automatic MIME type detection
  - Original filename preservation
- **Storage**: Files stored in user-specific directories
- **Metadata**: Size, MIME type, upload timestamp stored in database

### File Operations
- **List Files**: `GET /api/files`
  - Retrieve all files for authenticated user
  - Optional folder filtering
  - Pagination support
- **Download File**: `GET /api/files/:id/download`
  - Secure file serving
  - Content-Type headers set correctly
- **Delete File**: `DELETE /api/files/:id`
  - Soft delete in database
  - Physical file removal from storage
- **Rename File**: `PUT /api/files/:id/rename`
  - Update display name
  - Preserve original filename reference

### File Information
- **File Details**: `GET /api/files/:id`
  - File metadata (name, size, type, creation date)
  - Folder association
  - Download statistics (if implemented)

## Folder Management

### Folder Structure
- **Hierarchical Organization**: Support for nested folder structures
- **Root Level**: Each user has a personal root directory
- **Path Management**: Automatic path calculation and validation
- **Breadcrumb Support**: Full path tracking for navigation

### Folder Operations
- **Create Folder**: `POST /api/folders`
  - Create new folders in any accessible location
  - Automatic path generation
  - Icon type selection (default: folder icon)
- **List Folders**: `GET /api/folders`
  - Retrieve folder structure for user
  - Include subfolder and file counts
- **Rename Folder**: `PUT /api/folders/:id/rename`
  - Update folder name
  - Automatic path updates for all contents
- **Delete Folder**: `DELETE /api/folders/:id`
  - Recursive deletion of contents
  - Soft delete support
  - Confirmation required for non-empty folders
- **Move Folder**: `PUT /api/folders/:id/move`
  - Change parent folder
  - Path recalculation for all contents

### Folder Navigation
- **Browse Folder**: `GET /api/folders/:id/contents`
  - List all files and subfolders in a folder
  - Sorting options (name, date, size, type)
  - View modes (list, grid)

## Search Functionality

### File Search
- **Endpoint**: `GET /api/search`
- **Search Criteria**:
  - File name (partial matching)
  - File type/extension
  - Date range filtering
  - Size range filtering
- **Scope**: User's accessible files only
- **Results**: Paginated results with relevance scoring

### Advanced Search
- **Content Search**: Text file content indexing (if implemented)
- **Metadata Search**: Search by file properties
- **Folder Search**: Find folders by name or path

## Administrative Functions

### User Management (Admin Only)
- **List All Users**: `GET /api/admin/users`
  - View all registered users
  - User statistics and activity
- **User Details**: `GET /api/admin/users/:id`
  - Detailed user information
  - File and folder counts
  - Storage usage statistics
- **Manage User**: `PUT /api/admin/users/:id`
  - Update user roles
  - Enable/disable accounts
  - Reset passwords
- **Delete User**: `DELETE /api/admin/users/:id`
  - Remove user account
  - Handle file ownership transfer or deletion

### System Administration
- **System Statistics**: `GET /api/admin/stats`
  - Total users, files, storage usage
  - System health metrics
- **Storage Management**:
  - Monitor disk usage
  - Clean up orphaned files
  - Storage quota enforcement

### Bulk Operations
- **Endpoint**: `POST /api/admin/bulk`
- **Operations**:
  - Bulk user creation
  - Mass file operations
  - Batch user notifications

## Security Features

### Access Control
- **Authentication Required**: All API endpoints except registration and login
- **Role-Based Authorization**: Admin-only functions protected
- **User Isolation**: Users can only access their own files and folders
- **Path Validation**: Prevent directory traversal attacks

### Data Protection
- **Password Security**: Bcrypt hashing with automatic salt generation
- **JWT Security**: Configurable token expiration and secret rotation
- **Input Validation**: Sanitization of all user inputs
- **File Type Validation**: Configurable allowed file types

### Audit and Logging
- **Access Logs**: Track user authentication and file access
- **Operation Logs**: Record file and folder operations
- **Security Events**: Log failed authentication attempts

## User Interface Features

### File Explorer
- **Tree View**: Hierarchical folder navigation
- **File Grid/List**: Multiple view modes for files
- **Drag and Drop**: File upload via drag and drop
- **Context Menus**: Right-click operations for files and folders

### Upload Interface
- **Progress Tracking**: Real-time upload progress
- **Multiple Files**: Batch upload support
- **Upload Zones**: Designated drop zones for files

### Navigation
- **Breadcrumbs**: Current path display with navigation
- **Back/Forward**: Browser-style navigation
- **Bookmarks**: Favorite folder shortcuts (if implemented)

## Error Handling

### User-Friendly Messages
- **Upload Errors**: Clear messages for size/type restrictions
- **Authentication Errors**: Helpful login failure messages
- **Permission Errors**: Informative access denied messages

### System Error Recovery
- **Database Errors**: Graceful handling of database issues
- **Storage Errors**: Disk space and permission error handling
- **Network Errors**: Connection timeout and retry mechanisms

## Performance Features

### Optimization
- **Lazy Loading**: On-demand folder content loading
- **Caching**: Metadata caching for improved performance
- **Compression**: File compression for storage efficiency (if implemented)

### Scalability
- **Pagination**: Large dataset handling
- **Async Operations**: Non-blocking file operations
- **Load Balancing**: Support for horizontal scaling (if implemented)

## Configuration Options

### File Upload Settings
- `MAX_FILE_SIZE`: Maximum file upload size (default: 100MB)
- `ALLOWED_FILE_TYPES`: Permitted file extensions (default: all types)
- `UPLOAD_TIMEOUT`: Upload operation timeout

### Storage Settings
- `ROOT_DIRECTORY`: Base storage directory path
- `STORAGE_QUOTA`: Per-user storage limits (if implemented)
- `CLEANUP_SCHEDULE`: Automatic cleanup of deleted files

### Security Settings
- `JWT_SECRET`: Token signing secret (required for production)
- `JWT_EXPIRATION`: Token validity period
- `PASSWORD_REQUIREMENTS`: Minimum password complexity

## API Usage Examples

### Authentication Flow
```bash
# Register new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","email":"user1@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}'
```

### File Operations
```bash
# Upload file
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/file.pdf"

# List files
curl -X GET http://localhost:8080/api/files \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Folder Operations
```bash
# Create folder
curl -X POST http://localhost:8080/api/folders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Documents","parent_id":null}'

# List folders
curl -X GET http://localhost:8080/api/folders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Integration Points

### Frontend Integration
- RESTful API design for easy frontend consumption
- CORS configuration for cross-origin requests
- JSON response format with consistent error handling

### External Storage (Future)
- Plugin architecture for cloud storage providers
- S3-compatible API support
- Backup and sync capabilities

### Third-Party Authentication (Future)
- OAuth2 integration (Google, GitHub, etc.)
- LDAP/Active Directory support
- SSO capabilities

## Frontend User Experience

### Web Interface Features
For detailed information about the user interface and frontend functionality, see [Frontend Functional Documentation](FRONTEND_FUNCTIONAL_DOCUMENTATION.md), which covers:

- **Authentication Interface**: Login and registration pages
- **Dashboard Experience**: Main file management interface
- **File Explorer**: Grid/list views, navigation, and breadcrumbs
- **Upload Interface**: Drag and drop functionality and progress tracking
- **Search Interface**: Global search with filters and results
- **Admin Panel**: Administrative user interface
- **Mobile Support**: Responsive design for mobile and tablet devices

### User Interface Workflows
- **Daily Usage**: Login → Browse → Upload → Organize → Download
- **File Management**: Create folders → Upload files → Organize → Search
- **Administrative**: User management → System monitoring → File oversight

## Related Documentation

### Technical Implementation
- [Technical Documentation](TECHNICAL_DOCUMENTATION.md) - Full stack architecture
- [Frontend Technical Documentation](FRONTEND_TECHNICAL_DOCUMENTATION.md) - React implementation details

### Frontend Experience
- [Frontend Functional Documentation](FRONTEND_FUNCTIONAL_DOCUMENTATION.md) - Detailed UI features and workflows

### Project Management
- [Changelog](CHANGELOG.md) - Version history and updates
- [Documentation Update Guide](DOC_UPDATE_GUIDE.md) - Maintenance procedures

---

**Last Updated**: 2025-08-17
**Version**: 1.0.0