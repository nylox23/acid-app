package dto

import "time"

// Запрос для обновления заявки
type CarbonateUpdateRequest struct {
	Mass float32 `json:"mass" binding:"min=0"`
}

// Возвращаемые данные по заявке
type CarbonateResponse struct {
	ID         uint       `json:"id"`
	Status     string     `json:"status"`
	DateCreate time.Time  `json:"date_create"`
	DateUpdate time.Time  `json:"date_update"`
	DateFinish *time.Time `json:"date_finish,omitempty"`
	Creator    string     `json:"creator"`
	Moderator  string     `json:"moderator,omitempty"`
	Mass       float32    `json:"mass"`
}

// Возвращаемые данные по заявке со списком кислот
type CarbonateDetailResponse struct {
	CarbonateResponse
	Acids []CarbonateAcidResponse `json:"acids"`
}

// Возвращаемые данные по кислоте в составе заявки
type CarbonateAcidResponse struct {
	ID          uint         `json:"id"`
	AcidID      uint         `json:"acid_id"`
	CarbonateID uint         `json:"carbonate_id"`
	Mass        float32      `json:"mass"`
	Result      *float32     `json:"result"`
	Acid        AcidResponse `json:"acid"`
}

type CarbonateListEntry struct {
	CarbonateResponse
	Calculated int64 `json:"calculated"`
}

// Возвращаемый список заявок
type CarbonateListResponse struct {
	Carbonates []CarbonateListEntry `json:"carbonates"`
}

// Параметры для фильтра заявок
type CarbonateFilter struct {
	Status   string    `form:"status"`
	DateFrom time.Time `form:"date_from"`
	DateTo   time.Time `form:"date_to"`
}

// Возвращаемые данные для кнопки перехода к заявке
type CarbonateIconsResponse struct {
	CarbonateID uint  `json:"carbonate_id"`
	AcidCount   int64 `json:"acid_count"`
}
