param(
    [Parameter(Mandatory=$true)]
    [string]$ProjectName,
    
    [Parameter(Mandatory=$true)]
    [string]$ModulePath
)

# Renkli output için yardımcı fonksiyonlar
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

Write-ColorOutput Green "=== Yeni Go Projesi Olusturuluyor ==="
Write-ColorOutput Cyan "Proje Adi: $ProjectName"
Write-ColorOutput Cyan "Modul Path: $ModulePath"

# Kaynak ve hedef klasorleri belirle
$sourceFolder = "go-starter-example-project-master"
$targetFolder = $ProjectName

# Kaynak klasorun var olup olmadigini kontrol et
if (-not (Test-Path $sourceFolder)) {
    Write-ColorOutput Red "HATA: Kaynak klasor bulunamadi: $sourceFolder"
    Write-ColorOutput Yellow "Lutfen scripti go-starter-example-project-master klasorunun oldugu dizinde calistirin."
    exit 1
}

# Hedef klasor varsa uyari ver
if (Test-Path $targetFolder) {
    Write-ColorOutput Yellow "UYARI: '$targetFolder' klasoru zaten var!"
    $response = Read-Host "Devam etmek istiyor musunuz? Mevcut klasor silinecek! (y/n)"
    if ($response -ne 'y') {
        Write-ColorOutput Red "Islem iptal edildi."
        exit 0
    }
    Remove-Item -Path $targetFolder -Recurse -Force
    Write-ColorOutput Green "Mevcut klasor silindi."
}

# Klasoru kopyala
Write-ColorOutput Cyan "`n[1/5] Proje dosyalari kopyalaniyor..."
Copy-Item -Path $sourceFolder -Destination $targetFolder -Recurse
Write-ColorOutput Green "Dosyalar kopyalandi"

# Eski modul path'ini bul (go.mod dosyasindan)
$goModPath = Join-Path $targetFolder "go.mod"
$oldModulePath = ""

if (Test-Path $goModPath) {
    $goModContent = Get-Content $goModPath -Raw
    if ($goModContent -match 'module\s+(\S+)') {
        $oldModulePath = $Matches[1]
        Write-ColorOutput Cyan "`nEski modul path: $oldModulePath"
    }
}

if ([string]::IsNullOrEmpty($oldModulePath)) {
    Write-ColorOutput Red "HATA: go.mod dosyasinda modul path bulunamadi!"
    exit 1
}

# Degistirilecek dosya uzantilari
$fileExtensions = @("*.go", "*.mod", "*.yaml", "*.yml", "*.json", "*.md", "*.html")

Write-ColorOutput Cyan "`n[2/5] go.mod dosyasi guncelleniyor..."
$goModContent = Get-Content $goModPath -Raw
$goModContent = $goModContent -replace [regex]::Escape($oldModulePath), $ModulePath
Set-Content -Path $goModPath -Value $goModContent -NoNewline
Write-ColorOutput Green "go.mod guncellendi"

Write-ColorOutput Cyan "`n[3/5] Import path'leri guncelleniyor..."
$fileCount = 0
Get-ChildItem -Path $targetFolder -Recurse -Include $fileExtensions | ForEach-Object {
    $file = $_
    $content = Get-Content $file.FullName -Raw -ErrorAction SilentlyContinue
    
    if ($null -ne $content -and $content.Length -gt 0) {
        $newContent = $content -replace [regex]::Escape($oldModulePath), $ModulePath
        
        if ($newContent -ne $content) {
            Set-Content -Path $file.FullName -Value $newContent -NoNewline
            $fileCount++
            Write-ColorOutput Gray "  $($file.Name)"
        }
    }
}
Write-ColorOutput Green "$fileCount dosya guncellendi"

Write-ColorOutput Cyan "`n[4/5] Go dependencies cozumleniyor..."
Push-Location $targetFolder
try {
    $output = & go mod tidy 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput Green "Dependencies cozumlendi"
    } else {
        Write-ColorOutput Yellow "go mod tidy uyarisi:"
        Write-ColorOutput Gray $output
    }
} catch {
    Write-ColorOutput Red "HATA: go mod tidy calistirilamadi"
    Write-ColorOutput Gray $_.Exception.Message
} finally {
    Pop-Location
}

Write-ColorOutput Cyan "`n[5/5] Swagger dokumanlari yeniden olusturuluyor..."
Push-Location $targetFolder
try {
    # Swagger'in yuklu olup olmadigini kontrol et
    $swagInstalled = & { go list -f '{{.Target}}' github.com/swaggo/swag/cmd/swag 2>&1 }
    
    if ($LASTEXITCODE -ne 0) {
        Write-ColorOutput Yellow "swag bulunamadi, yukleniyor..."
        & go install github.com/swaggo/swag/cmd/swag@latest
    }
    
    # Swagger dokumanlarini olustur
    $swagOutput = & swag init 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-ColorOutput Green "Swagger dokumanlari olusturuldu"
    } else {
        Write-ColorOutput Yellow "Swagger olusturma uyarisi:"
        Write-ColorOutput Gray $swagOutput
    }
} catch {
    Write-ColorOutput Yellow "Swagger dokumanlari olusturulamadi"
    Write-ColorOutput Gray $_.Exception.Message
} finally {
    Pop-Location
}

Write-ColorOutput Green "`n=== Proje basariyla olusturuldu! ==="
Write-ColorOutput Cyan "`nProje klasoru: $targetFolder"
Write-ColorOutput Cyan "Modul path: $ModulePath"

Write-ColorOutput Yellow "`nSonraki adimlar:"
Write-ColorOutput Gray "  1. cd $targetFolder"
Write-ColorOutput Gray "  2. .env dosyasini duzenleyin"
Write-ColorOutput Gray "  3. go run main.go"

Write-ColorOutput Green "`nIyi kodlamalar!"
