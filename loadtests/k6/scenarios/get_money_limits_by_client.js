import { buildLimitByClientScenario } from '../lib/limitRead.js';
import { API_PREFIX, PATH_MONEY_LIMITS } from '../lib/quik.js';

const SCENARIO = 'get_money_limits_by_client';
const ENDPOINT = `${API_PREFIX}${PATH_MONEY_LIMITS}`;
const built = buildLimitByClientScenario(SCENARIO, PATH_MONEY_LIMITS, ENDPOINT);

export const options = built.options;
export default built.run;
export function handleSummary(data) {
  return built.summary(data);
}
