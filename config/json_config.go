package config

import "github.com/lrm25/moneyLeft/account"

type JsonConfig struct {
	inflationRate       float64
	stockReturnRate     float64
	bondReturnRate      float64
	socialSecurity      account.SocialSecurityAccount
	leisureMultiplier   float64
	moneyNeededPerMonth float64
}
