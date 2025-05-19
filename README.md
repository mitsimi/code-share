# Code Share

A web application for sharing and managing code snippets, built with Go and Vue.js.

## Project Structure

```
.
├── frontend/          # Vue.js frontend application
│   ├── src/          # Source files
│   ├── dist/         # Built frontend files
│   └── ...
├── internal/         # Backend application
│   ├── api/         # HTTP handlers
│   │   └── snippets.go
│   ├── models/      # Data models
│   │   └── snippet.go
│   ├── server/      # Server setup and routing
│   │   └── server.go
│   └── storage/     # Storage interfaces and implementations
│       ├── storage.go
│       └── memory.go
├── pkg/             # Shared packages
│   ├── middleware/  # Custom middleware
│   └── utils/       # Utility functions
└── main.go         # Application entry point
```

## Features

- Share code snippets with others
- Like and unlike snippets with real-time updates
- Modern, responsive UI with loading states
- RESTful API with proper error handling
- Client-side caching with TanStack Query

## API Endpoints

### Snippets

- `GET /api/snippets` - Get all snippets
- `GET /api/snippets/{id}` - Get a specific snippet
- `POST /api/snippets` - Create a new snippet
- `PUT /api/snippets/{id}` - Update a snippet
- `DELETE /api/snippets/{id}` - Delete a snippet
- `PATCH /api/snippets/{id}/like?action=like|unlike` - Like or unlike a snippet

### Request/Response Format

#### Create/Update Snippet

```json
{
  "title": "string",
  "content": "string",
  "author": "string"
}
```

#### Snippet Response

```json
{
  "id": "string",
  "title": "string",
  "content": "string",
  "author": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "likes": "number"
}
```

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

### Frontend

The frontend is built with Vue.js and uses:

- Vue 3 with Composition API
- TypeScript for type safety
- Vite for build tooling
- Tailwind CSS for styling

## Future Improvements

- [ ] Add database integration
- [ ] Add syntax highlighting

Optional:

- [ ] Add user authentication
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
