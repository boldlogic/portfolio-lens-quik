# Запуск: Run по ЭТОМУ файлу (.ps1), не scenarios/get_money_positions.js
param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [int]$TargetRps = 100,
    [int]$DurationSec = 30,
    [int]$PrepareClientCount = 20,
    [string]$FirmsFixturePath = "",
    [switch]$SkipPrepare
)

$runArgs = @{
    Scenario           = "get_money_positions"
    PortfolioBaseUrl   = $PortfolioBaseUrl
    WriterBaseUrl      = $WriterBaseUrl
    TargetRps          = $TargetRps
    DurationSec        = $DurationSec
    PrepareClientCount = $PrepareClientCount
    FirmsFixturePath   = $FirmsFixturePath
    SkipPrepare        = $SkipPrepare
}

& (Join-Path $PSScriptRoot "run.ps1") @runArgs
