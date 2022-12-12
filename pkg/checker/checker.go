package checker

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alexliesenfeld/health"
)

type Checker struct {
	Websites  []string
	Databases []*sql.DB
}

func (c *Checker) CheckSites() {
	checker := health.NewChecker(
		// Set the time-to-live for our cache to 1 second (default).
		health.WithCacheDuration(1 * time.Second),
	)
	fmt.Println(checker)
}
