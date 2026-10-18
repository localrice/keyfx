# KeyFX

Keyboard FX

Play any sound on key presses.

development is going on for Linux only at the moment.

## Usage

Build the binary and install it to `/usr/local/bin`:

```bash
go build -o keyfx .
sudo install -Dm755 keyfx /usr/local/bin/keyfx
```

Run:

```bash
sudo -E keyfx
```

To run the daemon manually without building:

```bash
sudo XDG_RUNTIME_DIR=/run/user/1000 \
  PULSE_SERVER=unix:/run/user/1000/pulse/native \
  go run main.go
```