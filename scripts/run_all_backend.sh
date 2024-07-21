echo "Please Run Me on the root dir, not in scripts dir."
echo "If here is no build dir, please run build script first."

gateway_directory="build/gateway"
service_directory="build/services"

log_directory="logs"

mkdir -p "$log_directory"

if [ ! -d "build" ]; then
    echo "Output dir does not exist, please run build script first."
fi

for gateway_file in "$gateway_directory"/*; do
    if [[ -x "$gateway_file" && -f "$gateway_file" ]]; then
        echo "Running $gateway_file"
        gateway_log_file="$log_directory"/"$(basename gateway_file)".log
        ./"$gateway_file" >> "$gateway_log_file" 2>&1 &
    fi
done

for service in "$service_directory"/*; do
    for service_file in "$service"/*; do
        if [[ -x "$service_file" && -f "$service_file" ]]; then
            echo "Running $service_file"
            service_log_file="$log_directory"/"$(basename service_file)".log
            ./"$service_file" >> "$service_log_file" 2>&1 &
        fi
    done
done