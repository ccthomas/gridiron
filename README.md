# Gridiron

[GoLang](https://go.dev/learn/) Modular Monolith designed to track American Football data.

[ChatGPT 3.5](https://openai.com/blog/chatgpt) was heavily used for "busy" or "repetitive" work.

## Setup

### Offline - Locally

1. Configure Database

    Start docker container running Postgres local.
    ```bash
    docker compose -f ./deployments/docker-compose.yml up gridiron-db -d
    ```

1. Terraform DB

    Run terraform to configure database
    ```bash
    terraform -chdir=deployments/terraform/ init
    terraform -chdir=deployments/terraform/ plan
    terraform -chdir=deployments/terraform/ apply -auto-approve
    ```

    Run migrations
    ```
    migrate -path deployments/sql/migrations/ -database "postgresql://my_user:my_password@127.0.0.1:5432/my_db?sslmode=disable" -verbose up
    ```

1. Generate Auth Private Key

    **This is not necessarily a good practice. For the sake of local demo, this approach get's the job done**

    First, you need to generate a private key.
    Copy contents outputted from the file and assign it to `.env` variable `PRIVATE_KEY`

    ```bash
    openssl genpkey -algorithm RSA -out private_key.pem
    ```

    Then assign it to the environment 

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
* Note: Only some files have unit tests. Go makes it hard to mock or stub out functions. And the Go community seems to think unit tests are pointless in some areas, so there is no clear way to test that part of the stack.

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

### Migrations

This project uses [Migrate](https://github.com/golang-migrate/migrate) to manage migrations.

The following command can be used to generate a new migration.
```bash
migrate create -ext sql -dir deployments/sql/migrations/ -seq <name>
```

## License

[ISC License](LICENSE)
