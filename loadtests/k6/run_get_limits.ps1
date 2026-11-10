# Запуск: Run по ЭТОМУ файлу (.ps1), не scenarios/get_limits.js
# TargetRps = суммарный HTTP GET/s за итерацию (3 GET: money-limits + security-limits + otc-security-limits)
param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [int]$TargetRps = 1,
    [int]$DurationSec = 30,
    [int]$PrepareClientCount = 20,
    [string]$FirmsFixturePath = "",
    [switch]$SkipPrepare,
    [switch]$IncludeDate
)

$invokeArgs = @{
    Scenario           = "get_limits"
    TargetRps          = $TargetRps
    DurationSec        = $DurationSec
    PortfolioBaseUrl   = $PortfolioBaseUrl
    WriterBaseUrl      = $WriterBaseUrl
    PrepareClientCount = $PrepareClientCount
    FirmsFixturePath   = $FirmsFixturePath
    SkipPrepare        = $SkipPrepare
}
if ($IncludeDate) {
    $invokeArgs.IncludeDate = $true
}

& (Join-Path $PSScriptRoot "run_limit_test.ps1") @invokeArgs
