# keyfx

keyboard fx

play any sound on key presses

development is going on for linux only at the moment

# using it
```
sudo XDG_RUNTIME_DIR=/run/user/1000 \
  PULSE_SERVER=unix:/run/user/1000/pulse/native \
  go run main.go
```