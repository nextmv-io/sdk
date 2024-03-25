if false; then
go run main.go
fi
sleep 0.5
go run main.go > /dev/null 2>&1 &
sleep 3.5
PID2=$(lsof -i -P | grep LISTEN | grep :9003 | tr -s ' ' | cut -d ' ' -f 2)
curl -s -X POST "http://localhost:9003?duration=5000000000" -H 'Content-Type: application/json' -d '{"message":"Hello one"}' > /dev/null 2>&1 &
sleep 0.01
curl -s -X POST "http://localhost:9003?duration=5000000000" -H 'Content-Type: application/json' -d '{"message":"Hello two"}' > /dev/null 2>&1 &
sleep 1 && curl -s -X POST "http://localhost:9003?duration=500000000" -H 'Content-Type: application/json' -d '{"message":"Hello"}'
kill $PID2 > /dev/null 2>&1
sleep 2
exit 0
