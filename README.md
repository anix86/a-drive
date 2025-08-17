# 🚀 A-Drive - Complete File Management System

A professional Google Drive-like file management system built with modern technologies.

![A-Drive Dashboard](https://via.placeholder.com/800x400/4F46E5/FFFFFF?text=A-Drive+Dashboard)

## ✨ Features

### 🔐 **Authentication & Security**
- **JWT Authentication** with secure password hashing
- **Role-based Access Control** (User/Admin roles)
- **User Isolation** - Each user has their own secure file space
- **Auto-logout** and session management

### 📁 **File Management**
- **Drag & Drop Upload** with progress tracking
- **File Operations** - Upload, download, rename, delete
- **Folder Organization** - Create nested folder structures
- **ZIP Downloads** - Download entire folders as archives
- **File Type Support** - All file types supported with smart icons

### 🎨 **Modern UI/UX**
- **Grid & List Views** - Switch between visual layouts
- **Responsive Design** - Works on desktop and mobile
- **Real-time Updates** with React Query
- **Custom Folder Icons** - 8 different icon types
- **Breadcrumb Navigation** - Easy directory navigation

### 👑 **Admin Panel**
- **User Management** - View and create user accounts
- **System Overview** - Browse all user files
- **Role Management** - Assign admin privileges
- **User Analytics** - Registration dates and activity

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| **Backend** | Go 1.21, Gin Framework, SQLite + GORM |
| **Frontend** | React 18, TypeScript, Tailwind CSS |
| **Authentication** | JWT Tokens, bcrypt |
| **File Storage** | Disk-based with user isolation |
| **Deployment** | Docker & Docker Compose |
| **API** | RESTful with JSON responses |

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)
```bash
# Clone and start the application
git clone <repository-url>
cd a-drive
docker-compose up -d

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### Option 2: Manual Development Setup
```bash
# Backend setup
cd backend
go mod tidy
go run main.go

# Frontend setup (in another terminal)
cd frontend
npm install
npm start
```

## 🔑 Default Admin Account

When you first start A-Drive, an admin account is automatically created:

- **Username:** `admin`
- **Password:** `admin123`

⚠️ **Important:** Change this password immediately in production!

## 📋 Testing

Run the comprehensive API test suite:
```bash
# Make sure the backend is running on port 8080
./test-api.sh
```

This tests all major functionality including:
- User registration and authentication
- File upload/download operations
- Folder management
- Admin panel access

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [FEATURES.md](FEATURES.md) | Complete feature overview |
| [API.md](API.md) | Full API documentation |
| [DEPLOYMENT.md](DEPLOYMENT.md) | Production deployment guide |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | Common issues and solutions |

## 🎯 Key Features Showcase

### File Upload with Progress
```javascript
// Drag & drop or click to upload
// Real-time progress tracking
// Batch upload support
// Error handling with retry
```

### Folder Management
```javascript
// Create nested folder structures
// Custom folder icons (📁 📷 🎵 🎬 💼 👤)
// Rename and organize
// Download as ZIP archives
```

### Admin Dashboard
```javascript
// View all users and their files
// Create new user accounts
// System-wide file management
// User role management
```

## 🔧 Configuration

### Backend Environment (.env)
```env
DATABASE_PATH=../storage/database/database.db
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ROOT_DIRECTORY=../storage/files
MAX_FILE_SIZE=104857600  # 100MB
ALLOWED_FILE_TYPES=*
PORT=8080
```

### Frontend Environment (.env)
```env
REACT_APP_API_URL=http://localhost:8080
REACT_APP_MAX_UPLOAD_SIZE=104857600
```

## 📊 Project Structure

```
a-drive/
├── backend/                 # Go API server
│   ├── config/             # Configuration management
│   ├── database/           # Database initialization
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # JWT auth, CORS, etc.
│   ├── models/             # Data models (User, File, Folder)
│   ├── routes/             # API route definitions
│   ├── utils/              # Helper functions
│   └── main.go             # Application entry point
├── frontend/               # React TypeScript app
│   ├── src/
│   │   ├── components/     # Reusable UI components
│   │   ├── hooks/          # Custom React hooks
│   │   ├── pages/          # Main application pages
│   │   ├── services/       # API integration
│   │   └── types/          # TypeScript definitions
│   └── public/             # Static assets
├── storage/                # Data persistence
│   ├── database/           # SQLite database
│   └── files/              # User file storage
├── docker-compose.yml      # Container orchestration
└── docs/                   # Documentation
```

## 🛡️ Security Features

- **Password Hashing** with bcrypt
- **JWT Token Authentication** with expiration
- **User Isolation** - Users can only access their own files
- **Input Validation** and sanitization
- **SQL Injection Prevention** via GORM
- **Path Traversal Protection**
- **Role-based Authorization**

## 🌟 Production Ready

A-Drive is designed for production use with:

- **Docker containerization** for easy deployment
- **Health check endpoints** for monitoring
- **Structured logging** for debugging
- **Environment-based configuration**
- **Horizontal scaling support**
- **Database migration handling**

## 📈 Performance

- **React Query** for efficient data caching
- **Optimistic UI updates** for better UX
- **Chunked file uploads** for large files
- **SQLite with WAL mode** for concurrent access
- **Lazy loading** for improved performance

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 API Quick Reference

### Authentication
```bash
POST /api/auth/register    # Create account
POST /api/auth/login       # Sign in
GET  /api/auth/me          # Get user info
```

### Files & Folders
```bash
GET    /api/files                    # List files/folders
POST   /api/files/upload             # Upload file
GET    /api/files/{id}/download      # Download file
DELETE /api/files/{id}               # Delete file
POST   /api/folders                  # Create folder
POST   /api/folders/{id}/zip         # Download folder as ZIP
```

### Admin (requires admin role)
```bash
GET  /api/admin/users                # List all users
POST /api/admin/users                # Create user
GET  /api/admin/files?user_id={id}   # Browse user files
```

## 📞 Support

If you encounter any issues:

1. Check the [Troubleshooting Guide](TROUBLESHOOTING.md)
2. Run the test script: `./test-api.sh`
3. Check application logs: `docker-compose logs`
4. Review the [API Documentation](API.md)

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Built with ❤️ using Go, React, and modern web technologies.**

*A-Drive provides enterprise-grade file management with the simplicity of cloud storage services.*