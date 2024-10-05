# Merchant-Bank API

This API provides a simple simulation for transactions between customers and merchants. It includes functionalities for customer login, payment, and logout, with a focus on secure JWT authentication and logging all activities to a history file. Data is stored in JSON files for simulation purposes.

Kindly check the documentation listed below.
- Postman documentation: https://documenter.getpostman.com/view/32334876/2sAXxMftR5

## Features
- Login: Customers can log in using a username and password, generating a JWT token for authentication.
- Payment: Logged-in customers can make payments to merchants. The transaction is logged, and the merchant's balance is updated.
- Logout: Logs out the customer and invalidates the JWT token.
- Register Customer: Create a new customer with a hashed password for security.
- Register Merchant: Create a new merchant with a unique ID and initial balance.

## Technology
- Golang: Backend API
- JWT: Authentication
- bcrypt: Password hashing
- gorilla/mux: HTTP routing

## Code Explanation
Main.go is the entry point of the application. It sets up the HTTP routes and starts the server.
- Router Setup: Using gorilla/mux to define routes such as /login, /payment, /logout, /create-customer, and /create-merchant.
- Middleware: Includes JWT authentication middleware to protect certain routes.
- Starting the Server: The server listens on port 8081.

## Getting Started
### Prerequisites
Install Go

### Installation
1. Clone the repository:
```shell
git clone https://github.com/your-username/merchant-bank-api.git
cd merchant-bank-api
```
2. Install Go modules and dependencies:
```shell
go mod tidy
```

### Running the Application
Start Go server:
```shell
go run main.go
```

The server will start at http://localhost:8081