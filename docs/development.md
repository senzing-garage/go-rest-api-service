# go-rest-api-service development

## Install Go

1. See Go's [Download and install](https://go.dev/doc/install)

## Install Senzing C library

Since the Senzing library is a prerequisite, it must be installed first.

1. Verify Senzing C shared objects, configuration, and SDK header files are installed.
    1. `/opt/senzing/g2/lib`
    1. `/opt/senzing/g2/sdk/c`
    1. `/etc/opt/senzing`

1. If not installed, see
   [How to Install Senzing for Go Development](https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/install-senzing-for-go-development.md).

## Install Git repository

1. Identify git repository.

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=go-rest-api-service
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Using the environment variables values just set, follow steps in
   [clone-repository](https://github.com/senzing-garage/knowledge-base/blob/main/HOWTO/clone-repository.md) to install the Git repository.

## Make a test database

1. Install
   [senzing-tools](https://github.com/senzing-garage/senzing-tools#install).
1. Create database.
   **Note:** The database location in the following example matches what's in the `Makefile`.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools init-database --database-url sqlite3://na:na@/tmp/sqlite/G2C.db
    ```

## Development cycle

Instructions are at
[Ogen QuickStart](https://ogen.dev/docs/intro/).

1. Get latest version of `ogen`

    ```console
    cd ${GIT_REPOSITORY_DIR}
    go get -d github.com/ogen-go/ogen
    ```

1. View version.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    go list -m github.com/ogen-go/ogen
    ```

1. Down latest version of
   [senzing-rest-api.yaml](https://raw.githubusercontent.com/senzing-garage/senzing-rest-api-specification/main/senzing-rest-api.yaml)
   to
   [restapiservice/openapi.yaml](https://github.com/senzing-garage/go-rest-api-service/blob/main/restapiservice/openapi.yaml).

1. Create `generate.go`

    ```console
    cd ${GIT_REPOSITORY_DIR}
    go generate ./...
    ```

## Test

1. Run Go tests.
   Example:

     ```console
     cd ${GIT_REPOSITORY_DIR}
     make test

     ```

## Documentation

1. Start `godoc` documentation server.
   Example:

    ```console
     cd ${GIT_REPOSITORY_DIR}
     godoc

    ```

1. Visit [localhost:6060](http://localhost:6060)
1. Senzing documentation will be in the "Third party" section.
   `github.com` > `senzing` > `go-rest-api-service`

1. When a versioned release is published with a `v0.0.0` format tag,
the reference can be found by clicking on the following badge at the top of the README.md page:
[![Go Reference](https://pkg.go.dev/badge/github.com/senzing-garage/go-rest-api-service.svg)](https://pkg.go.dev/github.com/senzing-garage/go-rest-api-service)
