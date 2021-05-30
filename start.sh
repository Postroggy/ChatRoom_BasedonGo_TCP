go build .
PID=$(lsof -t -i:8888)
echo $PID
kill -9 $PID
./tcp
