Set-Location 'c:\Users\joran\Desarrollos\TramaTex\apps\tramatex-api'

$srcFile = 'internal\sales\application\sales_service.go'
$lines = Get-Content $srcFile
$total = $lines.Count
Write-Host "Loaded $total lines"

$assign = @{
    '106'='sales'; '131'='sales'; '136'='sales'; '141'='sales'; '146'='sales'
    '152'='sales'; '162'='sales'
    '191'='quote'; '260'='quote'; '314'='quote'
    '366'='order'
    '418'='quote'; '442'='quote'; '457'='quote'; '489'='quote'
    '526'='order'; '627'='order'; '681'='order'; '735'='order'
    '766'='order'; '777'='order'; '789'='order'; '800'='order'; '848'='order'; '903'='order'
    '952'='delivery'
    '1066'='billing'
    '1204'='quote'; '1218'='quote'; '1227'='quote'
    '1249'='order'; '1264'='order'; '1275'='order'; '1313'='order'; '1335'='order'
    '1362'='delivery'; '1385'='delivery'; '1414'='delivery'
    '1439'='billing'; '1454'='billing'; '1478'='billing'
    '1515'='sales'; '1544'='order'
    '1565'='quote'; '1668'='order'; '1687'='order'
    '1795'='sales'; '1799'='sales'
    '1803'='quote'; '1817'='order'; '1831'='billing'
    '1848'='quote'; '1853'='order'; '1858'='delivery'; '1863'='billing'
    '1868'='order'; '1877'='order'; '1881'='order'
    '1890'='delivery'; '1907'='delivery'; '1923'='billing'
    '1944'='billing'; '1968'='billing'; '2017'='billing'; '2073'='billing'
    '2199'='sales'; '2210'='quote'; '2216'='order'; '2222'='delivery'; '2228'='billing'
}

$sortedStarts = [int[]]($assign.Keys | Sort-Object { [int]$_ })
$fileContent = @{
    sales   = [System.Collections.Generic.List[string]]::new()
    quote   = [System.Collections.Generic.List[string]]::new()
    order   = [System.Collections.Generic.List[string]]::new()
    delivery = [System.Collections.Generic.List[string]]::new()
    billing = [System.Collections.Generic.List[string]]::new()
}

for ($i=0; $i -lt $sortedStarts.Count; $i++) {
    $start = $sortedStarts[$i]
    $end = if ($i+1 -lt $sortedStarts.Count) { $sortedStarts[$i+1]-1 } else { $total }
    $target = $assign[$start.ToString()]
    if ($null -eq $target) { Write-Host "WARNING: no target for start=$start"; continue }
    for ($j=$start-1; $j -le $end-1; $j++) { $fileContent[$target].Add($lines[$j]) }
    $fileContent[$target].Add('')
}

Write-Host "quote: $($fileContent['quote'].Count), order: $($fileContent['order'].Count), delivery: $($fileContent['delivery'].Count), billing: $($fileContent['billing'].Count), sales: $($fileContent['sales'].Count)"

$allImports = @(
    'package application',
    '',
    'import (',
    '	"context"',
    '	"fmt"',
    '	"strings"',
    '	"time"',
    '',
    '	"github.com/google/uuid"',
    '	pricing_app "github.com/joran-cortez/tramatex/internal/pricing/application"',
    '	"github.com/joran-cortez/tramatex/internal/sales/domain"',
    ')',
    ''
)

$fileMap = @{
    quote    = 'internal\sales\application\quote_service.go'
    order    = 'internal\sales\application\order_service.go'
    delivery = 'internal\sales\application\delivery_note_service.go'
    billing  = 'internal\sales\application\billing_service.go'
}

foreach ($key in $fileMap.Keys) {
    $outFile = $fileMap[$key]
    $all = $allImports + $fileContent[$key].ToArray()
    [System.IO.File]::WriteAllLines($outFile, $all, [System.Text.UTF8Encoding]::new($false))
    Write-Host "Wrote $outFile ($($fileContent[$key].Count) lines of methods)"
}

# Rebuild sales_service.go: header (1-105) + sales-assigned functions
$salesHeader = $lines[0..104]  # 0-indexed, lines 1-105
$newSalesLines = $salesHeader + @('') + $fileContent['sales'].ToArray()
[System.IO.File]::WriteAllLines($srcFile, $newSalesLines, [System.Text.UTF8Encoding]::new($false))
Write-Host "Rewrote $srcFile ($($newSalesLines.Count) lines)"
