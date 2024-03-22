# Gridiron

[GoLang](https://go.dev/learn/) Modular Monolith designed to track American Football data.

## Setup

### Offline - Locally

1. Run Gridiron 

    Run the Gridiron docker container with
    ```bash
    docker compose -f ./deployments/docker-compose.yml up -d gridiron-app
    ```

    You can rebuild the and run the app with the following
    ```bash
    docker compose -f ./offline/docker-compose.yml up -d --no-deps --build gridiron-app
    ```


#### Clean Up

You can easily clean up your local environment with the following...

1. Mark the bash script as executable.

    ```bash
    chmod +x scripts/docker-nuke.sh
    ```

2. Run script to nuke environment

    ```bash
    ./scripts/docker-nuke.sh
    ```

## Technical Design

Gridiron follows the [Project Layout](https://github.com/golang-standards/project-layout) specified by the golang standards.

### Flaws in Design
* Dockerfile should be placed at `./build` and main.go should be placed at `./cmd/gridiron/gridiron.go`
    * Files are instead at root.
    * During initial development, there were problems getting things working. Instead of investing the time in this small item, I chose to move it to the root dir and move forward with the project.

## License

[ISC License](LICENSE)
