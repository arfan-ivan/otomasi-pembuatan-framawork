# go.mod
module github.com/yourusername/arvia

go 1.21

require (
    github.com/fsnotify/fsnotify v1.7.0
    github.com/gorilla/websocket v1.5.1
)

require (
    golang.org/x/net v0.17.0 // indirect
    golang.org/x/sys v0.13.0 // indirect
)

---

# Makefile
.PHONY: build install clean dev

# Build the arvia binary
build:
	go build -o bin/arvia cmd/arvia/main.go

# Install arvia globally
install: build
	go install cmd/arvia/main.go

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Development mode (for framework development)
dev:
	go run cmd/arvia/main.go serve

# Test the framework
test:
	go test ./...

help:
	@echo "Arvia Framework Build Commands:"
	@echo "  make build    - Build arvia binary"
	@echo "  make install  - Install arvia globally"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make dev      - Run in development mode"
	@echo "  make test     - Run tests"

---

# install.sh
#!/bin/bash

echo "🚀 Installing Arvia Framework..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed"

# Build the binary
echo "📦 Building Arvia..."
go build -o arvia cmd/arvia/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    
    # Make executable
    chmod +x arvia
    
    # Optionally move to system PATH
    if [ -w "/usr/local/bin" ]; then
        echo "🔧 Installing to /usr/local/bin..."
        sudo mv arvia /usr/local/bin/
        echo "✅ Arvia installed globally!"
        echo ""
        echo "Usage:"
        echo "  arvia init my-app"
        echo "  cd my-app"  
        echo "  arvia serve"
    else
        echo "✅ Binary created: ./arvia"
        echo ""
        echo "Usage:"
        echo "  ./arvia init my-app"
        echo "  cd my-app"
        echo "  ../arvia serve"
        echo ""
        echo "To install globally, run:"
        echo "  sudo mv arvia /usr/local/bin/"
    fi
else
    echo "❌ Build failed!"
    exit 1
fi

---

# README.md
# Arvia Framework

Simple, fast Go web development framework for static sites with live reload.

## Features

- 🚀 **Simple Setup**: Create projects in seconds
- 🔄 **Live Reload**: Automatic browser refresh during development  
- 📦 **Easy Build**: One command to build for production
- 🎨 **Modern UI**: Beautiful default styling
- 📱 **Responsive**: Mobile-first design
- ⚡ **Fast**: Built with Go for speed

## Installation

### Option 1: Install Script
```bash
curl -sSL https://raw.githubusercontent.com/yourusername/arvia/main/install.sh | bash
```

### Option 2: Manual Install
```bash
git clone https://github.com/yourusername/arvia.git
cd arvia
go build -o arvia cmd/arvia/main.go
sudo mv arvia /usr/local/bin/  # Optional: install globally
```

### Option 3: Go Install
```bash
go install github.com/yourusername/arvia/cmd/arvia@latest
```

## Quick Start

1. **Create new project:**
   ```bash
   arvia init my-website
   cd my-website
   ```

2. **Start development:**
   ```bash
   arvia serve
   ```
   Browser opens at `http://localhost:8080` with live reload

3. **Build for production:**
   ```bash
   arvia build
   ```
   Static files ready in `dist/` folder

## Commands

- `arvia init [name]` - Create new project
- `arvia serve` - Start development server with live reload
- `arvia build` - Build project for production  
- `arvia preview` - Preview built project
- `arvia version` - Show version
- `arvia help` - Show help

## Project Structure

```
my-website/
├── src/
│   └── index.html          # Main HTML file
├── assets/
│   ├── css/style.css       # Stylesheets
│   ├── js/app.js          # JavaScript
│   └── img/               # Images
├── dist/                   # Build output (auto-generated)
└── arvia.json             # Project configuration
```

## Configuration

Edit `arvia.json`:

```json
{
  "name": "my-website",
  "version": "1.0.0", 
  "source": "src",
  "build": "dist",
  "assets": "assets",
  "port": 8080
}
```

## Development Workflow

1. Edit HTML files in `src/`
2. Add CSS to `assets/css/`
3. Add JavaScript to `assets/js/`
4. Add images to `assets/img/`
5. Browser auto-reloads on file changes
6. Run `arvia build` when ready to deploy

## Deployment

After `arvia build`, upload the `dist/` folder to:
- **Static hosts**: Netlify, Vercel, GitHub Pages
- **CDN**: AWS S3 + CloudFront 
- **Traditional hosting**: Any web server

## Why Arvia?

- **No complex tooling** - Just Go binary
- **No framework lock-in** - Standard HTML/CSS/JS
- **Fast development** - Live reload built-in
- **Easy deployment** - Static files work everywhere
- **Simple learning curve** - Familiar web technologies

## Examples

### Basic HTML Page
```html
<!DOCTYPE html>
<html>
<head>
    <title>My Site</title>
    <link rel="stylesheet" href="../assets/css/style.css">
</head>
<body>
    <h1>Hello World!</h1>
    <script src="../assets/js/app.js"></script>
</body>
</html>
```

### Multi-page Site
```
src/
├── index.html
├── about.html
├── contact.html
└── blog/
    ├── index.html
    └── post-1.html
```

Access at:
- `http://localhost:8080/` (index.html)
- `http://localhost:8080/about.html`
- `http://localhost:8080/blog/`

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

MIT License - see LICENSE file

---

Built with ❤️ in Go