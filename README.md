# Code Share

A modern web application for sharing and managing code snippets, built with Go and Vue.js.

## Project Structure

```
.
├── frontend/          # Vue.js frontend application
│   ├── src/          # Source files
│   │   ├── components/  # Vue components
│   │   ├── stores/     # Pinia stores
│   │   ├── views/      # Page components
│   │   └── ...
│   ├── dist/         # Built frontend files
│   └── ...
├── internal/         # Backend application
│   ├── api/         # HTTP handlers
│   ├── models/      # Data models
│   ├── server/      # Server setup and routing
│   └── storage/     # Storage interfaces and implementations
├── data/            # Database migrations and seeds
├── yaak/           # Yaak configuration files
└── main.go         # Application entry point
```

## Features

- User authentication and authorization
- Share code snippets with others
- Like and unlike snippets with real-time updates
- Modern, responsive UI with loading states
- RESTful API with proper error handling
- Client-side caching with TanStack Query
- Secure API responses that protect user data
- Form validation with Zod and VeeValidate
- Toast notifications with Vue Sonner
- Database integration with SQLC

## API Endpoints

### Authentication

- `POST /api/auth/login` - User login
- `POST /api/auth/signup` - User registration
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh access token
- `GET /api/auth/me` - Get current user

### Snippets

- `GET /api/snippets` - Get all snippets
- `GET /api/snippets/{id}` - Get a specific snippet
- `POST /api/snippets` - Create a new snippet
- `PUT /api/snippets/{id}` - Update a snippet
- `DELETE /api/snippets/{id}` - Delete a snippet
- `PATCH /api/snippets/{id}/like?action=like|unlike` - Like or unlike a snippet

## Development

### Prerequisites

- Go 1.24 or later
- Node.js 22.14 or later
- pnpm (recommended) or npm

### Backend Setup

1. Install Go dependencies:

   ```bash
   go mod download
   ```

2. Run the server:
   ```bash
   # Using Air for hot reload
   air
   # Or directly
   go run main.go
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

## Project Architecture

### Backend

The backend follows a clean architecture pattern:

- **Models**: Data structures and validation
- **Storage**: Data persistence interface and implementations
- **API**: HTTP handlers and request/response handling
- **Server**: Routing and middleware setup
- **Database**: PostgreSQL with SQLC for type-safe queries

### Frontend

The frontend is built with Vue.js and uses:

- Vue 3 with Composition API
- TypeScript for type safety
- Vite for build tooling
- Tailwind CSS for styling
- Pinia for state management
- TanStack Query for data fetching
- VeeValidate with Zod for form validation
- Vue Sonner for toast notifications

## Technologies Used

### Backend

- Go
- PostgreSQL
- SQLC
- Air (for hot reload)
- Docker

### Frontend

- Vue 3
- TypeScript
- Vite
- Tailwind CSS
- Pinia
- TanStack Query
- VeeValidate
- Zod
- Vue Sonner

## Future Improvements

- [x] Add database integration
- [x] Add user authentication
- [ ] Add syntax highlighting

Optional:

- [ ] Add snippet categories/tags
- [ ] Add search functionality
- [ ] Add user profiles
- [ ] Add comments on snippets

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
