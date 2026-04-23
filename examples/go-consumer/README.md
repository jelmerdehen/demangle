# go-consumer example

Minimal Go program showing how to import and use the `demangle`
library. Run:

```
cd examples/go-consumer
go run .
```

Expected output:

```
  _ZN4llvm5Value4dumpEv                       [cpp-itanium] llvm::Value::dump()
  $s4main3FooV                                [swift-stable] main.Foo
  Java_com_example_Foo_bar                    [jni] com.example.Foo.bar
  _RNvCshIBIgx2Am2k_3std4open                 [rust] std::open
  ?foo@@YAXXZ                                 [cpp-msvc] void __cdecl foo(void)
  com.example.Foo$default                     [kotlin] com.example.Foo

Candidates for `_ZN4llvm5Value4dumpEv`:
  cpp-itanium (92)

Registered schemes:
  android-dex       java      dex-any
  cpp-itanium       cpp       itanium-abi
  …
```

## What's happening

- `import _ "github.com/jelmerdehen/demangle/scheme/all"` blank-
  imports every in-process scheme. Each subpackage has an `init()`
  that registers its `Scheme` value on `demangle.Default`.
- `demangle.Default.Demangle(ctx, input, nil)` auto-detects the
  scheme and dispatches.
- Per-scheme import is supported for minimal binaries: replace the
  `scheme/all` import with specific subpackages.
