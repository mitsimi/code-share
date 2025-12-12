# Code Share

A full‑stack code snippet platform built as an **engineering‑focused project** to explore
backend architecture, authentication, and frontend state management beyond typical course
requirements.

This project was originally developed for a university course and later expanded and refined
to experiment with **clean architecture, type‑safe database access, and modern frontend
patterns**.

## Why this project exists

Instead of building a minimal app to satisfy the course, I intentionally **overshot the scope**
to treat it as a realistic system:

- non‑trivial authentication and session handling
- clear separation between domain logic and transport
- type safety across the database, backend, and frontend
- real client‑side state management and caching

The goal was not to “ship a product”, but to **practice designing and reasoning about a system**
that stays understandable as it grows.

## High‑level overview

- **Backend**: Go REST API following a clean architecture approach
- **Frontend**: Vue 3 + TypeScript SPA with modern state and data‑fetching patterns
- **Database**: SQLite with SQLC for compile‑time checked queries
- **Auth**: JWT access tokens with refresh tokens and server‑side session tracking

## Backend architecture

The backend is structured into explicit layers:

- **Domain**  
  Core business entities and rules, independent of HTTP or storage.

- **Repositories**  
  Interfaces and implementations for persistence, isolating database concerns.

- **API / Handlers**  
  Thin HTTP layer responsible for request validation, authentication, and response mapping.

- **Infrastructure**  
  Routing, middleware, logging, and database connections.

Key design decisions:

- Used **SQLC** to eliminate an entire class of runtime SQL errors.
- Kept handlers small to prevent business logic from leaking into HTTP code.
- Explicit error types mapped consistently to HTTP responses.

## Frontend architecture

- **Vue 3** Composition API with strict TypeScript usage
- **TanStack Query** for server state and caching
- **Pinia** only for true client‑side state
- Form validation via **Zod** schemas, integrated with VeeValidate

The focus was on **predictable data flow** and minimizing implicit or duplicated state.

## Technical challenges explored

Some problems I deliberately spent time on:

- Designing authentication flows with refresh tokens and session invalidation
- Maintaining type safety from database queries to frontend components
- Avoiding tight coupling between API responses and UI state
- Structuring the project so new features do not require touching unrelated code

## What I would improve next

If I were to continue working on this project:

- Add integration and property‑based tests for critical backend paths
- Introduce explicit database migrations instead of schema initialization
- Implement rate‑limiting and more granular authorization rules
- Expand real‑time features beyond simple like/save interactions

## Tech stack

**Backend**
- Go
- Chi router
- SQLite + SQLC
- JWT, bcrypt
- Zap logging

**Frontend**
- Vue 3 + TypeScript
- Vite
- Tailwind CSS
- TanStack Query
- Pinia

## Running the project locally

```bash
# backend
go mod download
go run .

# frontend
cd frontend
pnpm install
pnpm dev
```

During development, the backend proxies requests to the frontend, so the application is
accessible via:

http://localhost:8080

## Project status

The project fully meets and exceeds the original course requirements and serves as a
long‑term playground for experimenting with architectural ideas and incremental
improvements.

While it is not developed against a fixed roadmap, I occasionally extend and refine it
(e.g. improving language support and syntax highlighting) when exploring new ideas or
techniques.
