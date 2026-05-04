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
			"_$ss20_SwiftNewtypeWrapperPss21_ObjectiveCBridgeable8RawValueRpzrlE016_forceBridgeFromD1C_6resultyAD_01_D5CTypeQZ_xSgztFZ",
			"static (extension in Swift):Swift._SwiftNewtypeWrapper< where A.RawValue: Swift._ObjectiveCBridgeable>._forceBridgeFromObjectiveC(_: A.RawValue._ObjectiveCType, result: inout A?) -> ()",
		},
		{
			"_$sSlsSIyxG7IndicesRtzrlE7indicesAAvpMV",
			"property descriptor for (extension in Swift):Swift.Collection< where A.Indices == Swift.DefaultIndices<A>>.indices : Swift.DefaultIndices<A>",
		},
		{
			"_$sSksSx5IndexRpzSnyABG7IndicesRtzSiAA_6StrideRTzrlE7indicesACvpMV",
			"property descriptor for (extension in Swift):Swift.RandomAccessCollection< where A.Index: Swift.Strideable, A.Indices == Swift.Range<A.Index>, A.Index.Stride == Swift.Int>.indices : Swift.Range<A.Index>",
		},
		{
			"_$s10Foundation11FormatStylePA2A07IntegerbC0VySiGRszrlE6numberAFvpZMV",
			"property descriptor for static (extension in Foundation):Foundation.FormatStyle< where A == Foundation.IntegerFormatStyle<Swift.Int>>.number : Foundation.IntegerFormatStyle<Swift.Int>",
		},
		{
			"_$s10Foundation11MeasurementVAASo11NSDimensionCRbzrlE11FormatStyleV6localeAA6LocaleVvpMV",
			"property descriptor for (extension in Foundation):Foundation.Measurement< where A: __C.NSDimension>.FormatStyle.locale : Foundation.Locale",
		},
		{
			"_$s10Foundation13CustomNSErrorPAASYRzs17FixedWidthInteger8RawValueSYRpzrlE9errorCodeSivpMV",
			"property descriptor for (extension in Foundation):Foundation.CustomNSError< where A: Swift.RawRepresentable, A.Swift.RawRepresentable.RawValue: Swift.FixedWidthInteger>.errorCode : Swift.Int",
		},
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
		tag := s.sym
		if len(tag) > 40 {
			tag = tag[:40]
		}
		result, err := cat.Demangle(context.Background(), s.sym, nil)
		if err != nil {
			fmt.Printf("ERROR[%s...]: %v\n", tag, err)
			continue
		}
		got := result.Output
		if got == s.want {
			fmt.Printf("PASS[%s...]\n", tag)
		} else {
			fmt.Printf("FAIL[%s...]\n  GOT:  %s\n  WANT: %s\n", tag, got, s.want)
		}
	}
}
