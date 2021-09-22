run below curl requests in worker folder. if it needs change port according cmd argument. 

send request to 3rd-party service. as the response you will get request key (RK). it needs for check of status:   

`curl -X POST "127.0.0.1:9080/client/request" -H "Content-Type: application/json" -d "@testdata/client_request_all_ok.json"`

send request to fetch result. "request" is RK (see above)

`curl -X POST "127.0.0.1:9080/client/status" -H "Content-Type: application/json" -d '{"request": "f229fcff-012d-4ad7-8b1a-88bf0ce2aba9"}'`