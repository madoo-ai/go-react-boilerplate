import { describe, expect, it } from "vitest"
import { cn, initials } from "./utils"

describe("utils", () => {
  it("merges conflicting tailwind classes", () => {
    expect(cn("px-2", "px-4")).toBe("px-4")
    expect(cn("text-sm", false && "hidden", "font-medium")).toBe("text-sm font-medium")
  })

  it("builds initials from a display name", () => {
    expect(initials("Super Admin")).toBe("SA")
    expect(initials("  single   name ")).toBe("SN")
    expect(initials("")).toBe("")
  })
})
