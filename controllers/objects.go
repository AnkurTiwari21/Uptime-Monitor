package controllers

type SignUpRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"required"`
	LastName  string `json:"last_name,omitempty" validate:"required"`
	Email     string `json:"email,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
}

type SignInRequest struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type RegisterWebsiteRequest struct {
	WebsiteURL string `json:"website_url" validate:"required"`
}

type RegisterWebsiteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type WebsiteLivelinessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UpdateAlertConfigRequest struct {
	WebsiteId        uint `json:"website_id,omitempty" validate:"required"`
	IsEnabled        bool `json:"is_enabled,omitempty"`
	FailureThreshold uint `json:"failure_threshold,omitempty"`
	LatencyThreshold uint `json:"latency_threshold,omitempty"`
}
type UpdateAlertConfigResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
