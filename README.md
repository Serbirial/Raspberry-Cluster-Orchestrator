![Go](https://img.shields.io/badge/Language-Go-00ADD8?logo=go&logoColor=white)
![Raspberry Pi](https://img.shields.io/badge/Platform-Raspberry%20Pi-green?logo=raspberry-pi&logoColor=white)
![GitHub stars](https://img.shields.io/github/stars/Serbirial/rco?style=social)

# 🚀 RCO – Raspberry Cluster Orchestrator

**RCO** is a lightweight automation and orchestration system for Raspberry Pi (or similar) clusters.  
It lets you execute commands, sync code, build/run binaries, and (not implemented YET) coordinate distributed devices via WLAN or GPIO – all from a single controller node.

---

### ✨ Features

- 📡 Send shell commands or JSON-defined task trees to all known nodes  
- 🧱 Build and execute binaries remotely  
- 🔗 Bridge nodes via GPIO or WLAN for distributed IO and communication (NOT IMPLEMENTED YET!!!)  
- 📁 Remote deployment across clusters of Raspberry Pis  
- 🪶 Minimal dependencies, built from scratch with simplicity in mind  

---

### 📡 Designed for...

- Raspberry Pi clusters  
- Low-power or offline-capable systems  

---

### 🧪 Example JSON Task

```json
{
  "PiWorker2": {
    "dir": "/home/USERNAME/PROJECT",
    "cmd": ["pkill -f './PROJECT -flag1' || true", "git pull origin main", "go build"],
    "bin": ["./BINARY", "-flag1"]
  }
}
```

## 🛠️ Installation

### Prerequisites

Before setting up RCO, ensure you have the following:

- **Go 1.24+**: Install from [golang.org](https://golang.org/dl/).
- **Git**: Install from [git-scm.com](https://git-scm.com/).
- **Raspberry Pi Devices**: At least one Raspberry Pi (or similar, will need GPIO functionality in future) running a Linux-based OS (e.g., Raspberry Pi OS).
- **Network Connectivity**: All Raspberry Pis in the cluster should be connected to the same network, with one acting as a controller.

### Clone the Repository

```bash
git clone https://github.com/Serbirial/Raspberry-Cluster-Orchestrator.git
cd Raspberry-Cluster-Orchestrator
```

### Build 
```bash
# A controller will need two binaries, watchdog and master

# Build the master binary and copy
cd master && go build && cp ./master ../master
cd ../
# Build the watchdog binary and copy
cd watchdog && go build && cp ./watchdog ../watchdog


# A worker/slave will need the slave binary

# Build the slave binary and copy
cd slave && go build && cp ./slave ../slave

```

## 🚀 Usage

### Configuration files
* `workers.txt`: A text file full of all known worker/slave nodes (one per line)
* `commands.json`: JSON file defining tasks to be executed on worker nodes.

### Example `workers.txt`

```
PiWorker1> 192.168.0.5
PiWorker2> 192.168.0.6
PiWorker3> 192.168.0.7
```

### Example `commands.json`
```json
{
  "PiWorker1": {
    "dir": "/home/pi/ascension-go",
    "cmd": ["pkill -f './ascension -remote-ws -ws-url=\"ws://localhost:8182/ws\" || true", "git pull origin main", "go build"],
    "bin": ["./ascension", "-remote-ws", "-ws-url=\"ws://localhost:8182/ws\""]
  },
  "PiWorker2": {
    "dir": "/home/summers/ascension-go",
    "cmd": ["pkill -f './ascension -ws-only' || true", "git pull origin main", "go build"],
    "bin": ["./ascension", "-ws-only"]
  }
}
```