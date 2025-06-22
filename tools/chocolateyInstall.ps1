$toolsDir = "$(Split-Path -Parent $MyInvocation.MyCommand.Definition)"
$url = 'https://github.com/qwertycodeqc/locc/releases/download/v1.1/locc.exe'
$exe = Join-Path $toolsDir 'locc.exe'

Invoke-WebRequest -Uri $url -OutFile $exe
Install-BinFile -Name 'locc' -Path $exe
