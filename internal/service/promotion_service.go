package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/repository"
)

var (
	posValues      = map[string]bool{"above": true, "below": true}
	alignValues    = map[string]bool{"left": true, "center": true, "right": true}
	styleValues    = map[string]bool{"bullets": true, "lines": true}
	densityValues  = map[string]bool{"compact": true, "comfortable": true, "spacious": true}
	borderValues   = map[string]bool{"none": true, "subtle": true, "strong": true}
	fontValues     = map[string]bool{"handwritten": true, "kawaii": true, "serif": true, "sans": true}
	countdownSizes = map[string]bool{"sm": true, "md": true, "lg": true}
	hexColorRegex  = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
)

// ValidatePromotion garante que o payload atende às regras
func ValidatePromotion(p *model.Promotion) error {
	// Global percentage
	if p.GlobalPercentage != nil {
		if *p.GlobalPercentage < 0 || *p.GlobalPercentage > 100 {
			return errors.New("globalPercentage deve estar entre 0 e 100")
		}
	}

	// Progressive rules
	for _, r := range p.ProgressiveRules {
		if r.Threshold <= 0 {
			return errors.New("progressiveRules.threshold deve ser > 0")
		}
		if r.Percentage < 0 || r.Percentage > 100 {
			return errors.New("progressiveRules.percentage deve estar entre 0 e 100")
		}
	}

	// Tempo: start <= end
	if p.StartAt != nil && p.EndAt != nil {
		if p.EndAt.Before(*p.StartAt) {
			return errors.New("endAt deve ser maior ou igual a startAt")
		}
	}

	// Enums
	checkEnum := func(val *string, allowed map[string]bool, field string) error {
		if val != nil && !allowed[*val] {
			return errors.New(field + " possui valor inválido")
		}
		return nil
	}

	if err := checkEnum(p.BannerConditionsPosition, posValues, "bannerConditionsPosition"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerCountdownPosition, posValues, "bannerCountdownPosition"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerAlignment, alignValues, "bannerAlignment"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerTitlePosition, posValues, "bannerTitlePosition"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerConditionsStyle, styleValues, "bannerConditionsStyle"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerDensity, densityValues, "bannerDensity"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerBorderStyle, borderValues, "bannerBorderStyle"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerTitleFont, fontValues, "bannerTitleFont"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerMessageFont, fontValues, "bannerMessageFont"); err != nil {
		return err
	}
	if err := checkEnum(p.BannerCountdownSize, countdownSizes, "bannerCountdownSize"); err != nil {
		return err
	}

	// Cores hex (aceita #RGB e #RRGGBB)
	checkHex := func(val *string, field string) error {
		if val != nil && !hexColorRegex.MatchString(*val) {
			return errors.New(field + " deve ser uma cor hex válida")
		}
		return nil
	}

	if err := checkHex(p.HighlightColor, "highlightColor"); err != nil {
		return err
	}
	if err := checkHex(p.BannerTitleColor, "bannerTitleColor"); err != nil {
		return err
	}
	if err := checkHex(p.BannerConditionsColor, "bannerConditionsColor"); err != nil {
		return err
	}
	if err := checkHex(p.BannerGlobalColor, "bannerGlobalColor"); err != nil {
		return err
	}
	if err := checkHex(p.BannerProgressiveColor, "bannerProgressiveColor"); err != nil {
		return err
	}
	if err := checkHex(p.BannerCountdownBgColor, "bannerCountdownBgColor"); err != nil {
		return err
	}
	if err := checkHex(p.BannerCountdownTextColor, "bannerCountdownTextColor"); err != nil {
		return err
	}

	return nil
}

// SavePromotion salva ou atualiza a promoção atual
func SavePromotion(p *model.Promotion) (*model.Promotion, error) {
	if err := ValidatePromotion(p); err != nil {
		return nil, err
	}

	existing, err := repository.GetLatestPromotion()
	if err != nil {
		return nil, err
	}

	if existing != nil {
		p.ID = existing.ID
		if err := repository.SavePromotion(p); err != nil {
			return nil, err
		}
	} else {
		if err := repository.CreatePromotion(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// GetActivePromotion retorna a promoção ativa (enabled e dentro da janela), senão nil
func GetActivePromotion() (*model.Promotion, error) {
	p, err := repository.GetLatestPromotion()
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	if !p.Enabled {
		return nil, nil
	}

	now := time.Now().UTC()
	if p.StartAt != nil && now.Before(*p.StartAt) {
		return nil, nil
	}
	if p.EndAt != nil && now.After(*p.EndAt) {
		return nil, nil
	}

	return p, nil
}
