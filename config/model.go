package config

import (
	"auth-api/internal/models/custom"
	"time"
)

const (
	EnvProduction  = "prod"
	EnvDevelopment = "dev"
	EnvStage       = "stage"
)

type Config struct {
	ENV        string `envconfig:"APP_ENV"`
	ServerPORT string `envconfig:"SERVER_PORT"`

	APIs
}

type APIs struct {
	Employ struct {
		Host           string        `envconfig:"EMPLOY_HOST" default:"http://employ.wbjobs.svc.k8s.dldevel"`
		Token          custom.Token  `envconfig:"EMPLOY_TOKEN"`
		RequestTimeout time.Duration `envconfig:"EMPLOY_REQUEST_TIMEOUT" default:"40s"`
	}

	UserInfo struct {
		Host           string        `envconfig:"USER_INFO_HOST" default:"https://api.user-info.svc.k8s.dldevel"`
		AuthID         string        `envconfig:"USER_INFO_AUTH_ID"`
		AuthSign       string        `envconfig:"USER_INFO_AUTH_SIGN"`
		RequestTimeout time.Duration `envconfig:"USER_INFO_REQUEST_TIMEOUT" default:"40s"`
	}

	Tarificator struct {
		Host           string        `envconfig:"TARIFICATOR_HOST" default:"https://wh-finance-salary.wildberries.ru"`
		Token          custom.Token  `envconfig:"TARIFICATOR_TOKEN"`
		RequestTimeout time.Duration `envconfig:"TARIFICATOR_REQUEST_TIMEOUT" default:"40s"`
	}

	HR struct {
		API struct {
			Host           string        `envconfig:"HR_HOST" default:"http://hr-employees.hr.svc.k8s.prod-dl"`
			Token          custom.Token  `envconfig:"HR_TOKEN"`
			RequestTimeout time.Duration `envconfig:"HR_REQUEST_TIMEOUT" default:"40s"`
		}
		Employee struct {
			Host           string        `envconfig:"HR_EMPLOYEE_HOST" default:"http://hr-employee-api.hr.svc.k8s.prod-dl"`
			Token          custom.Token  `envconfig:"HR_EMPLOYEE_TOKEN"`
			RequestTimeout time.Duration `envccnfig:"HR_REQUEST_TIMEOUT" default:"40s"`
		}
		ApplicationsRegistration struct {
			Host           string        `envconfig:"HR_APPLICATIONS_REGISTRATION_HOST" default:"http://hr-applications-registration.hr.svc.k8s.prod-dl"`
			Token          custom.Token  `envconfig:"HR_APPLICATIONS_REGISTRATION_TOKEN"`
			RequestTimeout time.Duration `envccnfig:"HR_APPLICATIONS_REGISTRATION_REQUEST_TIMEOUT" default:"40s"`
		}
	}
}
