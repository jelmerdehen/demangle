package stable_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

func TestCheckFailingSymbols(t *testing.T) {
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})

	syms := []struct{ sym, want string }{
		{
			"_$s10Foundation17_CalendarProtocolP4date8byAdding2to18wrappingComponentsAA4DateVSgAA0jI0V_AISbtFTj",
			"dispatch thunk of Foundation._CalendarProtocol.date(byAdding: Foundation.DateComponents, to: Foundation.Date, wrappingComponents: Swift.Bool) -> Foundation.Date?",
		},
		{
			"_$s10Foundation17_CalendarProtocolP4copy14changingLocale0E8TimeZone0E12FirstWeekday0e13MinimumDaysInI4WeekAaB_pAA0F0VSg_AA0gH0VSgSiSgAOtFTj",
			"dispatch thunk of Foundation._CalendarProtocol.copy(changingLocale: Foundation.Locale?, changingTimeZone: Foundation.TimeZone?, changingFirstWeekday: Swift.Int?, changingMinimumDaysInFirstWeek: Swift.Int?) -> Foundation._CalendarProtocol",
		},
		{
			"_$ss10_UTFParserP15_bufferedScalar8bitCount8Encoding_07EncodedC0QZs5UInt8V_tFTj",
			"dispatch thunk of Swift._UTFParser._bufferedScalar(bitCount: Swift.UInt8) -> A.Encoding.EncodedScalar",
		},
		{
			"_$ss20_ArrayBufferProtocolP13_copyContents8subRange12initializingSpy7ElementQzGSnySiG_AHtFTj",
			"dispatch thunk of Swift._ArrayBufferProtocol._copyContents(subRange: Swift.Range<Swift.Int>, initializing: Swift.UnsafeMutablePointer<A.Element>) -> Swift.UnsafeMutablePointer<A.Element>",
		},
		{
			"_$ss15WritableKeyPathC22_projectMutableAddress4fromSpyq_G7pointer_yXlSg5ownertSPyxG_tFTj",
			"dispatch thunk of Swift.WritableKeyPath._projectMutableAddress(from: Swift.UnsafePointer<A>) -> (pointer: Swift.UnsafeMutablePointer<B>, owner: Swift.AnyObject?)",
		},
	}

	for _, s := range syms {
		result, err := cat.Demangle(context.Background(), s.sym, nil)
		if err != nil {
			fmt.Printf("ERROR[%s...]: %v\n", s.sym[:40], err)
			continue
		}
		got := result.Output
		if got == s.want {
			fmt.Printf("PASS[%s...]\n", s.sym[:40])
		} else {
			fmt.Printf("FAIL[%s...]\n  GOT:  %s\n  WANT: %s\n", s.sym[:40], got, s.want)
		}
	}
}
