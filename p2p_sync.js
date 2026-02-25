const P2P_PORT = process.env.P2P_PORT || 5001;
const peers = process.env.PEERS ? process.env.PEERS.split(',') : [];

// This logic allows your 54-nation nodes to find each other automatically
class P2PServer {
    constructor(blockchain) {
        this.blockchain = blockchain;
        this.sockets = [];
    }

    listen() {
        const server = new WebSocket.Server({ port: P2P_PORT });
        server.on('connection', socket => this.connectSocket(socket));

        // Connect to known 'Lighthouse' seed nodes
        this.connectToPeers();
        console.log(`📡 Sovereign Discovery listening on port: ${P2P_PORT}`);
    }

    connectToPeers() {
        peers.forEach(peer => {
            const socket = new WebSocket(peer);
            socket.on('open', () => this.connectSocket(socket));
        });
    }

    // Sync the "Ancient" ledger across the global network
    syncChain(socket) {
        socket.send(JSON.stringify(this.blockchain.chain));
    }
}

module.exports = P2PServer;