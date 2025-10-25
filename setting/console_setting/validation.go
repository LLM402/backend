package console_setting

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	urlRegex       = regexp.MustCompile(`^https?://(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?|(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))(?:\:[0-9]{1,5})?(?:/.*)?$`)
	dangerousChars = []string{"<script", "<iframe", "javascript:", "onload=", "onerror=", "onclick="}
	validColors    = map[string]bool{
		"blue": true, "green": true, "cyan": true, "purple": true, "pink": true,
		"red": true, "orange": true, "amber": true, "yellow": true, "lime": true,
		"light-green": true, "teal": true, "light-blue": true, "indigo": true,
		"violet": true, "grey": true,
	}
	slugRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func parseJSONArray(jsonStr string, typeName string) ([]map[string]interface{}, error) {
	var list []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &list); err != nil {
		return nil, fmt.Errorf("%s format error: %s", typeName, err.Error())
	}
	return list, nil
}

func validateURL(urlStr string, index int, itemType string) error {
	if !urlRegex.MatchString(urlStr) {
		return fmt.Errorf("The URL format of the %dth %s is incorrect.", index, itemType)
	}
	if _, err := url.Parse(urlStr); err != nil {
		return fmt.Errorf("The URL of the %dth %s cannot be resolved: %s", index, itemType, err.Error())
	}
	return nil
}

func checkDangerousContent(content string, index int, itemType string) error {
	lower := strings.ToLower(content)
	for _, d := range dangerousChars {
		if strings.Contains(lower, d) {
			return fmt.Errorf("The %dth %s contains disallowed content.", index, itemType)
		}
	}
	return nil
}

func getJSONList(jsonStr string) []map[string]interface{} {
	if jsonStr == "" {
		return []map[string]interface{}{}
	}
	var list []map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &list)
	return list
}

func ValidateConsoleSettings(settingsStr string, settingType string) error {
	if settingsStr == "" {
		return nil
	}

	switch settingType {
	case "ApiInfo":
		return validateApiInfo(settingsStr)
	case "Announcements":
		return validateAnnouncements(settingsStr)
	case "FAQ":
		return validateFAQ(settingsStr)
	case "UptimeKumaGroups":
		return validateUptimeKumaGroups(settingsStr)
	default:
		return fmt.Errorf("Unknown setting type: %s", settingType)
	}
}

func validateApiInfo(apiInfoStr string) error {
	apiInfoList, err := parseJSONArray(apiInfoStr, "API Information")
	if err != nil {
		return err
	}

	if len(apiInfoList) > 50 {
		return fmt.Errorf("The number of API information cannot exceed 50.")
	}

	for i, apiInfo := range apiInfoList {
		urlStr, ok := apiInfo["url"].(string)
		if !ok || urlStr == "" {
			return fmt.Errorf("The URL field is missing for API information #%d.", i+1)
		}
		route, ok := apiInfo["route"].(string)
		if !ok || route == "" {
			return fmt.Errorf("The %dth API information is missing the line description field.", i+1)
		}
		description, ok := apiInfo["description"].(string)
		if !ok || description == "" {
			return fmt.Errorf("The description field is missing for API #%d.", i+1)
		}
		color, ok := apiInfo["color"].(string)
		if !ok || color == "" {
			return fmt.Errorf("The color field is missing for API #%d.", i+1)
		}

		if err := validateURL(urlStr, i+1, "API Information"); err != nil {
			return err
		}

		if len(urlStr) > 500 {
			return fmt.Errorf("The URL length of the %dth API information cannot exceed 500 characters.", i+1)
		}
		if len(route) > 100 {
			return fmt.Errorf("The line description length of API information %d cannot exceed 100 characters.", i+1)
		}
		if len(description) > 200 {
			return fmt.Errorf("The description length of the %dth API information cannot exceed 200 characters.", i+1)
		}

		if !validColors[color] {
			return fmt.Errorf("The color value of the %dth API information is invalid.", i+1)
		}

		if err := checkDangerousContent(description, i+1, "API Information"); err != nil {
			return err
		}
		if err := checkDangerousContent(route, i+1, "API Information"); err != nil {
			return err
		}
	}
	return nil
}

func GetApiInfo() []map[string]interface{} {
	return getJSONList(GetConsoleSetting().ApiInfo)
}

func validateAnnouncements(announcementsStr string) error {
	list, err := parseJSONArray(announcementsStr, "System Announcement")
	if err != nil {
		return err
	}
	if len(list) > 100 {
		return fmt.Errorf("The number of system announcements cannot exceed 100.")
	}
	validTypes := map[string]bool{
		"default": true, "ongoing": true, "success": true, "warning": true, "error": true,
	}
	for i, ann := range list {
		content, ok := ann["content"].(string)
		if !ok || content == "" {
			return fmt.Errorf("The %dth announcement is missing the content field.", i+1)
		}
		publishDateAny, exists := ann["publishDate"]
		if !exists {
			return fmt.Errorf("The %dth announcement is missing the publication date field.", i+1)
		}
		publishDateStr, ok := publishDateAny.(string)
		if !ok || publishDateStr == "" {
			return fmt.Errorf("The publication date of announcement %d cannot be empty.", i+1)
		}
		if _, err := time.Parse(time.RFC3339, publishDateStr); err != nil {
			return fmt.Errorf("The publication date format of announcement #%d is incorrect.", i+1)
		}
		if t, exists := ann["type"]; exists {
			if typeStr, ok := t.(string); ok {
				if !validTypes[typeStr] {
					return fmt.Errorf("The type value of the %dth announcement is invalid.", i+1)
				}
			}
		}
		if len(content) > 500 {
			return fmt.Errorf("The content length of announcement %d cannot exceed 500 characters.", i+1)
		}
		if extra, exists := ann["extra"]; exists {
			if extraStr, ok := extra.(string); ok && len(extraStr) > 200 {
				return fmt.Errorf("The description length of announcement %d cannot exceed 200 characters.", i+1)
			}
		}
	}
	return nil
}

func validateFAQ(faqStr string) error {
	list, err := parseJSONArray(faqStr, "FAQ Information")
	if err != nil {
		return err
	}
	if len(list) > 100 {
		return fmt.Errorf("The number of FAQs cannot exceed 100.")
	}
	for i, faq := range list {
		question, ok := faq["question"].(string)
		if !ok || question == "" {
			return fmt.Errorf("The %dth FAQ is missing the question field.", i+1)
		}
		answer, ok := faq["answer"].(string)
		if !ok || answer == "" {
			return fmt.Errorf("The answer field is missing for FAQ #%d.", i+1)
		}
		if len(question) > 200 {
			return fmt.Errorf("The length of the %dth FAQ question cannot exceed 200 characters.", i+1)
		}
		if len(answer) > 1000 {
			return fmt.Errorf("The answer length for FAQ #%d cannot exceed 1000 characters.", i+1)
		}
	}
	return nil
}

func getPublishTime(item map[string]interface{}) time.Time {
	if v, ok := item["publishDate"]; ok {
		if s, ok2 := v.(string); ok2 {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t
			}
		}
	}
	return time.Time{}
}

func GetAnnouncements() []map[string]interface{} {
	list := getJSONList(GetConsoleSetting().Announcements)
	sort.SliceStable(list, func(i, j int) bool {
		return getPublishTime(list[i]).After(getPublishTime(list[j]))
	})
	return list
}

func GetFAQ() []map[string]interface{} {
	return getJSONList(GetConsoleSetting().FAQ)
}

func validateUptimeKumaGroups(groupsStr string) error {
	groups, err := parseJSONArray(groupsStr, "Uptime Kuma group configuration")
	if err != nil {
		return err
	}

	if len(groups) > 20 {
		return fmt.Errorf("The number of Uptime Kuma groups cannot exceed 20.")
	}

	nameSet := make(map[string]bool)

	for i, group := range groups {
		categoryName, ok := group["categoryName"].(string)
		if !ok || categoryName == "" {
			return fmt.Errorf("The %dth group is missing the category name field.", i+1)
		}
		if nameSet[categoryName] {
			return fmt.Errorf("The category name of group %d duplicates with other groups.", i+1)
		}
		nameSet[categoryName] = true
		urlStr, ok := group["url"].(string)
		if !ok || urlStr == "" {
			return fmt.Errorf("The %dth group is missing the URL field.", i+1)
		}
		slug, ok := group["slug"].(string)
		if !ok || slug == "" {
			return fmt.Errorf("The %dth group is missing the Slug field.", i+1)
		}
		description, ok := group["description"].(string)
		if !ok {
			description = ""
		}

		if err := validateURL(urlStr, i+1, "Group"); err != nil {
			return err
		}

		if len(categoryName) > 50 {
			return fmt.Errorf("The classification name length of group %d cannot exceed 50 characters.", i+1)
		}
		if len(urlStr) > 500 {
			return fmt.Errorf("The URL length of group %d cannot exceed 500 characters.", i+1)
		}
		if len(slug) > 100 {
			return fmt.Errorf("The Slug length of group %d cannot exceed 100 characters.", i+1)
		}
		if len(description) > 200 {
			return fmt.Errorf("The description length of group %d cannot exceed 200 characters.", i+1)
		}

		if !slugRegex.MatchString(slug) {
			return fmt.Errorf("The slug for group %d can only contain letters, numbers, underscores, and hyphens.", i+1)
		}

		if err := checkDangerousContent(description, i+1, "Group"); err != nil {
			return err
		}
		if err := checkDangerousContent(categoryName, i+1, "Group"); err != nil {
			return err
		}
	}
	return nil
}

func GetUptimeKumaGroups() []map[string]interface{} {
	return getJSONList(GetConsoleSetting().UptimeKumaGroups)
}
