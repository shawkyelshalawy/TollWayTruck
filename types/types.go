package types

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type Distance struct {
	OBUID int     `json:"obuiD"`
	Value float64 `json:"value"`
	Unix  int64   `json:"unixTime"`
}

type Invoice struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
}

type GetInvoiceRequest struct {
	ObuID int32 `json:"obuID"`
}

type AggregateRequest struct {
	ObuID int32   `json:"obuID"`
	Value float64 `json:"value"`
	Unix  int64   `json:"unixTime"`
}
