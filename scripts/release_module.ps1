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

$tag = "$ModulePath/$Version"
git fetch --tags $Remote
if (-not (git tag --list $tag)) {
  git tag -a $tag -m "release $tag"
}
git push $Remote $Branch
git push $Remote $tag

Write-Host "Released $tag"
