package serializers

import (
	"godoor/models"
)

// UserPreferenceSerializer matching Tandoor UserPreferenceSerializer
type UserPreferenceSerializer struct {
	User                     UserSerializer `json:"user"`
	Image                    interface{}    `json:"image"` // TODO: implement UserFileViewSerializer
	Theme                    string         `json:"theme"`
	NavBgColor               string         `json:"nav_bg_color"`
	NavTextColor             string         `json:"nav_text_color"`
	NavShowLogo              bool           `json:"nav_show_logo"`
	DefaultUnit              string         `json:"default_unit"`
	DefaultPage              string         `json:"default_page"`
	UseFractions             bool           `json:"use_fractions"`
	UseKj                    bool           `json:"use_kj"`
	PlanShare                []UserSerializer `json:"plan_share"`
	NavSticky                bool           `json:"nav_sticky"`
	IngredientDecimals       int            `json:"ingredient_decimals"`
	Comments                 bool           `json:"comments"`
	ShoppingAutoSync         int            `json:"shopping_auto_sync"`
	MealplanAutoaddShopping  bool           `json:"mealplan_autoadd_shopping"`
	FoodInheritDefault       []interface{}  `json:"food_inherit_default"` // TODO: implement FoodInheritFieldSerializer
	DefaultDelay              float64        `json:"default_delay"`
	MealplanAutoincludeRelated bool          `json:"mealplan_autoinclude_related"`
	MealplanAutoexcludeOnhand bool          `json:"mealplan_autoexclude_onhand"`
	ShoppingShare             []UserSerializer `json:"shopping_share"`
	ShoppingRecentDays       int            `json:"shopping_recent_days"`
	CsvDelim                 string         `json:"csv_delim"`
	CsvPrefix                string         `json:"csv_prefix"`
	ShoppingUpdateFoodLists  bool           `json:"shopping_update_food_lists"`
	FilterToSupermarket      bool           `json:"filter_to_supermarket"`
	ShoppingAddOnhand        bool           `json:"shopping_add_onhand"`
	LeftHanded               bool           `json:"left_handed"`
	ShowStepIngredients      bool           `json:"show_step_ingredients"`
	FoodChildrenExist        bool           `json:"food_children_exist"` // TODO: implement logic
}

// SerializeUserPreference converts UserPreference model to serializer
func SerializeUserPreference(preference *models.UserPreference) UserPreferenceSerializer {
	return UserPreferenceSerializer{
		User:                     SerializeUser(&preference.User),
		Image:                    nil, // TODO: implement UserFileViewSerializer
		Theme:                    preference.Theme,
		NavBgColor:               preference.NavBgColor,
		NavTextColor:             preference.NavTextColor,
		NavShowLogo:              preference.NavShowLogo,
		DefaultUnit:              preference.DefaultUnit,
		DefaultPage:              preference.DefaultPage,
		UseFractions:             preference.UseFractions,
		UseKj:                   preference.UseKj,
		PlanShare:                []UserSerializer{}, // TODO: implement plan_share relationships
		NavSticky:                preference.NavSticky,
		IngredientDecimals:       preference.IngredientDecimals,
		Comments:                 preference.Comments,
		ShoppingAutoSync:         preference.ShoppingAutoSync,
		MealplanAutoaddShopping:  preference.MealplanAutoaddShopping,
		FoodInheritDefault:       []interface{}{}, // TODO: implement food inherit defaults
		DefaultDelay:             preference.DefaultDelay,
		MealplanAutoincludeRelated: preference.MealplanAutoincludeRelated,
		MealplanAutoexcludeOnhand: preference.MealplanAutoexcludeOnhand,
		ShoppingShare:            []UserSerializer{}, // TODO: implement shopping_share relationships
		ShoppingRecentDays:       preference.ShoppingRecentDays,
		CsvDelim:                 preference.CsvDelim,
		CsvPrefix:                preference.CsvPrefix,
		ShoppingUpdateFoodLists:  preference.ShoppingUpdateFoodLists,
		FilterToSupermarket:      preference.FilterToSupermarket,
		ShoppingAddOnhand:        preference.ShoppingAddOnhand,
		LeftHanded:               preference.LeftHanded,
		ShowStepIngredients:      preference.ShowStepIngredients,
		FoodChildrenExist:        false, // TODO: implement logic
	}
}
