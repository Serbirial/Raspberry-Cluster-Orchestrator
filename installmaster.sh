git clone https://github.com/Serbirial/Raspberry-Cluster-Orchestrator.git
cd Raspberry-Cluster-Orchestrator

cd master && go build && cp ./master ../master && cd ../
cd watchdog && go build && cp ./watchdog ../watchdog && cd ../
echo "Installed! Please configure the watchdog binary to run on start (service file), use the 'master' binary to send commands to workers"

echo "Running help command for master..."
./master -help