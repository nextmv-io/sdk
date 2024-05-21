if false; then
go run main.go
fi
sleep 0.5
go run main.go > /dev/null 2>&1 &
sleep 3.5
PID2=$(lsof -i -P | grep LISTEN | grep :9001 | tr -s ' ' | cut -d ' ' -f 2)
curl -s -X POST "http://localhost:9001?duration=500000000" -H 'Content-Type: application/json' -d '{"message":"Hello"}'
sleep 2
kill $PID2 > /dev/null 2>&1
jq . callback.txt > callback.json
exit 0
