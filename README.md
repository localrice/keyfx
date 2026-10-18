# KeyFX

<p align="center">
    <img src="./docs/images/keyfx.png" alt="mefetch" border="0">
</p>

Keyboard FX

play any sound on key presses.
development is going on for Linux only at the moment.
(happy to welcome devs to create windows support)

## Usage

Build the binary and install it to `/usr/local/bin`:

```bash
go build -o keyfx .
sudo install -Dm755 keyfx /usr/local/bin/keyfx
```

## Configuration

Create the configuration file at `~/.config/keyfx/config.yaml`:
```yaml
sound: "/absolute/path/to/your/sound.mp3"
```
Replace the example path with the absolute path to your sound file.

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

## future plans

- [ ] windows support
- [ ] multiple sound support for multiple keys
- [ ] make it cooler  