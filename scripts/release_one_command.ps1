param(
  [Parameter(Mandatory=$true)][string]$ModulePath,
  [Parameter(Mandatory=$true)][string]$Version,
  [string]$Remote = "origin",
  [string]$Branch = "main"
)

$ErrorActionPreference = "Stop"

if ($Version -notmatch '^v\d+\.\d+\.\d+$') {
  throw "Version must be semver format, for example v1.2.3"
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $repoRoot

& "$PSScriptRoot/bump_module_version.ps1" -ModulePath $ModulePath -Version $Version

git add tools/module_versions.json
$hasChanges = git diff --cached --name-only
if (-not $hasChanges) {
  Write-Host "No manifest changes to commit."
} else {
  git commit -m "chore(release): bump $ModulePath to $Version"
}

git push $Remote $Branch
& "$PSScriptRoot/release_module.ps1" -ModulePath $ModulePath -Version $Version -Remote $Remote -Branch $Branch

Write-Host "Release completed: $ModulePath $Version"
