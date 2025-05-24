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