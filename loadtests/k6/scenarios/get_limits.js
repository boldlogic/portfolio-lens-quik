import { buildLimitTripleReadScenario } from '../lib/limitRead.js';

const SCENARIO = 'get_limits';
const built = buildLimitTripleReadScenario(SCENARIO);

export const options = built.options;
export default built.run;
export function handleSummary(data) {
  return built.summary(data);
}
