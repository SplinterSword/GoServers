# GoServers

A RESTful API server built with Go, implementing user authentication, authorization, and chirp (microblog) functionality. This project serves as a learning exercise for building web servers in Go.

## Features

- **User Authentication**: JWT-based authentication system
- **Chirps**: Create, read, and delete short messages (chirps)
- **Database**: File-based JSON database for data persistence
- **API Metrics**: Track and monitor API usage
- **Webhooks**: Integration with external services via webhooks

## API Endpoints

### Health Check
- `GET /api/healthz`
  - Check if the server is running
  - Returns: 200 OK with status message

### Authentication
- `POST /api/users`
  - Register a new user
  - Required fields: `email`, `password`
  - Returns: User object (without password)

- `POST /api/login`
  - Authenticate a user and get JWT tokens
  - Required fields: `email`, `password`
  - Returns: User object, access token, and refresh token

- `POST /api/refresh`
  - Refresh an access token using a refresh token
  - Requires: Valid refresh token in Authorization header
  - Returns: New access token

- `POST /api/revoke`
  - Revoke a refresh token
  - Requires: Valid refresh token in Authorization header
  - Returns: 200 OK on success

- `PUT /api/users`
  - Update user information
  - Requires: Valid access token in Authorization header
  - Returns: Updated user object

### Chirps
- `POST /api/chirps`
  - Create a new chirp
  - Requires: Valid access token in Authorization header
  - Required fields: `body` (max 140 characters)
  - Returns: Created chirp object

- `GET /api/chirps`
  - Get all chirps
  - Optional query params:
    - `author_id`: Filter chirps by author ID
    - `sort`: Sort order (`asc` or `desc`)
  - Returns: Array of chirp objects

- `GET /api/chirps/{chirpID}`
  - Get a specific chirp by ID
  - Returns: Chirp object or 404 if not found

- `DELETE /api/chirps/{chirpID}`
  - Delete a chirp
  - Requires: Valid access token in Authorization header
  - User must be the author of the chirp
  - Returns: 200 OK on success

### Admin
- `GET /admin/metrics`
  - Get server metrics
  - Returns: Number of hits to the server

- `GET /api/reset`
  - Reset the server metrics counter
  - Returns: 200 OK with reset confirmation

### Webhooks
- `POST /api/polka/webhooks`
  - Handle Polka webhook events
  - Requires: API key in Authorization header
  - Currently supports user upgrade/downgrade events

## Setup

1. Clone the repository
2. Install Go dependencies:
   ```
   go mod download
   ```
3. Create a `.env` file with the following variables:
   ```
   JWT_SECRET=your_jwt_secret_here
   POLKA_KEY=your_polka_key_here
   ```
4. Run the server:
   ```
   go run .
   ```

## Development

- The server uses a file-based JSON database (`database.json`) that is automatically created on first run
- Use the `/api/reset` endpoint to reset the database:
  ```
  curl -X GET http://localhost:8080/api/reset
  ```

## License

This project is for learning purposes only.
