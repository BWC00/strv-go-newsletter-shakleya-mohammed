# Go Newsletter API

This repository contains the source code for a Go Newsletter platform API. The API allows registered users to create and manage their own newsletters, which can be subscribed to by other users. The API serves both mobile apps and websites and can be accessed using REST API.

## Design Considerations

- The idiomatic structure based on the resource-oriented design.
- The usage of Docker, Docker compose, Alpine images, and linters on development.
- APIs are implemented according to the OpenAPI specifications.
- The usage of different storage services like [Firebase](https://github.com/firebase/firebase-admin-go) and [Postgres](https://www.postgresql.org/) to maintain the data.
- [Goose](https://github.com/pressly/goose) is utilized for smooth database migrations, while [GORM](https://gorm.io/) serves as the database ORM.
- User authentication uses stateless authorization with JWT.
- The usage of [SendGrid](https://sendgrid.com/) as mailing service for sending email confirmations.
- The usage of [Zerolog](https://github.com/rs/zerolog) as the centralized Syslog logger.
- The usage of [Validator.v10](https://github.com/go-playground/validator) as the form validator.
- The usage of GitHub actions to run tests and linters, generate OpenAPI specifications, and build and push production images to the Docker registry.

## Features

- Sign up and sign in with email and password.
- Create and manage newsletters.
- Subscribe to newsletters using email addresses.
- Receive confirmation emails upon subscription with an unsubscribe link.
- Unsubscribe from newsletters.

## Endpoints

| Name               | HTTP Method | Route                     |
|--------------------|-------------|---------------------------|
| Health             | GET         | /live                     |
|                    |             |                           |
| Login user         | POST        | /api/v1/users/login       |
| Register user      | POST        | /api/v1/users/register    |
|                    |             |                           |
| Subscribe          | POST        | /api/v1/subscriptions     |
| Unsubscribe        | GET         | /api/v1/subscriptions     |
| List Subscriptions | GET         | /api/v1/subscriptions/all |
|                    |             |                           |
| List Newsletters   | GET         | /api/v1/newsletters       |
| Create Newsletter  | POST        | /api/v1/newsletters       |
| Read Newsletter    | GET         | /api/v1/newsletters/{id}  |
| Update Newsletter  | PUT         | /api/v1/newsletters/{id}  |
| Delete Newsletter  | DELETE      | /api/v1/newsletters/{id}  |

💡 [swaggo/swag](https://github.com/swaggo/swag) : `swag init -g cmd/api/main.go -o .swagger -ot yaml`

## 🗄️ Database design

User:
| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | SERIAL    | X        | X           |
| firstname      | VARCHAR   | X        |             |
| lastname       | VARCHAR   | X        |             |
| email          | VARCHAR   | X        |             |
| password       | VARCHAR   | X        |             |
| description    | TEXT      |          |             |
| created_at     | TIMESTAMP | X        |             |

Newsletter:
| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | SERIAL    | X        | X           |
| editor_id      | INT       | X        |             |
| name           | VARCHAR   | X        |             |
| description    | TEXT      |          |             |
| created_at     | TIMESTAMP | X        |             |

## Container image sizes

- DB: 400MB
- API
    - Development environment: 655MB
    - Production environment: 28MB ; 💡`docker build -f prod.Dockerfile . -t app`

## 📁 Project structure

```shell
strv-go-newsletter-shakleya-mohammed
├── go.mod
├── go.sum
├── README.md
├── Dockerfile
├── prod.Dockerfile
├── docker-compose.yml
├── api
│   ├── middleware
│   │   ├── api_version_ctx.go
│   │   ├── authentication.go
│   │   ├── content_type_json.go
│   │   ├── content_type_json_test.go
│   │   ├── middleware.go
│   │   ├── validate_newsletter.go
│   │   ├── validate_subscription.go
│   │   └── validate_user.go
│   ├── requestlog
│   │   ├── handler.go
│   │   └── log_entry.go
│   └── resource
│       ├── newsletter
│       │   ├── handler
│       │   │   ├── handler.go
│       │   │   └── register.go
│       │   ├── model.go
│       │   └── repository
│       │       └── repository.go
│       ├── subscription
│       │   ├── handler
│       │   │   ├── handler.go
│       │   │   └── register.go
│       │   ├── model.go
│       │   └── repository
│       │       └── repository.go
│       └── user
│           ├── handler
│           │   ├── handler.go
│           │   └── register.go
│           ├── model.go
│           └── repository
│               └── repository.go
├── bin
│   ├── gofumpt
│   ├── staticcheck
│   └── swag
├── cmd
│   ├── api
│   │   └── main.go
│   └── migrate
│       └── main.go
├── config
│   ├── config.go
│   ├── database.go
│   ├── email.go
│   └── server.go
├── database
│   ├── firebase.go
│   └── postgres.go
├── migrations
│   ├── 00001_create_users_table.sql
│   └── 00002_create_newsletters_table.sql
├── server
│   ├── audit.go
│   ├── database.go
│   ├── email.go
│   ├── router.go
│   └── server.go
└── util
    ├── auth
    │   ├── hash.go
    │   └── token.go
    ├── email
    │   └── email.go
    ├── err
    │   └── err.go
    ├── logger
    │   ├── logger.go
    │   └── logger_test.go
    └── validator
        ├── validator.go
        └── validator_test.go
```
