# Audit Trail Client Agent (go)
This is a client agent of Enterprise IT Audit Trail Platform

## How it works
1. Clone the repo
2. go to cmd folder
3. build and running with terminal
    ```bash
    go build -o audit-agent -race

    ./audit-agent -port=<your-UDP-port|default:8321> -host=<audit-trail-host> -key=<your-audit-trail-key> -secret=<your-audit-trail-secret> -name=<your-service-name> -max_in_flight=<setup_max_in_flight_request> -time_in_flight=<setup_time_limit_in_second> -timeout=<audit-trail-timeout> -client_spawn=<audit-trail_client_spawn>
    ```