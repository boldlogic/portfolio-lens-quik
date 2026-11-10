import { buildLimitPagedScenario } from '../lib/limitRead.js';
import { API_PREFIX, PATH_OTC_SECURITY_LIMITS } from '../lib/quik.js';

const SCENARIO = 'get_security_limits_otc_paged';
const ENDPOINT = `${API_PREFIX}${PATH_OTC_SECURITY_LIMITS}`;
const built = buildLimitPagedScenario(SCENARIO, PATH_OTC_SECURITY_LIMITS, ENDPOINT);

export const options = built.options;
export default built.run;
export function handleSummary(data) {
  return built.summary(data);
}
