param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [string]$WriterApiKey = "",
    [int]$TargetRps = 1000,
    [int]$PrepareFirmCount = 5,
    [int]$PrepareClientCount = 1
)

. (Join-Path $PSScriptRoot "_common.ps1")
$state = Build-QuikPrepareState @PSBoundParameters -SkipWriterLimits -FixtureFirmsCodesOnly
Write-QuikState -State $state -Scenario "write_limits_mixed"
