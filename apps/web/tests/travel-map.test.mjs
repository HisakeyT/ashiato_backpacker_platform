import assert from "node:assert/strict";
import test from "node:test";
import { groupMapPlans, mapPlans } from "../src/features/travelPlan/mapPlans.ts";
import { formatPeriod } from "../src/features/travelPlan/format.ts";
const makePlan = (id, sortOrder, [latitude, longitude]) => ({ id, sortOrder, latitude, longitude });

test("sorts by manual order then ID, leaves input intact and preserves zero coordinates", () => {
  const input = [makePlan(4, 2, [10, 20]), makePlan(3, 1, [null, 0]), makePlan(2, 2, [0, 0])];
  const original = structuredClone(input);
  const { located, missing } = mapPlans(input);
  assert.deepEqual(located.map(({ plan, number }) => [plan.id, number]), [[2, 2], [4, 3]]);
  assert.deepEqual(missing.map(p => p.id), [3]);
  assert.deepEqual(input, original);
});
test("empty, single and invalid coordinates", () => {
  assert.deepEqual(mapPlans([]), { located: [], missing: [] });
  assert.equal(mapPlans([makePlan(1, 1, [0, 0])]).located.length, 1);
  const invalid = [makePlan(1, 1, [91, 0]), makePlan(2, 2, [0, -181]), makePlan(3, 3, [NaN, 0]), makePlan(4, 4, [0, Infinity]), makePlan(5, 5, [0, null])];
  assert.equal(mapPlans(invalid).missing.length, 5);
});
test("dateline crossings use adjacent world positions and co-located plans remain reachable", () => {
  const { located } = mapPlans([makePlan(1, 1, [10, 179]), makePlan(2, 2, [10, -179]), makePlan(3, 3, [10, 179])]);
  assert.equal(Math.abs(located[0].position[1] - located[1].position[1]), 2);
  assert.deepEqual(groupMapPlans(located).map(group => group.map(item => item.plan.id)), [[1, 3], [2]]);
});
test("refreshed data reflects edits, deletion and order", () => {
  const updated = mapPlans([makePlan(2, 1, [20, 30]), makePlan(1, 2, [0, 0])]).located;
  assert.deepEqual(updated.map(p => p.plan.id), [2, 1]);
  assert.deepEqual(mapPlans([makePlan(2, 1, [20, 30])]).located.map(p => p.plan.id), [2]);
});
test("detail periods handle undated, single day and multiple days", () => {
  assert.equal(formatPeriod(null, null), "日付未定");
  assert.equal(formatPeriod("2026-10-01", null), "2026-10-01");
  assert.equal(formatPeriod("2026-10-01", "2026-10-01"), "2026-10-01");
  assert.equal(formatPeriod("2026-10-01", "2026-10-03"), "2026-10-01 〜 2026-10-03");
});
