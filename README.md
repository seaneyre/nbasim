# nbasim
Simulate historical or hypothetical NBA matches and output play-by-play events to a Websocket API.

## Quick Start

Build

```
go build -o nbasim cmd/nbasim
```

Start a Websocket server

```
./nbasim server --host localhost --port 8000
```

Start a simulation
```
./nbasim simulate --game-id 0022000180 --time-factor 4.00 --url localhost:8000
```

Connect to the Websocket server and watch the game events roll in
```
wscat -c ws:localhost:8000/ws/game/0022000180
```