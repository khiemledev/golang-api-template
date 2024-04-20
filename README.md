# Golang API Template

![pre-commit](https://github.com/khiemledev/golang-api-template/actions/workflows/pre-commit.yaml/badge.svg)

## Quick start

Prerequisite:

- [Golang](https://golang.org/doc/install) (>= 1.22.2)
- [NodeJs](https://nodejs.org/en/download/) (>= v20.6.1)
- [Make](https://www.gnu.org/software/make/) (>= 3.81)
- [Pre-commit](https://pre-commit.com/) (>= 3.5.0)

Install nodemon to auto reload when the code changed:

```bash
npx install -g nodemon

# Verify installation version
nodemon --version
```

Clone this repo and install dependencies:

1. Clone this repo

    ```bash
    git clone https://github.com/khiemledev/golang-api-template.git
    ```

2. Change directory into `golang-api-template`

    ```bash
    cd golang-api-template
    ```

3. Install dependencies Go dependencies

    ```bash
    go mod download
    ```

4. Start dev server

    ```bash
    make dev_server
    ```

5. Now you can vist [http://127.0.0.1:8080](http://127.0.0.1:8080) to check the server is running
    or using cURL

    ```bash
    curl http://127.0.0.1:8080
    ```


You can modify environment in the file `cmd/api/app.env`.



## Generate Swagger documentation

First, you need to install [swag](https://github.com/swaggo/swag#how-to-use-it-with-gin):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

You can generate Swagger documentation by running:

```bash
make docs
# or
swag init
```

## Before commit

1. First, you need to install [pre-commit](https://pre-commit.com/):

    ```bash
    brew install pre-commit
    # or
    pip install pre-commit
    ```

2. Install Go dev dependencies
    1. Install [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)

        ```bash
        go install https://pkg.go.dev/golang.org/x/tools/cmd/goimports
        ```

    2. Install [golangci-lint](https://golangci-lint.run/usage/install/)

        ```bash
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
        ```

3. Install pre-commit hooks

    ```bash
    pre-commit install
    ```

4. Verify installation

    ```bash
    pre-commit run --all-files
    ```


## Docker

To start the server using docker, using the following command:

```bash
# Create network
docker network create golang-api-template-network

# Start database
docker compose -f docker/docker-compose-postgres.yaml up -d

# Start the server
docker compose -f docker/docker-compose.yaml up -d
```

You can modify environment in the file `docker/app.env`.


## Todo

- [ ] Add cache for login session and user
- [ ] Add rate limit
- [ ] Revoke token
- [ ] Refresh token
