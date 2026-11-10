param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [int]$TargetRps = 1000,
    [int]$PrepareFirmCount = 0,
    [int]$PrepareClientCount = 20,
    [string]$FirmsFixturePath = "",
    [switch]$IncludeDate
)

. (Join-Path $PSScriptRoot "_common.ps1")
$state = Build-QuikPrepareState @PSBoundParameters -FixtureFirmsCodesOnly -WriterSeedKinds SecuritiesOtc
Write-QuikState -State $state -Scenario "get_security_limits_otc_by_client"
