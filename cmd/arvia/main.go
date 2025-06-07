// cmd/arvia/main.go
package main

import (
	"encoding/json"
	// "flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

// Config represents the project configuration
type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`
	Build   string `json:"build"`
	Assets  string `json:"assets"`
	Port    int    `json:"port"`
}

// Default configuration
var defaultConfig = Config{
	Name:    "arvia-app",
	Version: "1.0.0",
	Source:  "src",
	Build:   "dist",
	Assets:  "assets",
	Port:    8080,
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients = make(map[*websocket.Conn]bool)
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInit()
	case "serve", "dev":
		handleServe()
	case "build":
		handleBuild()
	case "preview":
		handlePreview()
	case "version":
		fmt.Println("Arvia Framework v2.0.0")
	case "help":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func handleInit() {
	var projectName string
	if len(os.Args) > 2 {
		projectName = os.Args[2]
	} else {
		projectName = "my-arvia-app"
	}

	createProject(projectName)
}

func createProject(name string) {
	fmt.Printf("🚀 Creating Arvia project: %s\n", name)

	// Create directory structure
	dirs := []string{
		filepath.Join(name, "src"),
		filepath.Join(name, "assets", "css"),
		filepath.Join(name, "assets", "js"),
		filepath.Join(name, "assets", "img"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("Error creating directory:", err)
		}
	}

	// Create config file
	config := defaultConfig
	config.Name = name

	configData, _ := json.MarshalIndent(config, "", "  ")
	writeFile(filepath.Join(name, "arvia.json"), string(configData))

	// Create main HTML file
	htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Arvia App</title>
    <link rel="stylesheet" href="../assets/css/style.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>🚀 Welcome to Arvia</h1>
            <p>Simple Go Web Framework</p>
        </header>
        
        <main>
            <div class="card">
                <h2>Getting Started</h2>
                <p>Edit <code>src/index.html</code> to start building your app.</p>
                <button onclick="showMessage()">Click Me!</button>
                <p id="message"></p>
            </div>
        </main>
    </div>
    
    <script src="../assets/js/app.js"></script>
</body>
</html>`

	// Create CSS file
	cssContent := `/* Arvia Default Styles */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    line-height: 1.6;
    color: #333;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
}

header {
    text-align: center;
    margin-bottom: 3rem;
    color: white;
}

header h1 {
    font-size: 3rem;
    margin-bottom: 0.5rem;
    text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

header p {
    font-size: 1.2rem;
    opacity: 0.9;
}

.card {
    background: rgba(255, 255, 255, 0.95);
    border-radius: 15px;
    padding: 2rem;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(10px);
    text-align: center;
}

.card h2 {
    color: #667eea;
    margin-bottom: 1rem;
}

.card p {
    margin-bottom: 1rem;
}

button {
    background: linear-gradient(45deg, #667eea, #764ba2);
    color: white;
    border: none;
    padding: 0.75rem 2rem;
    border-radius: 25px;
    cursor: pointer;
    font-size: 1rem;
    transition: transform 0.2s;
}

button:hover {
    transform: translateY(-2px);
}

code {
    background: #f4f4f4;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    font-family: 'Courier New', monospace;
}

#message {
    margin-top: 1rem;
    font-weight: bold;
    color: #667eea;
}

@media (max-width: 768px) {
    .container {
        padding: 1rem;
    }
    
    header h1 {
        font-size: 2rem;
    }
}`

	// Create JS file
	jsContent := `// Arvia App JavaScript
console.log('🚀 Arvia app loaded!');

function showMessage() {
    const messageEl = document.getElementById('message');
    messageEl.textContent = '🎉 Hello from Arvia Framework!';
    
    // Add animation
    messageEl.style.transform = 'scale(1.1)';
    setTimeout(() => {
        messageEl.style.transform = 'scale(1)';
    }, 200);
}

// Live reload (injected in development)
if (window.location.hostname === 'localhost') {
    console.log('🔄 Development mode - live reload enabled');
}`

	// Write files
	writeFile(filepath.Join(name, "src", "index.html"), htmlContent)
	writeFile(filepath.Join(name, "assets", "css", "style.css"), cssContent)
	writeFile(filepath.Join(name, "assets", "js", "app.js"), jsContent)

	fmt.Printf("✅ Project '%s' created successfully!\n", name)
	fmt.Println("📁 Structure:")
	fmt.Printf("   %s/\n", name)
	fmt.Println("   ├── src/index.html")
	fmt.Println("   ├── assets/css/style.css")
	fmt.Println("   ├── assets/js/app.js")
	fmt.Println("   └── arvia.json")
	fmt.Println("")
	fmt.Println("🔧 Get started:")
	fmt.Printf("   cd %s\n", name)
	fmt.Println("   arvia serve")
}

func handleServe() {
	config := loadConfig()
	if config == nil {
		return
	}

	if _, err := os.Stat(config.Source); os.IsNotExist(err) {
		fmt.Printf("❌ Source directory '%s' not found\n", config.Source)
		return
	}

	fmt.Println("🔄 Starting Arvia development server...")
	fmt.Printf("📁 Serving: %s\n", config.Source)
	fmt.Printf("🌐 URL: http://localhost:%d\n", config.Port)
	fmt.Println("🔄 Live reload: enabled")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop")

	// Start file watcher
	go startFileWatcher(config.Source, config.Assets)

	// Setup routes
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveDevelopmentFiles(w, r, config)
	})

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

func handleBuild() {
	config := loadConfig()
	if config == nil {
		return
	}

	fmt.Println("📦 Building Arvia project...")

	// Clean build directory
	if err := os.RemoveAll(config.Build); err != nil {
		fmt.Printf("Warning: Could not clean build directory: %v\n", err)
	}

	// Copy source files
	if err := copyDir(config.Source, config.Build); err != nil {
		log.Fatal("Error copying source files:", err)
	}
	fmt.Println("✅ Copied source files")

	// Copy assets
	if _, err := os.Stat(config.Assets); err == nil {
		assetsDestination := filepath.Join(config.Build, "assets")
		if err := copyDir(config.Assets, assetsDestination); err != nil {
			log.Fatal("Error copying assets:", err)
		}
		fmt.Println("✅ Copied assets")
	}

	fmt.Printf("🎉 Build completed: %s\n", config.Build)
	fmt.Println("📤 Ready for deployment!")
}

func handlePreview() {
	config := loadConfig()
	if config == nil {
		return
	}

	if _, err := os.Stat(config.Build); os.IsNotExist(err) {
		fmt.Println("❌ Build not found. Run 'arvia build' first.")
		return
	}

	port := config.Port + 1 // Use different port
	fmt.Println("🔍 Starting preview server...")
	fmt.Printf("📁 Serving: %s\n", config.Build)
	fmt.Printf("🌐 URL: http://localhost:%d\n", port)
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop")

	http.Handle("/", http.FileServer(http.Dir(config.Build)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func serveDevelopmentFiles(w http.ResponseWriter, r *http.Request, config *Config) {
	// Handle root path
	if r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}

	// Try to serve from source directory
	sourcePath := filepath.Join(config.Source, r.URL.Path)
	if _, err := os.Stat(sourcePath); err == nil {
		// Inject live reload script for HTML files
		if strings.HasSuffix(r.URL.Path, ".html") {
			content, err := os.ReadFile(sourcePath)
			if err != nil {
				http.Error(w, "File not found", 404)
				return
			}

			// Inject live reload script
			htmlContent := string(content)
			liveReloadScript := `
<script>
    const ws = new WebSocket('ws://localhost:` + fmt.Sprintf("%d", config.Port) + `/ws');
    ws.onmessage = function(event) {
        if (event.data === 'reload') {
            location.reload();
        }
    };
</script>`

			if idx := strings.LastIndex(htmlContent, "</body>"); idx != -1 {
				htmlContent = htmlContent[:idx] + liveReloadScript + htmlContent[idx:]
			}

			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Cache-Control", "no-cache")
			w.Write([]byte(htmlContent))
			return
		}

		http.ServeFile(w, r, sourcePath)
		return
	}

	// Try to serve from assets directory
	assetsPath := filepath.Join(config.Assets, strings.TrimPrefix(r.URL.Path, "/assets/"))
	if _, err := os.Stat(assetsPath); err == nil {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, assetsPath)
		return
	}

	http.Error(w, "File not found", 404)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	// Keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func startFileWatcher(dirs ...string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("File watcher error:", err)
		return
	}
	defer watcher.Close()

	// Add directories to watch
	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					watcher.Add(path)
				}
				return nil
			})
		}
	}

	// Debounce file changes
	var lastEvent time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				now := time.Now()
				if now.Sub(lastEvent) > 100*time.Millisecond {
					fmt.Printf("🔄 File changed: %s\n", event.Name)
					notifyClients()
					lastEvent = now
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("File watcher error:", err)
		}
	}
}

func notifyClients() {
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func loadConfig() *Config {
	if _, err := os.Stat("arvia.json"); os.IsNotExist(err) {
		fmt.Println("❌ Error: arvia.json not found. Run 'arvia init' first.")
		return nil
	}

	data, err := os.ReadFile("arvia.json")
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal("Error parsing config:", err)
	}

	return &config
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		log.Fatal("Error writing file:", err)
	}
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	os.MkdirAll(filepath.Dir(dst), 0755)

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func showHelp() {
	fmt.Println("🚀 Arvia Framework - Simple Go Web Development")
	fmt.Println("")
	fmt.Println("Usage: arvia <command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  init [name]     Create new project")
	fmt.Println("  serve           Start development server with live reload")
	fmt.Println("  build           Build project for production")
	fmt.Println("  preview         Preview built project")
	fmt.Println("  version         Show version")
	fmt.Println("  help            Show this help")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  arvia init my-app")
	fmt.Println("  cd my-app")
	fmt.Println("  arvia serve")
	fmt.Println("  arvia build")
}