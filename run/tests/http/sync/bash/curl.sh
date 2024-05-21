if false; then
go run main.go
fi
sleep 0.5
go run main.go > /dev/null 2>&1 &
sleep 3.5
PID2=$(lsof -i -P | grep LISTEN | grep :9000 | tr -s ' ' | cut -d ' ' -f 2)
curl -s -X POST "http://localhost:9000?duration=500000000" -H 'Content-Type: application/json' -d '{"message":"Hello"}' | jq
kill $PID2 > /dev/null 2>&1
exit 0
