package handler

import (
	"net/http"
	"strconv"
	"time"
	"web_service_auth/internal/app/dto"
	"web_service_auth/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GetCurrentCarbonateAPI godoc
// @Summary      Получить данные текущей заявки
// @Description  Возвращает данные по текущей заявке-черновику пользователя для кнопки перехода
// @Tags         carbonates
// @Produce      json
// @Success      200  {object}  dto.CarbonateIconsResponse
// @Router       /api/carbonates/current [get]
func (h *Handler) GetCurrentCarbonateAPI(c *gin.Context) {
	userID := h.GetCurrentUserID(c)

	carbonateID, _ := h.Repository.GetDraftCarbonate(userID)
	acidCount := h.Repository.GetAcidCount(userID)

	response := dto.CarbonateIconsResponse{
		CarbonateID: carbonateID,
		AcidCount:   acidCount,
	}

	c.JSON(http.StatusOK, response)
}

// GetCarbonatesAPI godoc
// @Summary      Получить список заявок
// @Description  Возвращает список заявок, отфильтрованный по статусу и дате
// @Tags         carbonates
// @Produce      json
// @Param        status     query     string  false  "Фильтр по статусу"
// @Param        date_from  query     string  false  "Фильтр даты от"
// @Param        date_to    query     string  false  "Фильтр даты до"
// @Success      200        {object}  dto.CarbonateListResponse
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /api/carbonates [get]
func (h *Handler) GetCarbonatesAPI(c *gin.Context) {
	var filter dto.CarbonateFilter
	var sortID uuid.UUID

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.GetCurrentUserRole(c) == role.Admin {
		sortID = uuid.Nil
	} else {
		sortID = h.GetCurrentUserID(c)
	}

	carbonates, err := h.Repository.GetCarbonatesWithFilter(filter.Status, filter.DateFrom, filter.DateTo, sortID)
	if err != nil {
		logrus.Error("Failed to get carbonates:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonates"})
		return
	}

	response := dto.CarbonateListResponse{
		Carbonates: make([]dto.CarbonateListEntry, len(carbonates)),
	}

	for i, carbonate := range carbonates {
		response.Carbonates[i] = dto.CarbonateListEntry{
			CarbonateResponse: dto.CarbonateResponse{
				ID:         carbonate.ID,
				Status:     carbonate.Status,
				DateCreate: carbonate.DateCreate,
				DateUpdate: carbonate.DateUpdate,
				Creator:    carbonate.Creator.Login,
				Moderator:  carbonate.Moderator.Login,
				Mass:       carbonate.Mass,
			},
			Calculated: h.Repository.GetCalculated(carbonate.ID),
		}

		if carbonate.DateFinish.Valid {
			response.Carbonates[i].DateFinish = &carbonate.DateFinish.Time
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetCarbonateAPI godoc
// @Summary      Получить заявку по ID
// @Description  Возвращает заявку со списком добавленных в нее кислот
// @Tags         carbonates
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success      200  {object}  dto.CarbonateDetailResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/carbonates/{id} [get]
func (h *Handler) GetCarbonateAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid carbonate ID"})
		return
	}

	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	if (h.GetCurrentUserRole(c) != role.Admin) && (h.GetCurrentUserID(c) != carbonate.CreatorID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Получение списка кислот
	acids, err := h.Repository.GetCarbonateAcids(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate acids:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate acids"})
		return
	}

	response := dto.CarbonateDetailResponse{
		CarbonateResponse: dto.CarbonateResponse{
			ID:         carbonate.ID,
			Status:     carbonate.Status,
			DateCreate: carbonate.DateCreate,
			DateUpdate: carbonate.DateUpdate,
			Creator:    carbonate.Creator.Login,
			Moderator:  carbonate.Moderator.Login,
			Mass:       carbonate.Mass,
		},
		Acids: make([]dto.CarbonateAcidResponse, len(acids)),
	}

	if carbonate.DateFinish.Valid {
		response.DateFinish = &carbonate.DateFinish.Time
	}

	for i, acid := range acids {
		response.Acids[i] = dto.CarbonateAcidResponse{
			ID:          acid.ID,
			AcidID:      acid.AcidID,
			CarbonateID: acid.CarbonateID,
			Mass:        acid.Mass,
			Result:      acid.Result,
			Acid: dto.AcidResponse{
				ID:        acid.Acid.ID,
				NameExt:   acid.Acid.NameExt,
				Info:      acid.Acid.Info,
				Name:      acid.Acid.Name,
				Hplus:     acid.Acid.Hplus,
				MolarMass: acid.Acid.MolarMass,
				Img:       acid.Acid.Img,
			},
		}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCarbonateAPI godoc
// @Summary      Обновить заявку
// @Description  Обновляет поле "масса карбоната" в текущей заявке-черновика пользователя
// @Tags         carbonates
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CarbonateUpdateRequest  true  "Обновляемые данные заявки"
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/carbonates [put]
func (h *Handler) UpdateCarbonateAPI(c *gin.Context) {
	userID := h.GetCurrentUserID(c)

	id, _ := h.Repository.GetDraftCarbonate(userID)

	var req dto.CarbonateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Заявка существует?
	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	// Обновление полей
	updates := make(map[string]interface{})
	if req.Mass > 0 {
		updates["mass"] = req.Mass
	}
	updates["date_update"] = time.Now()

	if err := h.Repository.UpdateCarbonate(uint(id), updates); err != nil {
		logrus.Error("Failed to update carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate updated successfully"})
}

// FormCarbonateAPI godoc
// @Summary      Сформировать заявку
// @Description  Меняет статус текущей заявки-черновика пользователя на "сформирован"
// @Tags         carbonates
// @Produce      json
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/carbonates/form [put]
func (h *Handler) FormCarbonateAPI(c *gin.Context) {
	userID := h.GetCurrentUserID(c)

	carbonateID, _ := h.Repository.GetDraftCarbonate(userID)

	carbonate, err := h.Repository.GetCarbonateByID(uint(carbonateID))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	// Проверка обязательного поля
	if carbonate.Mass <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mass is required and must be greater than 0"})
		return
	}

	// Обновление статуса заявки
	updates := map[string]interface{}{
		"status":      "сформирован",
		"date_update": time.Now(),
	}

	if err := h.Repository.UpdateCarbonate(uint(carbonateID), updates); err != nil {
		logrus.Error("Failed to form carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to form carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate formed successfully"})
}

// SetCarbonateStatusAPI godoc
// @Summary      Установить статус заявки
// @Description  Устанавливает статус заявки на "завершен" или "отклонен"
// @Tags         carbonates
// @Accept       json
// @Produce      json
// @Param        id       path      int  true  "ID заявки"
// @Param        request  body      object  true  "Новый статус"  SchemaExample({"status": "завершен"})
// @Success      200      {object}  dto.MessageResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/carbonates/{id}/status [put]
func (h *Handler) SetCarbonateStatusAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid carbonate ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=завершен отклонен"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Заявка существует?
	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	// Завершение заявки
	if carbonate.Status != "сформирован" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only formed carbonates can be finalized or rejected"})
		return
	}

	userID := h.GetCurrentUserID(c)

	// Обновление статуса заявки
	updates := map[string]interface{}{
		"status":       req.Status,
		"moderator_id": userID,
		"date_finish":  time.Now(),
		"date_update":  time.Now(),
	}

	if err := h.Repository.UpdateCarbonate(uint(id), updates); err != nil {
		logrus.Error("Failed to update carbonate status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update carbonate status"})
		return
	}

	if req.Status == "завершен" {
		// Вычисление полей по теме в м-м
		if err := h.calculateCarbonateResults(uint(id)); err != nil {
			logrus.Error("Failed to calculate carbonate results:", err)
		}
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate status updated successfully"})
}

// DeleteCarbonateAPI godoc
// @Summary      Удалить заявку
// @Description  Удаляет заявку (создателем - только черновики или отклоненные заявки)
// @Tags         carbonates
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/carbonates/{id} [delete]
func (h *Handler) DeleteCarbonateAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid carbonate ID"})
		return
	}

	// Заявка существует?
	carbonate, err := h.Repository.GetCarbonateByID(uint(id))
	if err != nil {
		logrus.Error("Failed to get carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carbonate"})
		return
	}

	if carbonate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carbonate not found"})
		return
	}

	userID := h.GetCurrentUserID(c)

	// Только создатель может удалить заявку
	if (carbonate.CreatorID != userID) && (h.GetCurrentUserRole(c) != role.Admin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own carbonates"})
		return
	}

	// Заявка - черновик/отклонена?
	if carbonate.Status != "черновик" && carbonate.Status != "отклонен" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft and rejected carbonates can be deleted"})
		return
	}

	// Удаление заявки
	if err := h.Repository.DeleteCarbonate(uint(id)); err != nil {
		logrus.Error("Failed to delete carbonate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete carbonate"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Carbonate deleted successfully"})
}

// Вычисление результатов в м-м
func (h *Handler) calculateCarbonateResults(carbonateID uint) error {
	acids, err := h.Repository.GetCarbonateAcids(carbonateID)
	if err != nil {
		return err
	}

	for _, acid := range acids {
		result := 22.4 * min(acid.Carbonate.Mass/100, acid.Mass*float32(acid.Acid.Hplus)/2/acid.Acid.MolarMass)

		if err := h.Repository.UpdateCarbonateAcidResult(acid.ID, result); err != nil {
			return err
		}
	}

	return nil
}
