package constants

const (
	GENERIC_SUCCESS_RESPONSE = "SUCCESS"
	GENERIC_FAILURE_RESPONSE = "FAILED"
)

const (
	WEBISTE_TYPE        = "WEBSITE"
	USER_TYPE           = "USER"
	INCIDENT_EVENT_TYPE = "INCIDENT_EVENT"
)

type Error struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}
