import http from 'k6/http';
import {
  buildOptions,
  DEFAULT_TARGET_RPS,
  resolveDurationPlan,
  resolveTestDurationSec,
} from '../lib/common.js';
import { createErrorRate, isJsonArray, runChecks, standardGetChecks } from '../lib/checks.js';
import { makeSummary } from '../lib/report.js';
import { readQuikState, referenceUrl } from '../lib/quik.js';

const SCENARIO = 'get_firms';
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
  const response = http.get(referenceUrl(STATE.referenceBaseUrl, '/firms/'), {
    headers: { Accept: 'application/json' },
    tags: { name: `GET ${ENDPOINT}` },
  });

  runChecks(response, {
    ...standardGetChecks(200),
    'response is json array': isJsonArray,
  }, errors);
}

export function handleSummary(data) {
  return makeSummary(data, {
    scenario: SCENARIO,
    scenarioId: 'get_firms_list',
    method: 'GET',
    endpoint: ENDPOINT,
    targetRps: TARGET_RPS,
    testDurationSec: TEST_DURATION_SEC,
    durationPlan: DURATION_PLAN,
    extraRows: [['reference base', STATE.referenceBaseUrl]],
  });
}
