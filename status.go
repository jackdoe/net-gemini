package gemini

type Status int

const (
	StatusInput                              = 10
	StatusSuccess                            = 20
	StatusSuccessEndClientCertificateSession = 21
	StatusRedirectTemporary                  = 30
	StatusRedirectPermanent                  = 31
	StatusTemporaryFailure                   = 40
	StatusServerUnavailable                  = 41
	StatusCGIError                           = 42
	StatusProxyError                         = 43
	StatusSlowDown                           = 44
	StatusPermanentFailure                   = 50
	StatusNotFound                           = 51
	StatusGone                               = 52
	StatusProxyRequestRefused                = 53
	StatusBadRequest                         = 59
	StatusClientCertRequired                 = 60
	StatusTransientCertRequested             = 61
	StatusAuthorisedCertRequired             = 62
	StatusCertNotAccepted                    = 63
	StatusFutureCertRejected                 = 64
	StatusExpiredCertRejected                = 65
)
