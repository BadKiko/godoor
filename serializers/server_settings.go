package serializers

import (
	"godoor/config"
)

// ServerSettingsSerializer matching Tandoor ServerSettingsSerializer
type ServerSettingsSerializer struct {
	ShoppingMinAutosyncInterval int    `json:"shopping_min_autosync_interval"`
	EnablePdfExport             bool   `json:"enable_pdf_export"`
	DisableExternalConnectors   bool   `json:"disable_external_connectors"`
	TermsURL                    string `json:"terms_url"`
	PrivacyURL                  string `json:"privacy_url"`
	ImprintURL                  string `json:"imprint_url"`
	Hosted                      bool   `json:"hosted"`
	Debug                       bool   `json:"debug"`
	Version                     string `json:"version"`

	UnauthenticatedThemeFromSpace int    `json:"unauthenticated_theme_from_space"`
	ForceThemeFromSpace           int    `json:"force_theme_from_space"`

	// Theme-related fields (simplified - no file handling yet)
	LogoColor32    interface{} `json:"logo_color_32"`
	LogoColor128   interface{} `json:"logo_color_128"`
	LogoColor144   interface{} `json:"logo_color_144"`
	LogoColor180   interface{} `json:"logo_color_180"`
	LogoColor192   interface{} `json:"logo_color_192"`
	LogoColor512   interface{} `json:"logo_color_512"`
	LogoColorSvg   interface{} `json:"logo_color_svg"`
	CustomSpaceTheme interface{} `json:"custom_space_theme"`
	NavLogo        interface{} `json:"nav_logo"`
	NavBgColor     interface{} `json:"nav_bg_color"`
}

// SerializeServerSettings creates ServerSettings from config
func SerializeServerSettings() ServerSettingsSerializer {
	return ServerSettingsSerializer{
		ShoppingMinAutosyncInterval:  config.ShoppingMinAutosyncInterval,
		EnablePdfExport:             config.EnablePDFExport,
		DisableExternalConnectors:   config.DisableExternalConnectors,
		TermsURL:                    config.TermsURL,
		PrivacyURL:                  config.PrivacyURL,
		ImprintURL:                  config.ImprintURL,
		Hosted:                      config.Hosted,
		Debug:                       config.Debug,
		Version:                     config.TandoorVersion,
		UnauthenticatedThemeFromSpace: config.UnauthenticatedThemeFromSpace,
		ForceThemeFromSpace:          config.ForceThemeFromSpace,

		// Theme fields - for now nil, can be extended later
		LogoColor32:     nil,
		LogoColor128:    nil,
		LogoColor144:    nil,
		LogoColor180:    nil,
		LogoColor192:    nil,
		LogoColor512:    nil,
		LogoColorSvg:    nil,
		CustomSpaceTheme: nil,
		NavLogo:         nil,
		NavBgColor:      nil,
	}
}
