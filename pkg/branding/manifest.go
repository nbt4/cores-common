package branding

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ManifestOptions contains stable service-specific PWA metadata and bundled
// icon fallbacks. Uploaded assets in Config take precedence.
type ManifestOptions struct {
	Name             string
	Description      string
	StartURL         string
	Scope            string
	ThemeColor       string
	BackgroundColor  string
	FallbackIcon192  string
	FallbackIcon512  string
	FallbackMaskable string
}

// Manifest returns a web app manifest that follows the active branding.
func Manifest(cfg Config, options ManifestOptions) map[string]any {
	name := firstNonEmpty(options.Name, cfg.ProductName)
	if cfg.ProductName != "" {
		name = cfg.ProductName
	}
	description := options.Description
	operator := firstNonEmpty(cfg.BrandName, cfg.CompanyName)
	if description == "" && operator != "" && operator != name {
		description = fmt.Sprintf("%s von %s", name, operator)
	}
	startURL := firstNonEmpty(options.StartURL, "/")
	scope := firstNonEmpty(options.Scope, "/")
	theme := firstNonEmpty(options.ThemeColor, "#0b0b0b")
	background := firstNonEmpty(options.BackgroundColor, theme)

	icons := make([]map[string]string, 0, 3)
	if cfg.Assets.AppIcon != "" {
		icons = append(icons, manifestIcon(cfg.Assets.AppIcon, "512x512", "any"))
	} else {
		if options.FallbackIcon192 != "" {
			icons = append(icons, manifestIcon(options.FallbackIcon192, "192x192", "any"))
		}
		if options.FallbackIcon512 != "" {
			icons = append(icons, manifestIcon(options.FallbackIcon512, "512x512", "any"))
		}
	}
	maskable := firstNonEmpty(cfg.Assets.MaskableIcon, options.FallbackMaskable)
	if maskable != "" {
		icons = append(icons, manifestIcon(maskable, "512x512", "maskable"))
	}

	return map[string]any{
		"id":               startURL,
		"name":             name,
		"short_name":       name,
		"description":      description,
		"lang":             "de",
		"start_url":        startURL,
		"scope":            scope,
		"display":          "standalone",
		"display_override": []string{"standalone"},
		"orientation":      "any",
		"background_color": background,
		"theme_color":      theme,
		"categories":       []string{"business", "productivity"},
		"icons":            icons,
	}
}

func manifestIcon(path, sizes, purpose string) map[string]string {
	contentType := "image/png"
	if strings.EqualFold(filepath.Ext(strings.SplitN(path, "?", 2)[0]), ".svg") {
		contentType = "image/svg+xml"
		sizes = "any"
	}
	return map[string]string{
		"src": path, "sizes": sizes, "type": contentType, "purpose": purpose,
	}
}
