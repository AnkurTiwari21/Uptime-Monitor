package jobs

type SQSIncidentEventType struct {
	WebsiteURL      string `json:"website_url"`
	Phone           string `json:"phone_number"`
	Email           string `json:"email"`
	Status          string `json:"status"`
	IncidentEventID string `json:"incident_event_id"`
}
