package revenueModel

import (
	"fmt"
	"lab2/internal/app/ds"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func calculateEMA(values []float64) float64 {
	n := float64(len(values))
	if n == 0 {
		return 0
	}

	alpha := 2 / (n + 1)
	ema := values[0]
	for i := 1; i < len(values); i++ {
		ema = alpha*values[i] + (1-alpha)*ema
	}

	pow := math.Pow(10, float64(3))
	return math.Round(ema*pow) / pow
}

func (r *RevenueModel) CreatePeriod(period *ds.Period) error {
	period.CreatedAt = time.Now()
	period.UpdatedAt = time.Now()
	return r.db.Create(period).Error
}

func (r *RevenueModel) UpdatePeriod(id int, updateData map[string]interface{}) error {
	updateData["updated_at"] = time.Now()
	result := r.db.Model(&ds.Period{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("period not found")
	}
	return nil
}

func (r *RevenueModel) UpdatePeriodImage(id int, imagePath string) error {
	return r.db.Model(&ds.Period{}).Where("id = ?", id).Updates(map[string]interface{}{
		"img":        imagePath,
		"updated_at": time.Now(),
	}).Error
}

func (r *RevenueModel) DeletePeriod(id int) error {
	result := r.db.Where("id = ?", id).Delete(&ds.Period{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("period not found")
	}
	return nil
}

func (r *RevenueModel) GetPeriodsApplications(status, startDate, endDate string) ([]map[string]interface{}, error) {
	var applications []ds.PeriodsApplication

	// Базовый запрос
	query := r.db.Where("status != ? AND status != ?", "deleted", "draft")

	// Фильтрация по статусу
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Фильтрация по диапазону дат
	if startDate != "" {
		query = query.Where("formed_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("formed_at <= ?", endDate)
	}

	err := query.Find(&applications).Error
	if err != nil {
		return nil, err
	}

	var creatorIDs []int
	var moderatorIDs []int

	for _, app := range applications {
		creatorIDs = append(creatorIDs, app.CreatorID)
		if app.ModeratorID != nil {
			moderatorIDs = append(moderatorIDs, *app.ModeratorID)
		}
	}

	var creators []ds.User
	if len(creatorIDs) > 0 {
		r.db.Where("id IN ?", creatorIDs).Find(&creators)
	}

	var moderators []ds.User
	if len(moderatorIDs) > 0 {
		r.db.Where("id IN ?", moderatorIDs).Find(&moderators)
	}

	creatorLogins := make(map[int]string)
	moderatorLogins := make(map[int]string)

	for _, creator := range creators {
		creatorLogins[creator.ID] = creator.Login
	}

	for _, moderator := range moderators {
		moderatorLogins[moderator.ID] = moderator.Login
	}

	var result []map[string]interface{}
	for _, app := range applications {
		item := map[string]interface{}{
			"id":            app.ID,
			"year":          app.Year,
			"full_duration": app.FullDuration,
			"status":        app.Status,
			"creator_login": creatorLogins[app.CreatorID],
			"company_name":  app.CompanyName,
			"created_at":    app.CreatedAt,
			"formed_at":     app.FormedAt,
			"completed_at":  app.CompletedAt,
		}

		if app.ModeratorID != nil {
			item["moderator_login"] = moderatorLogins[*app.ModeratorID]
		} else {
			item["moderator_login"] = nil
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *RevenueModel) GetPeriods() ([]ds.Period, error) {
	var periods []ds.Period
	err := r.db.Find(&periods).Error
	if err != nil {
		return nil, err
	}
	if len(periods) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return periods, nil
}

func (r *RevenueModel) GetPeriod(id int) (ds.Period, error) {
	period := ds.Period{}
	err := r.db.Where("id = ?", id).First(&period).Error
	if err != nil {
		return ds.Period{}, err
	}
	return period, nil
}

func (r *RevenueModel) GetPeriodsByTitle(title string) ([]ds.Period, error) {
	var periods []ds.Period
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&periods).Error
	if err != nil {
		return nil, err
	}
	return periods, nil
}

func (r *RevenueModel) GetSelectedPeriodsCount(periodsApplicationID int) (int, error) {
	var count int64
	if err := r.db.Table("selected_periods").
		Where("application_id = ?", periodsApplicationID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("ошибка подсчёта выбранных периодов:", err)
	}
	return int(count), nil
}

func (r *RevenueModel) CreateEmptyDraftPeriodsApplication(userID int) (*ds.PeriodsApplication, error) {
	newApplication := &ds.PeriodsApplication{
		Year:        0,
		CreatorID:   userID,
		Status:      "draft",
		CompanyName: "",
		CreatedAt:   time.Now(),
		FormedAt:    time.Now(),
		ModeratorID: nil,
		CompletedAt: nil,
	}

	err := r.db.Create(newApplication).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пустого черновика заявки периодов: %w", err)
	}

	return newApplication, nil
}

func (r *RevenueModel) UpdatePeriodsApplication(id int, updateData map[string]interface{}) (*ds.PeriodsApplication, error) {
	result := r.db.Model(&ds.PeriodsApplication{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("periods application not found")
	}

	var updatedApplication ds.PeriodsApplication
	if err := r.db.Where("id = ?", id).First(&updatedApplication).Error; err != nil {
		return nil, err
	}

	return &updatedApplication, nil
}

func (r *RevenueModel) CompletePeriodsApplication(id int, status string, moderatorID int) (map[string]interface{}, error) {
	// Получаем заявку
	var application ds.PeriodsApplication
	if err := r.db.Where("id = ?", id).First(&application).Error; err != nil {
		return nil, err
	}

	updateData := map[string]interface{}{
		"status":       status,
		"moderator_id": moderatorID,
		"completed_at": time.Now(),
	}

	var selectedPeriodsWithForecast []map[string]interface{}

	// Если заявка завершена, рассчитываем ForecastedRevenue для каждого периода
	if status == "completed" {
		periodsWithForecast, err := r.CalculateForecastedRevenueForApplication(id)
		if err != nil {
			return nil, err
		}
		selectedPeriodsWithForecast = periodsWithForecast
	} else {
		// Для отклоненной заявки просто получаем периоды без перерасчета
		periods, err := r.GetSelectedPeriodsWithDetails(id)
		if err != nil {
			return nil, err
		}
		selectedPeriodsWithForecast = periods
	}

	// Обновляем заявку
	result := r.db.Model(&ds.PeriodsApplication{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return nil, result.Error
	}

	// Получаем обновленную заявку
	var updatedApplication ds.PeriodsApplication
	if err := r.db.Where("id = ?", id).First(&updatedApplication).Error; err != nil {
		return nil, err
	}

	// Формируем результат с заявкой и периодами
	response := map[string]interface{}{
		"application":      updatedApplication,
		"selected_periods": selectedPeriodsWithForecast,
	}

	return response, nil
}

func (r *RevenueModel) CalculateForecastedRevenueForApplication(applicationID int) ([]map[string]interface{}, error) {
	var selectedPeriods []ds.SelectedPeriod
	if err := r.db.Where("application_id = ?", applicationID).Find(&selectedPeriods).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	// Для каждого периода рассчитываем ForecastedRevenue через EMA
	for _, sp := range selectedPeriods {
		// Парсим PreviousRevenue и вычисляем EMA
		parts := strings.Split(sp.PreviousRevenue, ";")
		var values []float64
		for _, p := range parts {
			v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
			if err == nil {
				values = append(values, v)
			}
		}

		// Используем EMA для прогноза
		forecast := calculateEMA(values)

		// Обновляем ForecastedRevenue в базе
		err := r.db.Model(&ds.SelectedPeriod{}).
			Where("application_id = ? AND period_id = ?", sp.ApplicationID, sp.PeriodID).
			Update("forecasted_revenue", forecast).Error
		if err != nil {
			return nil, err
		}

		// Получаем информацию о периоде
		var period ds.Period
		r.db.Where("id = ?", sp.PeriodID).First(&period)

		// Добавляем в результат
		result = append(result, map[string]interface{}{
			"period_id":          sp.PeriodID,
			"application_id":     sp.ApplicationID,
			"title":              period.Title,
			"description":        period.Description,
			"duration":           period.Duration,
			"previous_revenue":   sp.PreviousRevenue,
			"forecasted_revenue": forecast,
			"created_at":         sp.CreatedAt,
			"updated_at":         time.Now(),
		})
	}

	return result, nil
}

func (r *RevenueModel) GetSelectedPeriodsWithDetails(applicationID int) ([]map[string]interface{}, error) {
	var selectedPeriods []ds.SelectedPeriod
	if err := r.db.Where("application_id = ?", applicationID).Find(&selectedPeriods).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for _, sp := range selectedPeriods {
		var period ds.Period
		r.db.Where("id = ?", sp.PeriodID).First(&period)

		result = append(result, map[string]interface{}{
			"period_id":          sp.PeriodID,
			"application_id":     sp.ApplicationID,
			"title":              period.Title,
			"description":        period.Description,
			"duration":           period.Duration,
			"previous_revenue":   sp.PreviousRevenue,
			"forecasted_revenue": sp.ForecastedRevenue,
			"created_at":         sp.CreatedAt,
			"updated_at":         sp.UpdatedAt,
		})
	}

	return result, nil
}

func (r *RevenueModel) GetOrCreateDraftApplication(userID int) (*ds.PeriodsApplication, error) {
	var draftApplication ds.PeriodsApplication
	err := r.db.Where("creator_id = ? AND status = ?", userID, "draft").First(&draftApplication).Error
	if err == nil {
		return &draftApplication, nil
	}

	newApplication := &ds.PeriodsApplication{
		Year:      time.Now().Year(),
		CreatorID: userID,
		Status:    "draft",
		CreatedAt: time.Now(),
		FormedAt:  time.Now(),
	}

	err = r.db.Create(newApplication).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании черновика: %w", err)
	}

	return newApplication, nil
}

func generatePreviousRevenueFromPeriod(period string) string {
	parts := strings.Fields(period)
	if len(parts) == 0 {
		return ""
	}

	n, err := strconv.Atoi(parts[0])
	if err != nil || n <= 0 {
		n = 4
	}

	var values []string
	for i := 0; i < n; i++ {
		val := rand.Intn(6000) + 15000
		values = append(values, strconv.Itoa(val))
	}

	return strings.Join(values, "; ")
}

func (r *RevenueModel) AddPeriodToApplication(periodsApplicationID, periodID int) error {
	var existingSelectedPeriod ds.SelectedPeriod
	err := r.db.Where("application_id = ? AND period_id = ?", periodsApplicationID, periodID).
		First(&existingSelectedPeriod).Error

	if err == nil {
		return fmt.Errorf("период уже добавлен в заявку")
	}

	var period ds.Period
	if err := r.db.Where("id = ?", periodID).First(&period).Error; err != nil {
		return fmt.Errorf("не удалось получить период: %w", err)
	}
	previousRevenue := generatePreviousRevenueFromPeriod(period.Title)

	selectedPeriod := &ds.SelectedPeriod{
		ApplicationID:   periodsApplicationID,
		PeriodID:        periodID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		PreviousRevenue: previousRevenue,
	}

	return r.db.Create(selectedPeriod).Error
}

func (r *RevenueModel) DeleteApplication(periodsApplicationID int) error {
	result := r.db.Exec("UPDATE periods_applications SET status = 'deleted' WHERE id = ? AND status != 'deleted'", periodsApplicationID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("заявка не найдена или уже удалена")
	}

	return nil
}

func (r *RevenueModel) CanFormPeriodsApplication(id int) (bool, []string, error) {
	var application ds.PeriodsApplication
	if err := r.db.Where("id = ?", id).First(&application).Error; err != nil {
		return false, nil, err
	}

	var missingFields []string

	if application.CompanyName == "" {
		missingFields = append(missingFields, "company_name")
	}
	if application.Year == 0 {
		missingFields = append(missingFields, "year")
	}

	count, err := r.GetSelectedPeriodsCount(id)
	if err != nil {
		return false, nil, err
	}
	if count == 0 {
		missingFields = append(missingFields, "periods")
	}

	return len(missingFields) == 0, missingFields, nil
}

func (r *RevenueModel) DeletePeriodsApplication(id int) error {
	result := r.db.Model(&ds.PeriodsApplication{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    "deleted",
			"formed_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("periods application not found")
	}
	return nil
}

func (r *RevenueModel) DeletePeriodFromApplication(applicationID, periodID int) error {
	result := r.db.Where("application_id = ? AND period_id = ?", applicationID, periodID).
		Delete(&ds.SelectedPeriod{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("period not found in application")
	}
	return nil
}

func (r *RevenueModel) UpdatePeriodPreviousRevenue(applicationID, periodID int, previousRevenue string) (*ds.SelectedPeriod, error) {
	result := r.db.Model(&ds.SelectedPeriod{}).
		Where("application_id = ? AND period_id = ?", applicationID, periodID).
		Updates(map[string]interface{}{
			"previous_revenue": previousRevenue,
			"updated_at":       time.Now(),
		})

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("period not found in application")
	}

	// Получаем обновленную запись
	var updatedSelectedPeriod ds.SelectedPeriod
	if err := r.db.Where("application_id = ? AND period_id = ?", applicationID, periodID).
		First(&updatedSelectedPeriod).Error; err != nil {
		return nil, err
	}

	return &updatedSelectedPeriod, nil
}

func (r *RevenueModel) CreateUser(login, password string, isModerator bool) (*ds.User, error) {
	// Проверяем, нет ли уже пользователя с таким логином
	var existingUser ds.User
	if err := r.db.Where("login = ?", login).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("user with login '%s' already exists", login)
	}

	user := &ds.User{
		Login:       login,
		Password:    password,
		IsModerator: isModerator,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (r *RevenueModel) GetUserByID(userID int) (*ds.User, error) {
	var user ds.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *RevenueModel) UpdateUser(userID int, updateData map[string]interface{}) (*ds.User, error) {
	// Проверяем, что пользователь существует
	var existingUser ds.User
	if err := r.db.Where("id = ?", userID).First(&existingUser).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Если обновляется логин, проверяем что он уникальный
	if login, exists := updateData["login"]; exists {
		var userWithSameLogin ds.User
		if err := r.db.Where("login = ? AND id != ?", login, userID).First(&userWithSameLogin).Error; err == nil {
			return nil, fmt.Errorf("user with login '%s' already exists", login)
		}
	}

	// Добавляем время обновления
	updateData["updated_at"] = time.Now()

	// Обновляем пользователя
	result := r.db.Model(&ds.User{}).Where("id = ?", userID).Updates(updateData)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}

	// Получаем обновленного пользователя
	var updatedUser ds.User
	if err := r.db.Where("id = ?", userID).First(&updatedUser).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *RevenueModel) AuthenticateUser(login, password string) (*ds.User, error) {
	var user ds.User

	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if user.Password != password {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

func (r *RevenueModel) GetOrCreateApplication(applicationID, userID int) (*ds.PeriodsApplication, error) {
	// Сначала проверяем, существует ли заявка с таким ID
	var existingApplication ds.PeriodsApplication
	err := r.db.Where("id = ?", applicationID).First(&existingApplication).Error

	if err == nil {
		// Заявка существует, возвращаем её
		return &existingApplication, nil
	}

	// Заявка не существует, создаем новую
	newApplication := &ds.PeriodsApplication{
		ID:          applicationID,
		Year:        time.Now().Year(),
		CreatorID:   userID,
		Status:      "draft",
		CompanyName: "",
		CreatedAt:   time.Now(),
		FormedAt:    time.Now(),
		ModeratorID: nil,
		CompletedAt: nil,
	}

	err = r.db.Create(newApplication).Error
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании заявки: %w", err)
	}

	return newApplication, nil
}

func (r *RevenueModel) FormPeriodsApplication(id int, status string) (*ds.PeriodsApplication, error) {
	// Получаем текущую заявку
	var application ds.PeriodsApplication
	if err := r.db.Where("id = ?", id).First(&application).Error; err != nil {
		return nil, err
	}

	// Логируем текущие даты для отладки
	fmt.Printf("Application %d dates: created_at=%v, formed_at=%v\n",
		application.ID, application.CreatedAt, application.FormedAt)

	// Определяем правильное время для formed_at
	var formedAt time.Time
	now := time.Now()

	if application.CreatedAt.After(now) {
		// Если created_at в будущем, используем его
		formedAt = application.CreatedAt
	} else {
		// Иначе используем текущее время, но не раньше created_at
		formedAt = now
		if formedAt.Before(application.CreatedAt) {
			formedAt = application.CreatedAt
		}
	}

	fmt.Printf("Setting formed_at to: %v\n", formedAt)

	// Обновляем заявку
	updateData := map[string]interface{}{
		"status":    status,
		"formed_at": formedAt,
	}

	result := r.db.Model(&ds.PeriodsApplication{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating application: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("periods application not found")
	}

	// Получаем обновленную заявку
	var updatedApplication ds.PeriodsApplication
	if err := r.db.Where("id = ?", id).First(&updatedApplication).Error; err != nil {
		return nil, err
	}

	return &updatedApplication, nil
}

func (r *RevenueModel) GetPeriodsApplicationDates(id int) (map[string]interface{}, error) {
	var application ds.PeriodsApplication
	if err := r.db.Where("id = ?", id).First(&application).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         application.ID,
		"created_at": application.CreatedAt,
		"formed_at":  application.FormedAt,
		"status":     application.Status,
	}, nil
}
