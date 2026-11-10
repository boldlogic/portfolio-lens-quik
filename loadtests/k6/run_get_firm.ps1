# Запуск: Run по ЭТОМУ файлу (.ps1), не scenarios/get_firm.js
param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [int]$TargetRps = 1000,
    [int]$DurationSec = 30
)

$runArgs = @{
    Scenario             = "get_firm"
    PortfolioBaseUrl     = $PortfolioBaseUrl
    WriterBaseUrl        = $WriterBaseUrl
    ReferenceDataBaseUrl = $ReferenceDataBaseUrl
    TargetRps            = $TargetRps
    DurationSec          = $DurationSec
    PrepareFirmCount     = 5
}

& (Join-Path $PSScriptRoot "run.ps1") @runArgs
