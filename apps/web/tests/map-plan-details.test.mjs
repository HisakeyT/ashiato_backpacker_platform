import assert from "node:assert/strict";
import test from "node:test";
import { formatPeriod } from "../src/features/travelPlan/format.ts";

test("detail periods handle undated, single day and multiple days", () => {
  assert.equal(formatPeriod(null, null), "日付未定");
  assert.equal(formatPeriod("2026-10-01", null), "2026-10-01");
  assert.equal(formatPeriod("2026-10-01", "2026-10-01"), "2026-10-01");
  assert.equal(formatPeriod("2026-10-01", "2026-10-03"), "2026-10-01 〜 2026-10-03");
});
