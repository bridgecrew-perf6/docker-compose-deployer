# Docker Compose Deployer

Deployer for your Docker Compose services

## Requirements

1.  Use variable as image tag in the docker compose file

    ```yaml
    alpine:
      image: alpine:$ALPINE_IMG_TAG
    ```

2.  Create `.env` file in the same folder with the docker compose file, and declare variables in the `.env` file

    ```sh
    ALPINE_IMG_TAG=3.15.0
    ```

## Use cases

1.  Run command in your CD pipeline

    ```sh
    Usage:
      docker-compose-deployer deploy <service name> <image tag> [flags]

    Flags:
          --compose-v2            Use Docker Compose v2
      -f, --file stringArray      Compose configuration files
      -h, --help                  help for deploy
      -p, --project-name string   Compose project name
          --sudo                  Execute commands with sudo
      -v, --verbose               Show more output
      -w, --workdir string        Specify an alternate working directory (default ".")
    ```

2.  Start a http server, and make http request in your CD pipeline

    ```sh
    Usage:
      docker-compose-deployer server [flags]

    Flags:
          --compose-v2            Use Docker Compose v2
      -f, --file stringArray      Compose configuration files
      -h, --help                  help for server
          --port int              The server port (default 8000)
      -p, --project-name string   Compose project name
          --secret string         JWT secret
          --sudo                  Execute commands with sudo
      -v, --verbose               Show more output
      -w, --workdir string        Specify an alternate working directory (default ".")
    ```

    The server uses JWT for authorization. The payload should has following fields at least.

    ```json
    {
        "sub": "alpine",
        "jti": "3.14.3",
        "iat": 1516239000
    }
    ```

    `sub` represents the service name and `jti` represents the image tag need to be deployed.

    If you are not familiar with JWT, visit [jwt.io](https://jwt.io/). It will generate JWT automatically after you input your secret and payload.

    **Q:** Why no https?

    **A:** We can use Nginx with SSL as a reverse proxy.

## Command options

You can also see this information by running `docker-compose-deployer --help` from the command line.

Several environment variables are available for you to configure the command-line behavior.

| Flag                   | Environment Variable  | Description                            |
|------------------------|-----------------------|----------------------------------------|
| --compose-v2           | `DEPLOYER_COMPOSE_V2` | Use Docker Compose v2                  |
| -f, --file stringArray | `DEPLOYER_FILE`       | Compose configuration files            |
| -p, --port int         | `DEPLOYER_PORT`       | The server port                        |
| -s, --secret string    | `DEPLOYER_SECRET`     | JWT secret                             |
| --sudo                 | `DEPLOYER_SUDO`       | Execute commands with sudo             |
| -v, --verbose          | `DEPLOYER_VERBOSE`    | Show more output                       |
| -w, --workdir string   | `DEPLOYER_WORKDIR`    | Specify an alternate working directory |

## [Demos](./demo/)

## References

- [Compose V2](https://docs.docker.com/compose/cli-command/)
- [Environment variables in Compose](https://docs.docker.com/compose/environment-variables/)