# Устарело: GET /portfolios заменён на GET /positions. Обёртка для совместимости.
Write-Warning "run_get_portfolio.ps1 uстарел; ispolzuetsya run_get_positions.ps1 (GET /api/v1/quik/positions)"
& (Join-Path $PSScriptRoot "run_get_positions.ps1") @PSBoundParameters
