# Getting Started with AnChat

This guide walks you through setting up, running, and developing AnChat.

## Quick Start (30 seconds)

### 1. Clone and Build

```bash
git clone https://github.com/huyng/anchat.git
cd anchat

# Build the server
go build -o anchat-server ./cmd/server
```

### 2. Run the Server

```bash
./anchat-server
```

That's it! The server is now running on `http://localhost:8080`.

**Check it's working:**
```bash
curl http://localhost:8080/health
# Returns: {"status":"ok","time":"..."}
```

---

## Configuration (Optional)

To customize settings, create `anchat.toml`:

```toml
[server]
port = ":8080"

[database]
path = "anchat.db"

[tls]
enabled = false
# cert = "cert.pem"
# key = "key.pem"
```

**Available config options:**
```toml
[server]
port = ":8080"              # Port to listen on

[database]
path = "anchat.db"          # SQLite database file path

[security]
storage_key = "..."         # Optional: for future at-rest encryption

[tls]
enabled = false             # Enable TLS (requires cert and key)
cert = "path/to/cert.pem"   # TLS certificate path (if enabled)
key = "path/to/key.pem"     # TLS private key path (if enabled)
```

Run with custom config:
```bash
./anchat-server --config /path/to/config.toml
```

---

## Using TLS (Optional)

For production or public deployment, enable TLS:

### Option 1: Self-signed (for testing)

Generate a self-signed certificate:
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
```

Create `anchat.toml`:
```toml
[tls]
enabled = true
cert = "cert.pem"
key = "key.pem"
```

### Option 2: Let's Encrypt (for production)

```bash
# Install certbot
sudo apt install certbot  # Linux
brew install certbot      # macOS

# Get certificate
sudo certbot certonly --standalone -d chat.example.com
```

Update `anchat.toml`:
```toml
[tls]
enabled = true
cert = "/etc/letsencrypt/live/chat.example.com/fullchain.pem"
key = "/etc/letsencrypt/live/chat.example.com/privkey.pem"
```

### Option 3: Reverse Proxy (recommended for production)

Use Caddy or nginx to handle TLS, run AnChat without TLS:

**Caddyfile:**
```
yourdomain.com {
    reverse_proxy localhost:8080
}
```

Then run AnChat with `tls.enabled = false`.

---

## Development Guide

### Project Structure

```
anchat/
├── cmd/
│   └── server/
│       └── main.go          # Server entry point
├── internal/
│   ├── auth/                # Authentication logic
│   ├── config/              # Configuration handling
│   ├── crypto/               # Cryptography operations
│   ├── db/                   # Database layer
│   ├── models/               # Data models
│   └── server/               # HTTP/2 server
├── pkg/
│   └── protocol/            # Protocol definitions
├── PROTOCOL.md              # Protocol specification
└── README.md                # Project overview
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### API Usage Examples

#### Authentication

```bash
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{
    "cmd": "auth",
    "user": "alice",
    "password": "base64_password_proof",
    "pubkey_ed25519": "base64_encoded_pubkey",
    "pubkey_x25519": "base64_encoded_pubkey"
  }'
```

#### Listen to Events (SSE)

```bash
curl -N http://localhost:8080/api/v1/listen \
  -H "Authorization: Bearer your_session_token"
```

#### Send a Command

```bash
curl -X POST http://localhost:8080/api/v1/command \
  -H "Authorization: Bearer your_session_token" \
  -H "Content-Type: application/json" \
  -d '{"cmd": "status", "state": "online"}'
```

---

## Troubleshooting

### "dial tcp: connection refused"

Check server is running:
```bash
curl localhost:8080/health
```

### Database is locked

Another process may have the database open:
```bash
lsof anchat.db
```

### TLS errors

- Verify certificate paths in config
- Check certificate exists: `ls -la cert.pem key.pem`

---

## Next Steps

1. **Implement stubbed handlers** — Complete message routing and channel operations
2. **Build a client** — CLI, web, or mobile client
3. **Add tests** — Unit tests for crypto, integration tests for server
4. **Deploy** — Use Caddy or nginx for production TLS

---

## Resources

- **[PROTOCOL.md](PROTOCOL.md)** — Complete protocol specification
- **[README.md](README.md)** — Project overview and architecture
- **Go Documentation** — https://golang.org/doc/

---

**Happy coding! 💜**
