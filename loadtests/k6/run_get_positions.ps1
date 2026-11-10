# Запуск: Run по ЭТОМУ файлу (.ps1), не scenarios/get_positions.js
param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [int]$TargetRps = 50,
    [int]$DurationSec = 60,
    [int]$PrepareClientCount = 20,
    [string]$FirmsFixturePath = "",
    [switch]$SkipPrepare
)

$runArgs = @{
    Scenario             = "get_positions"
    PortfolioBaseUrl     = $PortfolioBaseUrl
    WriterBaseUrl        = $WriterBaseUrl
    ReferenceDataBaseUrl = $ReferenceDataBaseUrl
    TargetRps            = $TargetRps
    DurationSec          = $DurationSec
    PrepareClientCount   = $PrepareClientCount
    FirmsFixturePath       = $FirmsFixturePath
    SkipPrepare          = $SkipPrepare
}

& (Join-Path $PSScriptRoot "run.ps1") @runArgs
