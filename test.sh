go run .  & { sleep 2; go build && ./guthi_network -port 7000 -range 1; } 

trap "killall background" EXIT