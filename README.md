# Fiber Bootstrapped - Go

> Fiber Bootstrapped: Your Comprehensive Toolkit for Go Projects, with a Single Codebase.

## Features

- [x] Authentication (JWT Auth)
- [x] User Management
  - [x] Email Verification
  - [x] Password Reset
- [x] Role-based access control (user/admin)
- [ ] Email Notifications (Custom Mailer)
- [ ] Distributed Tasks Queue, Scheduled Jobs ([asynq](https://github.com/hibiken/asynq))

## Project Overview

This project consists of the following main components:

1. **Main File**: `main.go` - Entry point of the application.

2. **Source Code Directory** (`src`):

   - This directory contains the source code of the application.
   - It's organized into subdirectories based on different modules and components.

3. **App Module** (`app`):

   - Contains the core functionalities of the application.
   - Subdirectories:
     - `events`: Event handling related code.
     - `helpers`: Utility functions for error handling and middleware.
     - `hooks`: Hooks for service functionalities.
     - `modules`: Modules like mailer.
     - `schemas`: Schemas for defining data structures, organized by entity types.

4. **Services** (`services`):

   - Contains business logic for various services.
   - Subdirectories:
     - `auth`: Authentication related services.
     - `users`: User related services.

5. **Core Components** (`core`):
   - Core functionalities of the application.
   - Subdirectories:
     - `app`: Custom app functionalities.
     - `configuration`: Configuration handling.
     - `database`: Database initialization.
     - `events`: Event handling core.
     - `server`: Server setup and initialization.
     - `service`: Core service functionalities.

## Project Directory Structure

```
.
├── go.mod
├── go.sum
├── main.go
└── src
├── app
│ ├── app.go
│ ├── events
│ │ └── service.events.go
│ ├── helpers
│ │ ├── error.helper.go
│ │ └── middleware.helper.go
│ ├── hooks
│ │ └── service.hooks.go
│ ├── modules
│ │ └── mailer.module.go
│ ├── schemas
│ │ ├── auth
│ │ │ ├── auth.schema.go
│ │ │ └── manage
│ │ │ └── auth_manage.schema.go
│ │ └── users
│ │ └── users.schema.go
│ ├── services
│ │ ├── auth
│ │ │ ├── build
│ │ │ │ └── auth.build.go
│ │ │ ├── controllers
│ │ │ │ └── auth.controller.go
│ │ │ └── utils
│ │ │ └── auth.utils.go
│ │ ├── services.go
│ │ └── users
│ │ ├── build
│ │ │ └── users.build.go
│ │ └── controllers
│ │ └── users.controller.go
│ └── utils
│ └── shared.util.go
└── core
├── app.core.go
├── configuration.core.go
├── database.core.go
├── events.core.go
├── server.core.go
└── service.core.go
```

## Todo

- [ ] Add more data validation ([validator](https://pkg.go.dev/github.com/go-playground/validator/v10)).
- [ ] Support for logging to files, databases or external services.
- [ ] Publish Create/Read/Update/Delete events on service method calls.
- [ ] Support for bulk Create/Update/Delete operations.
- [ ] Support for MongoDB Aggregation Queries via Service interface.
- [ ] WebSockets or Server-Sent Events (SSE) support for real-time communication.
- [ ] Unit tests and end-to-end tests.
- [ ] Dockerize project.

## Usage

1. Make sure you have [Go](https://go.dev/) (and [MongoDB](https://www.mongodb.com/) for local instances) installed.

2. Install your dependencies.

   ```bash
   go mod vendor
   ```

3. Configuring the server with environment variables

   - Create a `.env` file in the root
   - Copy the values from `.env.sample` into the `.env` file and populate it accordingly.

4. Start your server.

```bash
 go run main.go
```

## Testing

_Implement Tests_

## Contributing

Contributions are welcome. Please follow the existing code style and conventions.

## Credits

### Acknowledgements

The project architecture and codebase is heavily inspired by [feathersjs](https://www.feathersjs.com/).

## License

This project is licensed under the [MIT License](LICENSE).
