# CVWO Project Backend

This is the backend for the CVWO (Complete Virtual Work Environment) project, built with Go and Gin.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Running the Server](#running-the-server)
  - [Endpoints](#endpoints)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

Make sure you have the following software installed:

- Go: [Installation Guide](https://golang.org/doc/install)
- PostgreSQL: [Download Page](https://www.postgresql.org/download/)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/cvwo-project.git
   ```

2. Navigate to the project directory:

   ```bash
   cd cvwo-project
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Set up the PostgreSQL database:
   - Create a new database and update the `database/db.go` file with your database credentials.

5. Build the project:

   ```bash
   go build -tags netgo -ldflags '-s -w' -o app
   ```

## Usage

### Running the Server

Run the server with the following command:

```bash
./app
```

The server will start on `http://0.0.0.0:8082`.

### Endpoints

## Post Handlers
- **POST /posts** Create a post
- **GET /posts/:id** Get a post by ID
- **GET /posts** Get all posts
- **DELETE /posts/:id** Delete a post by ID
- **PATCH /posts/:id**  Update a post by ID

## User Handlers
- **POST /register"** Register a user
- **POST /login"** Login for user
- **GET /userprofile/:id"** Get user by ID
- **DELETE /deleteuser/:id"** Delete user by ID

## Like Handlers
- **POST /likes** Create a like
- **DELETE /likes** Delete a like
- **GETlikecount /:post_id** Get like count by post ID

## Comment Handlers
- **POST /createcomment** Create a comment
- **PATCH /updatecomment/:comment_id** Update a comment by comment ID
- **DELETE /deletecomment/:comment_id** Delete a comment by comment ID
- **GET /getcommentsbypostid/:post_id** Get comment count by post ID

For more detailed information about each endpoint, please refer to the source code or use API documentation tools.

## Contributing

Contributions are welcome! If you'd like to contribute to the project, please follow the [contribution guidelines](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).

---

Feel free to customize the README according to your project's specific details. You can also add sections for API documentation, deployment instructions, and anything else relevant to your project.