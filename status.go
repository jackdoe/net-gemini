package gemini

type Status int

const (
	StatusInput               = 10
	StatusSuccess             = 20
	StautsRedirect            = 30
	StatusTemporaryFailure    = 40
	StatusPermanentFailure    = 50
	StatusNotFound            = 51
	StatusForbidden           = 52
	StatusCertificateRequired = 60
)
