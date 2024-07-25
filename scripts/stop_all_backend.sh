echo "Stopping all running services and gateway"

gateway_directory="build/gateway"
service_directory="build/services"

kill_process() {
    local process_name=$1
    pids=$(pgrep -f "$process_name")
    if [ -n "$pids" ]; then
        echo "Stopping $process_name..."
        kill $pids
    else
        echo "$process_name is not running."
    fi
}

for gateway_file in "$gateway_directory"/*; do
  if [[ -x "$gateway_file" && -f "$gateway_file" ]]; then
    gateway_name=$(basename "$gateway_file")
    kill_process "$gateway_name"
  fi
done


for service in "$service_directory"/*; do
    for service_file in "$service"/*; do
        if [[ -x "$service_file" && -f "$service_file" ]]; then
            service_name=$(basename "$service_file")
            kill_process "$service_name"
        fi
    done
done

rm -rf logs
mkdir logs

echo "All services and gateways have been stopped."