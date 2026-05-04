package main

import (
	"context"
	"fmt"

	"github.com/jelmerdehen/demangle"
	"github.com/jelmerdehen/demangle/scheme/swift/stable"
)

func main() {
	syms := []struct{ sym, expected string }{
		// Was passing, must remain passing
		{
			"$sSvSgA3ASbIetCyyd_SgSbIetCyyyd_SgD",
			"(@escaping @convention(thin) @convention(c) (@unowned Swift.UnsafeMutableRawPointer?, @unowned Swift.UnsafeMutableRawPointer?, @unowned (@escaping @convention(thin) @convention(c) (@unowned Swift.UnsafeMutableRawPointer?, @unowned Swift.UnsafeMutableRawPointer?) -> (@unowned Swift.Bool))?) -> (@unowned Swift.Bool))?",
		},
		// Was failing, needs to pass
		{
			"_$s10Foundation17_CalendarProtocolP4date8byAdding2to18wrappingComponentsAA4DateVSgAA0jI0V_AISbtFTj",
			"dispatch thunk of Foundation._CalendarProtocol.date(byAdding: Foundation.DateComponents, to: Foundation.Date, wrappingComponents: Swift.Bool) -> Foundation.Date?",
		},
		// Was passing, must remain passing
		{
			"_$s10Foundation17_TimeZoneProtocolP018nextDaylightSavingB10Transition5afterAA4DateVSgAG_tFTj",
			"dispatch thunk of Foundation._TimeZoneProtocol.nextDaylightSavingTimeTransition(after: Foundation.Date) -> Foundation.Date?",
		},
	}
	cat := demangle.NewCatalog()
	cat.Register(stable.Scheme{})
	ctx := context.Background()
	for _, s := range syms {
		result, err := cat.Demangle(ctx, s.sym, nil)
		if err != nil {
			fmt.Printf("Error:    %v\n", err)
		} else {
			match := "PASS"
			if result.Output != s.expected {
				match = "FAIL"
			}
			fmt.Printf("%s Got:  %s\n", match, result.Output)
		}
		fmt.Printf("     Exp: %s\n\n", s.expected)
	}
}
