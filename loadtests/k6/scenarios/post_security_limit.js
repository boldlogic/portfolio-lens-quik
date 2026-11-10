import http from 'k6/http';
import {
  buildOptions,
  DEFAULT_TARGET_RPS,
  resolveDurationPlan,
  resolveTestDurationSec,
} from '../lib/common.js';
import { createErrorRate, runChecks, writerNoContentChecks } from '../lib/checks.js';
import { PATH_WRITER_LIMITS } from '../lib/limitWrite.js';
import { makeSummary } from '../lib/report.js';
import {
  randomFrom,
  readQuikState,
  securityLimitBody,
  uniqueClientCode,
  writerHeaders,
  writerUrl,
} from '../lib/quik.js';

const SCENARIO = 'post_security_limit';
const STATE = readQuikState(SCENARIO);
const TARGET_RPS = Number(STATE.targetRps || DEFAULT_TARGET_RPS);
const TEST_DURATION_SEC = resolveTestDurationSec(STATE);
const DURATION_PLAN = resolveDurationPlan(TEST_DURATION_SEC);
const FIRM_CODES = STATE.firmCodes || [];
const SEC_CODES = STATE.secCodes || [];
const ENDPOINT = `/api/v1${PATH_WRITER_LIMITS}`;
const errors = createErrorRate();

if (!FIRM_CODES.length) {
  throw new Error(`loadtests/k6/state/${SCENARIO}.json must contain firm_codes`);
}

export const options = buildOptions(TARGET_RPS, {}, {
  state: STATE,
  testDurationSec: TEST_DURATION_SEC,
});

export default function () {
  const firmCode = randomFrom(FIRM_CODES);
  const clientCode = uniqueClientCode('K6S', __VU, __ITER);
  const secCode = randomFrom(SEC_CODES);
  const response = http.put(
    writerUrl(STATE.writerBaseUrl, PATH_WRITER_LIMITS),
    securityLimitBody(firmCode, clientCode, secCode),
    { headers: writerHeaders(STATE), tags: { name: `PUT ${ENDPOINT}` } },
  );

  runChecks(response, writerNoContentChecks(), errors);
}

export function handleSummary(data) {
  return makeSummary(data, {
    scenario: SCENARIO,
    scenarioId: 'post_security_limit',
    method: 'PUT',
    endpoint: ENDPOINT,
    targetRps: TARGET_RPS,
  });
}
