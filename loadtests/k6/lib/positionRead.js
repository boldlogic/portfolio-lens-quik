import http from 'k6/http';
import {
  buildOptions,
  jsonGetParams,
  resolveDurationPlan,
  resolveTargetRps,
  resolveTestDurationSec,
} from './common.js';
import { createErrorRate, isPortfolioResponse, runChecks, standardGetChecks } from './checks.js';
import { makeSummary } from './report.js';
import {
  abortIfUnreachable,
  portfolioQuery,
  portfolioUrl,
  positionsPathWithQuery,
  readQuikState,
} from './quik.js';
import { requireRunPs1Launch } from './launchGuard.js';

export function buildPositionScenario(scenario, positionsPath, apiPath, queryOpts = {}) {
  requireRunPs1Launch(scenario);
  const STATE = readQuikState(scenario);
  const ITER_RPS = resolveTargetRps(STATE.targetRps);
  const errors = createErrorRate();
  const targetCurrency = queryOpts.targetCurrency || 'RUB';
  const query = portfolioQuery({
    targetCurrency,
    clientCodes: queryOpts.clientCodes,
  });
  const requestPath = positionsPathWithQuery(positionsPath, query);
  const endpointLabel = `${apiPath}?targetCurrency=${targetCurrency}`;

  const testDurationSec = resolveTestDurationSec(STATE);
  const durationPlan = resolveDurationPlan(testDurationSec);
  const options = buildOptions(ITER_RPS, {}, {
    state: STATE,
    testDurationSec,
  });

  function run() {
    const url = portfolioUrl(STATE.portfolioBaseUrl, requestPath);
    const response = http.get(url, jsonGetParams({ tags: { name: `GET ${apiPath}` } }));
    abortIfUnreachable(response, `GET ${url}`);
    runChecks(response, {
      ...standardGetChecks(200),
      'response is portfolio object': isPortfolioResponse,
    }, errors);
  }

  function summary(data) {
    return makeSummary(data, {
      scenario,
      scenarioId: scenario,
      method: 'GET',
      endpoint: endpointLabel,
      targetRps: ITER_RPS,
      testDurationSec,
      durationPlan,
      extraRows: [
        ['targetCurrency', targetCurrency],
        ['clientCodes in query', queryOpts.clientCodes ? String(queryOpts.clientCodes) : '(omit)'],
        ['seed clients', String(STATE.seedClientCount || STATE.clientCodes.length)],
      ],
      footer: 'TargetRps v run_get_*_positions.ps1 = k6 iter/s (1 GET/iter). Kotirovki market.* vruchnuyu.',
    });
  }

  return { options, run, summary };
}
