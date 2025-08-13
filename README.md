# OTP Auth Service

A backend service in Golang that implements OTP-based login and registration, with rate limiting, user management, Swagger API documentation, and Docker support.

---

## Features

- **OTP Login & Registration**
  - Generate and verify one-time passwords (OTPs)
  - OTP expiration: 2 minutes
  - OTP rate limit: max 3 requests per phone number per 10 minutes
  - Stores OTP in Redis

- **User Management**
  - Retrieve a single user
  - Retrieve a paginated list of users with optional search by phone number
  - Stores user data in PostgreSQL

- **Security**
  - JWT authentication for protected routes
  - Swagger UI supports Bearer Token input

- **Documentation**
  - Swagger/OpenAPI docs auto-generated
  - API docs available at `/swagger/index.html`

- **Containerization**
  - Dockerized application with PostgreSQL and Redis in `docker-compose.yml`

---

## Tech Stack

- **Language**: Golang
- **Database**: PostgreSQL (GORM ORM)
- **Cache/OTP Storage**: Redis
- **Framework**: Gin
- **Auth**: JWT
- **Docs**: Swaggo/Swagger

---

## Project Structure

