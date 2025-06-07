
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
