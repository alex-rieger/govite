package govite

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
)

// Govite configuration object
type Config struct {
	// Run govite in development mode. This will include the vite websocket client
	// and point assets to vite dev server.
	// If false websocket client returns an empty string
	DevServerEnabled bool

	// Protocol for the vite dev server, default: http
	DevServerProtocol string

	// Host for vite dev server, default: "localhost"
	DevServerHost string

	// Port for vite dev server, default: "3001"
	DevServerPort string

	// Browser url for vite websocket client, default" @vite/client
	// check https://vitejs.dev/guide/backend-integration.html#backend-integration for reference
	WebSocketClientUrl string

	// Todo: i don't really know
	AssetsPath string

	// Filepath the where the generated vite manifest file is located.
	ManifestPath string
}

var ConfigDefault = Config{
	DevServerEnabled:   false,
	DevServerProtocol:  "http",
	DevServerHost:      "localhost",
	DevServerPort:      "3001",
	WebSocketClientUrl: "@vite/client",
	AssetsPath:         "./web/app/assets",
	ManifestPath:       "./web/app/dist/manifest.json",
}

// sets defaults for passed configurations
func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.DevServerProtocol == "" {
		cfg.DevServerProtocol = ConfigDefault.DevServerProtocol
	}
	if cfg.DevServerHost == "" {
		cfg.DevServerHost = ConfigDefault.DevServerHost
	}
	if cfg.DevServerPort == "" {
		cfg.DevServerPort = ConfigDefault.DevServerPort
	}
	if cfg.WebSocketClientUrl == "" {
		cfg.WebSocketClientUrl = ConfigDefault.WebSocketClientUrl
	}
	if cfg.AssetsPath == "" {
		cfg.AssetsPath = ConfigDefault.AssetsPath
	}
	if cfg.ManifestPath == "" {
		cfg.ManifestPath = ConfigDefault.ManifestPath
	}

	return cfg
}

type ViteAsset struct {
	File    string
	Src     string
	IsEntry bool
	Imports []string
	Css     []string
	Assets  []string
}

type ViteManifest map[string]ViteAsset

// get asset by name from vite manifest file
func lookupAsset(asset string, manifest ViteManifest) (ViteAsset, bool) {
	assetData, found := manifest[asset]
	return assetData, found
}

// check "Imports" and "Css" properties of vite asset
// add modulepreload tag for "Imports" entries
// add css tag for "Css" entries
func handleAssetDeps(asset ViteAsset) string {
	result := "<script type=\"module\" crossorigin src=\"" + asset.File + "\"></script>\r\n"

	if len(asset.Css) > 0 {
		for _, assetName := range asset.Css {
			result += "<link rel=\"stylesheet\" href=\"" + assetName + "\">\r\n"
		}
	}

	if len(asset.Imports) > 0 {
		for _, importName := range asset.Imports {
			result += "<link rel=\"modulepreload\" href=\"" + importName + "\">\r\n"
		}
	}

	return result
}

type Instance struct {
	// adds vite hot reload client to template in dev mode
	// returns empty string in production mode
	TemplateTagViteClient func() template.HTML

	// adds asset import to template
	// returns script tag in dev mode
	// in production mode also returns css / modulepreload tags
	TemplateTagAsset func(entry string) template.HTML
}

// load file and parse to ViteManifest
func toManifest(manifestPath string) ViteManifest {
	mfst := &ViteManifest{}
	file, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		panic("govite: manifest file not found")
	}
	err = json.Unmarshal(file, mfst)
	if err != nil {
		panic("govite: failed to unmarshal manifest")
	}
	return *mfst
}

// create new instance of govite
func New(config ...Config) *Instance {
	cfg := configDefault(config...)
	mfst := toManifest(cfg.ManifestPath)

	return &Instance{
		TemplateTagViteClient: func() template.HTML {
			if cfg.DevServerEnabled {
				return template.HTML("<script type=\"module\" src=\"" + cfg.DevServerProtocol + "://" + cfg.DevServerHost + ":" + cfg.DevServerPort + "/" + cfg.WebSocketClientUrl + "\"></script>")
			}
			return ""
		},
		TemplateTagAsset: func(asset string) template.HTML {
			if cfg.DevServerEnabled {
				return template.HTML("<script type=\"module\" src=\"" + cfg.DevServerProtocol + "://" + cfg.DevServerHost + ":" + cfg.DevServerPort + "/" + asset + "\"></script>")
			}
			assetData, found := lookupAsset(asset, mfst)
			if !found {
				return ""
			}

			result := handleAssetDeps(assetData)
			return template.HTML(result)
		},
	}
}
