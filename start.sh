mkdir Bin/Windows
mkdir Bin/Linux
mkdir Bin/Darwin

# compile windows version
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./Bin/Windows/TCPchatClient main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./Bin/Windows/TerminalClient ./ConsoleClient/terminalClient.go
# compile linux version
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./Bin/Linux/TCPchatClient main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./Bin/Linux/TerminalClient ./ConsoleClient/terminalClient.go
# compile darwin version
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./Bin/Darwin/TCPchatClient main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./Bin/Darwin/TerminalClient ./ConsoleClient/terminalClient.go

PID=$(lsof -t -i:8888)
echo $PID
kill -9 $PID
