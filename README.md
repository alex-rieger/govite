# govite

## dev mode setup with module reloading

### initialize the instance
```go
vite := govite.New(govite.Config{
    DevServerEnabled:   false,
	DevServerProtocol:  "http",
	DevServerHost:      "localhost",
	DevServerPort:      "3001",
	WebSocketClientUrl: "@vite/client",
	AssetsPath:         "./web/app/assets",
	ManifestPath:       "./web/app/dist/manifest.json",
})
```

### add template tags to template/html engine
```go
myFuncs := template.FuncMap{
	"vite": TemplateTagViteClient,
	"asset": TemplateTagAsset,
}
```

### use in templates
todo: only head tag supported
```html
<head>
	<!-- adds vite hmr client -->
	{{ vite }}

	<!-- adds entrypoint (+ css / module imports in production) -->
	{{ asset "src/main.ts" }}
</head>
```