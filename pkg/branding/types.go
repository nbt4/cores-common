package branding

import "time"

// Config is the central branding configuration singleton (table branding_config).
// Each service reads this table directly from the shared PostgreSQL database.
// Only cores-dashboard writes to it via the admin Branding UI.
type Config struct {
	ProductName     string   `json:"productName"`
	CompanyName     string   `json:"companyName"`
	BrandName       string   `json:"brandName"`
	Assets          AssetSet `json:"assets"`
	CompanyAssets   AssetSet `json:"companyAssets"`
	LogoSidebar     string   `json:"sidebarLogo"`
	LogoLogin       string   `json:"loginLogo"`
	FaviconPath     string   `json:"faviconPath"`
	LogoSizeSidebar int16    `json:"logoSizeSidebar"`
	LogoSizeLogin   int16    `json:"logoSizeLogin"`
}

// AssetSet describes logo files by their meaning rather than by a single UI
// position. "OnDark" assets are light artwork intended for dark surfaces;
// "OnLight" assets are dark artwork intended for light surfaces.
//
// SidebarLogo, LoginLogo and FaviconPath remain in Config as compatibility
// aliases while clients migrate to this structure.
type AssetSet struct {
	MarkOnDark        string `json:"markOnDark,omitempty"`
	MarkOnLight       string `json:"markOnLight,omitempty"`
	HorizontalOnDark  string `json:"horizontalOnDark,omitempty"`
	HorizontalOnLight string `json:"horizontalOnLight,omitempty"`
	StackedOnDark     string `json:"stackedOnDark,omitempty"`
	StackedOnLight    string `json:"stackedOnLight,omitempty"`
	Favicon           string `json:"favicon,omitempty"`
	AppIcon           string `json:"appIcon,omitempty"`
	MaskableIcon      string `json:"maskableIcon,omitempty"`
	Print             string `json:"print,omitempty"`
}

// Record is the GORM model for the branding_config table.
type Record struct {
	ID                     uint                `gorm:"primaryKey" json:"id"`
	CreatedAt              time.Time           `json:"createdAt"`
	UpdatedAt              time.Time           `json:"updatedAt"`
	CompanyName            string              `gorm:"column:company_name" json:"companyName"`
	BrandName              string              `gorm:"column:brand_name" json:"brandName"`
	LogoCoresSidebar       *string             `gorm:"column:logo_cores_sidebar" json:"logoCoresSidebar"`
	LogoCoresLogin         *string             `gorm:"column:logo_cores_login" json:"logoCoresLogin"`
	LogoRentalSidebar      *string             `gorm:"column:logo_rental_sidebar" json:"logoRentalSidebar"`
	LogoRentalLogin        *string             `gorm:"column:logo_rental_login" json:"logoRentalLogin"`
	LogoWarehouseSidebar   *string             `gorm:"column:logo_warehouse_sidebar" json:"logoWarehouseSidebar"`
	LogoWarehouseLogin     *string             `gorm:"column:logo_warehouse_login" json:"logoWarehouseLogin"`
	LogoPlannerSidebar     *string             `gorm:"column:logo_planner_sidebar" json:"logoPlannerSidebar"`
	LogoPlannerLogin       *string             `gorm:"column:logo_planner_login" json:"logoPlannerLogin"`
	LogoProcurementSidebar *string             `gorm:"column:logo_procurement_sidebar" json:"logoProcurementSidebar"`
	LogoProcurementLogin   *string             `gorm:"column:logo_procurement_login" json:"logoProcurementLogin"`
	FaviconCores           *string             `gorm:"column:favicon_cores" json:"faviconCores"`
	FaviconRental          *string             `gorm:"column:favicon_rental" json:"faviconRental"`
	FaviconWarehouse       *string             `gorm:"column:favicon_warehouse" json:"faviconWarehouse"`
	FaviconPlanner         *string             `gorm:"column:favicon_planner" json:"faviconPlanner"`
	FaviconProcurement     *string             `gorm:"column:favicon_procurement" json:"faviconProcurement"`
	FaviconPath            *string             `gorm:"column:favicon_path" json:"faviconPath"`
	Assets                 map[string]AssetSet `gorm:"column:assets_json;serializer:json;type:jsonb" json:"assets"`
	LogoSizeSidebar        int16               `gorm:"column:logo_size_sidebar" json:"logoSizeSidebar"`
	LogoSizeLogin          int16               `gorm:"column:logo_size_login" json:"logoSizeLogin"`
}

func (Record) TableName() string { return "branding_config" }
