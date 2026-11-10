import { buildLimitByClientScenario } from '../lib/limitRead.js';
import { API_PREFIX, PATH_SECURITY_LIMITS } from '../lib/quik.js';

const SCENARIO = 'get_security_limits_by_client';
const ENDPOINT = `${API_PREFIX}${PATH_SECURITY_LIMITS}`;
const built = buildLimitByClientScenario(SCENARIO, PATH_SECURITY_LIMITS, ENDPOINT);

export const options = built.options;
export default built.run;
export function handleSummary(data) {
  return built.summary(data);
}
