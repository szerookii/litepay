# LitePay

**LitePay** is a lightweight Litecoin payment processor that allows developers to seamlessly integrate Litecoin payments into their applications. Built with a **Golang backend** and a **SvelteKit frontend**, LitePay offers a fast, modern, and easy-to-deploy solution. The frontend is embedded into the backend for simplicity.

---

## 📋 Features

- **Backend**: Built with **Golang**, leveraging the Litecoin JSON-RPC API for blockchain interactions.
- **Frontend**: Developed using **SvelteKit** for fast, responsive UI.
- **Single Binary Deployment**: The SvelteKit frontend is embedded into the Golang backend for simplicity—no need for separate servers.
- **Integration with Litecoin Node**: Requires hosting a Litecoin node with RPC enabled.

---

## 🚀 Prerequisites

To run **LitePay**, you need the following:

1. **Golang** (1.20 or higher): [Download Go](https://golang.org/dl/)
2. **Litecoin Node**:
    - Install and run [Litecoin Core](https://litecoin.org/).
    - Enable RPC by configuring your `litecoin.conf` file.
3. **Node.js** (only required to build the frontend):
    - [Download Node.js](https://nodejs.org/).

---

## ⚙️ Setting Up Litecoin Node

Below is an example configuration for your `litecoin.conf` file (located in `~/.litecoin/` for Linux):

```conf
server=1
rpcuser=YOUR_RPC_USERNAME
rpcpassword=YOUR_RPC_PASSWORD
rpcallowip=YOUR_IP_ADDRESS
rpcbind=127.0.0.1
rpcport=9332
```

**Note:** This is just an example configuration. Adjust the parameters (e.g., rpcuser, rpcpassword) according to your security needs and network setup. Never share your RPC credentials publicly.

Start your Litecoin node with the following command:
```bash
litecoind -daemon
```

--- 
## 🛠 Installation

SOON...

---

## ⚠️ Disclaimer

LitePay is a simple payment processor and is not intended for high-value transactions or production use without additional security measures. Use at your own risk.
