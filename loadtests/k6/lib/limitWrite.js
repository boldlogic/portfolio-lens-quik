import http from 'k6/http';
import {
  moneyLimitBody,
  securityLimitBody,
  securityLimitOtcBody,
  uniqueClientCode,
  writerHeaders,
  writerUploadHeaders,
  writerUrl,
  randomFrom,
} from './quik.js';

export const PATH_WRITER_LIMITS = '/limits';
export const PATH_WRITER_LIMITS_UPLOAD = '/limits/upload';

export const CSV_HEADER = 'limit_type;client_code;ticker;settle_code;firm_code;balance;position_code;trade_account;acquisition_currency;isin';

export const LIMIT_KINDS = ['money', 'security', 'security_otc'];

const SETTLE_CODES = ['T0', 'T1', 'T2', 'Tx'];

export function randomLimitKind() {
  return randomFrom(LIMIT_KINDS);
}

export function randomSettleCode() {
  return randomFrom(SETTLE_CODES);
}

function moneyCsvRow(firmCode, clientCode, currency, settleCode) {
  return `money;${clientCode};${currency};${settleCode};${firmCode};100.50;EQTV;;;`;
}

function securityCsvRow(firmCode, clientCode, secCode, settleCode, acqCcy) {
  return `security;${clientCode};${secCode};${settleCode};${firmCode};10;;TRADEACC01;${acqCcy};`;
}

function securityOtcCsvRow(firmCode, clientCode, secCode, settleCode, acqCcy) {
  return `security_otc;${clientCode};${secCode};${settleCode};${firmCode};5;;;${acqCcy};`;
}

export function buildCsvContent(firmCode, rowCount, kind, secCodes, otcSecCodes, moneyCurrencies, vu, iter) {
  const lines = [CSV_HEADER];
  const settleCode = randomSettleCode();
  const count = rowCount === 2 ? 2 : 1;

  for (let i = 0; i < count; i += 1) {
    const prefix = kind === 'money' ? 'K6M' : (kind === 'security_otc' ? 'K6O' : 'K6S');
    const clientCode = uniqueClientCode(prefix, vu, `${iter}${i}`);
    if (kind === 'money') {
      const currency = randomFrom(moneyCurrencies);
      lines.push(moneyCsvRow(firmCode, clientCode, currency, settleCode));
    } else if (kind === 'security_otc') {
      const secCode = randomFrom(otcSecCodes);
      lines.push(securityOtcCsvRow(firmCode, clientCode, secCode, settleCode, 'RUB'));
    } else {
      const secCode = randomFrom(secCodes);
      lines.push(securityCsvRow(firmCode, clientCode, secCode, settleCode, 'RUB'));
    }
  }

  return lines.join('\n');
}

export function putLimitBody(firmCode, kind, secCodes, otcSecCodes, moneyCurrencies, vu, iter) {
  const clientCode = uniqueClientCode(
    kind === 'money' ? 'K6M' : (kind === 'security_otc' ? 'K6O' : 'K6S'),
    vu,
    iter,
  );
  const settleCode = randomSettleCode();

  if (kind === 'money') {
    const currency = randomFrom(moneyCurrencies);
    return moneyLimitBody(firmCode, clientCode, '100.50', currency, settleCode);
  }
  if (kind === 'security_otc') {
    const secCode = randomFrom(otcSecCodes);
    return securityLimitOtcBody(firmCode, clientCode, secCode, '5', 'RUB', settleCode);
  }
  const secCode = randomFrom(secCodes);
  return securityLimitBody(firmCode, clientCode, secCode, '10', 'RUB', settleCode);
}

export function putLimitRequest(state, firmCode, kind, secCodes, otcSecCodes, moneyCurrencies, tags) {
  const body = putLimitBody(firmCode, kind, secCodes, otcSecCodes, moneyCurrencies, __VU, __ITER);
  return http.put(
    writerUrl(state.writerBaseUrl, PATH_WRITER_LIMITS),
    body,
    {
      headers: writerHeaders(state),
      tags: tags || { name: `PUT /api/v1${PATH_WRITER_LIMITS}` },
    },
  );
}

export function uploadLimitsCsvRequest(state, firmCode, rowCount, kind, secCodes, otcSecCodes, moneyCurrencies, tags) {
  const csv = buildCsvContent(
    firmCode,
    rowCount,
    kind,
    secCodes,
    otcSecCodes,
    moneyCurrencies,
    __VU,
    __ITER,
  );
  return http.post(
    writerUrl(state.writerBaseUrl, PATH_WRITER_LIMITS_UPLOAD),
    {
      file: http.file(csv, 'limits.csv', 'text/csv'),
    },
    {
      headers: writerUploadHeaders(state),
      tags: tags || { name: `POST /api/v1${PATH_WRITER_LIMITS_UPLOAD}` },
    },
  );
}

export function runWriterMixedIteration(state, errors, runChecks, writerNoContentChecks, options = {}) {
  const firmCodes = state.firmCodes || [];
  const secCodes = state.secCodes || [];
  const otcSecCodes = state.otcSecCodes || [];
  const moneyCurrencies = state.moneyCurrencies || ['RUB'];
  const csvShare = options.csvShare ?? 0.5;

  const firmCode = randomFrom(firmCodes);
  const kind = randomLimitKind();

  let response;
  if (Math.random() < csvShare) {
    const rowCount = Math.random() < 0.5 ? 1 : 2;
    response = uploadLimitsCsvRequest(
      state,
      firmCode,
      rowCount,
      kind,
      secCodes,
      otcSecCodes,
      moneyCurrencies,
    );
  } else {
    response = putLimitRequest(
      state,
      firmCode,
      kind,
      secCodes,
      otcSecCodes,
      moneyCurrencies,
    );
  }

  runChecks(response, writerNoContentChecks(), errors);
}
