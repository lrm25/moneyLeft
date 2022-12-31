package account

const SS_EARLY_AGE = 62
const SS_NORMAL_AGE = 67
const SS_LATE_AGE = 70

type SocialSecuritySelection int

const (
	Early SocialSecuritySelection = iota + 1
	Normal
	Late
)

type SocialSecurityAccount struct {
	selection SocialSecuritySelection
	perMonth  float64
}

type SSAccount struct {
	account *Account
}
