package gemini

type Status int

const (
	StatusInput               = 10
	StatusSuccess             = 20
	StautsRedirect            = 30
	StatusTemporaryFailure    = 40
	StatusNotFound            = 41
	StatusForbidden           = 42
	StatusPermanentFailure    = 50
	StatusCertificateRequired = 60
)
