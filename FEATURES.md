# A-Drive Features Overview

## 🔐 Authentication & Authorization

### User Management
- **User Registration**: Create new accounts with username, email, password
- **Secure Login**: JWT-based authentication with encrypted passwords
- **Role-based Access**: User and Admin roles with different permissions
- **Auto-logout**: Automatic session expiry and token validation

### Security Features
- **Password Hashing**: bcrypt encryption for secure password storage
- **JWT Tokens**: Secure token-based authentication with expiration
- **Path Traversal Protection**: Prevents access to unauthorized directories
- **User Isolation**: Each user has their own isolated file space
- **Admin Verification**: Admin-only endpoints protected by role checks

## 📁 File Management

### File Operations
- **Upload Files**: Drag & drop or click to upload any file type
- **Download Files**: Direct download with original filenames
- **Delete Files**: Remove files with confirmation dialogs
- **Rename Files**: Change file names while preserving extensions
- **File Preview**: Visual icons based on file types (images, documents, etc.)
- **Favorites**: Mark files and folders as favorites for quick access
- **Recent Files**: View recently accessed files and folders
- **Photo Gallery**: Dedicated view for all image files across directories

### Upload Features
- **Drag & Drop**: Intuitive drag and drop interface
- **Progress Tracking**: Real-time upload progress with percentage
- **Batch Upload**: Upload multiple files simultaneously
- **Size Validation**: Configurable file size limits (default: 100MB)
- **Error Handling**: Clear error messages for failed uploads

## 📂 Folder Management

### Folder Operations
- **Create Folders**: Organize files into hierarchical folder structure
- **Rename Folders**: Change folder names with real-time updates
- **Delete Folders**: Remove folders and all contents with confirmation
- **Nested Structure**: Unlimited folder depth support
- **Move Items**: Drag files and folders between directories

### Folder Customization
- **Custom Icons**: Choose from 8 different folder icons:
  - 📁 Default Folder
  - 📁 Documents
  - 📷 Images
  - 🎵 Music
  - 🎬 Videos
  - ⬇️ Downloads
  - 💼 Work
  - 👤 Personal

## 🎨 User Interface

### View Modes
- **Grid View**: Visual card-based layout with large icons
- **List View**: Detailed table view with file information
- **Responsive Design**: Mobile-friendly interface that adapts to screen size
- **Dark/Light Theme**: Clean, modern design with Tailwind CSS

### Navigation
- **Breadcrumb Navigation**: Clear path indication with clickable breadcrumbs
- **Context Menus**: Right-click menus for quick actions
- **Keyboard Shortcuts**: Common operations (Ctrl+A, Delete, Enter)
- **Search**: Real-time search within current directory
- **Favorites Panel**: Quick access to favorited files and folders
- **Recent Files Panel**: View and access recently used items
- **Photos View**: Dedicated gallery view for all image files

### User Experience
- **Loading States**: Skeleton loaders and spinners for better UX
- **Error Messages**: User-friendly error notifications
- **Success Feedback**: Confirmation messages for completed actions
- **Drag Feedback**: Visual indicators during drag operations

## 🗜️ Archive Management

### ZIP Downloads
- **Folder Archives**: Download entire folders as ZIP files
- **On-demand Generation**: ZIP files created dynamically when requested
- **Recursive Compression**: Includes all subfolders and files
- **Automatic Cleanup**: Temporary files cleaned up after download

## 👑 Admin Panel

### User Management
- **User List**: View all registered users with roles
- **Create Users**: Admin can create new user accounts
- **User Analytics**: View user registration dates and activity
- **Role Management**: Assign admin or user roles

### File System Access
- **Browse User Files**: Admins can view any user's files and folders
- **User Statistics**: File counts and storage usage per user
- **System Overview**: Global file system statistics
- **User Directories**: Access to all user isolation directories

## ⭐ Advanced Features

### Favorites System
- **Mark Favorites**: Add files and folders to favorites with star icon
- **Favorites View**: Dedicated panel showing all favorited items
- **Quick Access**: One-click access to frequently used files
- **Persistent Storage**: Favorites saved across sessions
- **Smart UI**: Visual indicators for favorited items throughout the interface

### Recent Files Tracking
- **Access Tracking**: Automatically track file and folder access
- **Recent Panel**: View 20 most recently accessed items
- **Timestamp Display**: See when items were last accessed
- **Mixed Content**: Shows both files and folders in chronological order
- **Quick Navigation**: Jump directly to recently used items

### Photo Gallery
- **Image Detection**: Automatically identifies all image files (JPEG, PNG, GIF, etc.)
- **Gallery View**: Beautiful grid layout for photo browsing
- **Cross-Directory**: Shows images from all folders in one view
- **Metadata Display**: File size, upload date, and location information
- **Favorite Photos**: Photos can be marked as favorites within gallery view

## 🔧 Technical Features

### Performance
- **React Query**: Efficient data fetching with caching
- **Lazy Loading**: Components loaded on demand
- **Optimistic Updates**: UI updates before server confirmation
- **Connection Pooling**: Efficient database connections
- **MIME Type Filtering**: Efficient image queries using database indexing

### Storage
- **SQLite Database**: Lightweight, file-based database
- **File System Storage**: Direct file storage on disk
- **User Isolation**: Separate directories per user (`/root/{user_id}/`)
- **Metadata Tracking**: File information stored in database
- **Favorites Database**: Separate table for user favorites with constraints
- **Access Logging**: Recent file access tracking with timestamps

### API Design
- **RESTful API**: Clean, predictable endpoint structure
- **JSON Responses**: Consistent data format
- **Error Handling**: Structured error responses
- **Rate Limiting**: Protection against API abuse
- **Favorites API**: Complete CRUD operations for favorites management
- **Recent Files API**: Access tracking and retrieval endpoints
- **Photos API**: Efficient image file filtering and retrieval

## 🐳 Deployment

### Docker Support
- **Multi-stage Builds**: Optimized Docker images
- **Docker Compose**: Single-command deployment
- **Volume Persistence**: Data survives container restarts
- **Environment Configuration**: Easy environment variable setup

### Production Ready
- **Health Checks**: Built-in health monitoring endpoints
- **Logging**: Structured application logging
- **CORS Support**: Cross-origin request handling
- **SSL Ready**: HTTPS support for production deployment

## 📱 Cross-Platform

### Browser Support
- **Modern Browsers**: Chrome, Firefox, Safari, Edge
- **Mobile Responsive**: Touch-friendly interface
- **Progressive Web App**: Can be installed on mobile devices
- **Keyboard Navigation**: Full keyboard accessibility

### File Types
- **Universal Support**: No restrictions on file types
- **MIME Type Detection**: Automatic file type identification
- **Icon Association**: File type-specific icons
- **Download Preservation**: Original filenames and extensions maintained

## 🛡️ Security Best Practices

### Data Protection
- **Input Validation**: All user inputs validated and sanitized
- **SQL Injection Prevention**: Parameterized queries via GORM
- **XSS Protection**: React's built-in XSS prevention
- **CSRF Protection**: Token-based request validation

### Access Control
- **Authentication Required**: All operations require valid login
- **Authorization Checks**: Endpoint-level permission validation
- **User Ownership**: Users can only access their own files
- **Admin Verification**: Admin operations require admin role

This comprehensive feature set makes A-Drive a complete, production-ready file management solution comparable to commercial cloud storage services.