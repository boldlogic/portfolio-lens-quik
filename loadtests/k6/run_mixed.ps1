# Запуск: Run по ЭТОМУ файлу (.ps1), не scenarios/mixed.js
param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [int]$TargetRps = 1000,
    [int]$DurationSec = 30,
    [int]$PrepareFirmCount = 20,
    [int]$PrepareClientCount = 20
)

$runArgs = @{
    Scenario             = "mixed"
    PortfolioBaseUrl     = $PortfolioBaseUrl
    WriterBaseUrl        = $WriterBaseUrl
    ReferenceDataBaseUrl = $ReferenceDataBaseUrl
    TargetRps            = $TargetRps
    DurationSec          = $DurationSec
    PrepareFirmCount     = $PrepareFirmCount
    PrepareClientCount   = $PrepareClientCount
}

& (Join-Path $PSScriptRoot "run.ps1") @runArgs
