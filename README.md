
# Basic Shop Backend with Golang

This project is a basic backend for an online shop, built with Go. It provides RESTful API endpoints for product management and order processing, with access control for admin operations.

## Tech Stack

- **Golang**: The main programming language used to build this backend.
- **Gin**: A web framework for creating HTTP servers in Go, used to manage routes and middleware.
- **PostgreSQL**: The relational database for storing products, orders, and other data.
- **pgx**: PostgreSQL driver for Go, used for database connections.
- **Docker**: Used to run the PostgreSQL instance locally for development and testing.

## Usage

After setting up the environment variables and starting the PostgreSQL container, follow these steps to run the application.

Start the Go application by running:

```bash
go run main.go
```

The server will be accessible at `http://localhost:8080`.

## API Endpoints

Here is a list of available API endpoints and their descriptions:

### Public Endpoints

- **Get All Products**
  - **Endpoint**: `GET /api/v1/products`
  - **Description**: Retrieves a list of all available products.

- **Get Product by ID**
  - **Endpoint**: `GET /api/v1/products/:id`
  - **Description**: Retrieves detailed information about a specific product by ID.

- **Checkout Order**
  - **Endpoint**: `POST /api/v1/checkout`
  - **Description**: Processes an order checkout with the provided order data.

- **Confirm Order**
  - **Endpoint**: `POST /api/v1/orders/:id/confirm`
  - **Description**: Confirms an order by its ID.

- **Get Order by ID**
  - **Endpoint**: `GET /api/v1/orders/:id`
  - **Description**: Retrieves details about a specific order by ID.

### Admin Endpoints

The following endpoints are restricted to admin users and require the `AdminOnly` middleware to ensure only authorized users can perform these actions.

- **Create Product**
  - **Endpoint**: `POST /admin/products`
  - **Description**: Adds a new product to the inventory.

- **Update Product**
  - **Endpoint**: `PUT /admin/products/:id`
  - **Description**: Updates information for a specific product.

- **Delete Product**
  - **Endpoint**: `DELETE /admin/products/:id`
  - **Description**: Removes a product from the inventory.

## Project Structure

```
├── handler             # Contains route handlers for various endpoints
│   ├── order.go        # Order-related endpoints
│   └── product.go      # Product-related endpoints
├── middleware          # Middleware for route protection
│   └── admin.go        # Middleware to restrict access to admin routes
├── model               # Database models for products and orders
│   ├── order.go
│   └── product.go
├── main.go             # Entry point of the application
├── migrate.go          # Handles database migrations
├── go.mod              # Go module file
├── go.sum              # Go dependencies file
└── README.md           # Project documentation
```

## Setup and Installation

1. **Run PostgreSQL**

   Start a PostgreSQL instance (using Docker or a local installation) for database management.

2. **Export Environment Variables**

   Set up the environment variables for database connection and admin secret:

   ```bash
   export DB_URI=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
   export ADMIN_SECRET=secret
   ```

3. **Run the Application**

   Run the application with:

   ```bash
   go run main.go
   ```

