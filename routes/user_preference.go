package routes

import (
	"godoor/models"
	"godoor/serializers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserPreference returns current user's preferences as array
func GetUserPreference(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	preference, err := models.GetOrCreateUserPreference(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user preferences"})
		return
	}

	response := []serializers.UserPreferenceSerializer{
		serializers.SerializeUserPreference(preference),
	}
	c.JSON(http.StatusOK, response)
}

// UpdateUserPreference updates user's preferences
func UpdateUserPreference(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	preference, err := models.GetOrCreateUserPreference(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user preferences"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Update allowed fields
	if theme, exists := updates["theme"]; exists {
		if themeStr, ok := theme.(string); ok {
			preference.Theme = themeStr
		}
	}
	if navBgColor, exists := updates["nav_bg_color"]; exists {
		if navBgColorStr, ok := navBgColor.(string); ok {
			preference.NavBgColor = navBgColorStr
		}
	}
	if navTextColor, exists := updates["nav_text_color"]; exists {
		if navTextColorStr, ok := navTextColor.(string); ok {
			preference.NavTextColor = navTextColorStr
		}
	}
	if navShowLogo, exists := updates["nav_show_logo"]; exists {
		if navShowLogoBool, ok := navShowLogo.(bool); ok {
			preference.NavShowLogo = navShowLogoBool
		}
	}
	if defaultUnit, exists := updates["default_unit"]; exists {
		if defaultUnitStr, ok := defaultUnit.(string); ok {
			preference.DefaultUnit = defaultUnitStr
		}
	}
	if defaultPage, exists := updates["default_page"]; exists {
		if defaultPageStr, ok := defaultPage.(string); ok {
			preference.DefaultPage = defaultPageStr
		}
	}
	if useFractions, exists := updates["use_fractions"]; exists {
		if useFractionsBool, ok := useFractions.(bool); ok {
			preference.UseFractions = useFractionsBool
		}
	}
	if useKj, exists := updates["use_kj"]; exists {
		if useKjBool, ok := useKj.(bool); ok {
			preference.UseKj = useKjBool
		}
	}
	if navSticky, exists := updates["nav_sticky"]; exists {
		if navStickyBool, ok := navSticky.(bool); ok {
			preference.NavSticky = navStickyBool
		}
	}
	if ingredientDecimals, exists := updates["ingredient_decimals"]; exists {
		switch v := ingredientDecimals.(type) {
		case float64:
			preference.IngredientDecimals = int(v)
		case int:
			preference.IngredientDecimals = v
		}
	}
	if comments, exists := updates["comments"]; exists {
		if commentsBool, ok := comments.(bool); ok {
			preference.Comments = commentsBool
		}
	}
	if shoppingAutoSync, exists := updates["shopping_auto_sync"]; exists {
		switch v := shoppingAutoSync.(type) {
		case float64:
			preference.ShoppingAutoSync = int(v)
		case int:
			preference.ShoppingAutoSync = v
		}
	}
	if mealplanAutoaddShopping, exists := updates["mealplan_autoadd_shopping"]; exists {
		if mealplanAutoaddShoppingBool, ok := mealplanAutoaddShopping.(bool); ok {
			preference.MealplanAutoaddShopping = mealplanAutoaddShoppingBool
		}
	}
	if defaultDelay, exists := updates["default_delay"]; exists {
		if defaultDelayFloat, ok := defaultDelay.(float64); ok {
			preference.DefaultDelay = defaultDelayFloat
		}
	}
	if mealplanAutoincludeRelated, exists := updates["mealplan_autoinclude_related"]; exists {
		if mealplanAutoincludeRelatedBool, ok := mealplanAutoincludeRelated.(bool); ok {
			preference.MealplanAutoincludeRelated = mealplanAutoincludeRelatedBool
		}
	}
	if mealplanAutoexcludeOnhand, exists := updates["mealplan_autoexclude_onhand"]; exists {
		if mealplanAutoexcludeOnhandBool, ok := mealplanAutoexcludeOnhand.(bool); ok {
			preference.MealplanAutoexcludeOnhand = mealplanAutoexcludeOnhandBool
		}
	}
	if shoppingRecentDays, exists := updates["shopping_recent_days"]; exists {
		switch v := shoppingRecentDays.(type) {
		case float64:
			preference.ShoppingRecentDays = int(v)
		case int:
			preference.ShoppingRecentDays = v
		}
	}
	if csvDelim, exists := updates["csv_delim"]; exists {
		if csvDelimStr, ok := csvDelim.(string); ok {
			preference.CsvDelim = csvDelimStr
		}
	}
	if csvPrefix, exists := updates["csv_prefix"]; exists {
		if csvPrefixStr, ok := csvPrefix.(string); ok {
			preference.CsvPrefix = csvPrefixStr
		}
	}
	if shoppingUpdateFoodLists, exists := updates["shopping_update_food_lists"]; exists {
		if shoppingUpdateFoodListsBool, ok := shoppingUpdateFoodLists.(bool); ok {
			preference.ShoppingUpdateFoodLists = shoppingUpdateFoodListsBool
		}
	}
	if filterToSupermarket, exists := updates["filter_to_supermarket"]; exists {
		if filterToSupermarketBool, ok := filterToSupermarket.(bool); ok {
			preference.FilterToSupermarket = filterToSupermarketBool
		}
	}
	if shoppingAddOnhand, exists := updates["shopping_add_onhand"]; exists {
		if shoppingAddOnhandBool, ok := shoppingAddOnhand.(bool); ok {
			preference.ShoppingAddOnhand = shoppingAddOnhandBool
		}
	}
	if leftHanded, exists := updates["left_handed"]; exists {
		if leftHandedBool, ok := leftHanded.(bool); ok {
			preference.LeftHanded = leftHandedBool
		}
	}
	if showStepIngredients, exists := updates["show_step_ingredients"]; exists {
		if showStepIngredientsBool, ok := showStepIngredients.(bool); ok {
			preference.ShowStepIngredients = showStepIngredientsBool
		}
	}

	// Find existing preference
	var existingPreference models.UserPreference
	if err := models.DB.Where("user_id = ?", user.ID).First(&existingPreference).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user preferences"})
		return
	}

	// Create updates map with only changed fields
	updates = make(map[string]interface{})
	if preference.Theme != existingPreference.Theme {
		updates["theme"] = preference.Theme
	}
	if preference.UseFractions != existingPreference.UseFractions {
		updates["use_fractions"] = preference.UseFractions
	}

	if len(updates) > 0 {
		if err := models.DB.Model(&existingPreference).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user preferences"})
			return
		}
	}

	preference = &existingPreference

	response := []serializers.UserPreferenceSerializer{
		serializers.SerializeUserPreference(preference),
	}
	c.JSON(http.StatusOK, response)
}
