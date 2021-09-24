run proxy-server:

`clear; ./pswm` - (for mac os, by default will be use port 9080)

`clear; ./pswl` – for linux (arm64, by default will be use port 9080)

run 3d-party service emulationL:

`clear; ./lsem` - (for mac os, by default will be use port 9080)

`clear; ./lsel` – for linux (arm64, by default will be use port 9080)

run below curl requests in proxy-server folder.

send request to 3rd-party service. as the response you will get request key (RK). it needs for check of status:   

`curl -X POST "127.0.0.1:9080/client/request" -H "Content-Type: application/json" -d "@worker/testdata/client_request_all_ok.json"`

send request to fetch result. "request" is RK (see above)

`curl -X POST "127.0.0.1:9080/client/status" -H "Content-Type: application/json" -d '{"request": "__use_generated_before_code__"}'`