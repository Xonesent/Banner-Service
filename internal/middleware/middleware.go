package middleware

import "avito/assignment/config"

type MDWManager struct {
	cfg *config.Config
}

func NewOfficiantMiddleware(cfg *config.Config) *MDWManager {
	return &MDWManager{cfg: cfg}
}
