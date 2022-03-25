package govite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefault_DefaultValuesMatchSnapshot(t *testing.T) {
	result := configDefault()

	assert.Equal(t, result.DevServerEnabled, false, "DevServerEnabled matches snapshot")
	assert.Equal(t, result.DevServerProtocol, "http", "DevServerProtocol matches snapshot")
	assert.Equal(t, result.DevServerHost, "localhost", "DevServerHost matches snapshot")
	assert.Equal(t, result.DevServerPort, "3001", "DevServerPort matches snapshot")
	assert.Equal(t, result.WebSocketClientUrl, "@vite/client", "WebSocketClientUrl matches snapshot")
	assert.Equal(t, result.AssetsPath, "./web/app/assets", "AssetsPath matches snapshot")
	assert.Equal(t, result.ManifestPath, "./web/app/dist/manifest.json", "ManifestPath matches Snapshot")
}

func TestConfigDefault_OptionsMergingWorks(t *testing.T) {
	result := configDefault(Config{
		DevServerEnabled:  true,
		DevServerProtocol: "https",
	})

	assert.Equal(t, result.DevServerEnabled, true, "default value gets overridden")
	assert.Equal(t, result.DevServerProtocol, "https", "default value gets overridden")
	assert.Equal(t, result.DevServerHost, "localhost", "not passed value gets defaulted")
}

func TestLookupAsset_ExistingAssetIsReturned(t *testing.T) {
	manifest := ViteManifest{
		"assetFilename.suffix": ViteAsset{},
	}
	resultAsset, found := lookupAsset("assetFilename.suffix", manifest)

	assert.Equal(t, found, true, "found assets returns true")
	assert.NotEqual(t, resultAsset, nil, "found asset is not nil")
}

func TestLookupAsset_NotExistingAsset(t *testing.T) {
	manifest := ViteManifest{}
	resultAsset, found := lookupAsset("notExistingAsset.suffix", manifest)

	assert.Equal(t, found, false, "not found asset returns false")
	assert.Equal(t, resultAsset, ViteAsset{}, "not found asset is empty asset struct")
}

func TestHandleAssetDeps_ExpectedScriptTagsAreReturned(t *testing.T) {
	table := []struct {
		asset       ViteAsset
		expected    string
		description string
	}{
		{
			asset: ViteAsset{
				File:    "file.js",
				Src:     "something/file.js",
				IsEntry: false,
				Imports: []string{},
				Css:     []string{},
				Assets:  []string{},
			},
			expected:    "<script type=\"module\" crossorigin src=\"file.js\"></script>\r\n",
			description: "asset without deps generated correct script tag",
		},
		{
			asset: ViteAsset{
				File:    "file.js",
				Src:     "something/file.js",
				IsEntry: false,
				Imports: []string{"otherfile.js"},
				Css:     []string{},
				Assets:  []string{},
			},
			expected:    "<script type=\"module\" crossorigin src=\"file.js\"></script>\r\n<link rel=\"modulepreload\" href=\"otherfile.js\">\r\n",
			description: "asset with imports deps generated correct script tag",
		},
		{
			asset: ViteAsset{
				File:    "file.js",
				Src:     "something/otherfile.js",
				IsEntry: false,
				Imports: []string{},
				Css:     []string{"othercss.css"},
				Assets:  []string{},
			},
			expected:    "<script type=\"module\" crossorigin src=\"file.js\"></script>\r\n<link rel=\"stylesheet\" href=\"othercss.css\">\r\n",
			description: "asset with css deps generated correct script tag",
		},
		{
			asset: ViteAsset{
				File:    "file.js",
				Src:     "something/otherfile.js",
				IsEntry: false,
				Imports: []string{"otherfile.js"},
				Css:     []string{"othercss.css"},
				Assets:  []string{},
			},
			expected:    "<script type=\"module\" crossorigin src=\"file.js\"></script>\r\n<link rel=\"stylesheet\" href=\"othercss.css\">\r\n<link rel=\"modulepreload\" href=\"otherfile.js\">\r\n",
			description: "asset with imports and css deps generated correct script tag",
		},
	}

	for _, tc := range table {
		result := handleAssetDeps(tc.asset)
		assert.Equal(t, tc.expected, result, tc.description)
	}
}
