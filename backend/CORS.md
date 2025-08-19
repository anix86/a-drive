# CORS Configuration

The A-Drive backend now supports configurable CORS (Cross-Origin Resource Sharing) settings via environment variables.

## Environment Variables

### CORS_ORIGINS
Comma-separated list of allowed origins that can make requests to the API.
- **Default**: `http://localhost:3000,http://localhost:3001`
- **Example**: `CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com`

### CORS_METHODS
Comma-separated list of allowed HTTP methods.
- **Default**: `GET,POST,PUT,DELETE,OPTIONS`
- **Example**: `CORS_METHODS=GET,POST,PUT,DELETE,OPTIONS,PATCH`

### CORS_HEADERS
Comma-separated list of allowed request headers.
- **Default**: `Origin,Content-Type,Authorization`
- **Example**: `CORS_HEADERS=Origin,Content-Type,Authorization,X-Requested-With`

## Configuration Files

### .env.example
Template file showing all available environment variables with example values.

### .env.local
Local development configuration template.

### .env.production
Production configuration template with security considerations.

## API Endpoints

### GET /cors
Public endpoint that returns current CORS configuration.
```json
{
  "origins": ["http://localhost:3000", "http://localhost:3001"],
  "methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
  "headers": ["Origin", "Content-Type", "Authorization"],
  "message": "CORS configuration"
}
```

### GET /api/admin/cors
Admin-only endpoint with detailed CORS information.

### GET /api/admin/config
Admin-only endpoint returning full server configuration.

## Setup Instructions

1. **Copy Environment File**:
   ```bash
   cp .env.example .env
   ```

2. **Update CORS Origins**:
   ```bash
   # For development
   CORS_ORIGINS=http://localhost:3000,http://localhost:3001
   
   # For production
   CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
   ```

3. **Restart Server**:
   After updating .env, restart the Go server for changes to take effect.

## Security Notes

- **Production**: Only include necessary origins, avoid wildcards
- **Development**: Can include localhost and 127.0.0.1 variants  
- **HTTPS**: Use HTTPS origins in production
- **Credentials**: Server allows credentials by default

## Testing CORS

Test CORS configuration:
```bash
# Check current CORS settings
curl http://localhost:8080/cors

# Test preflight request
curl -X OPTIONS \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  http://localhost:8080/api/files
```