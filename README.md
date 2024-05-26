
This project is a Nakama module written in Go. It includes functions to interact with a database, read file contents, and perform RPC (Remote Procedure Call) functions. The project ensures that the necessary database tables exist and allows for data insertion via an RPC function.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Running the Application](#running-the-application)
- [Testing the Application](#testing-the-application)
- [Project Structure](#project-structure)
- [License](#license)

## Installation

### Prerequisites

- Go (1.19 or later)
- Docker
- Docker Compose
- Nakama Server
   ```

## Running the Application

1. Run Nakama server with Docker Compose:
    ```sh
    go mod vendor
    docker compose up --build
    ```

This command will start the Nakama server on localhost:7351 along with the postgres database, and load the Go plugin you have built.

## Usage

### RPC Function

The primary RPC function provided by this module is `my_rpc_function`. This function reads a file based on the provided payload, computes its hash, and stores the information in the database.

### Example Payload

```json
{
    "type": "core",
    "version": "1.0.0"
}
```
### Test via UI
To test the RCP function goto ui > localhost:7351 -> RUNTIME MODULES -> Registered RPC Functions -> select my_rpc_function -> click on API EXPLORE

### API Explorer

select `my_rpc_function` click on `send request`

### Run Unit Test
    ```sh
    go test
    
    ```

### Thoughts on the Task
This task provides a solid example of integrating custom server logic with Nakama using Go. It demonstrates how to handle file operations, database interactions, and RPC calls, which are common requirements in game server development. The use of Docker Compose simplifies the deployment process, making it easier to set up and run the Nakama server with the custom module.

### Ideas for Improvement
If I had more time, I would consider the following improvements:

Enhanced Error Handling: Improve error messages and add more granular error handling to cover edge cases.
Configuration Management: Externalize configuration parameters (like file paths and database credentials) to a configuration file or environment variables.
Logging Enhancements: Implement a more robust logging mechanism to capture detailed logs and possibly integrate with a logging service.
Unit Test Coverage: Increase test coverage to include more edge cases and different scenarios.
Continuous Integration: Set up a CI pipeline to automate the testing and deployment process.
Performance Optimization: Profile and optimize the code for performance, especially the file reading and hashing logic.
Documentation: Expand the documentation to include more detailed usage examples, API documentation, and troubleshooting tips.