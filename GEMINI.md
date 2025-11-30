# Gemini Code-Assist Instructions for traefik-auth-manager

This document provides AI-native developer documentation to help you get productive in this codebase.

## 🚀 Project Overview

This is a Go application to manage credentials for Traefik's forward authentication service. It's a self-contained web application where both the backend API and the frontend Progressive Web App (PWA) are written in Go.

## 🏗️ Architecture: A Dual-Build System

The core architectural concept is a single codebase in `cmd/webapp` that compiles into two distinct artifacts using Go build tags: a native backend server and a WebAssembly (WASM) frontend.

-   **Backend (Native Go)**: An API server built with the [Echo](https.://echo.labstack.com/) framework.
    -   Build tags: `windows` or `linux` (default)
    -   Entrypoint logic: The `Backend()` function in `cmd/webapp/main.go`.
    -   API handlers are located in `internal/handlers`.

-   **Frontend (WASM)**: A Progressive Web App (PWA) built with the [go-app](https://go-app.dev/) framework.
    -   Build tags: `GOOS=js GOARCH=wasm`
    -   Entrypoint logic: The `Frontend()` and `app.RunWhenOnBrowser()` functions in `cmd/webapp/main.go`.
    -   The UI is composed of components. Pages are defined in `internal/pages` and reusable components are in `internal/components`.

This dual-build approach allows for maximum code reuse between the frontend and backend, particularly for data models (see `internal/models`).

## 🛠️ Developer Workflow

All build and run commands are managed via the `Makefile`. The build output is placed in the `./tmp` directory.

-   **Build the app (backend + frontend)**:
    ```bash
    make build
    ```

-   **Build and run the server**:
    ```bash
    make run
    ```

-   **Live Reload (for development)**: The `README.md` recommends using `wgo` for live reloading. This will watch for file changes and automatically rebuild and restart the server.
    ```bash
    # Requires wgo: go install github.com/bokwoon95/wgo@latest
    wgo -xdir tmp -file .go -file .css -file .js make run
    ```

-   **Build Docker Image**:
    ```powershell
    $env:DOCKER_BUILDKIT=1; docker build . --network=host --tag apogee-dev/traefik-auth-manager:local
    ```

## ✨ Code Conventions

-   **Frontend Components**: UI components are Go structs in `internal/components` that implement the `go-app` rendering logic. They are composed to build pages in `internal/pages`.
-   **Backend Handlers**: Standard `Echo` handlers are located in `internal/handlers`. They are responsible for API logic and rendering the initial HTML that loads the WASM application.
-   **Styling**: The application uses Bootstrap. Static assets, including CSS and JS libraries, are located in the `/web` directory and copied to `/tmp/web` during the build.

## 🐍 Running Python Scripts via Docker

Since Python may not be installed on the local Windows machine, and we are using a remote Docker host, use the following workflow to run Python scripts (e.g., for image processing):

1.  **Start a container**:
    ```powershell
    docker run -d --name python-runner -w /app --entrypoint tail python:3.9-slim -f /dev/null
    ```

2.  **Copy files to the container**:
    ```powershell
    docker cp ./path/to/script.py python-runner:/app/
    docker cp ./path/to/input_file python-runner:/app/
    ```

3.  **Execute the script**:
    ```powershell
    docker exec python-runner python script.py
    ```

4.  **Copy results back**:
    ```powershell
    docker cp python-runner:/app/output_file ./path/to/output/
    ```

5.  **Cleanup**:
    ```powershell
    docker rm -f python-runner
    ```
