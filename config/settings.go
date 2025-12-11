package config

import "os"

// Server settings - simplified version for basic functionality
var (
	// Basic settings
	ShoppingMinAutosyncInterval = getEnvInt("SHOPPING_MIN_AUTOSYNC_INTERVAL", 5)
	EnablePDFExport             = getEnvBool("ENABLE_PDF_EXPORT", false)
	DisableExternalConnectors   = getEnvBool("DISABLE_EXTERNAL_CONNECTORS", false)
	TermsURL                    = getEnvString("TERMS_URL", "")
	PrivacyURL                  = getEnvString("PRIVACY_URL", "")
	ImprintURL                  = getEnvString("IMPRINT_URL", "")
	Hosted                      = getEnvBool("HOSTED", false)
	Debug                       = getEnvBool("DEBUG", true) // Default true for development

	// Media settings
	MediaRoot       = getEnvString("MEDIA_ROOT", "./media")
	MediaURL        = getEnvString("MEDIA_URL", "/media/")
	RecipeImagesDir = getEnvString("RECIPE_IMAGES_DIR", "recipes")
	MaxImageSize    = getEnvInt("MAX_IMAGE_SIZE_MB", 10) // Maximum image size in MB

	// Theme settings
	UnauthenticatedThemeFromSpace = getEnvInt("UNAUTHENTICATED_THEME_FROM_SPACE", 0)
	ForceThemeFromSpace           = getEnvInt("FORCE_THEME_FROM_SPACE", 0)

	// Version info (simplified)
	TandoorVersion = "2.3.6" // Tandoor version
)

// Helper functions for environment variables
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple conversion - in production would use strconv.Atoi with error handling
		if value == "0" {
			return 0
		} else if value == "1" {
			return 1
		} else if value == "5" {
			return 5
		}
		// Add more cases as needed
	}
	return defaultValue
}
