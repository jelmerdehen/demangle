// Swift 6.2+ inline arrays
// InlineArray<N, T> is a fixed-size stack-allocated array
func sumInlineArray(_ arr: borrowing InlineArray<5, Int>) -> Int {
    var total = 0
    for i in 0..<5 { total += arr[i] }
    return total
}

func makeInlineArray() -> InlineArray<3, Float> {
    InlineArray<3, Float>(repeating: 0.0)
}

struct InlineArrayWrapper {
    var data: InlineArray<8, UInt8>
    var small: InlineArray<2, String>

    func first() -> UInt8 { data[0] }
    func count() -> Int { 8 }
}

func processFixed<let N: Int, T>(_ arr: InlineArray<N, T>) -> Int { N }
