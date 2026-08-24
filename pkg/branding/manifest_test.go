package branding

import "testing"

func TestManifestUsesCustomApplicationIcons(t *testing.T) {
	cfg := Config{
		ProductName: "RentalCore",
		BrandName:   "Tsunami Events",
		Assets: AssetSet{
			AppIcon:      "/logos/rental-app.png?v=2",
			MaskableIcon: "/logos/rental-maskable.svg?v=2",
		},
	}
	manifest := Manifest(cfg, ManifestOptions{
		StartURL: "/rental/", Scope: "/rental/",
		FallbackIcon192: "/fallback-192.png",
	})

	if manifest["name"] != "RentalCore" || manifest["start_url"] != "/rental/" {
		t.Fatalf("unexpected identity fields: %#v", manifest)
	}
	icons := manifest["icons"].([]map[string]string)
	if len(icons) != 2 {
		t.Fatalf("expected app and maskable icon, got %#v", icons)
	}
	if icons[0]["src"] != cfg.Assets.AppIcon || icons[0]["type"] != "image/png" {
		t.Fatalf("unexpected app icon: %#v", icons[0])
	}
	if icons[1]["src"] != cfg.Assets.MaskableIcon || icons[1]["type"] != "image/svg+xml" || icons[1]["sizes"] != "any" {
		t.Fatalf("unexpected maskable icon: %#v", icons[1])
	}
}

func TestManifestKeepsBundledIconsWhenOnlyFaviconIsCustomized(t *testing.T) {
	manifest := Manifest(Config{ProductName: "Cores", Assets: AssetSet{Favicon: "/portrait.svg"}}, ManifestOptions{
		FallbackIcon192: "/app-icons/icon-192.png",
		FallbackIcon512: "/app-icons/icon-512.png",
	})
	icons := manifest["icons"].([]map[string]string)
	if len(icons) != 2 || icons[0]["src"] != "/app-icons/icon-192.png" {
		t.Fatalf("favicon must not become an install icon: %#v", icons)
	}
}
