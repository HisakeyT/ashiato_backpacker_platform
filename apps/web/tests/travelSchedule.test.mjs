import assert from "node:assert/strict";
import test from "node:test";

import { scheduleLabel } from "../src/features/travel/scheduleLabel.ts";
import { sortPlans } from "../src/features/travelPlan/sortPlans.ts";
import {
  movePlan,
  persistChangedPlanOrder,
} from "../src/features/travelPlan/TravelPlansSection/planOrderChanges.ts";
import { formatPeriod } from "../src/features/travelPlan/format.ts";

const plans = [
  { id: 1, sortOrder: 2, startDate: "2026-10-01", endDate: null },
  { id: 2, sortOrder: 1, startDate: null, endDate: null },
  { id: 3, sortOrder: 3, startDate: "2026-09-30", endDate: null },
];

test("planning keeps manual order even when dates differ", () => {
  assert.deepEqual(sortPlans(plans, "planning").map((plan) => plan.id), [2, 1, 3]);
});

test("recorded sorts by dates", () => {
  assert.deepEqual(sortPlans(plans.filter((plan) => plan.startDate), "recorded").map((plan) => plan.id), [3, 1]);
});

test("plan order is persisted once on completion and skipped when unchanged", async () => {
  const current = sortPlans(plans, "planning");
  const draft = movePlan(current, current[0], 2);
  assert.deepEqual(current.map((plan) => plan.id), [2, 1, 3]);
  assert.deepEqual(draft.map((plan) => plan.id), [1, 3, 2]);

  const requests = [];
  const persist = async (ids) => { requests.push(ids); };
  assert.equal(await persistChangedPlanOrder(current, current, persist), false);
  assert.equal(requests.length, 0);
  assert.equal(await persistChangedPlanOrder(current, draft, persist), true);
  assert.deepEqual(requests, [[1, 3, 2]]);
});

test("undated plans and published modes have clear labels", () => {
  assert.equal(formatPeriod(null, null), "日付未定");
  assert.equal(scheduleLabel("planning"), "計画");
  assert.equal(scheduleLabel("recorded"), "記録");
});
