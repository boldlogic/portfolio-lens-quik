import { buildPositionScenario } from '../lib/positionRead.js';
import { API_PREFIX, PATH_OTC_SECURITY_POSITIONS } from '../lib/quik.js';

const SCENARIO = 'get_otc_security_positions';
const ENDPOINT = `${API_PREFIX}${PATH_OTC_SECURITY_POSITIONS}`;
const built = buildPositionScenario(SCENARIO, PATH_OTC_SECURITY_POSITIONS, ENDPOINT);

export const options = built.options;
export default built.run;
export function handleSummary(data) {
  return built.summary(data);
}
