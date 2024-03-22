# Gridiron

[GoLang](https://go.dev/learn/) Modular Monolith designed to track American Football data.

[ChatGPT 3.5](https://openai.com/blog/chatgpt) was heavily used for "busy" or "repetitive" work.

## Setup

### Offline - Locally

1. Run Gridiron 

    Run the Gridiron docker container with
    ```bash
    docker compose -f ./deployments/docker-compose.yml up -d gridiron-app
    ```

    You can rebuild the and run the app with the following
    ```bash
    docker compose -f ./deployments/docker-compose.yml up -d --no-deps --build gridiron-app
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

## Testing

### Unit Test
`./internal` & `./pkg` contain unit tests. You can run them with the following.

```bash
go test ./pkg/... ./internal/...
```

### Integration Test
The `main.go` app and `./api` are tested with integration test.

```bash
go test ./test/...
```

## Technical Design

[Technical Design Document](docs/TECHNICAL_DESIGN.md)

## Contributing

This project follows [Feature branch workflow](https://docs.gitlab.com/ee/gitlab-basics/feature_branch_workflow.html)

## License

[ISC License](LICENSE)
