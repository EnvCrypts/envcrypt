# EnvCrypt Server

The backend service for EnvCrypt, a secure, end-to-end encrypted environment variable management system. This server acts as a blind storage and coordination layer, enforcing access control and storing encrypted blobs without ever having access to the raw secrets or encryption keys.

## Features

-   **Zero-Knowledge Architecture**: The server never sees plaintext secrets. It only stores encrypted data.
-   **Granular Access Control**: Manages user roles (Admin/Member) and project permissions.
-   **Service Role Management**: Supports machine identities (CI/CD) with delegated access.
-   **RESTful API**: Provides endpoints for user management, project coordination, and blob storage.
-   **Database**: Uses PostgreSQL for robust data persistence.

## Architecture

The server handles:
1.  **User Identity**: Storing public keys and encrypted private key backups.
2.  **Project Metadata**: Tracking project ownership and membership.
3.  **Key Exchange**: Facilitating the exchange of wrapped project keys between users.
4.  **Blob Storage**: Storing immutable versions of encrypted environment variables.

## Getting Started

### Prerequisites

-   Go 1.22+
-   PostgreSQL
-   Docker (optional, for containerized deployment)

### Configuration

The server is configured via environment variables. See `.env.example` (if available) or create a `.env` file with:

```bash
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/envcrypt?sslmode=disable
```

### Running Locally

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/envcrypt/main.go
```

### Docker Deployment

```bash
docker build -t envcrypt-server .
docker run -p 8080:8080 --env-file .env envcrypt-server
```


## Contributing

1.  Fork the repository
2.  Create your feature branch (`git checkout -b feature/amazing-feature`)
3.  Commit your changes (`git commit -m 'Add some amazing feature'`)
4.  Push to the branch (`git push origin feature/amazing-feature`)
5.  Open a Pull Request
