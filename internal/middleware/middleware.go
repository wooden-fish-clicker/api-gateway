package middleware

type Middleware struct {
	CheckJwtToken *CheckJwtToken
	RateLimiter   *RateLimiter
}

func NewMiddleware(checkJwtToken *CheckJwtToken, rateLimiter *RateLimiter) *Middleware {
	return &Middleware{
		CheckJwtToken: checkJwtToken,
		RateLimiter:   rateLimiter,
	}
}
