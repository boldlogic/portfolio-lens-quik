import http from 'k6/http';
import {
  buildOptions,
  DEFAULT_TARGET_RPS,
  resolveDurationPlan,
  resolveTestDurationSec,
} from '../lib/common.js';
import { createErrorRate, runChecks, standardGetChecks } from '../lib/checks.js';
import { makeSummary } from '../lib/report.js';
import {
  firmPatchBody,
  JSON_HEADERS,
  randomFrom,
  readQuikState,
  referenceUrl,
  uniqueSuffix,
} from '../lib/quik.js';
import { responseJson } from '../lib/common.js';

const SCENARIO = 'patch_firm';
const STATE = readQuikState(SCENARIO);
const TARGET_RPS = Number(STATE.targetRps || DEFAULT_TARGET_RPS);
const TEST_DURATION_SEC = resolveTestDurationSec(STATE);
const DURATION_PLAN = resolveDurationPlan(TEST_DURATION_SEC);
const FIRM_IDS = STATE.firmIds || [];
const ENDPOINT = '/api/v1/quik/firms/{id}';
const errors = createErrorRate();

if (!FIRM_IDS.length) {
  throw new Error(`loadtests/k6/state/${SCENARIO}.json must contain firm_ids`);
}

export const options = buildOptions(TARGET_RPS, {}, {
  state: STATE,
  testDurationSec: TEST_DURATION_SEC,
});

export default function () {
  const firmId = randomFrom(FIRM_IDS);
  const response = http.patch(
    referenceUrl(STATE.referenceBaseUrl, `/firms/${firmId}`),
    firmPatchBody(`Load Patch ${uniqueSuffix()}`),
    { headers: JSON_HEADERS, tags: { name: `PATCH ${ENDPOINT}` } },
  );

  runChecks(response, {
    ...standardGetChecks(200),
    'response id matches request': (r) => Number(responseJson(r, 'id')) === firmId,
  }, errors);
}

export function handleSummary(data) {
  return makeSummary(data, {
    scenario: SCENARIO,
    scenarioId: 'patch_firm_name',
    method: 'PATCH',
    endpoint: ENDPOINT,
    targetRps: TARGET_RPS,
    testDurationSec: TEST_DURATION_SEC,
    durationPlan: DURATION_PLAN,
    extraRows: [['firm ids', String(FIRM_IDS.length)]],
  });
}
