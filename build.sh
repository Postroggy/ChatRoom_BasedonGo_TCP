# compile windows version
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./Bin/Windows/TCPchatServer
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./Bin/Windows/TerminalClient ./ChatClient/Client.go
# compile linux version
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./Bin/Linux/TCPchatServer
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./Bin/Linux/TerminalClient ./ChatClient/Client.go
# compile darwin version
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./Bin/Darwin/TCPchatServer
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./Bin/Darwin/TerminalClient ./ChatClient/Client.go

PID=$(lsof -t -i:8888)
echo $PID
kill -9 $PID

