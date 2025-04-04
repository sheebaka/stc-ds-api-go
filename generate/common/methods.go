package common

type SharedAccount struct {
	DotNumber      int            `json:"dotNumber"`
	CrmNumber      int            `json:"crmNumber"`
	CompanyName    string         `json:"companyName"`
	CustomerStatus CustomerStatus `json:"customerStatus"`
}

type CustomerStatus struct {
	Factoring string `json:"factoring"`
}
