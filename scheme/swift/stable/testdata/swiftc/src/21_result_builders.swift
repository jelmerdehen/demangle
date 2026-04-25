// SPDX-License-Identifier: Apache-2.0
// Feature: @resultBuilder, buildBlock

@resultBuilder
public struct ArrayBuilder<T> {
    public static func buildBlock(_ components: [T]...) -> [T] {
        return components.flatMap { $0 }
    }
    public static func buildExpression(_ expression: T) -> [T] {
        return [expression]
    }
    public static func buildOptional(_ component: [T]?) -> [T] {
        return component ?? []
    }
    public static func buildEither(first component: [T]) -> [T] {
        return component
    }
    public static func buildEither(second component: [T]) -> [T] {
        return component
    }
    public static func buildArray(_ components: [[T]]) -> [T] {
        return components.flatMap { $0 }
    }
}

@resultBuilder
public struct StringBuidler {
    public static func buildBlock(_ parts: String...) -> String {
        return parts.joined()
    }
    public static func buildOptional(_ part: String?) -> String {
        return part ?? ""
    }
    public static func buildEither(first component: String) -> String {
        return component
    }
    public static func buildEither(second component: String) -> String {
        return component
    }
}

public func buildArray<T>(@ArrayBuilder<T> content: () -> [T]) -> [T] {
    return content()
}

public func buildString(@StringBuidler content: () -> String) -> String {
    return content()
}

public func makeIntArray(include: Bool) -> [Int] {
    return buildArray {
        1
        2
        if include { 3 }
    }
}

public func makeGreeting(name: String, formal: Bool) -> String {
    return buildString {
        if formal { "Dear " } else { "Hi " }
        name
        "!"
    }
}
