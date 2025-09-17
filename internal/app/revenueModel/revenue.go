package revenueModel

import (
    "fmt"
	"lab2/internal/app/ds"
    "time"
    "math/rand"
    "math"
    "strings"
    "strconv"
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

func (r *RevenueModel) GetPeriodsApplication(id int) (map[string]interface{}, error) {
    periodsApplication := &ds.PeriodsApplication{}
    if err := r.db.Where("id = ? AND status != ?", id, "deleted").First(periodsApplication).Error; err != nil {
        return nil, fmt.Errorf("ошибка получения заявки", err)
    }

	var clientCompany map[string]interface{}
    if periodsApplication.CompanyID != nil {
        company := &ds.Company{}
        if err := r.db.Where("id = ?", *periodsApplication.CompanyID).First(company).Error; err != nil {
            clientCompany = map[string]interface{}{"ID": nil, "Name": ""}
        } else {
            clientCompany = map[string]interface{}{"ID": company.ID, "Name": company.Name}
        }
    } else {
        clientCompany = map[string]interface{}{"ID": nil, "Name": ""}
    }

	var selectedPeriods []ds.SelectedPeriod
    err := r.db.Where("application_id = ?", id).Find(&selectedPeriods).Error
    if err != nil {
        return nil, fmt.Errorf("ошибка получения выбранных периодов", err)
    }

    periods, err := r.GetPeriods()
    if err != nil {
        periods = []ds.Period{}
    }


    var detailedSelectedPeriods []map[string]interface{}
    for _, sp := range selectedPeriods {
        var title, shortDesc, img, duration string
        for _, p := range periods {
            if p.ID == sp.PeriodID {
                title = p.Title
                shortDesc = p.ShortDescription
                img = p.Img
                duration = p.Duration
                break
            }
        }

        parts := strings.Split(sp.PreviousRevenue, ";")
        var values []float64
        for _, p := range parts {
            v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
            if err == nil {
                values = append(values, v)
            }
        }

        forecast := calculateEMA(values)

        detailedSelectedPeriods = append(detailedSelectedPeriods, map[string]interface{}{
            "PeriodID":          sp.PeriodID,
            "ApplicationID":     sp.ApplicationID,
            "Title":             title,
            "Description":       shortDesc,
            "Img":               img,
            "Duration":          duration,
            "PreviousRevenue":   sp.PreviousRevenue,
            "ForecastedRevenue": forecast,
            "CreatedAt":         sp.CreatedAt,
            "UpdatedAt":         sp.UpdatedAt,
        })
    }

    applicationInfo := map[string]interface{}{
        "ApplicationID":  periodsApplication.ID,
        "Year":           periodsApplication.Year,
        "FullDuration":   periodsApplication.FullDuration,
        "ClientCompany":  clientCompany,
        "SelectedPeriods": detailedSelectedPeriods,
        "CreatedAt":      periodsApplication.CreatedAt,
        "FormedAt":       periodsApplication.FormedAt,
        "CompletedAt":    periodsApplication.CompletedAt,
    }

    return applicationInfo, nil
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


func (r *RevenueModel) GetSelectedPeriodsCount(applicationID int) (int, error) {
    var count int64
    if err := r.db.Table("selected_periods").
        Where("application_id = ?", applicationID).Count(&count).Error; err != nil {
        return 0, fmt.Errorf("ошибка подсчёта выбранных периодов:", err)
    }
    return int(count), nil
}


func (r *RevenueModel) GetOrCreateDraftApplication(userID int) (*ds.PeriodsApplication, error) {
    var draftApplication ds.PeriodsApplication
    err := r.db.Where("creator_id = ? AND status = ?", userID, "draft").First(&draftApplication).Error
	if err == nil {
		return &draftApplication, nil
	}

    newApplication := &ds.PeriodsApplication{
        Year:         time.Now().Year(), 
        CreatorID:    userID,
        Status:       "draft",
        CreatedAt:    time.Now(),
        FormedAt:     time.Now(),
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

func (r *RevenueModel) AddPeriodToApplication(applicationID, periodID int) error {
    var existingSelectedPeriod ds.SelectedPeriod
    err := r.db.Where("application_id = ? AND period_id = ?", applicationID, periodID).
        First(&existingSelectedPeriod).Error
    
    if err == nil {
        return fmt.Errorf("период уже добавлен в заявку")
    }

    var period ds.Period
    if err := r.db.Where("id = ?", periodID).First(&period).Error; err != nil {
        return fmt.Errorf("не удалось получить период", err)
    }
    previousRevenue := generatePreviousRevenueFromPeriod(period.Title)
    
    selectedPeriod := &ds.SelectedPeriod{
        ApplicationID:   applicationID,
        PeriodID:        periodID,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
        PreviousRevenue: previousRevenue,
    }
    
    return r.db.Create(selectedPeriod).Error
}


func (r *RevenueModel) DeleteApplication(applicationID int) error {
    result := r.db.Exec("UPDATE periods_applications SET status = 'deleted' WHERE id = ? AND status != 'deleted'", applicationID)

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("заявка не найдена или уже удалена")
    }

    return nil
}
