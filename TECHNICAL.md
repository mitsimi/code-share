# Code Share

A modern web application for sharing and managing code snippets, built with Go and Vue.js.

## Features

- **User Authentication & Authorization**

  - JWT-based authentication with refresh tokens
  - Secure password hashing with bcrypt
  - Session management with automatic cleanup
  - Protected routes and middleware

- **Code Snippet Management**

  - Create, read, update, and delete snippets
  - Rich snippet metadata (title, content, language, author)

- **Social Features**

  - Like and unlike snippets with real-time updates
  - Save/bookmark snippets for later reference
  - View liked and saved snippets in user profiles

- **User Profiles**

  - Complete user profile management
  - Update username, email, and avatar
  - Change password with current password verification
  - View personal snippets, liked snippets, and saved snippets

- **Modern UI/UX**

  - Responsive design with mobile support
  - Dark/light theme switching
  - Loading states and error handling
  - Toast notifications with Vue Sonner
  - Form validation with Zod and VeeValidate
  - Client-side caching with TanStack Query

- **Backend Architecture**
  - Clean architecture
  - Type-safe database queries with SQLC
  - Comprehensive error handling and logging
  - RESTful API with proper HTTP status codes

## API Endpoints

### Authentication

- `POST /api/auth/login` - User login
- `POST /api/auth/signup` - User registration
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh access token

### Snippets

- `GET /api/snippets` - Get all snippets
- `GET /api/snippets/{id}` - Get a specific snippet
- `POST /api/snippets` - Create a new snippet
- `PUT /api/snippets/{id}` - Update a snippet
- `DELETE /api/snippets/{id}` - Delete a snippet
- `PATCH /api/snippets/{id}/like?action=like|unlike` - Like or unlike a snippet
- `PATCH /api/snippets/{id}/save?action=save|unsave` - Save or unsave a snippet

### Users

- `GET /api/users/{id}` - Get user by ID
- `GET /api/users/{id}/snippets` - Get user's snippets
- `GET /api/users/{id}/liked` - Get user's liked snippets
- `GET /api/users/{id}/saved` - Get user's saved snippets
- `PATCH /api/users/{id}` - Update user profile
- `PATCH /api/users/{id}/password` - Update user password
- `PATCH /api/users/{id}/avatar` - Update user avatar

### User Profile (Authenticated)

- `GET /api/users/me` - Get current user's profile
- `GET /api/users/me/snippets` - Get current user's snippets
- `GET /api/users/me/liked` - Get current user's liked snippets
- `GET /api/users/me/saved` - Get current user's saved snippets
- `PATCH /api/users/me` - Update current user's profile
- `PATCH /api/users/me/password` - Update current user's password
- `PATCH /api/users/me/avatar` - Update current user's avatar

## Project Architecture

### Backend

The backend follows a clean architecture pattern:

- **Domain Layer**: Core business logic and entities
- **Repository Layer**: Data access interfaces and implementations
- **API Layer**: HTTP handlers, DTOs, and request/response handling
- **Server Layer**: Routing, middleware, and server configuration
- **Storage Layer**: Database operations with SQLC for type-safe queries

**Key Design Patterns:**

- Repository pattern for data access
- Dependency injection for loose coupling
- Middleware pattern for cross-cutting concerns
- Clean separation of concerns between layers

### Frontend

The frontend is built with modern Vue.js practices:

- **Vue 3 Composition API**: Modern reactive programming
- **TypeScript**: Type safety throughout the application
- **Vite**: Fast build tooling and development server
- **Tailwind CSS**: Utility-first styling with custom design system
- **Pinia**: State management with TypeScript support
- **TanStack Query**: Server state management and caching
- **VeeValidate + Zod**: Form validation with schema validation
- **Vue Router**: Client-side routing with navigation guards

**Component Architecture:**

- Atomic design principles
- Reusable UI components with shadcn/ui
- Composition-based component logic
- Proper TypeScript interfaces and types

## Technologies Used

### Backend

- **Go 1.24+**: High-performance server language
- **SQLite**: Lightweight database with SQLC
- **Chi Router**: Lightweight HTTP router
- **Gorilla/WebSocket**: WebSocket implementation
- **JWT**: JSON Web Tokens for authentication
- **bcrypt**: Password hashing
- **Zap**: Structured logging
- **Docker**: Containerization

### Frontend

- **Vue 3**: Progressive JavaScript framework
- **TypeScript**: Type-safe JavaScript
- **Vite**: Build tool and dev server
- **Tailwind CSS**: Utility-first CSS framework
- **Pinia**: State management
- **TanStack Query**: Data fetching and caching
- **VeeValidate**: Form validation
- **Zod**: Schema validation
- **Lucide Vue**: Icon library

### Development Tools

- **SQLC**: Type-safe SQL code generation
- **Air**: Live reload backend for development
- **ESLint**: Code linting
- **Prettier**: Code formatting
- **pnpm**: Fast package manager

## Database Schema

The application uses SQLite with the following main tables:

- **users**: User accounts and profiles
- **snippets**: Code snippets with metadata
- **user_likes**: Many-to-many relationship for snippet likes
- **user_saves**: Many-to-many relationship for saved snippets
- **sessions**: User session management

## Development

### Prerequisites

- Go 1.24 or later
- Node.js 22.x or later
- pnpm (recommended) or npm
- Docker and Docker Compose (optional)

### Backend Setup

1. Install Go dependencies:

   ```bash
   go mod download
   ```

2. Prepare database
```bash
  mkdir data
  SEED=true go run .
```

3. Run the server:
   ```bash
   # Using Air for hot reload
   air
   # Or directly
   go run .
   ```

### Frontend Setup

1. Install dependencies:

   ```bash
   cd frontend
   pnpm install
   ```

2. Start development server:

   ```bash
   pnpm dev
   ```

3. Build for production:
   ```bash
   pnpm build
   ```

### Accessing the Application

In development, the backend server proxies requests to the Vite development server. You can access the application at:

- **Frontend**: http://localhost:8080
- **Backend API**: http://localhost:8080/api

The backend automatically forwards frontend requests to the Vite dev server running on port 5173, so you only need to access the application through the backend port (8080).


### Database Development Tips

When working on database schema changes or migrations, it's recommended to use an in-memory database instead of a file-based database for faster development iterations. You can set this up by:

1. Using SQLite in-memory mode:
   ```bash
   DB_PATH=":memory:" go run .
   ```
   ```bash
   DB_PATH=":memory:" air
   ```

2. Or modify the environment variable:
   ```bash
   export DB_PATH=":memory:"
   ```

This approach provides:
- Faster database operations
- No file I/O overhead
- Automatic cleanup between runs
- Easier testing and development

### Production Deployment

For production deployment, you can serve the built frontend files directly from the backend by setting the `SERVE_STATIC=true` environment variable. This eliminates the need for a separate frontend server and proxy configuration.
