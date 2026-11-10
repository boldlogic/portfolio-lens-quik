# Запуск: Run по ЭТОМУ файлу (.ps1), не scenarios/post_security_limit_otc.js
param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [string]$WriterApiKey = "",
    [int]$TargetRps = 1000,
    [int]$DurationSec = 30
)

$runArgs = @{
    Scenario             = "post_security_limit_otc"
    PortfolioBaseUrl     = $PortfolioBaseUrl
    WriterBaseUrl        = $WriterBaseUrl
    ReferenceDataBaseUrl = $ReferenceDataBaseUrl
    WriterApiKey         = $WriterApiKey
    TargetRps            = $TargetRps
    DurationSec          = $DurationSec
}

& (Join-Path $PSScriptRoot "run.ps1") @runArgs
