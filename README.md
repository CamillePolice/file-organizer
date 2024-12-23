# Go Development Container

This project provides a development container for building and running Go applications. It includes a pre-configured environment with all necessary tools and dependencies.

## Project Structure

```
go-dev-container
├── .devcontainer
│   ├── devcontainer.json       # Configuration for the development container
│   └── Dockerfile              # Docker image definition
├── src
│   └── main.go                 # Entry point of the Go application
└── README.md                   # Project documentation
```

## Getting Started

To get started with this project, follow these steps:

1. **Clone the Repository**
   ```
   git clone <repository-url>
   cd go-dev-container
   ```

2. **Open in Development Container**
   Open the project in your code editor and use the command to reopen in the development container. This will set up the environment as specified in the `.devcontainer` folder.

3. **Build the Application**
   Once inside the container, navigate to the `src` directory and build the application:
   ```
   cd src
   go build
   ```

4. **Run the Application**
   After building, you can run the application:
   ```
   go run main.go
   ```

## Usage

You can modify the `main.go` file to implement your application logic. The development container will provide all necessary tools for Go development, including linting and formatting.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.