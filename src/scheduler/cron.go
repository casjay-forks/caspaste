// This file is part of CasPaste.

// CasPaste is free software released under the MIT License.
// See LICENSE.md file for details.

// Cron expression parsing per AI.md PART 19
// Supports standard cron format and @every duration syntax
package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronExpr represents a parsed cron expression
type CronExpr struct {
	minute  []int
	hour    []int
	day     []int
	month   []int
	weekday []int
	isEvery bool
	every   time.Duration
}

// ParseCron parses a cron expression
// Supports:
// - Standard: "minute hour day month weekday" (e.g., "0 3 * * *")
// - Every: "@every 15m" (e.g., "@every 1h", "@every 30s")
// - Predefined: "@daily", "@hourly", "@weekly", "@monthly"
func ParseCron(expr string) (*CronExpr, error) {
	expr = strings.TrimSpace(expr)

	// Handle @every expressions
	if strings.HasPrefix(expr, "@every ") {
		durStr := strings.TrimPrefix(expr, "@every ")
		dur, err := time.ParseDuration(durStr)
		if err != nil {
			return nil, fmt.Errorf("invalid duration: %s", durStr)
		}
		return &CronExpr{isEvery: true, every: dur}, nil
	}

	// Handle predefined expressions
	switch expr {
	case "@yearly", "@annually":
		expr = "0 0 1 1 *"
	case "@monthly":
		expr = "0 0 1 * *"
	case "@weekly":
		expr = "0 0 * * 0"
	case "@daily", "@midnight":
		expr = "0 0 * * *"
	case "@hourly":
		expr = "0 * * * *"
	}

	// Parse standard cron expression
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(parts))
	}

	c := &CronExpr{}
	var err error

	c.minute, err = parseField(parts[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("invalid minute field: %w", err)
	}

	c.hour, err = parseField(parts[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("invalid hour field: %w", err)
	}

	c.day, err = parseField(parts[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("invalid day field: %w", err)
	}

	c.month, err = parseField(parts[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("invalid month field: %w", err)
	}

	c.weekday, err = parseField(parts[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("invalid weekday field: %w", err)
	}

	return c, nil
}

// parseField parses a single cron field
func parseField(field string, min, max int) ([]int, error) {
	// Handle wildcard
	if field == "*" {
		vals := make([]int, max-min+1)
		for i := range vals {
			vals[i] = min + i
		}
		return vals, nil
	}

	var result []int

	// Handle comma-separated values
	parts := strings.Split(field, ",")
	for _, part := range parts {
		// Handle step values (*/5 or 1-10/2)
		if strings.Contains(part, "/") {
			stepParts := strings.Split(part, "/")
			if len(stepParts) != 2 {
				return nil, fmt.Errorf("invalid step: %s", part)
			}

			step, err := strconv.Atoi(stepParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid step value: %s", stepParts[1])
			}

			var rangeVals []int
			if stepParts[0] == "*" {
				for i := min; i <= max; i++ {
					rangeVals = append(rangeVals, i)
				}
			} else {
				rangeVals, err = parseRange(stepParts[0], min, max)
				if err != nil {
					return nil, err
				}
			}

			for i := 0; i < len(rangeVals); i += step {
				result = append(result, rangeVals[i])
			}
			continue
		}

		// Handle range (1-5)
		if strings.Contains(part, "-") {
			rangeVals, err := parseRange(part, min, max)
			if err != nil {
				return nil, err
			}
			result = append(result, rangeVals...)
			continue
		}

		// Handle single value
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid value: %s", part)
		}
		if val < min || val > max {
			return nil, fmt.Errorf("value %d out of range [%d-%d]", val, min, max)
		}
		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("empty field")
	}

	return result, nil
}

// parseRange parses a range expression (e.g., "1-5")
func parseRange(r string, min, max int) ([]int, error) {
	parts := strings.Split(r, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range: %s", r)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid range start: %s", parts[0])
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid range end: %s", parts[1])
	}

	if start > end {
		return nil, fmt.Errorf("range start %d > end %d", start, end)
	}
	if start < min || end > max {
		return nil, fmt.Errorf("range out of bounds [%d-%d]", min, max)
	}

	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result, nil
}

// Next returns the next scheduled time after the given time
func (c *CronExpr) Next(after time.Time) time.Time {
	// Handle @every expressions
	if c.isEvery {
		return after.Add(c.every)
	}

	// Start from the next minute
	t := after.Truncate(time.Minute).Add(time.Minute)

	// Try to find the next matching time (max 4 years to prevent infinite loop)
	maxIterations := 366 * 4 * 24 * 60
	for i := 0; i < maxIterations; i++ {
		if c.matches(t) {
			return t
		}
		t = t.Add(time.Minute)
	}

	// No match found
	return time.Time{}
}

// matches checks if the given time matches the cron expression
func (c *CronExpr) matches(t time.Time) bool {
	return contains(c.minute, t.Minute()) &&
		contains(c.hour, t.Hour()) &&
		contains(c.day, t.Day()) &&
		contains(c.month, int(t.Month())) &&
		contains(c.weekday, int(t.Weekday()))
}

// contains checks if a slice contains a value
func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// String returns the string representation of the cron expression
func (c *CronExpr) String() string {
	if c.isEvery {
		return "@every " + c.every.String()
	}
	return fmt.Sprintf("%v %v %v %v %v", c.minute, c.hour, c.day, c.month, c.weekday)
}
