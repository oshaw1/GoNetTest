# GoNetTest
GoNetTest is a powerful all-in-one network test tool with an intuitive web interface. Written in Go and HTMX, it provides network analysis and helps detect regressions through its feature set:

- Real-time Performance Visualization: Generate interactive charts and graphs showcasing network metrics
- Scheduled Testing: Put together a task schedule to automate tests and chart generation - ensuring continuous network monitoring without manual intervention
- Historical Analytics: Access detailed historical data with customizable date ranges, helping identify patterns and troubleshoot network speed regressions
- Performance Insights: Get deep analytical insights into your network's performance over time, with customizable metrics and thresholds

## Types of tests:

- ICMP - Internet Control Message Protocol test that measures packet transmission 
between network hosts. this is just a small "healthcheck" request primarily used to validate connection. The "Jitter" test is a more advanced version of this

- Download - Tests download speeds over time by measuring the rate of data transfer 
from several different servers and data sizes to the client and calculates the average.

- Upload - Tests upload speeds over time by measuring the rate of data transfer 
from the client to serveral different servers and measures the average.

- Latency - Measures the variation in latency between successive packets. Helps identify 
network stability issues.

- Route - Traces the network path to a target, showing RTT for each hop. Helps identify 
routing bottlenecks and weak links.

- Bandwidth - Measures overall network capacity by testing maximum throughput at multiple different users to find the point at which performance suffers for x users

## Usage

To change any test/ui parameters such as Download/Upload urls or max requests please do so within `config/config.json`

Once the application is started you can access the dashboard via `{youripaddress/localhost}:7000/dashboard`

Alternatively you can view all accessable endpoints within the startup logs and view the specs within api/

### To Build -
- Bash/Go build tools:
```
go build -o GoNetTest cmd/main.go
```
- Docker:
```
docker build -t go-net-test . ;
docker run -p 7000:7000 go-net-test
```
**NOTE**: GoNetTest must be started from its root directory.

### Linux

For linux systems it is recomended to run GoNetTest as a system service as it is intended to run in the background. To do this first build the binary then configure gonettest.service to point to it.  
An example service configuration is:
```
[Unit]
Description=GoNetTest Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/home/code/projects/GoNetTest
ExecStart=/home/code/projects/GoNetTest/GoNetTest
Restart=always

[Install]
WantedBy=multi-user.target
```
Once this service has been configured it can be started with : ```sudo systemctl start gonettest```  

## Low Level System Architecture Diagram -
![Network Testing Webapp System ARCH](https://github.com/user-attachments/assets/d4563d27-be0a-4aad-b78e-80f2e9b19865)

## Application Flow Diagram -
![Network Testing Application Flow](https://github.com/user-attachments/assets/272cf95b-0e15-4226-8ec4-3e1b6e47495d)