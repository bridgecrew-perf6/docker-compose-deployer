# Demos

## Docker Compose demo

1.  Start the demo project

    ```sh
    $ cd demo
    $ echo "ALPINE_IMG_TAG=3.15.0" > .env
    $ docker-compose up -d
    # check OS version
    $ docker-compose exec alpine cat /etc/os-release
    NAME="Alpine Linux"
    ID=alpine
    VERSION_ID=3.15.0
    PRETTY_NAME="Alpine Linux v3.15"
    HOME_URL="https://alpinelinux.org/"
    BUG_REPORT_URL="https://bugs.alpinelinux.org/"
    ```

2.  Update the service through http api

    ```sh
    $ curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhbHBpbmUiLCJqdGkiOiIzLjE0LjMiLCJpYXQiOjE1MTYyMzkwMDB9.2rKVkCPtdmBUlZAzQLSW6IPfiO6_IuHrZpFiqdWl8oI" http://localhost:8000/v1/deploy
    {"code":0,"message":"OK","data":{"svc":"alpine","tag":"3.14.3"}}
    ```

    Check OS version again

    ```sh
    $ docker-compose exec alpine cat /etc/os-release
    NAME="Alpine Linux"
    ID=alpine
    VERSION_ID=3.14.3
    PRETTY_NAME="Alpine Linux v3.14"
    HOME_URL="https://alpinelinux.org/"
    BUG_REPORT_URL="https://bugs.alpinelinux.org/"
    ```

## CLI demo

1.  Start the demo project

    ```sh
    $ cd demo
    $ echo "ALPINE_IMG_TAG=3.15.0" > .env
    $ docker-compose up -d alpine
    # check OS version
    $ docker-compose exec alpine cat /etc/os-release
    NAME="Alpine Linux"
    ID=alpine
    VERSION_ID=3.15.0
    PRETTY_NAME="Alpine Linux v3.15"
    HOME_URL="https://alpinelinux.org/"
    BUG_REPORT_URL="https://bugs.alpinelinux.org/"
    ```

2.  Start the docker-compose-deployer

    ```sh
    $ docker-compose-deployer server &
    ```

3.  Update the service through http api

    ```sh
    $ curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhbHBpbmUiLCJqdGkiOiIzLjE0LjMiLCJpYXQiOjE1MTYyMzkwMDB9.2rKVkCPtdmBUlZAzQLSW6IPfiO6_IuHrZpFiqdWl8oI" http://localhost:8000/v1/deploy
    {"code":0,"message":"OK","data":{"svc":"alpine","tag":"3.14.3"}}
    ```

    Check OS version again

    ```sh
    $ docker-compose exec alpine cat /etc/os-release
    NAME="Alpine Linux"
    ID=alpine
    VERSION_ID=3.14.3
    PRETTY_NAME="Alpine Linux v3.14"
    HOME_URL="https://alpinelinux.org/"
    BUG_REPORT_URL="https://bugs.alpinelinux.org/"
    ```
