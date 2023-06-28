package core

import (
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ErrCaptcha = errors.New("Captcha detected")
var ErrSearchTimeout = errors.New("Timeout. Cannot find element on page")

type SearchResult struct {
	Rank        int    `json:"rank"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Query struct {
	Text         string
	LangCode     string // eg. EN, ES, RU...
	DateInterval string // format: YYYYMMDD..YYYMMDD - 20181010..20231010
	Filetype     string // File extension to search.
	Site         string // Search site
	Limit        int    // Limit the number of results
}

func (q Query) IsEmpty() bool {
	if q.Site == "" && q.Filetype == "" && q.Text == "" {
		return true
	}
	return false
}

func (q *Query) InitFromContext(c *fiber.Ctx) error {
	q.Text = c.Query("text")
	q.LangCode = c.Query("lang")
	q.DateInterval = c.Query("date")
	q.Filetype = c.Query("file")
	q.Site = c.Query("site")

	limit, err := strconv.Atoi(c.Query("limit", "25"))
	if err != nil {
		return err
	}
	q.Limit = limit

	if q.IsEmpty() {
		return errors.New("Query cannot be empty")
	}

	return nil
}

type SearchEngineOptions struct {
	RateRequests    int           `mapstructure:"rate_requests"`
	RateTime        time.Duration `mapstructure:"rate_seconds"`
	RateBurst       int           `mapstructure:"rate_burst"`
	SelectorTimeout time.Duration `mapstructure:"selector_timeout"`
}

func (o *SearchEngineOptions) Init() {
	if o.RateRequests == 0 {
		o.RateRequests = 1
	}
	if o.RateTime == 0 {
		o.RateTime = time.Second * 10
	}
	if o.RateBurst == 0 {
		o.RateBurst = 1
	}
	if o.SelectorTimeout == 0 {
		o.SelectorTimeout = time.Second * 5
	}
}
