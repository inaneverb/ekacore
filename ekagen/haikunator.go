package ekagen

// Ruby original: https://github.com/usmanbashir/haikunator
// Go ver of Ruby original: https://github.com/yelinaung/go-haikunator
// Not as package to easily add a new words.

// TODO: Add more words.
// ATM its ~ 400k-4kk diff variants with 4-digits tail.

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	adjectives = []string{
		"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry",
		"dark", "summer", "icy", "delicate", "quiet", "white", "cool",
		"spring", "winter", "patient", "twilight", "dawn", "crimson", "wispy",
		"weathered", "blue", "billowing", "broken", "cold", "damp", "falling",
		"frosty", "green", "long", "late", "lingering", "bold", "little",
		"morning", "muddy", "old", "red", "rough", "still", "small",
		"sparkling", "throbbing", "shy", "wandering", "withered", "wild",
		"black", "young", "holy", "solitary", "fragrant", "aged", "snowy",
		"proud", "floral", "restless", "divine", "polished", "ancient",
		"purple", "lively", "nameless",
	}

	nouns = []string{
		"waterfall", "river", "breeze", "moon", "rain", "wind", "sea",
		"morning", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn",
		"glitter", "forest", "hill", "cloud", "meadow", "sun", "glade", "bird",
		"brook", "butterfly", "bush", "dew", "dust", "field", "fire", "flower",
		"firefly", "feather", "grass", "haze", "mountain", "night", "pond",
		"darkness", "snowflake", "silence", "sound", "sky", "shape", "surf",
		"thunder", "violet", "water", "wildflower", "wave", "water",
		"resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost",
		"voice", "paper", "frog", "smoke", "star",
	}
)

//
var r *rand.Rand

//
func init() {
	r = rand.New(rand.New(rand.NewSource(99)))
	r.Seed(time.Now().UTC().Unix())
}

//
func Haikunate() string {
	return HaikunateWithRange(0, 9999)
}

//
func HaikunateWithRange(from, to uint) string {
	if from > to {
		from, to = to, from
	}
	return fmt.Sprintf(
		"%s-%s-%04d",
		adjectives[r.Intn(len(adjectives))],
		nouns[r.Intn(len(nouns))],
		from+(uint(r.Int())%(to-from)),
	)
}
