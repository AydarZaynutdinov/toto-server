package entity

type SkuConfig struct {
	ID            string `json:"id" db:"id"`
	Package       string `json:"package" db:"package"`
	CountryCode   string `json:"country_code" db:"country_code"`
	PercentileMin int    `json:"percentile_min" db:"percentile_min"`
	PercentileMax int    `json:"percentile_max" db:"percentile_max"`
	MainSku       string `json:"main_sku" db:"main_sku"`
}
