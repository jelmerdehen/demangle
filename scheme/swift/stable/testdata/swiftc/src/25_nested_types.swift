// SPDX-License-Identifier: Apache-2.0
// Feature: nested struct inside struct/class/enum

public struct Matrix {
    public struct Row {
        public let values: [Double]
        public init(_ values: [Double]) { self.values = values }
        public var count: Int { return values.count }
        public subscript(i: Int) -> Double { return values[i] }
    }

    public let rows: [Row]
    public init(_ rows: [[Double]]) { self.rows = rows.map { Row($0) } }
    public var rowCount: Int { return rows.count }
    public var colCount: Int { return rows.first?.count ?? 0 }
}

public class Node<T> {
    public struct Metadata {
        public let id: Int
        public let label: String
        public init(id: Int, label: String) { self.id = id; self.label = label }
    }

    public let value: T
    public let meta: Metadata
    public var children: [Node<T>] = []

    public init(value: T, meta: Metadata) { self.value = value; self.meta = meta }
    public func addChild(_ node: Node<T>) { children.append(node) }
}

public enum Network {
    public struct Request {
        public let url: String
        public let method: Method
        public enum Method { case get, post, put, delete }
        public init(url: String, method: Method) { self.url = url; self.method = method }
    }

    public struct Response {
        public let statusCode: Int
        public let body: String?
        public init(statusCode: Int, body: String? = nil) {
            self.statusCode = statusCode; self.body = body
        }
        public var isSuccess: Bool { return (200..<300).contains(statusCode) }
    }

    case request(Request)
    case response(Response)
}
