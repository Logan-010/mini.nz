# Mini.nz - Encrypted File Upload Service

Mini.nz is a lightweight file upload service written in Go that provides encrypted file storage.

## Features

- **Encryption**: Mini.nz ensures user privacy by employing file encryption.

- **Simple Setup**: Get started with Mini.nz in a few simple steps.

- **Gzip compression**: Lightweight executable and no runtime dependencies.

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
