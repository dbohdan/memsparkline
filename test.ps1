Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$env:CGO_ENABLED = "0"

$extension = ""
if (((Get-Variable 'IsWindows' -Scope 'Global' -ErrorAction 'Ignore') -and
     $IsWindows) -or
    $env:OS -eq "Windows_NT") {
    $extension = ".exe"
}
$build = @{
    "memsparkline" = "main.go"
    "test/sleep" = "test/sleep.go"
}

foreach ($dest in $build.Keys) {
    $executable = "$dest$extension"
    $source = $build[$dest]

    Remove-Item -Force -ErrorAction Ignore $executable
    go build -o $executable $source
}

go test