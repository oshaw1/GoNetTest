# GoNetTest
GoNetTest is a powerful, lightweight network testing tool built in Go. It provides network performance analysis with an intuitive web interface.

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
