param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [int]$TargetRps = 1000,
    [int]$PrepareFirmCount = 20,
    [int]$PrepareClientCount = 20
)

. (Join-Path $PSScriptRoot "_common.ps1")

$state = Build-QuikPrepareState @PSBoundParameters
Write-QuikState -State $state -Scenario "mixed"
