package economy

// Good identifies a trade commodity.
type Good int

const (
	GoodRum Good = iota
	GoodPowder
	GoodCloth
)

var AllGoods = []Good{GoodRum, GoodPowder, GoodCloth}

func (g Good) String() string {
	switch g {
	case GoodRum:
		return "rum"
	case GoodPowder:
		return "powder"
	case GoodCloth:
		return "cloth"
	default:
		return "?"
	}
}

func (g Good) Label() string {
	switch g {
	case GoodRum:
		return "Rum"
	case GoodPowder:
		return "Powder"
	case GoodCloth:
		return "Cloth"
	default:
		return "?"
	}
}

// CargoHold tracks goods with a fixed capacity.
type CargoHold struct {
	capacity int
	items    [3]int
}

func NewCargoHold(capacity int) CargoHold {
	return CargoHold{capacity: capacity}
}

func (c *CargoHold) Capacity() int {
	return c.capacity
}

func (c *CargoHold) Total() int {
	total := 0
	for _, n := range c.items {
		total += n
	}
	return total
}

func (c *CargoHold) FreeSpace() int {
	return c.capacity - c.Total()
}

func (c *CargoHold) Amount(g Good) int {
	if int(g) < 0 || int(g) >= len(c.items) {
		return 0
	}
	return c.items[g]
}

func (c *CargoHold) Set(g Good, amount int) {
	if int(g) >= 0 && int(g) < len(c.items) {
		c.items[g] = amount
	}
}

// Add returns how many units were actually added.
func (c *CargoHold) Add(g Good, amount int) int {
	if amount <= 0 {
		return 0
	}
	space := c.FreeSpace()
	if space <= 0 {
		return 0
	}
	added := amount
	if added > space {
		added = space
	}
	c.items[g] += added
	return added
}

// Remove returns how many units were actually removed.
func (c *CargoHold) Remove(g Good, amount int) int {
	if amount <= 0 {
		return 0
	}
	have := c.Amount(g)
	if have <= 0 {
		return 0
	}
	removed := amount
	if removed > have {
		removed = have
	}
	c.items[g] -= removed
	return removed
}

func (c CargoHold) Summary() string {
	return formatCargoSummary(c.items[:])
}

func formatCargoSummary(items []int) string {
	parts := make([]string, 0, len(AllGoods))
	for i, g := range AllGoods {
		if items[i] > 0 {
			parts = append(parts, g.Label()+":"+itoa(items[i]))
		}
	}
	if len(parts) == 0 {
		return "empty"
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out += ", " + p
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var digits [20]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		digits[i] = '-'
	}
	return string(digits[i:])
}
