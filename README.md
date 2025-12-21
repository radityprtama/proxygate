# ProxyGate

**Unified AI API gateway for OpenAI, Gemini, Claude, and Codex CLI tools.**

ProxyGate is a self-hosted proxy server that provides OpenAI-compatible API endpoints for multiple AI providers. Use your existing AI subscriptions with any compatible client, SDK, or development tool.

## Features

- **Multi-Provider Support** - OpenAI, Gemini, Claude, and Codex compatible API endpoints
- **OAuth Authentication** - Simple CLI-based OAuth login flows for all supported providers
- **Load Balancing** - Multi-account round-robin distribution with automatic failover
- **Streaming Support** - Full streaming and non-streaming response support
- **Function Calling** - Complete tool use and function calling support
- **Multimodal** - Text and image input support
- **Quota Management** - Automatic credential rotation and quota handling
- **Embeddable SDK** - Reusable Go SDK for embedding the proxy in your applications

## Installation

### Binary Release (Recommended)

Download the latest release from GitHub:

```bash
# Linux (amd64)
curl -fsSL https://github.com/radityprtama/proxygate/releases/latest/download/proxygate_linux_amd64.tar.gz | tar xz
sudo mv proxygate /usr/local/bin/
proxygate --help

# Linux (arm64)
curl -fsSL https://github.com/radityprtama/proxygate/releases/latest/download/proxygate_linux_arm64.tar.gz | tar xz
sudo mv proxygate /usr/local/bin/
```

### Docker

```bash
docker pull radityprtama/proxygate:latest
docker run -d -p 8317:8317 -v ./config.yaml:/proxygate/config.yaml radityprtama/proxygate
```

### Docker Compose

```bash
docker compose up -d
```

### Build from Source

```bash
git clone https://github.com/radityprtama/proxygate.git
cd proxygate
go build -o proxygate ./cmd/server
```

## Configuration

Create a `config.yaml` file:

```yaml
host: "127.0.0.1"
port: 8317
auth-dir: "~/.proxygate"
api-keys:
  - "your-api-key-here"

# Enable debug logging
debug: false

# Gemini API keys (optional)
# gemini-api-key:
#   - api-key: "AIzaSy..."

# Claude API keys (optional)
# claude-api-key:
#   - api-key: "sk-ant-..."

# OpenAI compatibility providers (optional)
# openai-compatibility:
#   - name: "openrouter"
#     base-url: "https://openrouter.ai/api/v1"
#     api-key-entries:
#       - api-key: "sk-or-..."
```

See `config.example.yaml` for the complete configuration reference.

## OAuth Login Flows

ProxyGate supports OAuth-based authentication for CLI tools:

### Gemini CLI

```bash
proxygate -login
```

### OpenAI Codex

```bash
proxygate -codex-login
```

### Claude Code

```bash
proxygate -claude-login
```

### Qwen Code

```bash
proxygate -qwen-login
```

### iFlow

```bash
proxygate -iflow-login
```

Use `-no-browser` flag if you want to manually open the OAuth URL:

```bash
proxygate -login -no-browser
```

## Running as a Systemd Service

Create `/etc/systemd/system/proxygate.service`:

```ini
[Unit]
Description=ProxyGate AI API Gateway
After=network.target

[Service]
Type=simple
User=proxygate
WorkingDirectory=/opt/proxygate
ExecStart=/usr/local/bin/proxygate -config /etc/proxygate/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable proxygate
sudo systemctl start proxygate
```

## Reverse Proxy with NGINX

```nginx
server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8317;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 86400;
    }
}
```

## API Endpoints

ProxyGate exposes the following API endpoints:

| Endpoint | Description |
|----------|-------------|
| `POST /v1/chat/completions` | OpenAI-compatible chat completions |
| `POST /v1/responses` | OpenAI Responses API |
| `POST /v1beta/models/{model}:generateContent` | Gemini-compatible endpoint |
| `POST /v1/messages` | Claude-compatible messages API |

## SDK Usage

ProxyGate includes an embeddable Go SDK:

```go
import (
    "github.com/radityprtama/proxygate/v6/sdk/cliproxy"
    "github.com/radityprtama/proxygate/v6/sdk/config"
)

cfg := &config.SDKConfig{
    // Configuration options
}

service := cliproxy.NewService(cfg)
// Use service...
```

See `docs/sdk-usage.md` for detailed SDK documentation.

## Security Notes

- **Use TLS in production** - Always deploy behind HTTPS
- **Restrict management API** - Keep management endpoints on localhost only
- **Secure API keys** - Use strong, unique API keys
- **Rotate credentials** - Regularly rotate OAuth tokens and API keys
- **Monitor access** - Enable logging to track API usage

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
