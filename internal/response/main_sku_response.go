package response

type (
	MainSkuResponse struct {
		MainSku string `json:"main_sku"`
	}

	ErrorResponse struct {
		Message string `json:"message"`
	}
)
