import http from 'k6/http';
import { DEFAULT_TARGET_RPS, jsonGetParams, resolveTargetRps } from './common.js';
import exec from 'k6/execution';

export const API_PREFIX = '/api/v1/quik';
export const WRITER_API_PREFIX = '/api/v1';

export const PATH_MONEY_LIMITS = '/money-limits';
export const PATH_SECURITY_LIMITS = '/security-limits';
export const PATH_OTC_SECURITY_LIMITS = '/otc-security-limits';
export const PATH_POSITIONS = '/positions';
export const PATH_MONEY_POSITIONS = '/money-positions';
export const PATH_SECURITY_POSITIONS = '/security-positions';
export const PATH_OTC_SECURITY_POSITIONS = '/otc-security-positions';

export const DEFAULT_SEC_CODES = [
  'GLD', 'EQMX', 'OBLG', 'DIVD', 'FXRW', 'BCSR',
  'SU26245RMFS9', 'SU26246RMFS7', 'F', 'C-RM',
];

export const LIMIT_READ_PATHS_TRIPLE = [
  PATH_MONEY_LIMITS,
  PATH_SECURITY_LIMITS,
  PATH_OTC_SECURITY_LIMITS,
];

export const DEFAULT_PORTFOLIO_BASE_URL = 'http://127.0.0.1:5030';
export const DEFAULT_WRITER_BASE_URL = 'http://127.0.0.1:5060';
export const DEFAULT_REFERENCE_BASE_URL = 'http://127.0.0.1:6030';

export const JSON_HEADERS = {
  'Content-Type': 'application/json',
  Accept: 'application/json',
};

export function writerHeaders(state, extra = {}) {
  const headers = { ...JSON_HEADERS, ...extra };
  const apiKey = (state && state.writerApiKey) || __ENV.WRITER_API_KEY || '';
  if (apiKey) {
    headers['X-API-KEY'] = apiKey;
  }
  return headers;
}

export function writerUploadHeaders(state, extra = {}) {
  const headers = { Accept: 'application/json', ...extra };
  const apiKey = (state && state.writerApiKey) || __ENV.WRITER_API_KEY || '';
  if (apiKey) {
    headers['X-API-KEY'] = apiKey;
  }
  return headers;
}

export function trimBaseUrl(url) {
  return String(url || '').replace(/\/+$/, '');
}

export function portfolioUrl(baseUrl, path) {
  return `${trimBaseUrl(baseUrl)}${API_PREFIX}${path}`;
}

export function writerUrl(baseUrl, path) {
  return `${trimBaseUrl(baseUrl)}${WRITER_API_PREFIX}${path}`;
}

export function referenceUrl(baseUrl, path) {
  return `${trimBaseUrl(baseUrl)}${API_PREFIX}${path}`;
}

export const DEFAULT_PAGE_LIMIT = 50;

export function todayISO() {
  const d = new Date();
  const y = d.getUTCFullYear();
  const m = String(d.getUTCMonth() + 1).padStart(2, '0');
  const day = String(d.getUTCDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

export function todayLocalISO() {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

function resolveIncludeDate(raw) {
  if (raw.include_date === true || raw.includeDate === true) {
    return true;
  }
  const mode = raw.date_mode || raw.dateMode || 'omit';
  return mode === 'include';
}

export function readQuikState(scenario) {
  const raw = JSON.parse(open(import.meta.resolve(`../state/${scenario}.json`)).replace(/^\uFEFF/, ''));
  const includeDate = resolveIncludeDate(raw);
  const clientCodes = raw.client_codes || raw.clientCodes || [];
  const primaryClientCode = raw.primary_client_code || raw.primaryClientCode
    || (clientCodes.length ? clientCodes[0] : '');
  return {
    portfolioBaseUrl: raw.portfolio_base_url || raw.portfolioBaseUrl || DEFAULT_PORTFOLIO_BASE_URL,
    writerBaseUrl: raw.writer_base_url || raw.writerBaseUrl || DEFAULT_WRITER_BASE_URL,
    referenceBaseUrl: raw.reference_base_url || raw.referenceBaseUrl || DEFAULT_REFERENCE_BASE_URL,
    targetRps: resolveTargetRps(raw.target_iter_rps || raw.targetIterRps || raw.target_rps || raw.targetRps),
    targetHttpRps: Number(raw.target_http_rps || raw.targetHttpRps || 0),
    targetIterRps: Number(raw.target_iter_rps || raw.targetIterRps || 0),
    pagesPerIter: Number(raw.pages_per_iter || raw.pagesPerIter || 0),
    date: raw.date || todayLocalISO(),
    includeDate,
    dateMode: raw.date_mode || raw.dateMode || 'omit',
    clientCodes,
    primaryClientCode,
    pageLimit: Number(raw.page_limit || raw.pageLimit || DEFAULT_PAGE_LIMIT),
    seedClientCount: Number(raw.seed_client_count || raw.seedClientCount || raw.prepare_client_count || 0),
    writerApiKey: raw.writer_api_key || raw.writerApiKey || __ENV.WRITER_API_KEY || '',
    writerSeedKinds: raw.writer_seed_kinds || raw.writerSeedKinds || [],
    secCodes: raw.sec_codes || raw.secCodes || DEFAULT_SEC_CODES,
    otcSecCodes: raw.otc_sec_codes || raw.otcSecCodes || DEFAULT_SEC_CODES,
    moneyCurrencies: raw.money_currencies || raw.moneyCurrencies || ['RUB'],
    firmIds: (raw.firm_ids || raw.firmIds || []).map(Number),
    firmCodes: raw.firm_codes || raw.firmCodes || [],
    prepareFirmCount: Number(raw.prepare_firm_count || raw.prepareFirmCount || 0),
    prepareClientCount: Number(raw.prepare_client_count || raw.prepareClientCount || 0),
    runId: raw.run_id || raw.runId || '',
  };
}

function buildQueryString(parts) {
  return parts
    .filter((p) => p.value !== undefined && p.value !== null && p.value !== '')
    .map((p) => `${encodeURIComponent(p.key)}=${encodeURIComponent(String(p.value))}`)
    .join('&');
}

function appendDateParts(parts, options = {}) {
  if (options.includeDate) {
    parts.push({ key: 'loadDate', value: options.date || todayLocalISO() });
  }
}

export function listQueryByClient(clientCode, options = {}) {
  const parts = [{ key: 'clientCodes', value: clientCode }];
  appendDateParts(parts, options);
  return buildQueryString(parts);
}

export function listQueryPaged(limit, offset, options = {}) {
  const parts = [
    { key: 'limit', value: limit },
    { key: 'offset', value: offset },
    { key: 'includeTotalCount', value: 'true' },
  ];
  appendDateParts(parts, options);
  return buildQueryString(parts);
}

export function listQueryLimitsAggregate(options = {}) {
  const parts = [];
  appendDateParts(parts, options);
  return buildQueryString(parts);
}

export function limitsPathWithQuery(path, query) {
  return query.length ? `${path}?${query}` : path;
}

export function portfolioQuery(options = {}) {
  const parts = [
    { key: 'targetCurrency', value: options.targetCurrency || 'RUB' },
  ];
  if (options.clientCodes && options.clientCodes.length) {
    const codes = Array.isArray(options.clientCodes)
      ? options.clientCodes.join(',')
      : String(options.clientCodes);
    parts.push({ key: 'clientCodes', value: codes });
  }
  return buildQueryString(parts);
}

export function positionsPathWithQuery(path, query) {
  return query.length ? `${path}?${query}` : path;
}

export function randomFrom(items) {
  return items[Math.floor(Math.random() * items.length)];
}

export function uniqueSuffix() {
  return `${__VU}-${__ITER}-${Math.random().toString(36).slice(2, 8)}`;
}

export function uniqueClientCode(prefix, vu, iter) {
  const core = `${vu}${iter}${Math.floor(Math.random() * 1000)}`;
  const code = `${prefix}${core}`.toUpperCase().replace(/[^A-Z0-9]/g, '');
  return code.slice(0, 12).padEnd(12, '0');
}

export function uniqueFirmCode(runId, vu, iter) {
  const suffix = `${runId}${vu}${iter}`.replace(/[^A-Z0-9]/gi, '').toUpperCase();
  const code = `LF${suffix}`;
  return code.slice(0, 12).padEnd(12, '0');
}

export function clientCodesQuery(codes) {
  if (!codes || !codes.length) {
    return '';
  }
  return `&clientCodes=${encodeURIComponent(codes.join(','))}`;
}

export function listQuery(date, limit, offset, clientCodes) {
  const parts = [
    { key: 'loadDate', value: date },
    { key: 'limit', value: limit },
    { key: 'offset', value: offset },
  ];
  if (clientCodes && clientCodes.length) {
    parts.push({ key: 'clientCodes', value: clientCodes.join(',') });
  }
  return buildQueryString(parts);
}

export function moneyLimitBody(firmCode, clientCode, balance = '100.00', currency = 'RUB', settleCode = 'Tx') {
  return JSON.stringify({
    limitType: 'money',
    clientCode,
    ticker: currency,
    positionCode: 'EQTV',
    settleCode,
    firmCode,
    balance,
  });
}

export function securityLimitBody(firmCode, clientCode, secCode, balance = '10', acquisitionCurrencyCode = 'RUB', settleCode = 'Tx') {
  return JSON.stringify({
    limitType: 'security',
    clientCode,
    ticker: secCode,
    tradeAccount: 'TRADEACC01',
    settleCode,
    firmCode,
    balance,
    acquisitionCurrencyCode,
  });
}

export function securityLimitOtcBody(firmCode, clientCode, secCode, balance = '5', acquisitionCurrencyCode = 'RUB', settleCode = 'Tx') {
  return JSON.stringify({
    limitType: 'security_otc',
    clientCode,
    ticker: secCode,
    settleCode,
    firmCode,
    balance,
    acquisitionCurrencyCode,
  });
}

export function firmCreateBody(firmCode, firmName) {
  return JSON.stringify({
    firmCode,
    firmName,
  });
}

export function firmPatchBody(firmName) {
  return JSON.stringify({
    firmName,
  });
}

export function abortIfUnreachable(response, context) {
  const unreachable = response.status === 0
    || (response.error && String(response.error).length > 0);
  if (!unreachable) {
    return;
  }
  console.error(`[k6] service unreachable, abort test: ${context} (status=${response.status})`);
  exec.test.abort(`service unreachable: ${context}`);
}

export function fetchAllLimitPages(baseUrl, path, pageLimit, queryOptions, onPage) {
  let offset = 0;
  let page = 0;
  while (true) {
    const query = listQueryPaged(pageLimit, offset, queryOptions);
    const url = portfolioUrl(baseUrl, limitsPathWithQuery(path, query));
    const response = http.get(url, jsonGetParams({
      tags: { name: `GET ${API_PREFIX}${path}`, page: String(page) },
    }));
    abortIfUnreachable(response, `GET ${url}`);
    const stop = onPage(response, page, offset);
    if (stop) {
      break;
    }

    let totalCount = 0;
    let limitsLen = 0;
    try {
      const body = response.json();
      totalCount = Number(body.totalCount);
      limitsLen = Array.isArray(body.limits) ? body.limits.length : 0;
    } catch (_) {
      break;
    }

    offset += pageLimit;
    page += 1;
    if (limitsLen === 0 || offset >= totalCount) {
      break;
    }
  }
}

export function allocateMixedRates(targetRps) {
  const positions = Math.min(targetRps, Math.round(targetRps * 0.55));
  let remaining = targetRps - positions;

  const getFirms = remaining > 0 ? Math.min(remaining, Math.round(targetRps * 0.12)) : 0;
  remaining -= getFirms;

  const getFirm = remaining > 0 ? Math.min(remaining, Math.max(1, Math.round(targetRps * 0.05))) : 0;
  remaining -= getFirm;

  const postMoney = remaining > 0 ? Math.min(remaining, Math.round(targetRps * 0.10)) : 0;
  remaining -= postMoney;

  const postSecurities = remaining > 0 ? Math.min(remaining, Math.round(targetRps * 0.10)) : 0;
  remaining -= postSecurities;

  const postSecuritiesOtc = remaining > 0 ? Math.min(remaining, Math.round(targetRps * 0.04)) : 0;
  remaining -= postSecuritiesOtc;

  const postFirm = remaining > 0 ? Math.min(remaining, Math.round(targetRps * 0.02)) : 0;
  remaining -= postFirm;

  const patchFirm = remaining;

  return {
    positions,
    getFirms,
    getFirm,
    postMoney,
    postSecurities,
    postSecuritiesOtc,
    postFirm,
    patchFirm,
  };
}
