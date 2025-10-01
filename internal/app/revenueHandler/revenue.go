package revenueHandler

import (
	"lab2/internal/app/ds"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *RevenueHandler) RegisterRevenueHandler(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/periods", h.GetPeriods)
		api.GET("/periods/:id", h.GetPeriod)
		api.POST("/periods", h.CreatePeriod)
		api.PUT("/periods", h.UpdatePeriod)
		api.POST("/periods/periods-application", h.AddPeriodToApplication)

		api.GET("/periods-cart-info", h.GetPeriodsCartInfo)
		api.GET("/periods-applications", h.GetPeriodsApplications)
		api.PUT("/periods-applications/:id", h.UpdatePeriodsApplication)
		api.PUT("/periods-applications/:id/form", h.FormPeriodsApplication)
		api.PUT("/periods-applications/:id/complete", h.CompletePeriodsApplication)
		api.DELETE("/periods-applications/:id", h.DeletePeriodsApplication)

		api.DELETE("/selected-periods/:application_id/:period_id", h.DeletePeriodFromApplication)
		api.PUT("/selected-periods/:application_id/:period_id/previous-revenue", h.UpdatePeriodPreviousRevenue)

		api.POST("/auth/signup", h.RegisterUser)
		api.GET("/user/checkauth/:user_id", h.CheckUserAuth)
		api.PUT("/user/profile/:user_id", h.UpdateUserProfile)
		api.POST("/auth/login", h.LoginUser)
		api.POST("/auth/logout", h.LogoutUser)
	}
}

func (h *RevenueHandler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *RevenueHandler) errorRevenueHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *RevenueHandler) GetPeriods(ctx *gin.Context) {
	title := ctx.Query("title")

	var periods []ds.Period
	var err error

	if title != "" {
		periods, err = h.RevenueModel.GetPeriodsByTitle(title)
	} else {
		periods, err = h.RevenueModel.GetPeriods()
	}

	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"periods": periods,
	})
}

func (h *RevenueHandler) UpdatePeriodsApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid id",
		})
		return
	}

	var request struct {
		CompanyName string `json:"company_name"`
		Year        int    `json:"year"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	updateData := make(map[string]interface{})
	if request.CompanyName != "" {
		updateData["company_name"] = request.CompanyName
	}
	if request.Year > 0 {
		updateData["year"] = request.Year
	}

	if len(updateData) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "no fields to update",
		})
		return
	}

	updatedApplication, err := h.RevenueModel.UpdatePeriodsApplication(id, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error updating application",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "application updated successfully",
		"data":    updatedApplication,
	})
}

func (h *RevenueHandler) GetPeriod(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusBadRequest, err)
		return
	}

	period, err := h.RevenueModel.GetPeriod(id)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"period": period,
	})
}

func (h *RevenueHandler) CreatePeriod(ctx *gin.Context) {
	var request struct {
		Title               string `json:"title" binding:"required"`
		Description         string `json:"description" binding:"required"`
		Duration            string `json:"duration" binding:"required"`
		ShortDescription    string `json:"short_description" binding:"required"`
		DetailedDescription string `json:"detailed_description" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	period := &ds.Period{
		Title:               request.Title,
		Description:         request.Description,
		Duration:            request.Duration,
		ShortDescription:    request.ShortDescription,
		DetailedDescription: request.DetailedDescription,
		IsActive:            true,
	}

	if err := h.RevenueModel.CreatePeriod(period); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error creating period",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "period created successfully",
		"data":    period,
	})
}

func (h *RevenueHandler) UpdatePeriod(ctx *gin.Context) {
	var request struct {
		ID                  int    `json:"id" binding:"required"`
		Title               string `json:"title"`
		Description         string `json:"description"`
		Duration            string `json:"duration"`
		ShortDescription    string `json:"short_description"`
		DetailedDescription string `json:"detailed_description"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	updateData := make(map[string]interface{})
	if request.Title != "" {
		updateData["title"] = request.Title
	}
	if request.Description != "" {
		updateData["description"] = request.Description
	}
	if request.Duration != "" {
		updateData["duration"] = request.Duration
	}
	if request.ShortDescription != "" {
		updateData["short_description"] = request.ShortDescription
	}
	if request.DetailedDescription != "" {
		updateData["detailed_description"] = request.DetailedDescription
	}

	if len(updateData) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "no fields to update",
		})
		return
	}

	if err := h.RevenueModel.UpdatePeriod(request.ID, updateData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error updating period",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "period updated successfully",
	})
}

func (h *RevenueHandler) CreateEmptyDraftPeriodsApplication(ctx *gin.Context) {
	userID := 1

	periodsApplication, err := h.RevenueModel.CreateEmptyDraftPeriodsApplication(userID)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "empty draft periods application created successfully",
		"data":    periodsApplication,
	})
}

func (h *RevenueHandler) GetPeriodsCartInfo(ctx *gin.Context) {
	userID := 1
	periodsApplication, err := h.RevenueModel.GetOrCreateDraftApplication(userID)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	count, err := h.RevenueModel.GetSelectedPeriodsCount(periodsApplication.ID)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"application_id": periodsApplication.ID,
			"items_count":    count,
		},
	})
}

func (h *RevenueHandler) GetPeriodsApplications(ctx *gin.Context) {
	status := ctx.Query("status")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	applications, err := h.RevenueModel.GetPeriodsApplications(status, startDate, endDate)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"applications": applications,
	})
}

func (h *RevenueHandler) GetOrCreateDraftApplication(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logrus.Error("error", err)
	}

	periodsApplication, err := h.RevenueModel.GetOrCreateDraftApplication(userID)
	if err != nil {
		logrus.Error("error", err)
	}

	ctx.HTML(http.StatusOK, "periodsApplication.html", gin.H{
		"periodsApplication": periodsApplication,
	})
}

func (h *RevenueHandler) AddPeriodToDraftApplication(ctx *gin.Context) {
	var request struct {
		PeriodID int `json:"period_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	userID := 1

	periodsApplication, err := h.RevenueModel.GetOrCreateDraftApplication(userID)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.RevenueModel.AddPeriodToApplication(periodsApplication.ID, request.PeriodID)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "period added to draft application successfully",
		"data": gin.H{
			"application_id": periodsApplication.ID,
			"period_id":      request.PeriodID,
		},
	})
}

func (h *RevenueHandler) DeleteApplication(ctx *gin.Context) {
	periodsApplicationIDStr := ctx.Param("id")
	periodsApplicationID, err := strconv.Atoi(periodsApplicationIDStr)
	if err != nil {
		logrus.Error("error", err)
	}

	err = h.RevenueModel.DeleteApplication(periodsApplicationID)
	if err != nil {
		logrus.Error("error", err)
	}

	ctx.Redirect(http.StatusFound, "/home")
}

func (h *RevenueHandler) FormPeriodsApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid id",
		})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	canForm, missingFields, err := h.RevenueModel.CanFormPeriodsApplication(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error checking application",
		})
		return
	}

	if !canForm {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":         "fail",
			"message":        "cannot form application - missing required fields",
			"missing_fields": missingFields,
		})
		return
	}

	updateData := map[string]interface{}{
		"status":    request.Status,
		"formed_at": time.Now(),
	}

	updatedApplication, err := h.RevenueModel.FormPeriodsApplication(id, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error forming application",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "application formed successfully",
		"data":    updatedApplication,
	})
}

func (h *RevenueHandler) CompletePeriodsApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid id",
		})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	if request.Status != "completed" && request.Status != "rejected" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "status must be 'completed' or 'rejected'",
		})
		return
	}

	moderatorID := 1
	result, err := h.RevenueModel.CompletePeriodsApplication(id, request.Status, moderatorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error completing application",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "application " + request.Status + " successfully",
		"data":    result,
	})
}

func (h *RevenueHandler) DeletePeriodsApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid id",
		})
		return
	}

	err = h.RevenueModel.DeletePeriodsApplication(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error deleting application",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "application deleted successfully",
	})
}

func (h *RevenueHandler) DeletePeriodFromApplication(ctx *gin.Context) {
	applicationIDStr := ctx.Param("application_id")
	periodIDStr := ctx.Param("period_id")

	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil || applicationID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid application id",
		})
		return
	}

	periodID, err := strconv.Atoi(periodIDStr)
	if err != nil || periodID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid period id",
		})
		return
	}

	err = h.RevenueModel.DeletePeriodFromApplication(applicationID, periodID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error deleting period from application",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "period deleted from application successfully",
	})
}

func (h *RevenueHandler) UpdatePeriodPreviousRevenue(ctx *gin.Context) {
	applicationIDStr := ctx.Param("application_id")
	periodIDStr := ctx.Param("period_id")

	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil || applicationID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid application id",
		})
		return
	}

	periodID, err := strconv.Atoi(periodIDStr)
	if err != nil || periodID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid period id",
		})
		return
	}

	var request struct {
		PreviousRevenue string `json:"previous_revenue" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	updatedSelectedPeriod, err := h.RevenueModel.UpdatePeriodPreviousRevenue(applicationID, periodID, request.PreviousRevenue)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error updating previous revenue",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "previous revenue updated successfully",
		"data":    updatedSelectedPeriod,
	})
}

func (h *RevenueHandler) RegisterUser(ctx *gin.Context) {
	var request struct {
		Login       string `json:"login" binding:"required"`
		Password    string `json:"password" binding:"required"`
		IsModerator bool   `json:"is_moderator"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	// Создаем пользователя
	user, err := h.RevenueModel.CreateUser(request.Login, request.Password, request.IsModerator)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error creating user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user registered successfully",
		"data": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
			"created_at":   user.CreatedAt,
		},
	})
}

func (h *RevenueHandler) CheckUserAuth(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")

	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "user ID is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid user ID",
		})
		return
	}

	user, err := h.RevenueModel.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user is authenticated",
		"data": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
			"created_at":   user.CreatedAt,
			"updated_at":   user.UpdatedAt,
		},
	})
}

func (h *RevenueHandler) UpdateUserProfile(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")

	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "user ID is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid user ID",
		})
		return
	}

	var request struct {
		Login       string `json:"login"`
		IsModerator *bool  `json:"is_moderator"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	updateData := make(map[string]interface{})
	if request.Login != "" {
		updateData["login"] = request.Login
	}
	if request.IsModerator != nil {
		updateData["is_moderator"] = *request.IsModerator
	}

	if len(updateData) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "no fields to update",
		})
		return
	}

	updatedUser, err := h.RevenueModel.UpdateUser(userID, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "error updating user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user profile updated successfully",
		"data": gin.H{
			"id":           updatedUser.ID,
			"login":        updatedUser.Login,
			"is_moderator": updatedUser.IsModerator,
			"created_at":   updatedUser.CreatedAt,
			"updated_at":   updatedUser.UpdatedAt,
		},
	})
}

func (h *RevenueHandler) LoginUser(ctx *gin.Context) {
	var request struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	user, err := h.RevenueModel.AuthenticateUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid credentials",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login successful",
		"data": gin.H{
			"user_id":      user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
		},
	})
}

func (h *RevenueHandler) LogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logout successful",
	})
}

func (h *RevenueHandler) AddPeriodToApplication(ctx *gin.Context) {
	var request struct {
		ApplicationID int `json:"application_id" binding:"required"`
		PeriodID      int `json:"period_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid input data",
		})
		return
	}

	userID := 1
	// Получаем или создаем заявку
	periodsApplication, err := h.RevenueModel.GetOrCreateApplication(request.ApplicationID, userID)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Добавляем период в заявку
	err = h.RevenueModel.AddPeriodToApplication(periodsApplication.ID, request.PeriodID)
	if err != nil {
		h.errorRevenueHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "period added to application successfully",
		"data": gin.H{
			"application_id": periodsApplication.ID,
			"period_id":      request.PeriodID,
			"status":         periodsApplication.Status,
		},
	})
}
