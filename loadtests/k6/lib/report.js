import {
  escapeHtml,
  formatLocalFilenameTimestamp,
  formatLocalTimestamp,
  formatMs,
  formatPercent,
  formatRate,
  PRE_ALLOCATED_VUS,
  resolveDurationPlan,
  resolveTestDurationSec,
} from './common.js';

// Проверяет, что все k6 thresholds из summary завершились успешно.
// От этого зависит итоговый статус отчета: ПРОЙДЕН или ПРОВАЛЕН.
function allThresholdsPassed(data) {
  const results = [];
  Object.keys(data.metrics || {}).forEach((metricName) => {
    const thresholds = data.metrics[metricName]?.thresholds || {};
    Object.keys(thresholds).forEach((threshold) => {
      results.push(thresholds[threshold].ok);
    });
  });
  return results.length > 0 && results.every((ok) => ok);
}

function tableRows(rows) {
  return rows.map(([label, value]) => `
    <tr>
      <th>${escapeHtml(label)}</th>
      <td>${escapeHtml(value)}</td>
    </tr>`).join('');
}

function thresholdRows(data) {
  const rows = [];
  Object.keys(data.metrics || {}).forEach((metricName) => {
    const thresholds = data.metrics[metricName]?.thresholds || {};
    Object.keys(thresholds).forEach((threshold) => {
      rows.push([
        `${metricName} ${threshold}`,
        thresholds[threshold].ok ? 'пройден' : 'провален',
      ]);
    });
  });
  return rows;
}

function statCard(label, value, hint = '') {
  return `
    <section class="card">
      <div class="label">${escapeHtml(label)}</div>
      <div class="value">${escapeHtml(value)}</div>
      ${hint ? `<div class="hint">${escapeHtml(hint)}</div>` : ''}
    </section>`;
}

function resolveReportTargets(config) {
  const httpTargetRps = Number(config.httpTargetRps || 0);
  const iterTargetRps = Number(config.iterTargetRps || config.targetRps || 0);
  const pagesPerIter = Number(config.pagesPerIter || 0);
  return {
    httpTargetRps: httpTargetRps > 0 ? httpTargetRps : 0,
    iterTargetRps,
    pagesPerIter,
    filenameRps: httpTargetRps > 0 ? httpTargetRps : iterTargetRps,
  };
}

function resolveReportDuration(config) {
  if (config.durationPlan) {
    return config.durationPlan;
  }
  const totalSec = config.testDurationSec || resolveTestDurationSec({});
  return resolveDurationPlan(totalSec);
}

// Собирает метрики для отчёта: только HTTP и задержки (без iter/RPS-разбивок).
function summaryValues(data, plan, duration) {
  const durationSec = (data.state?.testRunDurationMs / 1000) || 0;
  const httpReqs = data.metrics?.http_reqs?.values || {};
  const droppedIterations = data.metrics?.dropped_iterations?.values || {};
  const vus = data.metrics?.vus?.values || {};
  const vusMax = data.metrics?.vus_max?.values || {};
  const httpReqFailed = data.metrics?.http_req_failed?.values || {};
  const httpReqDuration = data.metrics?.http_req_duration?.values || {};

  const totalHttpRequests = httpReqs.count || 0;
  const totalDroppedIterations = droppedIterations.count || 0;
  const failedHttpRequests = Math.round(totalHttpRequests * (httpReqFailed.rate || 0));
  const successfulHttpRequests = totalHttpRequests - failedHttpRequests;

  const httpTargetRps = plan.httpTargetRps;
  const iterTargetRps = plan.iterTargetRps;
  const targetRps = httpTargetRps > 0 ? httpTargetRps : iterTargetRps;

  const plannedHttpRequests = (0.5 * targetRps * duration.warmupSec) + (targetRps * duration.holdSec);
  const httpCompletionRate = plannedHttpRequests > 0 ? totalHttpRequests / plannedHttpRequests : 0;
  const httpRpsByWallTime = durationSec > 0 ? totalHttpRequests / durationSec : 0;

  return {
    durationSec,
    vus,
    vusMax,
    httpReqFailed,
    httpReqDuration,
    totalHttpRequests,
    totalDroppedIterations,
    failedHttpRequests,
    successfulHttpRequests,
    plannedHttpRequests,
    httpCompletionRate,
    httpRpsByWallTime,
    httpTargetRps,
    iterTargetRps,
  };
}

// Единая точка формирования HTML-отчетов.
// Каждый сценарий передает сюда только свои названия, endpoint и дополнительные строки.
export function makeSummary(data, config) {
  const now = new Date();
  const timestamp = formatLocalTimestamp(now);
  const dateStr = formatLocalFilenameTimestamp(now);
  const targets = resolveReportTargets(config);
  const duration = resolveReportDuration(config);
  const values = summaryValues(data, targets, duration);
  const goalHint = values.httpTargetRps > 0
    ? (targets.pagesPerIter > 1 ? `${targets.pagesPerIter} GET за проход` : '')
    : `${duration.warmupSec} с разогрев + ${duration.holdSec} с удержание`;
  const droppedHint = values.totalDroppedIterations > 0
    ? `не успели запустить: ${values.totalDroppedIterations}`
    : '';
  const factHttpHint = [
    values.httpTargetRps > 0 ? `${formatPercent(values.httpCompletionRate)} от плана HTTP` : 'по длительности теста',
    droppedHint,
  ].filter(Boolean).join(' · ');
  const thresholdsPassed = allThresholdsPassed(data);
  const resultId = thresholdsPassed ? 'pass' : 'fail';
  const resultRu = thresholdsPassed ? 'ПРОЙДЕН' : 'ПРОВАЛЕН';
  // Только относительный путь: k6 на Windows не пишет handleSummary в абсолютный K6_REPORTS_DIR.
  // run.ps1 делает Push-Location в корень репо перед k6 run.
  const outputDir = 'docs/reports/k6';
  const filename = `${outputDir}/load_${config.scenario}__dur${duration.totalSec}s__rps${targets.filenameRps}__${resultId}__${dateStr}.html`;
  const lastFilename = `${outputDir}/last_${config.scenario}.html`;
  const extraRows = config.extraRows || [];
  const extraCards = config.extraCards || [];

  const html = `<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <title>k6 ${escapeHtml(config.method)} ${escapeHtml(config.endpoint)} ${escapeHtml(resultRu)}</title>
  <style>
    :root { color-scheme: light; --ok:#0f766e; --bad:#b91c1c; --ink:#172033; --muted:#5b6475; --line:#d9dee8; --bg:#f6f8fb; --card:#fff; }
    * { box-sizing: border-box; }
    body { margin: 0; padding: 32px; background: var(--bg); color: var(--ink); font: 14px/1.45 system-ui, -apple-system, Segoe UI, sans-serif; }
    main { max-width: 1120px; margin: 0 auto; }
    header { display: flex; justify-content: space-between; gap: 24px; align-items: flex-start; margin-bottom: 24px; }
    h1 { margin: 0 0 8px; font-size: 28px; line-height: 1.15; }
    h2 { margin: 28px 0 12px; font-size: 18px; }
    .meta { color: var(--muted); }
    .badge { padding: 8px 12px; border-radius: 6px; color: #fff; background: ${thresholdsPassed ? 'var(--ok)' : 'var(--bad)'}; font-weight: 700; }
    .grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
    .card { background: var(--card); border: 1px solid var(--line); border-radius: 8px; padding: 14px; }
    .label { color: var(--muted); font-size: 12px; text-transform: uppercase; letter-spacing: .04em; }
    .value { margin-top: 6px; font-size: 24px; font-weight: 750; }
    .hint { margin-top: 4px; color: var(--muted); font-size: 12px; }
    table { width: 100%; border-collapse: collapse; background: var(--card); border: 1px solid var(--line); border-radius: 8px; overflow: hidden; }
    th, td { padding: 10px 12px; border-bottom: 1px solid var(--line); text-align: left; }
    th { width: 42%; color: var(--muted); font-weight: 600; }
    tr:last-child th, tr:last-child td { border-bottom: 0; }
    .split { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
    footer { margin-top: 24px; color: var(--muted); font-size: 12px; }
    @media (max-width: 820px) { body { padding: 18px; } header, .split { display: block; } .grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } .badge { display: inline-block; margin-top: 12px; } }
  </style>
</head>
<body>
<main>
  <header>
    <div>
      <h1>${escapeHtml(config.method)} ${escapeHtml(config.endpoint)}</h1>
      <div class="meta">Сценарий: ${escapeHtml(config.scenarioId)} · ${escapeHtml(timestamp)}</div>
    </div>
    <div class="badge">${escapeHtml(resultRu)}</div>
  </header>

  <div class="grid">
    ${values.httpTargetRps > 0
    ? statCard('Цель HTTP', formatRate(values.httpTargetRps, 'запросов/с'), goalHint)
    : statCard('Цель RPS', formatRate(values.iterTargetRps, 'запросов/с'), goalHint)}
    ${statCard('Факт HTTP RPS', formatRate(values.httpRpsByWallTime, 'запросов/с'), factHttpHint)}
    ${statCard('HTTP ошибки', String(values.failedHttpRequests), `rate: ${formatPercent(values.httpReqFailed.rate || 0)}`)}
    ${statCard('P95', formatMs(values.httpReqDuration['p(95)']), 'порог: < 500 мс')}
    ${statCard('P99', formatMs(values.httpReqDuration['p(99)']), 'порог: < 1000 мс')}
    ${statCard('Max', formatMs(values.httpReqDuration.max), 'максимальная задержка')}
    ${statCard('Активные VU', String(values.vus.max ?? 'нет данных'), `доступно: ${values.vusMax.max ?? PRE_ALLOCATED_VUS}`)}
    ${extraCards.map((card) => statCard(card[0], card[1], card[2])).join('')}
  </div>

  <div class="split">
    <section>
      <h2>Нагрузка</h2>
      <table><tbody>${tableRows([
        ['Endpoint', `${config.method} ${config.endpoint}`],
        ['Плановая длительность', `${duration.totalSec} с (${duration.warmupSec} разогрев + ${duration.holdSec} удержание)`],
        ['Фактическая длительность', `${values.durationSec.toFixed(2)} с`],
        ['Плановое число HTTP', String(Math.round(values.plannedHttpRequests))],
        ...extraRows,
      ])}</tbody></table>
    </section>

    <section>
      <h2>HTTP и задержки</h2>
      <table><tbody>${tableRows([
        ['HTTP-запросов всего', String(values.totalHttpRequests)],
        ['Успешных HTTP-запросов', String(values.successfulHttpRequests)],
        ['Неуспешных HTTP-запросов', String(values.failedHttpRequests)],
        ['Среднее время ответа', formatMs(values.httpReqDuration.avg)],
        ['Медиана', formatMs(values.httpReqDuration.med)],
        ['Минимум', formatMs(values.httpReqDuration.min)],
        ['Максимум', formatMs(values.httpReqDuration.max)],
      ])}</tbody></table>
    </section>
  </div>

  <h2>Пороговые условия</h2>
  <table><tbody>${tableRows(thresholdRows(data))}</tbody></table>

  <footer>${escapeHtml(config.footer || 'Подготовка данных выполняется до запуска k6 и не входит в метрики этого отчета.')}</footer>
</main>
</body>
</html>`;

  const lastHtml = `<!doctype html>
<html lang="ru">
<head><meta charset="utf-8"><title>Последний отчет ${escapeHtml(config.scenario)}</title></head>
<body>
  <h1>${escapeHtml(config.method)} ${escapeHtml(config.endpoint)}: ${escapeHtml(resultRu)}</h1>
  <p>${escapeHtml(timestamp)} · HTTP ${escapeHtml(formatRate(values.httpRpsByWallTime, 'запросов/с'))}${values.httpTargetRps > 0 ? ` (цель ${formatRate(values.httpTargetRps, 'запросов/с')})` : ''}</p>
  <ul>
    <li>HTTP-запросы: ${values.totalHttpRequests}</li>
    <li>HTTP-ошибки: ${values.failedHttpRequests}</li>
    <li>P95: ${escapeHtml(formatMs(values.httpReqDuration['p(95)']))}</li>
    <li>P99: ${escapeHtml(formatMs(values.httpReqDuration['p(99)']))}</li>
    <li>Max: ${escapeHtml(formatMs(values.httpReqDuration.max))}</li>
  </ul>
</body>
</html>`;

  return {
    [filename]: html,
    [lastFilename]: lastHtml,
  };
}
