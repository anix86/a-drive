# A-Drive API Documentation

## Authentication

All API endpoints except authentication require a JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

### Auth Endpoints

#### POST /api/auth/register
Register a new user account.

**Request:**
```json
{
  "username": "john_doe",
  "email": "john@example.com", 
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "jwt_token_here",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

#### POST /api/auth/login
Login with existing credentials.

**Request:**
```json
{
  "username": "john_doe",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "jwt_token_here",
  "user": {
    "id": 1,
    "username": "john_doe", 
    "email": "john@example.com",
    "role": "user"
  }
}
```

#### GET /api/auth/me
Get current user information.

**Response:**
```json
{
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com", 
    "role": "user"
  }
}
```

## File Operations

#### GET /api/files?folder_id={id}
List files and folders in a directory.

**Parameters:**
- `folder_id` (optional): Folder ID to list contents. Use "root" for root directory.

**Response:**
```json
{
  "folders": [
    {
      "id": 1,
      "name": "Documents",
      "parent_id": null,
      "user_id": 1,
      "icon_type": "documents",
      "path": "Documents",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "files": [
    {
      "id": 1,
      "name": "document.pdf", 
      "original_name": "document.pdf",
      "folder_id": 1,
      "user_id": 1,
      "file_path": "/storage/files/root/1/1_document.pdf",
      "size": 1024576,
      "mime_type": "application/pdf",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /api/files/upload
Upload a file.

**Request:** `multipart/form-data`
- `file`: File to upload
- `folder_id` (optional): Target folder ID

**Response:**
```json
{
  "file": {
    "id": 1,
    "name": "document.pdf",
    "original_name": "document.pdf", 
    "folder_id": 1,
    "user_id": 1,
    "file_path": "/storage/files/root/1/1_document.pdf",
    "size": 1024576,
    "mime_type": "application/pdf",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /api/files/{id}/download
Download a file.

**Response:** File content with appropriate headers.

#### DELETE /api/files/{id}
Delete a file.

**Response:**
```json
{
  "message": "File deleted successfully"
}
```

#### PUT /api/files/{id}
Rename a file.

**Request:**
```json
{
  "name": "new_filename.pdf"
}
```

**Response:**
```json
{
  "file": {
    "id": 1,
    "name": "new_filename.pdf",
    // ... other file properties
  }
}
```

## Folder Operations

#### POST /api/folders
Create a new folder.

**Request:**
```json
{
  "name": "New Folder",
  "parent_id": 1,
  "icon_type": "folder"
}
```

**Response:**
```json
{
  "folder": {
    "id": 2,
    "name": "New Folder",
    "parent_id": 1,
    "user_id": 1,
    "icon_type": "folder",
    "path": "Documents/New Folder",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /api/folders/{id}
Get folder details with contents.

**Response:**
```json
{
  "folder": {
    "id": 1,
    "name": "Documents",
    "parent_id": null,
    "user_id": 1,
    "icon_type": "documents", 
    "path": "Documents",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "subfolders": [...],
    "files": [...]
  }
}
```

#### PUT /api/folders/{id}
Update folder (rename or change icon).

**Request:**
```json
{
  "name": "Renamed Folder",
  "icon_type": "work"
}
```

**Response:**
```json
{
  "folder": {
    "id": 1,
    "name": "Renamed Folder",
    "icon_type": "work",
    // ... other properties
  }
}
```

#### DELETE /api/folders/{id}
Delete a folder and all its contents.

**Response:**
```json
{
  "message": "Folder deleted successfully"
}
```

#### POST /api/folders/{id}/zip
Create and download a ZIP archive of the folder.

**Response:** ZIP file download.

## Favorites

#### GET /api/favorites
Get all favorites for the authenticated user.

**Response:**
```json
{
  "favorites": [
    {
      "id": 1,
      "user_id": 1,
      "item_type": "file",
      "item_id": 5,
      "created_at": "2024-01-01T00:00:00Z",
      "item": {
        "id": 5,
        "name": "document.pdf",
        "is_favorite": true,
        // ... other file/folder properties
      }
    }
  ]
}
```

#### POST /api/favorites
Add an item to favorites.

**Request:**
```json
{
  "item_type": "file",
  "item_id": 5
}
```

**Response:**
```json
{
  "favorite": {
    "id": 1,
    "user_id": 1,
    "item_type": "file",
    "item_id": 5,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### DELETE /api/favorites/{id}
Remove a favorite by favorite ID.

**Response:**
```json
{
  "message": "Favorite removed successfully"
}
```

#### DELETE /api/favorites/item
Remove a favorite by item type and ID.

**Request:**
```json
{
  "item_type": "file",
  "item_id": 5
}
```

**Response:**
```json
{
  "message": "Favorite removed successfully"
}
```

#### GET /api/favorites/check/{type}/{id}
Check if an item is favorited.

**Parameters:**
- `type`: "file" or "folder"
- `id`: Item ID

**Response:**
```json
{
  "is_favorite": true,
  "favorite_id": 1
}
```

## Recent Files

#### GET /api/recent-files
Get the 20 most recently accessed files and folders.

**Response:**
```json
{
  "recent_files": [
    {
      "id": 1,
      "user_id": 1,
      "item_type": "file",
      "item_id": 5,
      "accessed_at": "2024-01-01T12:00:00Z",
      "created_at": "2024-01-01T10:00:00Z",
      "item": {
        "id": 5,
        "name": "document.pdf",
        // ... other file/folder properties
      }
    }
  ]
}
```

#### POST /api/recent-files/track/file/{id}
Track access to a file.

**Response:**
```json
{
  "message": "Access tracked successfully"
}
```

#### POST /api/recent-files/track/folder/{id}
Track access to a folder.

**Response:**
```json
{
  "message": "Access tracked successfully"
}
```

## Photos

#### GET /api/photos
Get all image files for the authenticated user.

**Response:**
```json
{
  "files": [
    {
      "id": 3,
      "name": "photo.jpg",
      "original_name": "photo.jpg",
      "folder_id": null,
      "user_id": 1,
      "file_path": "../storage/files/root/1/1_photo.jpg",
      "size": 1048576,
      "mime_type": "image/jpeg",
      "current_version": 1,
      "versioning_enabled": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "is_favorite": false
    }
  ],
  "count": 1
}
```

## Admin Operations

*Requires admin role*

#### GET /api/admin/users
List all users.

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /api/admin/users
Create a new user (admin only).

**Request:**
```json
{
  "username": "new_user",
  "email": "newuser@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "user": {
    "id": 2,
    "username": "new_user",
    "email": "newuser@example.com", 
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /api/admin/files?user_id={id}&folder_id={id}
Browse files for any user.

**Parameters:**
- `user_id`: Required. User ID to browse files for.
- `folder_id` (optional): Folder to browse. Use "root" for root directory.

**Response:**
```json
{
  "user": {
    "id": 2,
    "username": "target_user",
    "email": "user@example.com",
    "role": "user"
  },
  "folders": [...],
  "files": [...]
}
```

## Error Responses

All endpoints return appropriate HTTP status codes with error messages:

```json
{
  "error": "Error description here"
}
```

Common status codes:
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (missing/invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource doesn't exist)
- `409` - Conflict (duplicate username/email)
- `500` - Internal Server Error