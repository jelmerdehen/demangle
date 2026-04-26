// SPDX-License-Identifier: Apache-2.0
// Feature: unicode identifier corpus for Swift stable remangler / Punycode encoder tests.
// Covers Latin extended, Greek, Cyrillic, CJK, Arabic, Hebrew, Devanagari, emoji, mixed.

public struct UnicodeCorpus {

    // ── Latin extended ──────────────────────────────────────────────────────
    public static func café() -> Int { return 1 }
    public static func résumé() -> String { return "" }
    public static func naïve() -> Bool { return true }
    public static func fiancée() -> Double { return 0.0 }
    public static func ñoño() -> Int { return 2 }
    public static func über() -> String { return "" }
    public static func façade() -> Bool { return false }
    public static func crème() -> String { return "" }
    public static func piñata() -> Int { return 3 }
    public static func jalapeño() -> String { return "" }

    // ── Greek ────────────────────────────────────────────────────────────────
    public static func αlpha() -> Int { return 10 }
    public static func βeta() -> Double { return 0.0 }
    public static func Δelta() -> Int { return 11 }
    public static func Σigma() -> Double { return 0.0 }
    public static func Ωmega() -> Int { return 12 }
    public static func γamma() -> Bool { return true }
    public static func λambda() -> Int { return 13 }
    public static func πi() -> Double { return 3.14159 }

    // ── Cyrillic ─────────────────────────────────────────────────────────────
    public static func функция() -> Int { return 20 }
    public static func переменная() -> String { return "" }
    public static func результат() -> Bool { return true }
    public static func данные() -> Double { return 0.0 }
    public static func список() -> Int { return 21 }
    public static func словарь() -> String { return "" }

    // ── CJK ──────────────────────────────────────────────────────────────────
    public static func 计算() -> Int { return 30 }
    public static func 函数() -> String { return "" }
    public static func 变量名() -> Bool { return false }
    public static func 数据() -> Double { return 0.0 }
    public static func 结果() -> Int { return 31 }
    public static func 列表() -> String { return "" }
    public static func 字典() -> Bool { return true }

    // ── Arabic ───────────────────────────────────────────────────────────────
    public static func متغير() -> Int { return 40 }
    public static func دالة() -> String { return "" }
    public static func نتيجة() -> Bool { return true }

    // ── Hebrew ───────────────────────────────────────────────────────────────
    public static func פונקציה() -> Int { return 50 }
    public static func משתנה() -> String { return "" }
    public static func תוצאה() -> Bool { return false }

    // ── Devanagari ───────────────────────────────────────────────────────────
    public static func फंक्शन() -> Int { return 60 }
    public static func चर() -> String { return "" }
    public static func परिणाम() -> Bool { return true }

    // ── Mathematical symbols (letter-like) ───────────────────────────────────
    public static func ℕatural() -> Int { return 70 }
    public static func ℤinteger() -> Int { return 71 }
    public static func ℝeal() -> Double { return 0.0 }

    // ── Emoji identifiers ────────────────────────────────────────────────────
    public static func 🎉celebrate() -> Int { return 80 }
    public static func 🚀launch() -> Bool { return true }
    public static func 💡idea() -> String { return "" }
    public static func 🔑key() -> Int { return 81 }
    public static func 🎯target() -> Bool { return false }

    // ── Mixed ASCII + Unicode ────────────────────────────────────────────────
    public static func user名前() -> String { return "" }
    public static func dataДата() -> Int { return 90 }
    public static func func関数() -> Bool { return true }
    public static func valueВалор() -> Double { return 0.0 }
    public static func item项目() -> Int { return 91 }
    public static func key键值() -> String { return "" }
}
