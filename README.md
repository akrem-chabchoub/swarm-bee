# swarm-bee — mini Bee node (learning implementation)

This repository is a focused, incremental reimplementation of the Bee node architecture (Swarm-inspired).
The goal is to build a production-minded, well-tested node in small phases: identity & crypto, config/CLI, storage, p2p (libp2p), Kademlia topology, core protocols, HTTP API, and CI/CD.

Key components:
- pkg/address, pkg/crypto — identities and signing
- pkg/storage, pkg/sharky, pkg/storer — persistent chunk storage
- pkg/p2p, pkg/topology — libp2p networking and Kademlia routing
- pkg/api — HTTP surface for uploads/downloads


Quick start:
- go build -o build/swarm-bee ./cmd/swarm-bee
- go test ./... 

License: MIT (or choose your preferred license)
