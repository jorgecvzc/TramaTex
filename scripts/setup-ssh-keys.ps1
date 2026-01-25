# Generate SSH keys for pcele access
# Run: .\scripts\setup-ssh-keys.ps1

Write-Host "=== SSH Key Generation for TramaTex ===" -ForegroundColor Green
Write-Host ""

$sshDir = "$env:USERPROFILE\.ssh"
$keyPath = "$sshDir\id_rsa"

# Create .ssh directory if it doesn't exist
if (-not (Test-Path $sshDir)) {
    New-Item -ItemType Directory -Path $sshDir -Force | Out-Null
    Write-Host "Created directory: $sshDir"
}

# Check if keys already exist
if (Test-Path "$keyPath.pub") {
    Write-Host "SSH keys already exist at: $keyPath" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Your public key:"
    Write-Host "---"
    Get-Content "$keyPath.pub"
    Write-Host "---"
} else {
    Write-Host "Generating RSA 4096 keys..." -ForegroundColor Cyan
    
    # Generate keys
    ssh-keygen -t rsa -b 4096 -f $keyPath -N "" -C "tramatex@pcele"
    
    if (Test-Path "$keyPath.pub") {
        Write-Host ""
        Write-Host "SSH keys generated successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Your public key:"
        Write-Host "---"
        Get-Content "$keyPath.pub"
        Write-Host "---"
    } else {
        Write-Host "Error: Keys were not created" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "NEXT STEP: Add this key to pcele" -ForegroundColor Yellow
Write-Host ""
Write-Host "1. SSH to pcele and add the key:"
Write-Host "   ssh ele@pcele"
Write-Host "   mkdir -p ~/.ssh"
Write-Host "   cat >> ~/.ssh/authorized_keys"
Write-Host "   (paste the public key above, then Ctrl+D)"
Write-Host "   chmod 600 ~/.ssh/authorized_keys"
Write-Host "   exit"
Write-Host ""
Write-Host "After that, SSH will work without password prompts!"
