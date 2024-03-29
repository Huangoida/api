package service

type AutoGenerated struct {
	APIs []APIs `json:"APIs"`
}

type Services struct {
	Host string `json:"Host"`
	Port int    `json:"Port"`
}

type APIs struct {
	ID         string             `json:"id"`
	BatchIndex int                `json:"BatchIndex"`
	ServiceID  string             `json:"serviceId"`
	Services   Services           `json:"Services"`
	APIID      string             `json:"apiId"`
	Name       string             `json:"Name"`
	Path       string             `json:"Path"`
	Protocol   string             `json:"Protocol"`
	Method     string             `json:"Method"`
	Parameter  []ParametersStruct `json:"Parameter,omitempty"`
	ParentID   int                `json:"-"`
	ParentAPI  []interface{}      `json:"ParentApi,omitempty"`
}
