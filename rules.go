package gol

// Colour values
type Colour struct {
	R, G, B float32
}

// Rule with alive status and transitions which
// represent what the rule changes to based on
// amount of adjacent alive cells (0-8)
type Rule struct {
	Alive       bool
	Transitions [9]uint8
	Colour      Colour
}

// Randomize a single Rule
func (ru *Rule) Randomize(RuleAmount int) {
	ru.Alive = randInt(2) == 0
	ru.Colour.R = float32(float32(randInt(255)) / 255.0)
	ru.Colour.G = float32(float32(randInt(255)) / 255.0)
	ru.Colour.B = float32(float32(randInt(255)) / 255.0)
	for idx := range ru.Transitions {
		ru.Transitions[idx] = uint8(randInt(RuleAmount))
	}
}

// Rules is an ordered array of Rule structs
type Rules struct {
	Array []Rule
}

// Randomize an array of Rules
func (rs *Rules) Randomize(RuleAmount int) {
	rs.Array = make([]Rule, RuleAmount)
	for idx := range rs.Array {
		rs.Array[idx].Randomize(RuleAmount)
	}
}
