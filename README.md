# Audit Trail Client Agent (go)
This is a client agent of Enterprise IT Audit Trail Platform

## How it works
1. Clone the repo
2. go to cmd folder
3. build and running with terminal
    ```bash
    go build -o audit-agent -race

    ./audit-agent <options>

    <options>:
        -client_spawn uint
            setup audit-trail client spawn (default 100)
        -host string
            setup audit-trail server host with port
        -key string
            key of audit-trail account
        -max_in_flight int
            limit request that should be sent to audit-trail server (default 1000)
        -name string
            your app name
        -port string
            setup UDP port (default "8321")
        -secret string
            secret of audit-trail account
        -time_in_flight int
            limit time request that should be sent to server (in second) (default 60)
        -timeout uint
            setup audit-trail timeout (default 10)
    ```
