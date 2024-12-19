# README

## Project Overview
This project is a basic **CRUD (Create, Read, Update, Delete)** web application built with Go. It provides a RESTful API for managing users and includes an example frontend that interacts with the backend. It uses **Gorilla Mux** for routing, **GORM** for database interaction, and **PostgreSQL** as the database.

## Features
- Create a new user
- Fetch all users
- Retrieve a user by ID
- Update user information
- Delete a user
- Basic frontend interface for managing users
- Handles JSON-based GET and POST requests on the root endpoint

## Prerequisites
1. Go 1.23 or higher
2. PostgreSQL database
3. Node.js (optional, for frontend-related modifications)

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/crud-go-app.git
   cd crud-go-app
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the database:
   - Ensure PostgreSQL is running.
   - Update the `dsn` in `initDB()` to match your database credentials:
     ```go
     dsn := "host=localhost user=postgres password=0000 dbname=postgres port=5433 sslmode=disable TimeZone=Asia/Almaty"
     ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. Access the application:
   - Backend API: [http://localhost:8080](http://localhost:8080)
   - Frontend UI: [http://localhost:8080](http://localhost:8080)

## API Endpoints
### Users Endpoints
- **POST `/users`**  
  Create a new user.  
  **Body (JSON):**
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com"
  }
  ```

- **GET `/users`**  
  Retrieve a list of all users.

- **GET `/users/{id}`**  
  Retrieve a user by ID.

- **PUT `/users/{id}`**  
  Update a user's information.  
  **Body (JSON):**
  ```json
  {
    "name": "Updated Name",
    "email": "updated.email@example.com"
  }
  ```

- **DELETE `/users/{id}`**  
  Delete a user by ID.

### Root Endpoints
- **GET `/`**  
  Returns a greeting message.
  
- **POST `/`**  
  Accepts a JSON message.  
  **Body (JSON):**
  ```json
  {
    "message": "Your custom message"
  }
  ```

## Database
- The application uses **GORM** for database interaction and automatically migrates the `User` table during startup.
- User model:
  ```go
  type User struct {
      ID    uint   `gorm:"primaryKey"`
      Name  string `json:"name"`
      Email string `json:"email"`
  }
  ```

## Frontend
The basic HTML file is located in the root directory and provides a simple interface for managing users. It uses JavaScript to interact with the API.

- To fetch all users: Calls `/users` endpoint.
- To create a user: Prompts for `name` and `email` and sends a POST request to `/users`.

## Dependencies
- [Gorilla Mux](https://github.com/gorilla/mux) - Router for handling HTTP requests.
- [GORM](https://gorm.io/) - ORM library for database operations.
- [PostgreSQL Driver](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL) - PostgreSQL support for GORM.

## File Structure
```
├── main.go         # Main application code
├── go.mod          # Module dependencies
├── go.sum          # Dependency checksums
├── index.html      # Basic frontend UI
```

## Running Tests
To test the API:
1. Use tools like [Postman](https://www.postman.com/) or [cURL](https://curl.se/).
2. Ensure the database is properly set up before testing.

## Troubleshooting
- **Database connection errors**: Verify the `dsn` string in the `initDB()` function.
- **Port conflicts**: Ensure port `8080` is not in use or change it in the `ListenAndServe` function.

## License
This project is licensed under the MIT License.
