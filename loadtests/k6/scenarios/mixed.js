import http from 'k6/http';
import {
  DEFAULT_TARGET_RPS,
  resolveDurationPlan,
  resolveTestDurationSec,
  responseJson,
} from '../lib/common.js';
import {
  createErrorRate,
  isJsonArray,
  isPortfolioResponse,
  postCreatedChecks,
  runChecks,
  standardGetChecks,
  writerNoContentChecks,
} from '../lib/checks.js';
import { PATH_WRITER_LIMITS } from '../lib/limitWrite.js';
import { makeSummary } from '../lib/report.js';
import {
  allocateMixedRates,
  firmCreateBody,
  firmPatchBody,
  moneyLimitBody,
  API_PREFIX,
  JSON_HEADERS,
  PATH_POSITIONS,
  portfolioQuery,
  portfolioUrl,
  positionsPathWithQuery,
  randomFrom,
  readQuikState,
  referenceUrl,
  securityLimitBody,
  securityLimitOtcBody,
  uniqueClientCode,
  uniqueFirmCode,
  uniqueSuffix,
  writerHeaders,
  writerUrl,
} from '../lib/quik.js';

const SCENARIO = 'mixed';
const STATE = readQuikState(SCENARIO);
const TARGET_RPS = Number(STATE.targetRps || DEFAULT_TARGET_RPS);
const DURATION = resolveDurationPlan(resolveTestDurationSec(STATE));
const CLIENT_CODES = STATE.clientCodes || [];
const FIRM_IDS = STATE.firmIds || [];
const FIRM_CODES = STATE.firmCodes || [];
const RATES = allocateMixedRates(TARGET_RPS);
const errors = createErrorRate();

if (!CLIENT_CODES.length || !FIRM_IDS.length || !FIRM_CODES.length) {
  throw new Error('loadtests/k6/state/mixed.json must contain client_codes, firm_ids and firm_codes');
}

function rampingScenario(execName, targetRate) {
  return {
    executor: 'ramping-arrival-rate',
    exec: execName,
    startRate: 0,
    timeUnit: '1s',
    preAllocatedVUs: Math.max(10, Math.ceil(targetRate / 10)),
    maxVUs: Math.max(100, Math.ceil(targetRate * 2)),
    stages: [
      { target: targetRate, duration: `${DURATION.warmupSec}s` },
      { target: targetRate, duration: `${DURATION.holdSec}s` },
    ],
  };
}

function addScenario(scenarios, name, execName, targetRate) {
  if (targetRate > 0) {
    scenarios[name] = rampingScenario(execName, targetRate);
  }
}

const mixedScenarios = {};
addScenario(mixedScenarios, 'positions_flow', 'positionsFlow', RATES.positions);
addScenario(mixedScenarios, 'get_firms_flow', 'getFirmsFlow', RATES.getFirms);
addScenario(mixedScenarios, 'get_firm_flow', 'getFirmFlow', RATES.getFirm);
addScenario(mixedScenarios, 'post_money_flow', 'postMoneyFlow', RATES.postMoney);
addScenario(mixedScenarios, 'post_securities_flow', 'postSecuritiesFlow', RATES.postSecurities);
addScenario(mixedScenarios, 'post_securities_otc_flow', 'postSecuritiesOtcFlow', RATES.postSecuritiesOtc);
addScenario(mixedScenarios, 'post_firm_flow', 'postFirmFlow', RATES.postFirm);
addScenario(mixedScenarios, 'patch_firm_flow', 'patchFirmFlow', RATES.patchFirm);

export const options = {
  scenarios: mixedScenarios,
  thresholds: {
    http_req_duration: ['p(50)<150', 'p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
    checks: ['rate>0.99'],
    errors: ['rate<0.01'],
  },
  summaryTrendStats: ['avg', 'min', 'med', 'p(95)', 'p(99)', 'max'],
};

export function positionsFlow() {
  const query = portfolioQuery({ targetCurrency: 'RUB' });
  const path = positionsPathWithQuery(PATH_POSITIONS, query);
  const response = http.get(portfolioUrl(STATE.portfolioBaseUrl, path), {
    headers: { Accept: 'application/json' },
    tags: { name: `GET ${API_PREFIX}${PATH_POSITIONS}` },
  });
  runChecks(response, { ...standardGetChecks(200), 'response is portfolio object': isPortfolioResponse }, errors);
}

export function getFirmsFlow() {
  const response = http.get(referenceUrl(STATE.referenceBaseUrl, '/firms/'), {
    headers: { Accept: 'application/json' },
    tags: { name: 'GET /api/v1/quik/firms/' },
  });
  runChecks(response, { ...standardGetChecks(200), 'response is json array': isJsonArray }, errors);
}

export function getFirmFlow() {
  const firmId = randomFrom(FIRM_IDS);
  const response = http.get(referenceUrl(STATE.referenceBaseUrl, `/firms/${firmId}`), {
    headers: { Accept: 'application/json' },
    tags: { name: 'GET /api/v1/quik/firms/{id}' },
  });
  runChecks(response, {
    ...standardGetChecks(200),
    'response id matches request': (r) => Number(responseJson(r, 'id')) === firmId,
  }, errors);
}

export function postMoneyFlow() {
  const firmCode = randomFrom(FIRM_CODES);
  const clientCode = uniqueClientCode('K6M', __VU, __ITER);
  const response = http.put(
    writerUrl(STATE.writerBaseUrl, PATH_WRITER_LIMITS),
    moneyLimitBody(firmCode, clientCode),
    { headers: writerHeaders(STATE), tags: { name: `PUT /api/v1${PATH_WRITER_LIMITS}` } },
  );
  runChecks(response, writerNoContentChecks(), errors);
}

export function postSecuritiesFlow() {
  const firmCode = randomFrom(FIRM_CODES);
  const clientCode = uniqueClientCode('K6S', __VU, __ITER);
  const secCode = randomFrom(STATE.secCodes || []);
  const response = http.put(
    writerUrl(STATE.writerBaseUrl, PATH_WRITER_LIMITS),
    securityLimitBody(firmCode, clientCode, secCode),
    { headers: writerHeaders(STATE), tags: { name: `PUT /api/v1${PATH_WRITER_LIMITS}` } },
  );
  runChecks(response, writerNoContentChecks(), errors);
}

export function postSecuritiesOtcFlow() {
  const firmCode = randomFrom(FIRM_CODES);
  const clientCode = uniqueClientCode('K6O', __VU, __ITER);
  const secCode = randomFrom(STATE.otcSecCodes || []);
  const response = http.put(
    writerUrl(STATE.writerBaseUrl, PATH_WRITER_LIMITS),
    securityLimitOtcBody(firmCode, clientCode, secCode),
    { headers: writerHeaders(STATE), tags: { name: `PUT /api/v1${PATH_WRITER_LIMITS}` } },
  );
  runChecks(response, writerNoContentChecks(), errors);
}

export function postFirmFlow() {
  const firmCode = uniqueFirmCode(STATE.runId || 'RUN', __VU, __ITER);
  const response = http.post(
    referenceUrl(STATE.referenceBaseUrl, '/firms/'),
    firmCreateBody(firmCode, `Load Firm ${firmCode}`),
    { headers: { 'Content-Type': 'application/json', Accept: 'application/json' }, tags: { name: 'POST /api/v1/quik/firms/' } },
  );
  runChecks(response, postCreatedChecks('firmCode'), errors);
}

export function patchFirmFlow() {
  const firmId = randomFrom(FIRM_IDS);
  const response = http.patch(
    referenceUrl(STATE.referenceBaseUrl, `/firms/${firmId}`),
    firmPatchBody(`Load Patch ${uniqueSuffix()}`),
    { headers: JSON_HEADERS, tags: { name: 'PATCH /api/v1/quik/firms/{id}' } },
  );
  runChecks(response, {
    ...standardGetChecks(200),
    'response id matches request': (r) => Number(responseJson(r, 'id')) === firmId,
  }, errors);
}

export function handleSummary(data) {
  const readRps = RATES.positions + RATES.getFirms + RATES.getFirm;
  const writeRps = RATES.postMoney + RATES.postSecurities + RATES.postSecuritiesOtc + RATES.postFirm + RATES.patchFirm;
  return makeSummary(data, {
    scenario: SCENARIO,
    scenarioId: 'mixed_quik_http',
    method: 'MIXED',
    endpoint: 'positions read + writer POST + reference-data',
    targetRps: TARGET_RPS,
    testDurationSec: DURATION.totalSec,
    durationPlan: DURATION,
    extraRows: [
      ['GET positions', `${RATES.positions} итераций/с`],
      ['GET firms', `${RATES.getFirms} итераций/с`],
      ['GET firm', `${RATES.getFirm} итераций/с`],
      ['POST money', `${RATES.postMoney} итераций/с`],
      ['POST securities', `${RATES.postSecurities} итераций/с`],
      ['POST securities OTC', `${RATES.postSecuritiesOtc} итераций/с`],
      ['POST firm', `${RATES.postFirm} итераций/с`],
      ['PATCH firm', `${RATES.patchFirm} итераций/с`],
      ['client codes', String(CLIENT_CODES.length)],
      ['firm ids', String(FIRM_IDS.length)],
    ],
    extraCards: [
      ['Read доля', `${readRps} RPS`, 'positions + reference (без GET лимитов)'],
      ['Write доля', `${writeRps} RPS`, 'writer + reference write'],
    ],
    footer: 'GET *-limits в mixed не входят. Подготовка firms и limits выполняется до k6.',
  });
}
