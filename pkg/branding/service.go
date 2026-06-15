package branding

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// cacheBuster appends a version query param to force browser re-fetch on change.
func cacheBuster(path string, updatedAt time.Time) string {
	if path == "" {
		return ""
	}
	return fmt.Sprintf("%s?v=%d", path, updatedAt.Unix())
}

// Service reads branding_config directly from the shared PostgreSQL DB.
// No caching — every call hits the DB for instant updates.
type Service struct {
	db      *gorm.DB
	service string
}

// NewService creates a new branding service for the specified service name
// ("cores", "rental", "warehouse", "planner").
func NewService(db *gorm.DB, service string) *Service {
	return &Service{db: db, service: service}
}

// GetConfig returns the branding configuration for the service.
func (s *Service) GetConfig() Config {
	var rec Record
	if err := s.db.First(&rec, 1).Error; err != nil {
		return Config{
			CompanyName:     s.defaultName(),
			LogoSizeSidebar: 100,
			LogoSizeLogin:   100,
		}
	}

	cfg := Config{
		CompanyName:     s.coalesceName(rec.CompanyName),
		BrandName:       rec.BrandName,
		FaviconPath:     cacheBuster(s.faviconFor(rec), rec.UpdatedAt),
		LogoSizeSidebar: s.coalesceSize(rec.LogoSizeSidebar),
		LogoSizeLogin:   s.coalesceSize(rec.LogoSizeLogin),
	}

	switch s.service {
	case "cores":
		cfg.LogoSidebar = cacheBuster(s.deref(rec.LogoCoresSidebar), rec.UpdatedAt)
		cfg.LogoLogin = cacheBuster(s.deref(rec.LogoCoresLogin), rec.UpdatedAt)
	case "rental":
		cfg.LogoSidebar = cacheBuster(s.deref(rec.LogoRentalSidebar), rec.UpdatedAt)
		cfg.LogoLogin = cacheBuster(s.deref(rec.LogoRentalLogin), rec.UpdatedAt)
	case "warehouse":
		cfg.LogoSidebar = cacheBuster(s.deref(rec.LogoWarehouseSidebar), rec.UpdatedAt)
		cfg.LogoLogin = cacheBuster(s.deref(rec.LogoWarehouseLogin), rec.UpdatedAt)
	case "planner":
		cfg.LogoSidebar = cacheBuster(s.deref(rec.LogoPlannerSidebar), rec.UpdatedAt)
		cfg.LogoLogin = cacheBuster(s.deref(rec.LogoPlannerLogin), rec.UpdatedAt)
	}

	return cfg
}

func (s *Service) faviconFor(rec Record) string {
	switch s.service {
	case "cores":
		return s.deref(rec.FaviconCores)
	case "rental":
		return s.deref(rec.FaviconRental)
	case "warehouse":
		return s.deref(rec.FaviconWarehouse)
	case "planner":
		return s.deref(rec.FaviconPlanner)
	}
	return s.deref(rec.FaviconPath)
}

func (s *Service) defaultName() string {
	switch s.service {
	case "rental":
		return "RentalCore"
	case "warehouse":
		return "WarehouseCore"
	case "planner":
		return "PlannerCore"
	default:
		return ""
	}
}

func (s *Service) coalesceName(dbName string) string {
	if dbName != "" {
		return dbName
	}
	return s.defaultName()
}

func (s *Service) coalesceSize(val int16) int16 {
	if val >= 50 && val <= 200 {
		return val
	}
	return 100
}

func (s *Service) deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
