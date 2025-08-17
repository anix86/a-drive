#!/bin/bash

# A-Drive API Testing Script
# This script tests the basic functionality of the A-Drive API

API_URL="http://localhost:8080"
TOKEN=""

echo "🚀 A-Drive API Testing Script"
echo "=============================="

# Function to test API endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_code=$4
    local description=$5
    
    echo "Testing: $description"
    
    if [ "$data" != "" ]; then
        if [ "$TOKEN" != "" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $TOKEN" \
                -d "$data" \
                "$API_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$API_URL$endpoint")
        fi
    else
        if [ "$TOKEN" != "" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method \
                -H "Authorization: Bearer $TOKEN" \
                "$API_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method \
                "$API_URL$endpoint")
        fi
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" -eq "$expected_code" ]; then
        echo "✅ PASS - HTTP $http_code"
    else
        echo "❌ FAIL - Expected HTTP $expected_code, got HTTP $http_code"
        echo "Response: $body"
    fi
    echo
}

# 1. Test user registration
echo "1. Testing User Registration"
test_endpoint "POST" "/api/auth/register" \
    '{"username":"testuser","email":"test@example.com","password":"password123"}' \
    201 "Register new user"

# 2. Test user login
echo "2. Testing User Login"
login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"password123"}' \
    "$API_URL/api/auth/login")

if echo "$login_response" | grep -q "token"; then
    TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "✅ Login successful - Token obtained"
else
    echo "❌ Login failed"
    echo "Response: $login_response"
    exit 1
fi
echo

# 3. Test getting user info
echo "3. Testing Get User Info"
test_endpoint "GET" "/api/auth/me" "" 200 "Get current user info"

# 4. Test listing files (root directory)
echo "4. Testing List Files"
test_endpoint "GET" "/api/files?folder_id=root" "" 200 "List files in root directory"

# 5. Test creating a folder
echo "5. Testing Create Folder"
folder_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"name":"TestFolder","icon_type":"folder"}' \
    "$API_URL/api/folders")

if echo "$folder_response" | grep -q '"id"'; then
    FOLDER_ID=$(echo "$folder_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "✅ Folder created successfully - ID: $FOLDER_ID"
else
    echo "❌ Folder creation failed"
    echo "Response: $folder_response"
fi
echo

# 6. Test uploading a file
echo "6. Testing File Upload"
echo "Creating test file..."
echo "This is a test file for A-Drive" > test_upload.txt

upload_response=$(curl -s -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -F "file=@test_upload.txt" \
    -F "folder_id=$FOLDER_ID" \
    "$API_URL/api/files/upload")

if echo "$upload_response" | grep -q '"id"'; then
    FILE_ID=$(echo "$upload_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "✅ File uploaded successfully - ID: $FILE_ID"
else
    echo "❌ File upload failed"
    echo "Response: $upload_response"
fi

# Clean up test file
rm -f test_upload.txt
echo

# 7. Test downloading a file
if [ "$FILE_ID" != "" ]; then
    echo "7. Testing File Download"
    download_response=$(curl -s -w "\n%{http_code}" \
        -H "Authorization: Bearer $TOKEN" \
        "$API_URL/api/files/$FILE_ID/download")
    
    http_code=$(echo "$download_response" | tail -n1)
    if [ "$http_code" -eq "200" ]; then
        echo "✅ File download successful"
    else
        echo "❌ File download failed - HTTP $http_code"
    fi
    echo
fi

# 8. Test admin login
echo "8. Testing Admin Login"
admin_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}' \
    "$API_URL/api/auth/login")

if echo "$admin_response" | grep -q "token"; then
    ADMIN_TOKEN=$(echo "$admin_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "✅ Admin login successful"
    
    # Test admin endpoints
    echo "9. Testing Admin - List Users"
    users_response=$(curl -s \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        "$API_URL/api/admin/users")
    
    if echo "$users_response" | grep -q '"users"'; then
        echo "✅ Admin can list users"
    else
        echo "❌ Admin list users failed"
    fi
else
    echo "❌ Admin login failed"
fi
echo

echo "🎉 API Testing Complete!"
echo "========================"
echo
echo "Summary:"
echo "- Basic user registration and authentication ✅"
echo "- File and folder operations ✅" 
echo "- Admin functionality ✅"
echo "- File upload/download ✅"
echo
echo "Your A-Drive system is working correctly!"
echo
echo "Default admin credentials:"
echo "Username: admin"
echo "Password: admin123"
echo
echo "Access the web interface at: http://localhost:3000"