package branding

import "time"

// Config is the central branding configuration singleton (table branding_config).
// Each service reads this table directly from the shared PostgreSQL database.
// Only cores-dashboard writes to it via the admin Branding UI.
type Config struct {
	CompanyName     string `json:"companyName"`
	BrandName       string `json:"brandName"`
	LogoSidebar     string `json:"sidebarLogo"`
	LogoLogin       string `json:"loginLogo"`
	FaviconPath     string `json:"faviconPath"`
	LogoSizeSidebar int16  `json:"logoSizeSidebar"`
	LogoSizeLogin   int16  `json:"logoSizeLogin"`
}

// Record is the GORM model for the branding_config table.
type Record struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	CompanyName          string    `gorm:"column:company_name" json:"companyName"`
	BrandName            string    `gorm:"column:brand_name" json:"brandName"`
	LogoCoresSidebar     *string   `gorm:"column:logo_cores_sidebar" json:"logoCoresSidebar"`
	LogoCoresLogin       *string   `gorm:"column:logo_cores_login" json:"logoCoresLogin"`
	LogoRentalSidebar    *string   `gorm:"column:logo_rental_sidebar" json:"logoRentalSidebar"`
	LogoRentalLogin      *string   `gorm:"column:logo_rental_login" json:"logoRentalLogin"`
	LogoWarehouseSidebar *string   `gorm:"column:logo_warehouse_sidebar" json:"logoWarehouseSidebar"`
	LogoWarehouseLogin   *string   `gorm:"column:logo_warehouse_login" json:"logoWarehouseLogin"`
	LogoPlannerSidebar   *string   `gorm:"column:logo_planner_sidebar" json:"logoPlannerSidebar"`
	LogoPlannerLogin     *string   `gorm:"column:logo_planner_login" json:"logoPlannerLogin"`
	FaviconCores         *string   `gorm:"column:favicon_cores" json:"faviconCores"`
	FaviconRental        *string   `gorm:"column:favicon_rental" json:"faviconRental"`
	FaviconWarehouse     *string   `gorm:"column:favicon_warehouse" json:"faviconWarehouse"`
	FaviconPlanner       *string   `gorm:"column:favicon_planner" json:"faviconPlanner"`
	FaviconPath          *string   `gorm:"column:favicon_path" json:"faviconPath"`
	LogoSizeSidebar      int16     `gorm:"column:logo_size_sidebar" json:"logoSizeSidebar"`
	LogoSizeLogin        int16     `gorm:"column:logo_size_login" json:"logoSizeLogin"`
}

func (Record) TableName() string { return "branding_config" }
