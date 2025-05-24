cd master && go build && cp ./master ../Master && cd ../
cd watchdog && go build && cp ./watchdog ../Watchdog && cd ../
echo "Installed! Please configure the watchdog binary to run on start (service file), use the 'master' binary to send commands to workers"

echo "Running help command for master..."
./master -help