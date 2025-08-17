# A-Drive Frontend Technical Documentation

## Overview
The A-Drive frontend is a modern React application built with TypeScript, providing a responsive and intuitive user interface for file storage and management. It features a component-based architecture with state management, authentication, and real-time file operations.

## Technology Stack

### Core Technologies
- **React 19.1.1**: Modern React with latest features and hooks
- **TypeScript 4.9.5**: Type-safe JavaScript development
- **React Router DOM 7.8.1**: Client-side routing and navigation
- **TanStack React Query 5.85.3**: Data fetching, caching, and synchronization

### UI and Styling
- **Tailwind CSS 3.3.6**: Utility-first CSS framework
- **@tailwindcss/forms 0.5.7**: Enhanced form styling
- **Custom Design System**: Primary color palette and component styles

### Development Tools
- **React Scripts 5.0.1**: Build tools and development server
- **PostCSS 8.4.31**: CSS processing and optimization
- **Autoprefixer 10.4.21**: CSS vendor prefixing

### HTTP and API
- **Axios 1.11.0**: HTTP client for API communication
- **Interceptors**: Request/response handling for authentication

### Drag & Drop and File Handling
- **React DnD 16.0.1**: Drag and drop functionality
- **React Dropzone 14.3.8**: File upload with drag and drop
- **HTML5 Backend**: Native drag and drop support

### Testing
- **Testing Library**: React, DOM, and Jest DOM testing utilities
- **Jest**: Test runner and framework
- **User Event Testing**: User interaction testing

## Architecture

### Project Structure
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
│   │   │   ├── FileItem.tsx   # Individual file component
│   │   │   ├── FileList.tsx   # File listing component
│   │   │   └── FolderItem.tsx # Folder component
│   │   ├── Modals/            # Modal dialog components
│   │   │   ├── CreateFolderModal.tsx
│   │   │   ├── RenameModal.tsx
│   │   │   └── IconSelectorModal.tsx
│   │   ├── Upload/            # File upload components
│   │   │   └── UploadZone.tsx # Drag and drop upload
│   │   ├── BulkOperations.tsx # Bulk file operations
│   │   ├── Profile.tsx        # User profile management
│   │   └── Search.tsx         # Search functionality
│   ├── pages/                 # Page-level components
│   │   ├── Dashboard.tsx      # Main user dashboard
│   │   ├── Login.tsx          # User authentication
│   │   ├── Register.tsx       # User registration
│   │   └── AdminPanel.tsx     # Administrative interface
│   ├── hooks/                 # Custom React hooks
│   │   ├── useAuth.tsx        # Authentication state management
│   │   └── useFiles.tsx       # File operations management
│   ├── services/              # API layer
│   │   └── api.ts             # HTTP client and API functions
│   ├── types/                 # TypeScript type definitions
│   │   └── index.ts           # Shared interfaces and types
│   └── styles/                # Global styles
│       └── index.css          # Global CSS and Tailwind imports
├── package.json               # Project dependencies and scripts
├── tsconfig.json              # TypeScript configuration
├── tailwind.config.js         # Tailwind CSS configuration
└── postcss.config.js          # PostCSS configuration
```

## Core Components Architecture

### Application Root (`App.tsx`)
```typescript
// Main application setup with providers and routing
<QueryClientProvider client={queryClient}>
  <AuthProvider>
    <Router>
      <Routes>
        {/* Route definitions with protection */}
      </Routes>
    </Router>
  </AuthProvider>
</QueryClientProvider>
```

#### Route Protection System
- **ProtectedRoute**: Requires authentication, redirects to login if not authenticated
- **PublicRoute**: Redirects to dashboard if already authenticated
- **Loading States**: Spinner components during authentication checks

### State Management

#### React Query Configuration (`App.tsx:10-17`)
```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,                      // Retry failed requests once
      refetchOnWindowFocus: false,   // Don't refetch on window focus
    },
  },
});
```

#### Custom Hooks Architecture

##### Authentication Hook (`hooks/useAuth.tsx`)
- **Features**:
  - JWT token management in localStorage
  - User session persistence
  - Automatic login/logout handling
  - Loading states during authentication

##### File Operations Hook (`hooks/useFiles.tsx`)
- **Features**:
  - File and folder state management
  - Upload progress tracking
  - Real-time updates with React Query
  - Optimistic updates for better UX

### API Layer Architecture (`services/api.ts`)

#### Base Configuration
```typescript
const api = axios.create({
  baseURL: API_URL,                    // Environment-configurable base URL
  headers: {
    'Content-Type': 'application/json',
  },
});
```

#### Authentication Interceptors
```typescript
// Request interceptor: Add JWT token to all requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor: Handle authentication errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

#### API Modules Structure
- **authApi**: Registration, login, user info
- **filesApi**: File CRUD operations, upload, download
- **foldersApi**: Folder operations, ZIP downloads
- **adminApi**: Administrative functions
- **Additional APIs**: Search, profile, bulk operations

### Component Design Patterns

#### File Explorer Components
- **FileList**: Container component for file grid/list view
- **FileItem**: Individual file representation with actions
- **FolderItem**: Folder component with navigation and actions
- **Hierarchical Navigation**: Breadcrumb and tree structure support

#### Modal System
- **CreateFolderModal**: Folder creation with icon selection
- **RenameModal**: Generic rename functionality for files/folders
- **IconSelectorModal**: Visual icon picker for folders

#### Upload System
- **UploadZone**: Drag and drop file upload with progress
- **Progress Tracking**: Real-time upload progress with visual feedback
- **Multiple File Support**: Batch upload handling

### TypeScript Integration

#### Core Type Definitions (`types/index.ts`)
```typescript
interface User {
  id: number;
  username: string;
  email: string;
  role: string;
}

interface File {
  id: number;
  name: string;
  original_name: string;
  folder_id?: number;
  user_id: number;
  file_path: string;
  size: number;
  mime_type: string;
  created_at: string;
  updated_at: string;
}

interface Folder {
  id: number;
  name: string;
  parent_id?: number;
  user_id: number;
  icon_type: string;
  path: string;
  created_at: string;
  updated_at: string;
  subfolders?: Folder[];
  files?: File[];
}
```

#### Request/Response Types
- **AuthResponse**: Login/registration response format
- **LoginRequest/RegisterRequest**: Authentication request formats
- **CreateFolderRequest**: Folder creation parameters
- **UploadProgress**: File upload progress tracking

### Styling Architecture

#### Tailwind CSS Configuration (`tailwind.config.js`)
```javascript
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),  // Enhanced form styling
  ],
}
```

#### Design System
- **Primary Color Palette**: Blue-based theme with multiple shades
- **Form Styling**: Enhanced form controls with Tailwind Forms plugin
- **Responsive Design**: Mobile-first responsive breakpoints
- **Component Classes**: Utility-first approach with custom components

### Performance Optimizations

#### React Query Features
- **Caching**: Automatic response caching with configurable TTL
- **Background Updates**: Stale-while-revalidate pattern
- **Request Deduplication**: Prevents duplicate API calls
- **Optimistic Updates**: Immediate UI feedback for better UX

#### Code Splitting
- **Route-Based Splitting**: Automatic code splitting by routes
- **Component Lazy Loading**: On-demand component loading
- **Bundle Optimization**: Tree shaking and dead code elimination

#### File Upload Optimization
- **Progress Tracking**: Real-time upload progress with axios
- **Chunked Uploads**: Support for large file uploads (configurable)
- **Concurrent Uploads**: Multiple file upload support

### Security Features

#### Client-Side Security
- **JWT Token Storage**: Secure token management in localStorage
- **Automatic Token Refresh**: Token validation and refresh handling
- **Route Protection**: Authentication-based route access control
- **XSS Prevention**: React's built-in XSS protection

#### Input Validation
- **Form Validation**: Client-side form validation with TypeScript
- **File Type Validation**: MIME type checking before upload
- **Size Limitations**: Client-side file size validation

### Error Handling

#### API Error Handling
- **Axios Interceptors**: Global error handling for API responses
- **User-Friendly Messages**: Meaningful error messages for users
- **Network Error Recovery**: Retry mechanisms for failed requests
- **Loading States**: Proper loading indicators during operations

#### Component Error Boundaries
- **Error Boundaries**: React error boundaries for component crashes
- **Fallback UI**: Graceful degradation when errors occur
- **Error Reporting**: Console logging for development debugging

### Development Workflow

#### Build and Development Scripts (`package.json:31-35`)
```json
{
  "scripts": {
    "start": "react-scripts start",    // Development server
    "build": "react-scripts build",   // Production build
    "test": "react-scripts test",     // Test runner
    "eject": "react-scripts eject"    // Eject from Create React App
  }
}
```

#### Development Server
- **Hot Reloading**: Automatic reload on file changes
- **Proxy Configuration**: API proxy for development
- **Environment Variables**: Support for `.env` files
- **Source Maps**: Development debugging support

#### Testing Setup
- **Jest Configuration**: Unit testing with Jest
- **React Testing Library**: Component testing utilities
- **User Event Testing**: User interaction testing
- **Coverage Reports**: Test coverage analysis

### Browser Compatibility

#### Target Browsers (`package.json:43-54`)
```json
{
  "production": [
    ">0.2%",           // Browsers with >0.2% market share
    "not dead",        // Browsers that are not dead
    "not op_mini all"  // Exclude Opera Mini
  ],
  "development": [
    "last 1 chrome version",   // Latest Chrome for development
    "last 1 firefox version",
    "last 1 safari version"
  ]
}
```

#### Progressive Web App Features
- **Manifest.json**: PWA configuration
- **Service Worker**: Offline support (if implemented)
- **Responsive Design**: Mobile-first responsive layout

### Environment Configuration

#### Environment Variables
- **REACT_APP_API_URL**: Backend API base URL (default: http://localhost:8080)
- **Build-time Variables**: Available during build process
- **Runtime Configuration**: Dynamic configuration support

#### Deployment Configuration
- **Static Build**: Optimized static files for deployment
- **Asset Optimization**: Minification and compression
- **Cache Headers**: Browser caching optimization

## Integration with Backend

### API Communication
- **RESTful API**: Consistent REST API communication
- **JSON Format**: Structured data exchange
- **Error Handling**: Standardized error response handling
- **Authentication**: JWT token-based authentication

### Real-time Features (Future)
- **WebSocket Support**: Real-time file updates
- **Live Collaboration**: Multi-user file operations
- **Notifications**: Real-time user notifications

### File Handling
- **Upload Streaming**: Efficient file upload handling
- **Download Management**: Secure file download with proper headers
- **Preview Support**: File preview capabilities (future enhancement)

---

**Last Updated**: 2025-08-17
**Version**: 1.0.0