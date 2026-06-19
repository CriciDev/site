$ErrorActionPreference = 'Stop'

$rootDir = Split-Path -Parent $PSScriptRoot
Set-Location $rootDir

$devProxyUpstream = if ($env:DEV_PROXY_UPSTREAM) { $env:DEV_PROXY_UPSTREAM } else { 'http://127.0.0.1:8080' }
$devProxyBind = if ($env:DEV_PROXY_BIND) { $env:DEV_PROXY_BIND } else { '127.0.0.1' }
$devProxyPort = if ($env:DEV_PROXY_PORT) { $env:DEV_PROXY_PORT } else { '7331' }

$air = Start-Process -FilePath 'air' -ArgumentList @('-c', '.air.toml') -NoNewWindow -PassThru
$templ = Start-Process -FilePath 'templ' -ArgumentList @(
	'generate',
	'-watch',
	"-proxy=$devProxyUpstream",
	"-proxybind=$devProxyBind",
	"-proxyport=$devProxyPort",
	'-watch-pattern=(.+\.go$)|(.+\.templ$)|(.+_templ\.txt$)|(.+\.css$)|(.+\.js$)|(.+\.(png|jpg|jpeg|svg|gif|webp)$)',
	'-ignore-pattern=(^|/)(tmp|data)(/|$)|.+_templ\.go$'
) -NoNewWindow -PassThru

$exitCode = 0
try {
	while ($true) {
		if ($air.HasExited) {
			$exitCode = $air.ExitCode
			break
		}

		if ($templ.HasExited) {
			$exitCode = $templ.ExitCode
			break
		}

		Start-Sleep -Seconds 1
	}
}
finally {
	foreach ($process in @($air, $templ)) {
		if ($process -and -not $process.HasExited) {
			Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
		}
	}
}

exit $exitCode
