go build 
count=$1
for (( i=1; i<=$count; i++ ))
do 
    port_no=$((7000+$i))   
    echo $port_no
    ./guthi_network -port $port_no -range 1 &
done

# trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT