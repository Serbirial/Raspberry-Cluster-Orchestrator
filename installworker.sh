git clone https://github.com/Serbirial/Raspberry-Cluster-Orchestrator.git
cd Raspberry-Cluster-Orchestrator

cd slave && go build && cp ./slave ../slave && cd ../
echo "Installed! Please configure the slave binary to run on startup through a service file"