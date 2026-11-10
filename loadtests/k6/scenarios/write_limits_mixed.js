import {
  buildOptions,
  DEFAULT_TARGET_RPS,
  resolveDurationPlan,
  resolveTestDurationSec,
} from '../lib/common.js';
import { createErrorRate, runChecks, writerNoContentChecks } from '../lib/checks.js';
import { PATH_WRITER_LIMITS, runWriterMixedIteration } from '../lib/limitWrite.js';
import { makeSummary } from '../lib/report.js';
import { readQuikState } from '../lib/quik.js';

const SCENARIO = 'write_limits_mixed';
const STATE = readQuikState(SCENARIO);
const TARGET_RPS = Number(STATE.targetRps || DEFAULT_TARGET_RPS);
const TEST_DURATION_SEC = resolveTestDurationSec(STATE);
const DURATION_PLAN = resolveDurationPlan(TEST_DURATION_SEC);
const FIRM_CODES = STATE.firmCodes || [];
const ENDPOINT = `PUT/POST /api/v1/limits (50% PUT, 50% CSV 1-2 rows)`;
const errors = createErrorRate();

if (!FIRM_CODES.length) {
  throw new Error(`loadtests/k6/state/${SCENARIO}.json must contain firm_codes`);
}
if (!STATE.writerApiKey) {
  throw new Error(`writer_api_key not set: pass -WriterApiKey, env WRITER_API_KEY, or api_key in configs/quik-portfolio-writer-config.yaml`);
}

export const options = buildOptions(TARGET_RPS, {}, {
  state: STATE,
  testDurationSec: TEST_DURATION_SEC,
});

export default function () {
  runWriterMixedIteration(STATE, errors, runChecks, writerNoContentChecks, { csvShare: 0.5 });
}

export function handleSummary(data) {
  return makeSummary(data, {
    scenario: SCENARIO,
    scenarioId: 'write_limits_mixed',
    method: 'MIXED',
    endpoint: ENDPOINT,
    targetRps: TARGET_RPS,
    testDurationSec: TEST_DURATION_SEC,
    durationPlan: DURATION_PLAN,
    extraRows: [
      ['firm codes', String(FIRM_CODES.length)],
      ['csv share', '50%'],
      ['csv rows', '1-2 per upload'],
      ['put share', '50%'],
    ],
  });
}
