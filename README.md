# Go GraphQL Backend Boilerplate

Welcome to the Go GraphQL Backend Boilerplate! This repository provides a solid foundation for building scalable and maintainable GraphQL applications using Go, following Clean Architecture and Domain-Driven Design (DDD) principles.

## Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.21 or later
- PostgreSQL database

## Getting Started

Follow these steps to get your local development environment running:

### 1. Clone the Repository

Clone this repository and navigate to the project directory:

```bash
git clone <your-repo-url>
cd gobase
```

### 2. Configuration

Copy the example configuration file and update it with your local settings:

```bash
cp config/file/config.local.yaml.example config/file/config.local.yaml
```

Edit `config/file/config.local.yaml` to customize the configuration, particularly the database connection string.

### 3. Build and Run

To build and run the service locally:

```bash
# Generate wire injection code
make wire

# Start the service
make start

# Alternatively, use the freshstart command to do all of the above
make freshstart
```

The GraphQL server will be available at `http://localhost:6080/query` by default.

## Project Architecture

This boilerplate is built using Clean Architecture principles, structured to maintain a strong separation of concerns and to support microservices pattern.

### Clean Architecture Overview

This project follows a modified Clean Architecture pattern with these key principles:

1.  **Dependency Rule**: Dependencies always point inward. Inner layers (domain) have no knowledge of outer layers (infrastructure, interface).
2.  **Domain-Centric**: All business logic is central to the application and is independent of any framework or external system.
3.  **Separation of Concerns**: Clear boundaries between different components of the system (e.g., database logic, business rules, API definitions).
4.  **Testability**: Business logic can be tested in isolation without needing a UI, database, or external dependencies.

### Project Structure

```
/gobase
├── cmd/                     # Application entrypoints
├── config/                  # Application configuration files
├── di/                      # Dependency injection setup (using Google Wire)
│   ├── container/           # Wire-generated DI container
│   ├── provider/            # Component providers for DI
│   └── registry/            # Registry types for DI
├── graphql/                 # GraphQL schema, generated code, and resolvers
├── internal/                # Private application code
│   ├── domain/              # Business domains (e.g., product, user)
│   │   ├── <domain_name>/   # A specific domain
│   │   │   ├── repository/  # Data access layer interface
│   │   │   ├── usecase/     # Business logic layer
│   │   │   └── dto/         # Data Transfer Objects for the domain
│   ├── db/                  # Database schema and entity definitions
│   └── pkg/                 # Internal shared packages
│       ├── helper/          # Utility helpers
│       └── service/         # Cross-cutting services (e.g., CRUD helpers)
└── transport/               # API transport layers (e.g., REST)
```

### Key Components

-   **Domain Layer**: Each business domain (e.g., `product`) contains its core logic, including Repositories (interfaces), Use Cases (business rules), and DTOs.
-   **Infrastructure Layer**: Contains implementations for external concerns. This includes repository implementations using `Bun ORM` for database operations.
-   **Interface Layer**: The `graphql` directory holds all GraphQL-related code, including schemas (`.graphql`), generated models, and resolver implementations that connect the schema to the use cases.
-   **Dependency Injection**: The application uses `Google Wire` to manage dependencies. Providers are defined in `di/provider` and assembled in `di/container`.

### Microservices Patterns

This boilerplate is designed as a well-structured microservices architecture:

1.  **Clear Boundaries**: Each domain is self-contained with well-defined interfaces, making it easier to extract into a separate service later.
2.  **Domain-Driven Design**: Business capabilities are organized into distinct domains, aligning the code structure with the business structure.

## License

[Your License Information]