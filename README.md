# keyfx

keyboard fx

play any sound on key presses

development is going on for linux only at the moment

# using it

build the binary and install it to /usr/local/bin:

```
go build -o keyfx .
sudo install -Dm755 keyfx /usr/local/bin/keyfx
```

To run the daemon manually:
```
sudo XDG_RUNTIME_DIR=/run/user/1000 \
  PULSE_SERVER=unix:/run/user/1000/pulse/native \
  go run main.go
```