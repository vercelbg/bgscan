<div align="center">

<br/>

```
██████╗  ██████╗ ███████╗ ██████╗ █████╗ ███╗   ██╗
██╔══██╗██╔════╝ ██╔════╝██╔════╝██╔══██╗████╗  ██║
██████╔╝██║  ███╗███████╗██║     ███████║██╔██╗ ██║
██╔══██╗██║   ██║╚════██║██║     ██╔══██║██║╚██╗██║
██████╔╝╚██████╔╝███████║╚██████╗██║  ██║██║ ╚████║
╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝
```

**Blazing-fast multi-protocol scanner with modular chain architecture**

<br/>

[![Go Version](https://img.shields.io/badge/Go-1.26.3+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-6366f1?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20|%20Windows%20|%20macOS%20|%20Termux-64748b?style=flat-square)](https://github.com/MohsenBg/bgscan/releases)
[![UI](https://img.shields.io/badge/UI-BubbleTea%20TUI-ec4899?style=flat-square)](https://github.com/charmbracelet/bubbletea)
[![Status](https://img.shields.io/badge/Status-Production%20Ready-22c55e?style=flat-square)](https://github.com/MohsenBg/bgscan/releases/latest)

<br/>

[📦 Installation](#-installation) · [📚 Documentation](#-documentation) · [🔹 Protocols](#-supported-protocols) · [💖 Support](#-support--donate)

<br/>

---

</div>

## About bgscan

**bgscan** is a high-performance scanning engine written in **Go**.
It runs multiple protocols concurrently and chains scan stages together to build advanced detection and reconnaissance workflows.

Built for developers and researchers who demand **speed, flexibility, and a modern scanning experience**.

---

<img width="1552" height="841" alt="bgscan-ui" src="https://github.com/user-attachments/assets/08a50bc0-186d-45a4-8c27-4bb52a2041ee" />

---

## ✨ Key Features

<table>
<tr>
<td width="50%">

### 🔗 Smart Chaining
Chain multiple scan stages with `Stream`, `Sequential`, and `Batch Size` execution modes for complete control over your data pipeline.

</td>
<td width="50%">

### 🖥️ Interactive TUI
A fully interactive **Bubble Tea** terminal UI — scan, monitor, and manage results without ever opening a browser.

</td>
</tr>
<tr>
<td width="50%">

### 📡 Broad Protocol Support
From ICMP and TCP to DNS Tunnel, Slipstream, and Xray — all in one tool.

</td>
<td width="50%">

### 💾 Save & Replay
Persist scan results and run new scans directly against previously saved data.

</td>
</tr>
<tr>
<td width="50%">

### 🛰️ Xray Integration
Full support for saving, managing, and validating Xray outbounds.

</td>
<td width="50%">

### 🌐 Advanced DNS
DNS Tunnel support with a fallback mechanism for complex resolution scenarios.

</td>
</tr>
</table>

---

## 🔹 Supported Protocols

| Protocol | Layer | Description |
|:--------:|:-----:|-------------|
| **ICMP** | 3 | Host discovery and reachability via Ping |
| **TCP** | 4 | Connection scanning and TCP handshake validation |
| **HTTP** | 7 | Full HTTP/1.1 and HTTP/2 support *(no QUIC / HTTP/3)* |
| **TLS** | 7 | TLS 1.0 through TLS 1.3 |
| **DNS** | 7 | Advanced DNS queries with fallback mechanism |
| **DNSTT** | 7 | DNS Tunnel validation *(SOCKS only, no auth)* |
| **Slipstream** | 7 | Slipstream validation *(SOCKS only, no auth)* |
| **Xray** | 7 | Xray outbound validation and testing |

---

## 📦 Installation

### Automatic (Recommended)

**🪟 Windows**
```powershell
irm https://raw.githubusercontent.com/MohsenBg/bgscan/refs/heads/main/scripts/install.ps1 | iex
```

**🐧 Linux / 🍎 macOS**
```bash
curl -fsSL https://raw.githubusercontent.com/MohsenBg/bgscan/refs/heads/main/scripts/install.sh | sh
```

**🤖 Android (Termux)**
```bash
pkg update -y && pkg install bash curl unzip -y
curl -fsSL https://raw.githubusercontent.com/MohsenBg/bgscan/refs/heads/main/scripts/install.sh | bash
```

---

### Manual Install

Download the binary for your OS and architecture from the Releases page:

👉 **[Download the latest release](https://github.com/MohsenBg/bgscan/releases/latest)**

**Linux / macOS / Termux**
```bash
unzip bgscan-*.zip
chmod +x bgscan
./bgscan
```

**Windows**

Extract the ZIP and run `bgscan.exe`.

---

## 📚 Documentation

### 👤 User Docs
Usage guide, configuration, supported protocols, execution modes, and practical examples:

📖 **[→ User Documentation](docs/user/README.md)**

---

### 🛠️ Developer Docs
Project architecture, module structure, the Pipeline engine, the Asset system, coding standards, and contribution guidelines:

📖 **[→ Developer Documentation](docs/developer/README.md)**

---

## 💖 Support / Donate

If bgscan has been useful to you and you'd like to support its development:

| Network | Currency | Address |
|:-------:|:--------:|---------|
| ₿ Bitcoin | `BTC` | `bc1qdwh57dm97nmx5jzdr7lrc9cxe5xh3zc59er7z9` |
| 🟣 Ethereum | `ETH` | `0x40Fd22Fff4E059e906A10747Fd0a45A1DB39c985` |
| 🟡 BNB Smart Chain | `BNB / BEP20` | `0x40Fd22Fff4E059e906A10747Fd0a45A1DB39c985` |
| 🔷 TRON | `TRX / TRC20` | `TNW6pbfY8zZVZezZWyYXo7h12MycRsVJK7` |

---

<div align="center">

Built with ❤️ and Go &nbsp;|&nbsp; MIT License

</div>
