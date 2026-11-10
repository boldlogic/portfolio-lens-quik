import http from 'k6/http';
import {
  buildOptions,
  jsonGetParams,
  resolveDurationPlan,
  resolveHttpTargetRps,
  resolveTargetRps,
  resolveTestDurationSec,
} from './common.js';
import { createErrorRate, isLimitsList, isLimitsPage, runChecks, standardGetChecks } from './checks.js';
import { makeSummary } from './report.js';
import {
  abortIfUnreachable,
  API_PREFIX,
  fetchAllLimitPages,
  LIMIT_READ_PATHS_TRIPLE,
  limitsPathWithQuery,
  listQueryByClient,
  portfolioUrl,
  readQuikState,
} from './quik.js';
import { requireRunPs1Launch } from './launchGuard.js';

function resolveLimitTargets(scenario) {
  const state = readQuikState(scenario);
  return {
    state,
    httpTargetRps: resolveHttpTargetRps(state),
    iterTargetRps: resolveTargetRps(state.targetIterRps || state.targetRps),
    pagesPerIter: Number(state.pages_per_iter || state.pagesPerIter || 0),
  };
}

export function buildLimitByClientScenario(scenario, limitsPath, apiPath) {
  requireRunPs1Launch(scenario);
  const { state: STATE, httpTargetRps: HTTP_TARGET_RPS, iterTargetRps: ITER_RPS } = resolveLimitTargets(scenario);
  const clientCode = STATE.primaryClientCode;
  const errors = createErrorRate();
  const queryOpts = { includeDate: STATE.includeDate, date: STATE.date };

  if (!clientCode) {
    throw new Error(`loadtests/k6/state/${scenario}.json must contain primary_client_code or client_codes`);
  }

  const testDurationSec = resolveTestDurationSec(STATE);
  const durationPlan = resolveDurationPlan(testDurationSec);
  const options = buildOptions(ITER_RPS, {}, {
    httpTargetRps: HTTP_TARGET_RPS > 0 ? HTTP_TARGET_RPS : 0,
    state: STATE,
    testDurationSec,
  });

  function run() {
    const query = listQueryByClient(clientCode, queryOpts);
    const url = portfolioUrl(STATE.portfolioBaseUrl, limitsPathWithQuery(limitsPath, query));
    const response = http.get(url, jsonGetParams({ tags: { name: `GET ${apiPath}` } }));
    abortIfUnreachable(response, `GET ${url}`);
    runChecks(response, {
      ...standardGetChecks(200),
      'response has limits list': isLimitsList,
    }, errors);
  }

  function summary(data) {
    return makeSummary(data, {
      scenario,
      scenarioId: scenario,
      method: 'GET',
      endpoint: `${apiPath}?clientCodes`,
      targetRps: HTTP_TARGET_RPS > 0 ? HTTP_TARGET_RPS : ITER_RPS,
      httpTargetRps: HTTP_TARGET_RPS,
      iterTargetRps: ITER_RPS,
      pagesPerIter: 1,
      testDurationSec,
      durationPlan,
      extraRows: [
        ['client code', clientCode],
        ['loadDate in query', STATE.includeDate ? STATE.date : '(omit)'],
        ['seed clients', String(STATE.seedClientCount || STATE.clientCodes.length)],
      ],
      footer: 'TargetRps v run_get_* = HTTP GET/s (1 GET/iter). Pri status 0 test ostanavlivaetsya.',
    });
  }

  return { options, run, summary };
}

export function buildLimitPagedScenario(scenario, limitsPath, apiPath) {
  requireRunPs1Launch(scenario);
  const {
    state: STATE,
    httpTargetRps: HTTP_TARGET_RPS,
    iterTargetRps: ITER_RPS,
    pagesPerIter,
  } = resolveLimitTargets(scenario);
  const errors = createErrorRate();
  const queryOpts = { includeDate: STATE.includeDate, date: STATE.date };
  const pageLimit = STATE.pageLimit;

  const testDurationSec = resolveTestDurationSec(STATE);
  const durationPlan = resolveDurationPlan(testDurationSec);
  const options = buildOptions(ITER_RPS, {}, {
    httpTargetRps: HTTP_TARGET_RPS > 0 ? HTTP_TARGET_RPS : 0,
    state: STATE,
    testDurationSec,
  });

  function run() {
    fetchAllLimitPages(
      STATE.portfolioBaseUrl,
      limitsPath,
      pageLimit,
      queryOpts,
      (response) => !runChecks(response, {
        ...standardGetChecks(200),
        'response has limits page': isLimitsPage,
      }, errors),
    );
  }

  function summary(data) {
    return makeSummary(data, {
      scenario,
      scenarioId: scenario,
      method: 'GET',
      endpoint: `${apiPath}?limit&offset`,
      targetRps: HTTP_TARGET_RPS > 0 ? HTTP_TARGET_RPS : ITER_RPS,
      httpTargetRps: HTTP_TARGET_RPS,
      iterTargetRps: ITER_RPS,
      pagesPerIter,
      testDurationSec,
      durationPlan,
      extraRows: [
        ['page limit', String(pageLimit)],
        ['GET за проход', pagesPerIter > 0 ? String(pagesPerIter) : '(probe)'],
        ['loadDate in query', STATE.includeDate ? STATE.date : '(omit)'],
        ['seed clients', String(STATE.seedClientCount || STATE.clientCodes.length)],
      ],
      footer: 'TargetRps v run_get_*_paged = HTTP GET/s; k6 iter/s = HTTP / pages. Pri status 0 test ostanavlivaetsya.',
    });
  }

  return { options, run, summary };
}

export function buildLimitTripleReadScenario(scenario) {
  requireRunPs1Launch(scenario);
  const { state: STATE, httpTargetRps: HTTP_TARGET_RPS, iterTargetRps: ITER_RPS, pagesPerIter } = resolveLimitTargets(scenario);
  const clientCode = STATE.primaryClientCode;
  const errors = createErrorRate();
  const queryOpts = { includeDate: STATE.includeDate, date: STATE.date };
  const httpPerIter = pagesPerIter > 0 ? pagesPerIter : LIMIT_READ_PATHS_TRIPLE.length;

  if (!clientCode) {
    throw new Error(`loadtests/k6/state/${scenario}.json must contain primary_client_code or client_codes`);
  }

  const testDurationSec = resolveTestDurationSec(STATE);
  const durationPlan = resolveDurationPlan(testDurationSec);
  const options = buildOptions(ITER_RPS, {}, {
    httpTargetRps: HTTP_TARGET_RPS > 0 ? HTTP_TARGET_RPS : 0,
    state: STATE,
    testDurationSec,
  });

  const readTargets = LIMIT_READ_PATHS_TRIPLE.map((path) => ({
    path,
    apiPath: `${API_PREFIX}${path}`,
  }));

  function run() {
    const query = listQueryByClient(clientCode, queryOpts);
    for (const { path, apiPath } of readTargets) {
      const url = portfolioUrl(STATE.portfolioBaseUrl, limitsPathWithQuery(path, query));
      const response = http.get(url, jsonGetParams({ tags: { name: `GET ${apiPath}` } }));
      abortIfUnreachable(response, `GET ${url}`);
      runChecks(response, {
        ...standardGetChecks(200),
        'response has limits list': isLimitsList,
      }, errors);
    }
  }

  function summary(data) {
    return makeSummary(data, {
      scenario,
      scenarioId: scenario,
      method: 'GET',
      endpoint: `${API_PREFIX}{money,security,otc}-limits?clientCodes`,
      targetRps: HTTP_TARGET_RPS > 0 ? HTTP_TARGET_RPS : ITER_RPS,
      httpTargetRps: HTTP_TARGET_RPS,
      iterTargetRps: ITER_RPS,
      pagesPerIter: httpPerIter,
      testDurationSec,
      durationPlan,
      extraRows: [
        ['client code', clientCode],
        ['GET за итерацию', String(httpPerIter)],
        ['loadDate in query', STATE.includeDate ? STATE.date : '(omit)'],
        ['seed clients', String(STATE.seedClientCount || STATE.clientCodes.length)],
      ],
      footer: 'TargetRps v run_get_limits.ps1 = HTTP GET/s (3 GET/iter: money + securities + otc). Pri status 0 test ostanavlivaetsya.',
    });
  }

  return { options, run, summary };
}
