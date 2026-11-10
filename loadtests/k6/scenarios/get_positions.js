import { buildPositionScenario } from '../lib/positionRead.js';
import { API_PREFIX, PATH_POSITIONS } from '../lib/quik.js';

const SCENARIO = 'get_positions';
const ENDPOINT = `${API_PREFIX}${PATH_POSITIONS}`;
const built = buildPositionScenario(SCENARIO, PATH_POSITIONS, ENDPOINT);

export const options = built.options;
export default built.run;
export function handleSummary(data) {
  return built.summary(data);
}
