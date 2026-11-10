import http from 'k6/http';
import {
  buildOptions,
  DEFAULT_TARGET_RPS,
  resolveDurationPlan,
  resolveTestDurationSec,
} from '../lib/common.js';
import { createErrorRate, postCreatedChecks, runChecks } from '../lib/checks.js';
import { makeSummary } from '../lib/report.js';
import {
  firmCreateBody,
  JSON_HEADERS,
  readQuikState,
  referenceUrl,
  uniqueFirmCode,
} from '../lib/quik.js';

const SCENARIO = 'post_firm';
const STATE = readQuikState(SCENARIO);
const TARGET_RPS = Number(STATE.targetRps || DEFAULT_TARGET_RPS);
const TEST_DURATION_SEC = resolveTestDurationSec(STATE);
const DURATION_PLAN = resolveDurationPlan(TEST_DURATION_SEC);
const ENDPOINT = '/api/v1/quik/firms/';
const errors = createErrorRate();

export const options = buildOptions(TARGET_RPS, {}, {
  state: STATE,
  testDurationSec: TEST_DURATION_SEC,
});

export default function () {
  const firmCode = uniqueFirmCode(STATE.runId || 'RUN', __VU, __ITER);
  const response = http.post(
    referenceUrl(STATE.referenceBaseUrl, '/firms/'),
    firmCreateBody(firmCode, `Load Firm ${firmCode}`),
    { headers: JSON_HEADERS, tags: { name: `POST ${ENDPOINT}` } },
  );

  runChecks(response, postCreatedChecks('firmCode'), errors);
}

export function handleSummary(data) {
  return makeSummary(data, {
    scenario: SCENARIO,
    scenarioId: 'post_firm',
    method: 'POST',
    endpoint: ENDPOINT,
    targetRps: TARGET_RPS,
  });
}
