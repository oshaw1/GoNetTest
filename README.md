# GoNetTest
GoNetTest is a powerful all-in-one network test tool with an intuitive web interface. Written in Go and HTMX, it provides network analysis and helps detect regressions through its feature set:

Real-time Performance Visualization: Generate interactive charts and graphs showcasing network metrics
Scheduled Testing: Set up automated test intervals - ensuring continuous network monitoring without manual intervention
Historical Analytics: Access detailed historical data with customizable date ranges, helping identify patterns and troubleshoot network speed regressions
Performance Insights: Get deep analytical insights into your network's performance over time, with customizable metrics and thresholds

## To Build -
```
go build cmd/main.go
./main.exe
```
### With Docker -
```
docker build -t go-net-test-app . ;
docker run --network host -p 7000:7000 go-net-test-app
```
## Low Level System Architecture Diagram -
![Network Testing Webapp System ARCH](https://github.com/user-attachments/assets/d4563d27-be0a-4aad-b78e-80f2e9b19865)

## Application Flow Diagram -
![Network Testing Application Flow](https://github.com/user-attachments/assets/272cf95b-0e15-4226-8ec4-3e1b6e47495d)

## Types of tests -
 
ICMP -
Download -
Upload -
