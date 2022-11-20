package response

type (
	GeoResponse struct {
		IP      string  `json:"ip"`
		City    City    `json:"city"`
		Region  Region  `json:"region"`
		Country Country `json:"country"`
	}

	City struct {
		ID    int     `json:"id"`
		Lat   float64 `json:"lat"`
		Lon   float64 `json:"lon"`
		Name  string  `json:"name_en"`
		Okato string  `json:"okato"`
	}

	Region struct {
		ID       int     `json:"id"`
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		Name     string  `json:"name_en"`
		Okato    string  `json:"okato"`
		Iso      string  `json:"iso"`
		TimeZone string  `json:"timezone"`
	}

	Country struct {
		ID        int     `json:"id"`
		Iso       string  `json:"iso"`
		Continent string  `json:"continent"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Name      string  `json:"name_en"`
		TimeZone  string  `json:"timezone"`
	}
)
