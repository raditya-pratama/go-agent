package entity

type ActivityLog struct {
	ElementID      string `json:"element_id"`
	NewData        string `json:"new_data"`
	OldData        string `json:"old_data"`
	DisplayMessage string `json:"display_message"`
	Payload        string `json:"payload"`
	UriPath        string `json:"uri_path"`
	ModuleURL      string `json:"module_url"`
	Message        string `json:"message"`
	Environment    string `json:"environment"`
	ActivityName   string `json:"activity_name"`
	ElementName    string `json:"element_name"`
	TribeName      string `json:"tribe_name"`
}

// AuditTrailConfig -
type AuditTrailConfig struct {
	Host        string
	ClientSpawn uint
	Timeout     uint
	Key         string
	Secret      string
	AppName     string
}
