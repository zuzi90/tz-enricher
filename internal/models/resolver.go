package models

type (
	AgeResolver struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
		Age   int    `json:"age"`
	}

	GenderResolver struct {
		Count       int     `json:"count"`
		Name        string  `json:"name"`
		Gender      string  `json:"gender"`
		Probability float64 `json:"probability"`
	}

	NationalityResolver struct {
		Count   int    `json:"count"`
		Name    string `json:"name"`
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}
)
