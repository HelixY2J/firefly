# Firefly beta

Current approach - P2P system with master node for sync

![firefly](./res/logo.svg)

## Design 

Each music node has
- db to track files
- grpc server
- inbuilt audio player
- web interface

master client model
- master node should have central databse with file metadata
- tracks all conn clients - syn playback
- client req missing file chinks from master, maybe other clietn too

file deduplication
- files are split into chunks and fingerprinted
- chunks are stored logically 
- clients only download mising chinks instead of full files

comm by grpc
- nodes comm over grpc, handling metadata transfer, file transfermplayback snyc
- heartbeats, lib updates,chunk req

## The system flow

1. **Nodes Join via Heartbeats**
    
    - Clients send a heartbeat to the master.
    - Master assigns a unique node ID and shares the group status.
    - Master removes inactive nodes after a timeout, notifying others.
2. **Library Synchronization**
    
    - Clients periodically sync their file list with the master.
    - Master compares the client’s list with the authoritative library and sends back only changes (deltas).
3. **File Download Process**
    
    - Clients request metadata from the master (to get chunk info).
    - They check which chunks they already hav and only request missing ones.
    - Chunks are downloaded using load balancing  across available nodes.
4. **Synchronized Playback**
    
    - Any node can request synchronized playback by the master.
    - Master sends control requests to all clients.
    - Clients respond once ready after downloading the file
    - Master waits for all active clients before issuing a final playback command

## Known Upcoming challenges
- How handling network failures like when master nodes crashes
- Latency !!!



## Unknown challenges
???