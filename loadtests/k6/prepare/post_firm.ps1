param(
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [int]$TargetRps = 1000
)

. (Join-Path $PSScriptRoot "_common.ps1")
$state = @{
    portfolio_base_url  = $PortfolioBaseUrl.TrimEnd("/")
    writer_base_url     = $WriterBaseUrl.TrimEnd("/")
    reference_base_url  = $ReferenceDataBaseUrl.TrimEnd("/")
    target_rps          = $TargetRps
    date                = (Get-Date).ToString("yyyy-MM-dd")
    run_id              = [guid]::NewGuid().ToString("N").Substring(0, 8)
    firm_ids            = @()
    firm_codes          = @()
    client_codes        = @()
    prepare_firm_count  = 0
    prepare_client_count = 0
}
Write-QuikState -State $state -Scenario "post_firm"
