package economy

import "fmt"

// Clock tracks world time in ticks.
type Clock struct {
	tick        int
	ticksPerDay int
}

func NewClock(ticksPerDay int) *Clock {
	return &Clock{ticksPerDay: ticksPerDay}
}

func (c *Clock) Tick() {
	c.tick++
}

func (c *Clock) CurrentTick() int {
	return c.tick
}

func (c *Clock) TicksPerDay() int {
	return c.ticksPerDay
}

// TimeOfDay returns HH:MM for the current tick within the day cycle.
func (c *Clock) TimeOfDay() string {
	if c.ticksPerDay <= 0 {
		return "00:00"
	}
	dayTick := c.tick % c.ticksPerDay
	if dayTick < 0 {
		dayTick += c.ticksPerDay
	}
	totalMinutes := dayTick * 24 * 60 / c.ticksPerDay
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
