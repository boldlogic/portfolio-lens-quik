export const DEFAULT_PORTFOLIO_BASE_URL = 'http://127.0.0.1:5030';
export const DEFAULT_WRITER_BASE_URL = 'http://127.0.0.1:5060';
export const DEFAULT_REFERENCE_BASE_URL = 'http://127.0.0.1:6030';
export const DEFAULT_TARGET_RPS = 1000;
export const DEFAULT_TEST_DURATION_SEC = 30;
export const DEFAULT_WARMUP_DURATION_SEC = 5;

function readSidecarDuration() {
  try {
    const value = Number(String(open(import.meta.resolve('../.last-test-duration-sec'))).trim());
    if (!Number.isNaN(value) && value >= 6) {
      return Math.floor(value);
    }
  } catch (_) {
    // sidecar создает run.ps1 перед k6 run
  }
  return 0;
}

// Приоритет: K6_TEST_DURATION_SEC, sidecar, state.test_duration_sec.
export function resolveTestDurationSec(state) {
  const env = __ENV.K6_TEST_DURATION_SEC;
  if (env !== undefined && String(env).length > 0) {
    const fromEnv = Number(env);
    if (!Number.isNaN(fromEnv) && fromEnv >= 6) {
      return Math.floor(fromEnv);
    }
  }
  const fromSidecar = readSidecarDuration();
  if (fromSidecar > 0) {
    return fromSidecar;
  }
  const raw = state || {};
  const fromState = Number(raw.test_duration_sec || raw.testDurationSec || 0);
  if (fromState >= 6) {
    return Math.floor(fromState);
  }
  return DEFAULT_TEST_DURATION_SEC;
}

// totalSec: полная длительность сценария; 5 с разогрев (или меньше при коротком тесте) + удержание.
export function resolveDurationPlan(totalSec) {
  const total = Math.max(6, Math.floor(Number(totalSec) || DEFAULT_TEST_DURATION_SEC));
  const warmupSec = Math.min(DEFAULT_WARMUP_DURATION_SEC, total - 1);
  const holdSec = total - warmupSec;
  return { totalSec: total, warmupSec, holdSec };
}
export const PRE_ALLOCATED_VUS = 250;
export const MAX_VUS = 2800;
export const RESPONSE_TIME_LIMIT_MS = 500;
export const HTTP_GET_TIMEOUT = '3s';
export const SCENARIO_GRACEFUL_STOP = '0s';

export function jsonGetParams(extra = {}) {
  return {
    headers: { Accept: 'application/json' },
    timeout: HTTP_GET_TIMEOUT,
    ...extra,
  };
}

function readSidecarRps(filename) {
  try {
    const value = Number(String(open(import.meta.resolve(`../${filename}`))).trim());
    if (!Number.isNaN(value) && value > 0) {
      return value;
    }
  } catch (_) {
    // sidecar создает run.ps1 перед k6 run
  }
  return 0;
}

// Приоритет: K6_TARGET_RPS, затем .last-k6-iter-rps, затем state.
export function resolveTargetRps(stateValue) {
  const env = __ENV.K6_TARGET_RPS;
  if (env !== undefined && String(env).length > 0) {
    const fromEnv = Number(env);
    if (!Number.isNaN(fromEnv) && fromEnv > 0) {
      return fromEnv;
    }
  }
  const fromSidecar = readSidecarRps('.last-k6-iter-rps');
  if (fromSidecar > 0) {
    return fromSidecar;
  }
  return Number(stateValue || DEFAULT_TARGET_RPS);
}

// Цель HTTP GET/с: K6_HTTP_TARGET_RPS, затем .last-http-target-rps, затем state.
export function resolveHttpTargetRps(state) {
  const env = __ENV.K6_HTTP_TARGET_RPS;
  if (env !== undefined && String(env).length > 0) {
    const fromEnv = Number(env);
    if (!Number.isNaN(fromEnv) && fromEnv > 0) {
      return fromEnv;
    }
  }
  const fromSidecar = readSidecarRps('.last-http-target-rps');
  if (fromSidecar > 0) {
    return fromSidecar;
  }
  const raw = state || {};
  const fromState = Number(raw.target_http_rps || raw.targetHttpRps || 0);
  if (fromState > 0) {
    return fromState;
  }
  return 0;
}

export function arrivalRateConfig(iterPerSec) {
  const rate = Math.max(0.001, Number(iterPerSec));
  if (rate >= 1 && Number.isInteger(rate)) {
    return { rate, timeUnit: '1s', effectiveIterPerSec: rate };
  }
  if (rate >= 1) {
    const scale = 10;
    const scaled = Math.round(rate * scale);
    return {
      rate: scaled,
      timeUnit: `${scale}s`,
      effectiveIterPerSec: scaled / scale,
    };
  }
  const periodSec = Math.max(1, Math.ceil(1 / rate));
  return { rate: 1, timeUnit: `${periodSec}s`, effectiveIterPerSec: 1 / periodSec };
}

function vuPoolForTarget(targetRps) {
  return {
    preAllocatedVUs: Math.max(2, Math.min(PRE_ALLOCATED_VUS, Math.ceil(targetRps * 2))),
    maxVUs: Math.max(10, Math.min(MAX_VUS, Math.ceil(targetRps * 4))),
  };
}

// Читает локальный state-файл, который заранее создает PowerShell prepare-скрипт.
// В state лежат base_url, target_rps и подготовленные id для конкретного сценария.
export function readState(scenario) {
  return JSON.parse(open(import.meta.resolve(`../state/${scenario}.json`)).replace(/^\uFEFF/, ''));
}

function buildThresholds(targetRps) {
  const latency = ['p(50)<150', 'p(95)<500', 'p(99)<1000'];
  const isLowRate = targetRps <= 20;
  return {
    http_req_duration: latency,
    http_req_failed: [{
      threshold: isLowRate ? 'rate<0.25' : 'rate<0.05',
      abortOnFail: true,
      delayAbortEval: isLowRate ? '12s' : '3s',
    }],
    checks: [isLowRate ? 'rate>0.90' : 'rate>0.99'],
    errors: [isLowRate ? 'rate<0.25' : 'rate<0.01'],
  };
}

// Limit read: быстрый abort по доле http_req_failed (test.abort не прерывает уже висящие http.get).
function buildLimitReadThresholds() {
  return {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: [{
      threshold: 'rate<0.01',
      abortOnFail: true,
      delayAbortEval: '0s',
    }],
    checks: ['rate>0.90'],
    errors: ['rate<0.05'],
  };
}

// iterPerSec: итерации k6/с. opts.httpTargetRps: для порогов/VU при paged (HTTP GET/с).
export function buildOptions(iterPerSec, extraThresholds = {}, opts = {}) {
  const httpTargetRps = Number(opts.httpTargetRps || 0);
  const thresholdBase = httpTargetRps > 0 ? httpTargetRps : iterPerSec;
  const vu = vuPoolForTarget(Math.max(thresholdBase, iterPerSec, 1));
  const arrival = arrivalRateConfig(iterPerSec);
  const duration = resolveDurationPlan(
    opts.testDurationSec || resolveTestDurationSec(opts.state || {}),
  );

  const useConstantArrival = iterPerSec < 1 || !Number.isInteger(iterPerSec);
  let stabilityTest;
  if (useConstantArrival) {
    stabilityTest = {
      executor: 'constant-arrival-rate',
      rate: arrival.rate,
      timeUnit: arrival.timeUnit,
      duration: `${duration.totalSec}s`,
      preAllocatedVUs: vu.preAllocatedVUs,
      maxVUs: vu.maxVUs,
      gracefulStop: SCENARIO_GRACEFUL_STOP,
    };
  } else {
    stabilityTest = {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: vu.preAllocatedVUs,
      maxVUs: vu.maxVUs,
      gracefulStop: SCENARIO_GRACEFUL_STOP,
      stages: [
        { target: iterPerSec, duration: `${duration.warmupSec}s` },
        { target: iterPerSec, duration: `${duration.holdSec}s` },
      ],
    };
  }

  const thresholds = httpTargetRps > 0
    ? buildLimitReadThresholds()
    : buildThresholds(thresholdBase);

  return {
    scenarios: {
      stability_test: stabilityTest,
    },
    thresholds: {
      ...thresholds,
      ...extraThresholds,
    },
    summaryTrendStats: ['avg', 'min', 'med', 'p(95)', 'p(99)', 'max'],
  };
}

// Безопасно достает поле из JSON-ответа. Если тело не JSON или поля нет, вернет undefined.
export function responseJson(response, path) {
  try {
    return response.json(path);
  } catch (_) {
    return undefined;
  }
}

export function escapeHtml(value) {
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

export function formatLocalTimestamp(date) {
  return [
    date.getFullYear(),
    String(date.getMonth() + 1).padStart(2, '0'),
    String(date.getDate()).padStart(2, '0'),
  ].join('-') + ' ' + [
    String(date.getHours()).padStart(2, '0'),
    String(date.getMinutes()).padStart(2, '0'),
    String(date.getSeconds()).padStart(2, '0'),
  ].join(':');
}

export function formatLocalFilenameTimestamp(date) {
  return [
    date.getFullYear(),
    String(date.getMonth() + 1).padStart(2, '0'),
    String(date.getDate()).padStart(2, '0'),
  ].join('-') + '_' + [
    String(date.getHours()).padStart(2, '0'),
    String(date.getMinutes()).padStart(2, '0'),
    String(date.getSeconds()).padStart(2, '0'),
  ].join('-');
}

export function formatMs(value) {
  return value === undefined ? 'нет данных' : `${value.toFixed(2)} мс`;
}

export function formatRate(value, unit) {
  return `${value.toFixed(2)} ${unit}`;
}

export function formatPercent(value) {
  return `${(value * 100).toFixed(2)}%`;
}
