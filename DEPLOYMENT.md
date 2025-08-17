# A-Drive Deployment Guide

## Quick Start with Docker Compose

1. **Clone and start the application:**
```bash
git clone <your-repo>
cd a-drive
docker-compose up -d
```

2. **Access the application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

3. **Default admin login:**
- Username: `admin`
- Password: `admin123`

## Development Setup

### Backend Development
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend Development
```bash
cd frontend
npm install
npm start
```

## Production Deployment

### Environment Variables

#### Backend (.env)
```env
DATABASE_PATH=/app/storage/database/database.db
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ROOT_DIRECTORY=/app/storage/files
MAX_FILE_SIZE=104857600  # 100MB
ALLOWED_FILE_TYPES=*
PORT=8080
```

#### Frontend (.env)
```env
REACT_APP_API_URL=https://your-api-domain.com
REACT_APP_MAX_UPLOAD_SIZE=104857600
```

### Security Considerations

1. **Change JWT Secret**: Update `JWT_SECRET` in production
2. **HTTPS**: Use SSL certificates for production
3. **File Size Limits**: Adjust `MAX_FILE_SIZE` as needed
4. **CORS**: Update CORS origins in backend for production domains
5. **Admin Password**: Change default admin password immediately

### Database Backup

```bash
# Backup SQLite database
cp ./storage/database/database.db ./backup/database-$(date +%Y%m%d).db

# Backup user files
tar -czf ./backup/files-$(date +%Y%m%d).tar.gz ./storage/files/
```

### Nginx Configuration (Production)

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://frontend:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        client_max_body_size 100M;
    }
}
```

### Health Checks

The backend includes health check endpoints:
- `GET /health` - Basic health check
- `GET /api/auth/me` - Authentication health check

### Monitoring

Consider adding:
- Application logs (structured logging implemented)
- File system monitoring for storage usage
- Database monitoring for SQLite performance
- Rate limiting monitoring