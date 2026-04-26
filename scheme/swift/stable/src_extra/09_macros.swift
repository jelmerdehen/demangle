// Swift 5.9+ macro declarations
// Expansion symbols (@__swiftmacro_) are emitted at compile time when macros are used;
// they cannot be compiled in isolation without the macro plugin.
// The fixture corpus (09_macros.txt) uses fMX (MacroExpansionLoc) symbols from Apple's
// swift/test/Macros/ test suite and grammar-derived synthetic examples.
//
// Grammar: @__swiftmacro_ <module> <file-punycode> fMX <line-1>_ <col-1>_ <macro-name> fMf <disc>_
// Covered roles: fMf (freestanding expression/declaration)
// Grammar gaps (not yet covered): fMp, fMm, fMr, fMe, fMa (attached roles)

@freestanding(expression)
public macro stringify<T>(_ value: T) -> (T, String) = #externalMacro(module: "MacroPlugin", type: "StringifyMacro")

@freestanding(expression)
public macro addOne(_ value: Int) -> Int = #externalMacro(module: "MacroPlugin", type: "AddOneMacro")

@freestanding(expression)
public macro coerceToInt<T: BinaryInteger>(_ value: T) -> Int = #externalMacro(module: "MacroPlugin", type: "CoerceMacro")

@freestanding(declaration, names: arbitrary)
public macro bitwidthNumbers() = #externalMacro(module: "MacroPlugin", type: "BitwidthNumbersMacro")

@freestanding(declaration, names: arbitrary)
public macro multiStatement() = #externalMacro(module: "MacroPlugin", type: "MultiStatementMacro")

@attached(member, names: named(init))
public macro MemberwiseInit() = #externalMacro(module: "MacroPlugin", type: "MemberwiseInitMacro")

@attached(peer)
public macro addCompletionHandler() = #externalMacro(module: "MacroPlugin", type: "AddCompletionHandlerMacro")

@attached(memberAttribute)
public macro myTypeWrapper() = #externalMacro(module: "MacroPlugin", type: "TypeWrapperMacro")

@attached(accessor)
public macro myPropertyWrapper() = #externalMacro(module: "MacroPlugin", type: "PropertyWrapperMacro")

@attached(extension, conformances: Equatable)
public macro DelegatedConformance() = #externalMacro(module: "MacroPlugin", type: "DelegatedConformanceMacro")
