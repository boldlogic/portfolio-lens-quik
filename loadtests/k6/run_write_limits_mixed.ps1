# Запуск: Run по ЭТОМУ файлу (.ps1), не scenarios/write_limits_mixed.js
param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [string]$WriterApiKey = "",
    [int]$TargetRps = 50,
    [int]$DurationSec = 120
)

$runArgs = @{
    Scenario             = "write_limits_mixed"
    PortfolioBaseUrl     = $PortfolioBaseUrl
    WriterBaseUrl        = $WriterBaseUrl
    ReferenceDataBaseUrl = $ReferenceDataBaseUrl
    WriterApiKey         = $WriterApiKey
    TargetRps            = $TargetRps
    DurationSec          = $DurationSec
}

& (Join-Path $PSScriptRoot "run.ps1") @runArgs
