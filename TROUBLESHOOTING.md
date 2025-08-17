# A-Drive Troubleshooting Guide

## Common Issues and Solutions

### 🐳 Docker Issues

#### "Port already in use" Error
```bash
Error: bind: address already in use
```

**Solutions:**
1. Check what's using the ports:
```bash
lsof -i :3000  # Frontend port
lsof -i :8080  # Backend port
```

2. Stop conflicting services or change ports in `docker-compose.yml`:
```yaml
services:
  frontend:
    ports:
      - "3001:80"  # Change from 3000 to 3001
  backend:
    ports:
      - "8081:8080"  # Change from 8080 to 8081
```

#### Docker Build Fails
```bash
docker-compose build --no-cache
docker-compose up -d
```

#### Container Exits Immediately
```bash
# Check logs
docker-compose logs backend
docker-compose logs frontend

# Check container status
docker-compose ps
```

### 🔧 Backend Issues

#### Database Connection Error
```
Failed to connect to database
```

**Solutions:**
1. Check database directory permissions:
```bash
mkdir -p ./storage/database
chmod 755 ./storage/database
```

2. Verify environment variables in `.env`:
```env
DATABASE_PATH=../storage/database/database.db
```

#### JWT Token Issues
```
Invalid token / Unauthorized
```

**Solutions:**
1. Clear browser localStorage:
```javascript
// In browser console
localStorage.clear()
```

2. Check JWT secret in backend `.env`:
```env
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

#### File Upload Fails
```
Failed to create file / Permission denied
```

**Solutions:**
1. Check storage directory permissions:
```bash
mkdir -p ./storage/files
chmod 755 ./storage/files
```

2. Verify file size limits:
```env
MAX_FILE_SIZE=104857600  # 100MB in bytes
```

3. Check disk space:
```bash
df -h .
```

### 🎨 Frontend Issues

#### "Module not found" Errors
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

#### Build Fails
```bash
cd frontend
npm run build
```

Common fixes:
1. Update TypeScript:
```bash
npm install typescript@latest
```

2. Clear npm cache:
```bash
npm cache clean --force
```

#### API Connection Issues
Check frontend `.env`:
```env
REACT_APP_API_URL=http://localhost:8080
```

For production, use your actual domain:
```env
REACT_APP_API_URL=https://your-api-domain.com
```

### 🔐 Authentication Issues

#### Can't Login as Admin
Default credentials:
- Username: `admin`
- Password: `admin123`

If admin user doesn't exist, check backend logs:
```bash
docker-compose logs backend | grep admin
```

#### Password Reset
Currently not implemented. To reset admin password:
1. Stop the application
2. Delete the database file:
```bash
rm ./storage/database/database.db
```
3. Restart the application (admin user will be recreated)

### 📁 File System Issues

#### Files Not Appearing
1. Check user permissions
2. Verify file paths in database vs filesystem
3. Check browser network tab for API errors

#### Upload Progress Stuck
1. Check file size limits
2. Verify network connection
3. Check backend logs for errors

### 🔍 Debugging Steps

#### Check Application Health
```bash
# Test backend health
curl http://localhost:8080/health

# Test authentication
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

#### View Application Logs
```bash
# Backend logs
docker-compose logs -f backend

# Frontend logs  
docker-compose logs -f frontend

# All logs
docker-compose logs -f
```

#### Database Inspection
```bash
# Access SQLite database
sqlite3 ./storage/database/database.db

# Useful queries
.tables
SELECT * FROM users;
SELECT * FROM folders;
SELECT * FROM files LIMIT 10;
```

#### File System Check
```bash
# Check storage structure
tree ./storage/

# Check user directories
ls -la ./storage/files/root/

# Check file permissions
ls -la ./storage/files/root/*/
```

### 🌐 Network Issues

#### CORS Errors
Check backend CORS configuration in `main.go`:
```go
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))
```

#### API Timeouts
1. Check network connectivity
2. Increase timeout values
3. Check server resources (CPU/Memory)

### 🔧 Performance Issues

#### Slow File Uploads
1. Check file size limits
2. Verify network bandwidth
3. Check server disk I/O

#### Database Performance
SQLite performance tips:
1. Enable WAL mode:
```sql
PRAGMA journal_mode=WAL;
```

2. Optimize database:
```sql
VACUUM;
ANALYZE;
```

### 🛠️ Development Issues

#### Hot Reload Not Working
```bash
cd frontend
rm -rf node_modules
npm install
npm start
```

#### Go Module Issues
```bash
cd backend
go clean -modcache
go mod tidy
go mod download
```

### 📊 Monitoring Commands

#### System Resources
```bash
# Check disk usage
df -h

# Check memory usage  
free -h

# Check running processes
ps aux | grep -E "(go|node|nginx)"

# Check port usage
netstat -tulpn | grep -E "(3000|8080)"
```

#### Application Status
```bash
# Docker container status
docker-compose ps

# Container resource usage
docker stats

# Container logs (last 100 lines)
docker-compose logs --tail=100
```

### 🆘 Emergency Recovery

#### Complete Reset
```bash
# Stop all containers
docker-compose down

# Remove all data (WARNING: This deletes all files and users)
rm -rf ./storage/

# Rebuild and restart
docker-compose build --no-cache
docker-compose up -d
```

#### Backup Before Reset
```bash
# Backup database
cp ./storage/database/database.db ./backup-$(date +%Y%m%d).db

# Backup files
tar -czf ./files-backup-$(date +%Y%m%d).tar.gz ./storage/files/
```

### 📞 Getting Help

If you're still experiencing issues:

1. **Check the logs** - Most issues show up in application logs
2. **Run the test script** - `./test-api.sh` to verify basic functionality
3. **Review configuration** - Double-check environment variables
4. **Check GitHub issues** - Look for similar problems in the repository
5. **Create detailed bug report** - Include logs, configuration, and steps to reproduce

#### Log Collection Script
```bash
#!/bin/bash
echo "Collecting A-Drive diagnostic information..."
echo "===========================================" > debug.log
echo "Date: $(date)" >> debug.log
echo "System: $(uname -a)" >> debug.log
echo >> debug.log

echo "Docker Status:" >> debug.log
docker-compose ps >> debug.log 2>&1
echo >> debug.log

echo "Backend Logs:" >> debug.log
docker-compose logs backend >> debug.log 2>&1
echo >> debug.log

echo "Frontend Logs:" >> debug.log  
docker-compose logs frontend >> debug.log 2>&1
echo >> debug.log

echo "Environment Variables:" >> debug.log
cat backend/.env >> debug.log 2>&1
cat frontend/.env >> debug.log 2>&1

echo "Debug information saved to debug.log"
```

This troubleshooting guide covers most common issues you might encounter with A-Drive. Keep it handy for quick reference!