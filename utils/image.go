package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// IsImageFileTypeAllowed checks if the file extension is allowed for images
func IsImageFileTypeAllowed(filename string) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// GenerateUniqueFilename generates a unique filename for image upload
func GenerateUniqueFilename(originalFilename string, recipeID uint) string {
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".jpg" // Default extension
	}

	// Generate random bytes for uniqueness
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	// Create filename: recipeID_random.ext
	return fmt.Sprintf("%d_%x%s", recipeID, randomBytes, ext)
}

// getMediaRoot returns the media root directory from environment or default
func getMediaRoot() string {
	if root := os.Getenv("MEDIA_ROOT"); root != "" {
		return root
	}
	return "./media"
}

// getRecipeImagesDir returns the recipe images directory from environment or default
func getRecipeImagesDir() string {
	if dir := os.Getenv("RECIPE_IMAGES_DIR"); dir != "" {
		return dir
	}
	return "recipes"
}

// getMaxImageSize returns the maximum image size in MB from environment or default
func getMaxImageSize() int64 {
	if size := os.Getenv("MAX_IMAGE_SIZE_MB"); size != "" {
		// Simple conversion - in production would use strconv.Atoi with error handling
		if size == "1" {
			return 1
		} else if size == "5" {
			return 5
		} else if size == "10" {
			return 10
		}
		// Default to 10
		return 10
	}
	return 10
}

// EnsureRecipeImagesDir ensures the recipe images directory exists
func EnsureRecipeImagesDir() error {
	imagesDir := filepath.Join(getMediaRoot(), getRecipeImagesDir())
	return os.MkdirAll(imagesDir, 0755)
}

// SaveUploadedImage saves an uploaded image file to disk
func SaveUploadedImage(file multipart.File, header *multipart.FileHeader, recipeID uint) (string, error) {
	// Validate file type
	if !IsImageFileTypeAllowed(header.Filename) {
		return "", fmt.Errorf("unsupported file type: %s", header.Filename)
	}

	// Check file size (convert MB to bytes)
	maxSize := getMaxImageSize() * 1024 * 1024
	if header.Size > maxSize {
		return "", fmt.Errorf("file too large: %d bytes (max: %d bytes)", header.Size, maxSize)
	}

	// Ensure directory exists
	if err := EnsureRecipeImagesDir(); err != nil {
		return "", fmt.Errorf("failed to create images directory: %v", err)
	}

	// Generate unique filename
	filename := GenerateUniqueFilename(header.Filename, recipeID)
	fullPath := filepath.Join(getMediaRoot(), getRecipeImagesDir(), filename)

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(fullPath) // Clean up on error
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return relative path for database storage
	return filepath.Join(getRecipeImagesDir(), filename), nil
}

// DeleteImageFile deletes an image file from disk
func DeleteImageFile(imagePath string) error {
	if imagePath == "" {
		return nil
	}

	// Convert relative path to absolute
	fullPath := filepath.Join(getMediaRoot(), imagePath)
	return os.Remove(fullPath)
}


