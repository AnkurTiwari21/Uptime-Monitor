package models

import (
	"gorm.io/gorm"
)

func InitUserRepo(DB *gorm.DB) IUser {
	return &userRepo{
		db: DB,
	}
}

func InitWebsiteRepo(DB *gorm.DB) IWebiste {
	return &websiteRepo{
		db: DB,
	}
}

func InitAlertConfigRepo(DB *gorm.DB) IAlertConfig {
	return &alertConfigRepo{
		db: DB,
	}
}

func InitAlertTargetRepo(DB *gorm.DB) IAlertTarget {
	return &alertTargetRepo{
		db: DB,
	}
}

func InitLogsRepo(DB *gorm.DB) ILog {
	return &logsRepo{
		db: DB,
	}
}

func InitIncidentsRepo(DB *gorm.DB) IIncident {
	return &incidentsRepo{
		db: DB,
	}
}

func InitIncidentEventsRepo(DB *gorm.DB) IIncidentEvent {
	return &incidentEventsRepo{
		db: DB,
	}
}
