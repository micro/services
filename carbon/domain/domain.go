package domain

type Project struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"splitPercentage"`
	Tonnes     float64 `json:"splitAmountTonnes"`
}

type OffsetRequest struct {
	Number int32  `json:"number"`
	Units  string `json:"units"`
}

type OffsetResponse struct {
	Number   int32     `json:"number"`
	Units    string    `json:"units"`
	Tonnes   float64   `json:"numberInTonnes"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
	Projects []Project `json:"projectDetails"`
}
