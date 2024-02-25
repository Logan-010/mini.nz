# Mini.nz - End-to-End Encrypted File Upload Service

Mini.nz is a lightweight file upload service written in Go that focuses on providing end-to-end encryption for user privacy. With Mini.nz, you can easily upload and share files while ensuring that your data remains confidential.

## Features

- **End-to-End Encryption**: Mini.nz ensures the privacy of your files by employing end-to-end encryption.

- **Simple Setup**: Get started with Mini.nz in a few simple steps. Just clone the repository, build the application, and run it. The provided Dockerfile also allows for easy deployment in containerized environments.

## Getting Started

Follow these steps to set up Mini.nz on your local machine:

1. Clone the repository:

    ```bash
    git clone https://github.com/Logan-010/mini.nz
    cd mini.nz
    ```

2. Compile the Tailwind CSS styles:

    ```bash
    npx tailwindcss -i ./assets/input.css -o ./assets/style.css
    ```

3. Build the Mini.nz executable:

    ```bash
    go build -ldflags "-s -w" -o ./mini.nz
    ```

4. Run Mini.nz:

    ```bash
    ./mini.nz
    ```

5. Visit http://localhost:8080 to access Mini.nz in your browser.

## Docker Support

Mini.nz comes with an included Dockerfile for easy containerization. Build and run Mini.nz in a Docker container using the following commands:

```bash
docker build -t mini.nz .
docker run -p 8080:8080 mini.nz
```

Visit http://localhost:8080 to access Mini.nz in your browser.
