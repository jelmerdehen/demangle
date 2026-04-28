package stable_test

import (
	"context"
	"fmt"
	"testing"
	_ "github.com/jelmerdehen/demangle/scheme/swift/stable"
	"github.com/jelmerdehen/demangle"
)

func TestStaticExt(t *testing.T) {
	syms := []string{
		"_$s10Foundation11FormatStylePA2A014DateComponentsV07ISO8601bC0VRszrlE7iso8601AGvgZ",
		"_$s10Foundation11FormatStylePA2A013FloatingPointbC0VySdGRszrlE6numberAFvpZMV",
		"_$s10Foundation11FormatStylePA2A013FloatingPointbC0VySfGRszrlE6numberAFvpZMV",
		"_$s10Foundation11FormatStylePA2A3URLVABVRszrlE3urlAFvgZ",
		"_$s10Foundation11FormatStylePA2A4DateV010ComponentsbC0VRszrlE12timeDurationAGvgZ",
	}
	wants := []string{
		`static (extension in Foundation):Foundation.FormatStyle< where A == Foundation.DateComponents.ISO8601FormatStyle>.iso8601.getter : Foundation.DateComponents.ISO8601FormatStyle`,
		`property descriptor for static (extension in Foundation):Foundation.FormatStyle< where A == Foundation.FloatingPointFormatStyle<Swift.Double>>.number : Foundation.FloatingPointFormatStyle<Swift.Double>`,
		`property descriptor for static (extension in Foundation):Foundation.FormatStyle< where A == Foundation.FloatingPointFormatStyle<Swift.Float>>.number : Foundation.FloatingPointFormatStyle<Swift.Float>`,
		`static (extension in Foundation):Foundation.FormatStyle< where A == Foundation.URL.FormatStyle>.url.getter : Foundation.URL.FormatStyle`,
		`static (extension in Foundation):Foundation.FormatStyle< where A == Foundation.Date.ComponentsFormatStyle>.timeDuration.getter : Foundation.Date.ComponentsFormatStyle`,
	}
	for i, sym := range syms {
		r, err := demangle.Default.Demangle(context.Background(), sym, nil)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			got := r.Output
			want := wants[i]
			if got == want {
				fmt.Printf("PASS: %s\n", got)
			} else {
				fmt.Printf("FAIL:\n  got:  %s\n  want: %s\n", got, want)
			}
		}
	}
}
