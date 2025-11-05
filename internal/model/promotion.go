package model

import (
	"time"

	"gorm.io/gorm"
)

type ProgressiveRule struct {
	Threshold  int `json:"threshold"`
	Percentage int `json:"percentage"`
}

type Promotion struct {
	gorm.Model
	Enabled                  bool              `json:"enabled"`
	GlobalPercentage         *int              `json:"globalPercentage"`
	ProgressiveRules         []ProgressiveRule `json:"progressiveRules" gorm:"type:jsonb;serializer:json"`
	StartAt                  *time.Time        `json:"startAt"`
	EndAt                    *time.Time        `json:"endAt"`
	MessageTemplate          *string           `json:"messageTemplate"`
	HighlightColor           *string           `json:"highlightColor"`
	BannerShowConditions     *bool             `json:"bannerShowConditions"`
	BannerConditionsPosition *string           `json:"bannerConditionsPosition"`
	BannerShowCountdown      *bool             `json:"bannerShowCountdown"`
	BannerCountdownPosition  *string           `json:"bannerCountdownPosition"`
	BannerAlignment          *string           `json:"bannerAlignment"`
	BannerShowTitle          *bool             `json:"bannerShowTitle"`
	BannerTitle              *string           `json:"bannerTitle"`
	BannerTitlePosition      *string           `json:"bannerTitlePosition"`
	BannerConditionsStyle    *string           `json:"bannerConditionsStyle"`
	BannerDensity            *string           `json:"bannerDensity"`
	BannerBorderStyle        *string           `json:"bannerBorderStyle"`
	BannerTitleFont          *string           `json:"bannerTitleFont"`
	BannerMessageFont        *string           `json:"bannerMessageFont"`
	BannerTitleColor         *string           `json:"bannerTitleColor"`
	BannerConditionsColor    *string           `json:"bannerConditionsColor"`
	BannerGlobalColor        *string           `json:"bannerGlobalColor"`
	BannerProgressiveColor   *string           `json:"bannerProgressiveColor"`
	BannerCountdownBgColor   *string           `json:"bannerCountdownBgColor"`
	BannerCountdownTextColor *string           `json:"bannerCountdownTextColor"`
	BannerCountdownSize      *string           `json:"bannerCountdownSize"`
}
