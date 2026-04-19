param(
	[string]$Owner,
	[string]$Repo,
	[string[]]$Branches = @("master", "develop", "staging"),
	[int]$RequiredApprovals = 1
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Assert-Command {
	param([string]$Name)
	if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
		throw "Missing required command: $Name"
	}
}

function Resolve-RepoFromOrigin {
	$originUrl = git remote get-url origin 2>$null
	if (-not $originUrl) {
		throw "Could not resolve 'origin' remote URL."
	}

	$pattern = 'github\.com[:/](?<owner>[^/]+)/(?<repo>[^/.]+)(\.git)?$'
	$match = [regex]::Match($originUrl, $pattern)
	if (-not $match.Success) {
		throw "Could not parse owner/repo from origin URL: $originUrl"
	}

	return @{
		Owner = $match.Groups["owner"].Value
		Repo  = $match.Groups["repo"].Value
	}
}

Assert-Command "gh"
Assert-Command "git"

gh auth status *> $null
if ($LASTEXITCODE -ne 0) {
	throw "GitHub CLI is not authenticated. Run: gh auth login -h github.com -s repo"
}

if (-not $Owner -or -not $Repo) {
	$resolved = Resolve-RepoFromOrigin
	if (-not $Owner) { $Owner = $resolved.Owner }
	if (-not $Repo) { $Repo = $resolved.Repo }
}

$payload = @{
	required_status_checks = @{
		strict   = $false
		contexts = @()
	}
	enforce_admins                = $true
	required_pull_request_reviews = @{
		dismiss_stale_reviews           = $false
		require_code_owner_reviews      = $false
		required_approving_review_count = $RequiredApprovals
		require_last_push_approval      = $false
	}
	restrictions      = $null
	allow_deletions   = $false
	allow_force_pushes = $false
}

$tmpFile = [System.IO.Path]::GetTempFileName()
try {
	# GitHub API rejects JSON with BOM in some environments, so write UTF-8 without BOM.
	$json = $payload | ConvertTo-Json -Depth 10
	[System.IO.File]::WriteAllText($tmpFile, $json, (New-Object System.Text.UTF8Encoding($false)))

	Write-Host "Applying branch protection in $Owner/$Repo" -ForegroundColor Cyan
	foreach ($branch in $Branches) {
		Write-Host "  -> $branch" -ForegroundColor Yellow
		gh api --method PUT "repos/$Owner/$Repo/branches/$branch/protection" `
			--input "$tmpFile" `
			-H "Accept: application/vnd.github+json" `
			-H "X-GitHub-Api-Version: 2022-11-28" *> $null

		if ($LASTEXITCODE -ne 0) {
			throw "Failed to protect branch '$branch'."
		}
	}

	Write-Host "Done. Protection applied to: $($Branches -join ', ')" -ForegroundColor Green
}
finally {
	if (Test-Path $tmpFile) {
		Remove-Item $tmpFile -Force
	}
}
