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
