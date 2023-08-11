# Go-DDD: Domain Driven Design Template in Golang

Welcome to `go-ddd`, a reference implementation/template repository demonstrating the Domain Driven Design (DDD) approach in Golang. This project aims to help developers and architects understand the DDD structure, especially in the context of Go, and how it can lead to cleaner, more maintainable, and scalable codebases.

## Overview

Domain-Driven Design is a methodology and design pattern used to build complex enterprise software by connecting the implementation to an evolving model. `go-ddd` showcases this by setting up a simple marketplace where `Sellers` can sell `Products`.

### Why DDD?

- **Ubiquitous Language**: Promotes a common language between developers and stakeholders.
- **Isolation of Domain Logic**: The domain logic is separate from the infrastructure and application layers, promoting SOLID principles.
- **Scalability**: Allows for easier microservices architecture transitions.

## Repository Structure

- `domain`: The heart of the software, representing business logic and rules.
    - `entities`: Fundamental objects within our system, like `Product` and `Seller`.
- `application`: Contains use-case specific operations that interact with the domain layer.
- `infrastructure`: Supports the higher layers with technical capabilities like database access.
    - `db`: Database access and models.
    - `mappers`: Converters between domain entities and database models.
    - `repositories`: Concrete implementations of our storage needs.
- `interface`: The external layer which interacts with the outside world, like API endpoints.
    - `api/rest`: Handlers or controllers for managing HTTP requests and responses.

## Getting Started

1. Clone this repository:
```bash
git clone https://github.com/sklinkert/go-ddd.git
cd go-ddd
go mod download
go run main.go
```

### Contributions
Contributions, issues, and feature requests are welcome! Feel free to check the issues page.

### License
Distributed under the MIT License. See LICENSE for more information.

### Acknowledgments
Eric Evans, for introducing Domain-Driven Design.