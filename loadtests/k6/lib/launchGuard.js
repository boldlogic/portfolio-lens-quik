export function requireRunPs1Launch(scenarioFileHint) {
  if (__ENV.K6_LAUNCHED_FROM_RUN_PS1 === '1') {
    return;
  }
  if (__ENV.K6_HTTP_TARGET_RPS !== undefined && String(__ENV.K6_HTTP_TARGET_RPS).length > 0) {
    const fromEnv = Number(__ENV.K6_HTTP_TARGET_RPS);
    if (!Number.isNaN(fromEnv) && fromEnv > 0) {
      return;
    }
  }
  throw new Error(
    `Запускайте loadtests/k6/run_${scenarioFileHint}.ps1 (или .cmd), не k6 run scenarios/.... `
    + 'В Code Runner откройте .ps1/.cmd, не .js. run.ps1 передаёт TargetRps в k6.',
  );
}
