import { check } from 'k6';
import { Rate } from 'k6/metrics';
import { RESPONSE_TIME_LIMIT_MS, responseJson } from './common.js';

export function createErrorRate(name = 'errors') {
  return new Rate(name);
}

export function runChecks(response, rules, errors) {
  const success = check(response, rules);
  errors.add(!success);
  return success;
}

export function standardGetChecks(expectedStatus = 200) {
  return {
    [`status is ${expectedStatus}`]: (r) => r.status === expectedStatus,
    [`response time < ${RESPONSE_TIME_LIMIT_MS}ms`]: (r) => r.timings.duration < RESPONSE_TIME_LIMIT_MS,
  };
}

export function isJsonArray(response) {
  try {
    const body = response.json();
    return Array.isArray(body);
  } catch (_) {
    return false;
  }
}

export function isLimitsList(response) {
  const limits = responseJson(response, 'limits');
  return Array.isArray(limits);
}

function isNonEmptyString(value) {
  return typeof value === 'string' && value.length > 0;
}

function isDecimalField(value) {
  if (value === null || value === undefined) {
    return false;
  }
  if (typeof value === 'number' && !Number.isNaN(value)) {
    return true;
  }
  return isNonEmptyString(value);
}

export function isPortfolioResponse(response) {
  try {
    const body = response.json();
    if (Array.isArray(body)) {
      return false;
    }
    if (!Array.isArray(body.positions)) {
      return false;
    }
    if (!isNonEmptyString(body.targetCurrency)) {
      return false;
    }
    if (!isDecimalField(body.marketValueTotal)) {
      return false;
    }
    if (body.positions.length === 0) {
      return true;
    }
    const item = body.positions[0];
    return isNonEmptyString(item.clientCode)
      && isNonEmptyString(item.firmCode)
      && isDecimalField(item.amount)
      && isDecimalField(item.marketValueInTargetCurrency);
  } catch (_) {
    return false;
  }
}

export function isLimitsPage(response) {
  const limits = responseJson(response, 'limits');
  const totalCount = responseJson(response, 'totalCount');
  return Array.isArray(limits) && typeof totalCount === 'number';
}

export function writerNoContentChecks(expectedStatus = 204) {
  return standardGetChecks(expectedStatus);
}

export function postCreatedChecks(entityField) {
  return {
    ...standardGetChecks(201),
    'response has entity field': (r) => {
      const value = responseJson(r, entityField);
      return value !== undefined && value !== null && String(value).length > 0;
    },
  };
}
