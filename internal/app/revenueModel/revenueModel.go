package revenueModel

import (
	"fmt"
	"strings"
)

type RevenueModel struct {
}

func NewRevenueModel() (*RevenueModel, error) {
	return &RevenueModel{}, nil
}

type Period struct {
	ID               int
	Title            string
	Description      string
	Img              string
	Duration         string
	ShortDescription string
	MoreDescription  string
}

type Company struct {
	ID   int
	Name string
}

type SelectedPeriod struct {
	ID                int
	PeriodID          int
	PreviousRevenue   string
	ForecastedRevenue int
}

type PeriodsApplication struct {
	ID            int
	Year          int
	ClientCompany Company
	Periods       []SelectedPeriod
}

func (r *RevenueModel) GetPeriodsApplication(id int) (map[string]interface{}, error) {
	periodsApplication := &PeriodsApplication{
		ID:   1,
		Year: 2025,
		ClientCompany: Company{
			ID:   1,
			Name: "Wildberries",
		},
		Periods: []SelectedPeriod{
			{
				ID:                1,
				PeriodID:          1,
				PreviousRevenue:   "17000; 17000; 17000; 17000",
				ForecastedRevenue: 17000,
			},
		},
	}

	periods, err := r.GetPeriods()
	if err != nil {
		return nil, err
	}

	var selectedPeriods []map[string]interface{}
	for _, p := range periodsApplication.Periods {
		for _, period := range periods {
			if p.PeriodID == period.ID {
				selectedPeriods = append(selectedPeriods,
					map[string]interface{}{
						"ID":                p.ID,
						"Title":             period.Title,
						"Description":       period.ShortDescription,
						"Img":               period.Img,
						"Duration":          period.Duration,
						"PreviousRevenue":   p.PreviousRevenue,
						"ForecastedRevenue": p.ForecastedRevenue,
					})
			}
		}
	}

	periodsApplicationInfo := map[string]interface{}{
		"ApplicationID": periodsApplication.ID,
		"Year":          periodsApplication.Year,
		"ClientCompany": map[string]interface{}{
			"ID":   periodsApplication.ClientCompany.ID,
			"Name": periodsApplication.ClientCompany.Name,
		},
		"SelectedPeriods": selectedPeriods,
	}

	return periodsApplicationInfo, nil
}

func (r *RevenueModel) GetPeriods() ([]Period, error) {
	periods := []Period{
		{
			ID:               1,
			Title:            "4 квартала",
			Description:      "Прогнозирование выручки за следующий квартал методом скользящей средней на основе данных за 4 квартала",
			ShortDescription: "Прогнозирование выручки за следующий квартал на основе данных за 4 квартала",
			MoreDescription:  "Методом скользящей средней вычисляется прогнозируемое значение выручки за следующий период (например, квартал). За основу берутся данные о выручке за предыдущие кварталы (не менее четырех). При расчете используется экспоненциальный метод скользящей средней, что увеличивает точность прогноза.",
			Img:              "http://localhost:9002/periods/4quarters.png",
			Duration:         "2 дня",
		},
		{
			ID:               2,
			Title:            "5 месяцев",
			Description:      "Прогнозирование выручки за следующий месяц методом скользящей средней на основе данных за 5 месяцев",
			ShortDescription: "Прогнозирование выручки за следующий месяц на основе данных за 5 месяцев",
			MoreDescription:  "Метод скользящей средней применяется для месячных данных. На основе выручки за предыдущие 5 месяцев строится прогноз на следующий месяц. Метод особенно эффективен для сезонных бизнесов с четко выраженными месячными циклами продаж и операций.",
			Img:              "http://localhost:9002/periods/5months.png",
			Duration:         "2 дня",
		},
		{
			ID:               3,
			Title:            "18 кварталов",
			Description:      "Прогнозирование выручки за следующий квартал методом скользящей средней на основе данных за 18 кварталов",
			ShortDescription: "Прогнозирование выручки за следующий квартал на основе данных за 18 кварталов",
			MoreDescription:  "Долгосрочное прогнозирование на основе данных за 4.5 года (18 кварталов). Позволяет учесть долгосрочные тренды и циклические колебания бизнеса. Подходит для компаний с устоявшейся бизнес-моделью и историей.",
			Img:              "http://localhost:9002/periods/18quarters.png",
			Duration:         "5 дней",
		},
		{
			ID:               4,
			Title:            "12 месяцев",
			Description:      "Прогнозирование выручки за следующий месяц методом скользящей средней на основе данных за 12 месяцев",
			ShortDescription: "Прогнозирование выручки за следующий месяц на основе данных за 12 месяцев",
			MoreDescription:  "Годовой анализ месячных данных позволяет учесть сезонность и годовые циклы бизнеса. Особенно полезен для розничной торговли, туризма и других отраслей с выраженной сезонной динамикой выручки.",
			Img:              "http://localhost:9002/periods/12months.png",
			Duration:         "4 дня",
		},
		{
			ID:               5,
			Title:            "12 кварталов",
			Description:      "Прогнозирование выручки за следующий квартал методом скользящей средней на основе данных за 12 кварталов",
			ShortDescription: "Прогнозирование выручки за следующий квартал на основе данных за 12 кварталов",
			MoreDescription:  "Трехлетний анализ квартальных данных обеспечивает высокую точность прогнозирования. Учитывает как краткосрочные колебания, так и среднесрочные тренды развития компании и рынка в целом.",
			Img:              "http://localhost:9002/periods/12quarters.png",
			Duration:         "5 дней",
		},
		{
			ID:               6,
			Title:            "24 месяца",
			Description:      "Прогнозирование выручки за следующий месяц методом скользящей средней на основе данных за 24 месяца",
			ShortDescription: "Прогнозирование выручки за следующий месяц на основе данных за 24 месяца",
			MoreDescription:  "Двухлетний анализ месячных данных предоставляет наиболее точный прогноз для среднесрочного планирования. Идеален для растущих компаний, позволяет учесть темпы роста и сезонные паттерны выручки.",
			Img:              "http://localhost:9002/periods/24months.png",
			Duration:         "6 дней",
		},
	}

	if len(periods) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return periods, nil
}

func (r *RevenueModel) GetPeriod(id int) (Period, error) {
	periods, err := r.GetPeriods()
	if err != nil {
		return Period{}, err
	}

	for _, period := range periods {
		if period.ID == id {
			return period, nil
		}
	}
	return Period{}, fmt.Errorf("заказ не найден")
}

func (r *RevenueModel) GetPeriodsByTitle(title string) ([]Period, error) {
	periods, err := r.GetPeriods()
	if err != nil {
		return []Period{}, err
	}

	var result []Period
	for _, period := range periods {
		if strings.Contains(strings.ToLower(period.Title), strings.ToLower(title)) {
			result = append(result, period)
		}
	}

	return result, nil
}

func (r *RevenueModel) GetSelectedPeriodsCount(id int) (int, error) {
	periodsApplication, err := r.GetPeriodsApplication(id)
	if err != nil {
		return 0, err
	}

	selectedPeriods, ok := periodsApplication["SelectedPeriods"].([]map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("неверный формат SelectedPeriods")
	}

	return len(selectedPeriods), nil
}
