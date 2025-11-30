#!/bin/bash
#
# ╔═══════════════════════════════════════════════════════════════════════════╗
# ║                         SERVERPANEL INSTALLER                             ║
# ║                    Tek Komutla Tam Kurulum Scripti                        ║
# ╚═══════════════════════════════════════════════════════════════════════════╝
#
# Kullanım:
#   curl -sSL https://raw.githubusercontent.com/asergenalkan/serverpanel/main/install.sh | bash
#

set -e

# ═══════════════════════════════════════════════════════════════════════════════
# RENK TANIMLARI
# ═══════════════════════════════════════════════════════════════════════════════
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
WHITE='\033[1;37m'
NC='\033[0m'
BOLD='\033[1m'

# ═══════════════════════════════════════════════════════════════════════════════
# YAPILANDIRMA
# ═══════════════════════════════════════════════════════════════════════════════
VERSION="1.0.0"
INSTALL_DIR="/opt/serverpanel"
DATA_DIR="/var/lib/serverpanel"
LOG_DIR="/var/log/serverpanel"
GITHUB_REPO="asergenalkan/serverpanel"
GITHUB_RAW="https://raw.githubusercontent.com/${GITHUB_REPO}/main"

# Sayaçlar
STEP_CURRENT=0
STEP_TOTAL=12
ERRORS=0
WARNINGS=0

# ═══════════════════════════════════════════════════════════════════════════════
# YARDIMCI FONKSİYONLAR
# ═══════════════════════════════════════════════════════════════════════════════

print_banner() {
    clear
    echo -e "${CYAN}"
    cat << "EOF"
   ___                          ___                 _ 
  / __| ___ _ ___ _____ _ _    | _ \__ _ _ _  ___  | |
  \__ \/ -_) '_\ V / -_) '_|   |  _/ _` | ' \/ -_) | |
  |___/\___|_|  \_/\___|_|     |_| \__,_|_||_\___| |_|
                                                      
EOF
    echo -e "${WHITE}  ════════════════════════════════════════════════════${NC}"
    echo -e "${WHITE}     Web Hosting Control Panel - Kurulum v${VERSION}${NC}"
    echo -e "${WHITE}  ════════════════════════════════════════════════════${NC}"
    echo ""
}

log_step() {
    ((STEP_CURRENT++))
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  [${STEP_CURRENT}/${STEP_TOTAL}] ${WHITE}${BOLD}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

log_info() {
    echo -e "  ${GREEN}✓${NC} $1"
}

log_warn() {
    echo -e "  ${YELLOW}⚠${NC} $1"
    ((WARNINGS++))
}

log_error() {
    echo -e "  ${RED}✗${NC} $1"
    ((ERRORS++))
}

log_detail() {
    echo -e "    ${CYAN}→${NC} $1"
}

log_progress() {
    echo -ne "  ${MAGENTA}◌${NC} $1...\r"
}

log_done() {
    echo -e "  ${GREEN}●${NC} $1    "
}

spinner() {
    local pid=$1
    local delay=0.1
    local spinstr='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local temp=${spinstr#?}
        printf "  ${CYAN}%c${NC} İşleniyor...\r" "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
    done
    printf "                    \r"
}

check_command() {
    if command -v "$1" &> /dev/null; then
        log_info "$1 mevcut $(command -v $1)"
        return 0
    else
        log_detail "$1 bulunamadı, kurulacak"
        return 1
    fi
}

run_silent() {
    "$@" > /dev/null 2>&1
}

run_with_log() {
    local logfile="/tmp/serverpanel_install_$$.log"
    if "$@" >> "$logfile" 2>&1; then
        return 0
    else
        log_error "Komut başarısız: $*"
        log_error "Detaylar için: cat $logfile"
        return 1
    fi
}

# ═══════════════════════════════════════════════════════════════════════════════
# KONTROL FONKSİYONLARI
# ═══════════════════════════════════════════════════════════════════════════════

check_root() {
    log_step "Yetki Kontrolü"
    
    if [[ $EUID -ne 0 ]]; then
        log_error "Bu script root yetkisi gerektirir!"
        echo ""
        echo -e "  ${YELLOW}Kullanım:${NC}"
        echo -e "    ${WHITE}sudo bash install.sh${NC}"
        echo -e "    ${WHITE}curl -sSL ... | sudo bash${NC}"
        echo ""
        exit 1
    fi
    
    log_info "Root yetkisi doğrulandı"
}

check_os() {
    log_step "İşletim Sistemi Kontrolü"
    
    if [[ ! -f /etc/os-release ]]; then
        log_error "Desteklenmeyen işletim sistemi!"
        exit 1
    fi
    
    source /etc/os-release
    
    log_detail "Dağıtım: $NAME"
    log_detail "Sürüm: $VERSION_ID"
    log_detail "Kod Adı: $VERSION_CODENAME"
    
    if [[ "$ID" != "ubuntu" && "$ID" != "debian" ]]; then
        log_error "Sadece Ubuntu ve Debian desteklenmektedir"
        log_error "Tespit edilen: $ID"
        exit 1
    fi
    
    if [[ "$ID" == "ubuntu" ]]; then
        if [[ "${VERSION_ID%%.*}" -lt 20 ]]; then
            log_error "Ubuntu 20.04 veya üzeri gereklidir"
            exit 1
        fi
    fi
    
    log_info "İşletim sistemi uyumlu: $PRETTY_NAME"
    
    # Mimari kontrolü
    ARCH=$(uname -m)
    if [[ "$ARCH" != "x86_64" ]]; then
        log_warn "x86_64 dışı mimari: $ARCH - sorunlar olabilir"
    else
        log_info "Mimari: $ARCH"
    fi
}

check_resources() {
    log_step "Sistem Kaynakları Kontrolü"
    
    # RAM
    local total_ram=$(free -m | awk '/^Mem:/{print $2}')
    local free_ram=$(free -m | awk '/^Mem:/{print $7}')
    
    if [[ $total_ram -lt 512 ]]; then
        log_error "Minimum 512MB RAM gerekli (Mevcut: ${total_ram}MB)"
        exit 1
    elif [[ $total_ram -lt 1024 ]]; then
        log_warn "1GB+ RAM önerilir (Mevcut: ${total_ram}MB)"
    else
        log_info "RAM: ${total_ram}MB toplam, ${free_ram}MB kullanılabilir"
    fi
    
    # Disk
    local free_disk=$(df -m / | awk 'NR==2 {print $4}')
    local total_disk=$(df -m / | awk 'NR==2 {print $2}')
    
    if [[ $free_disk -lt 2048 ]]; then
        log_error "Minimum 2GB boş alan gerekli (Mevcut: ${free_disk}MB)"
        exit 1
    else
        log_info "Disk: ${free_disk}MB boş / ${total_disk}MB toplam"
    fi
    
    # CPU
    local cpu_cores=$(nproc)
    log_info "CPU: ${cpu_cores} çekirdek"
}

check_ports() {
    log_step "Port Kontrolü"
    
    local ports=(80 443 8443 3306 53)
    local port_names=("HTTP" "HTTPS" "Panel" "MySQL" "DNS")
    
    for i in "${!ports[@]}"; do
        local port=${ports[$i]}
        local name=${port_names[$i]}
        
        if ss -tuln | grep -q ":$port "; then
            local proc=$(ss -tulnp | grep ":$port " | awk '{print $7}' | cut -d'"' -f2 | head -1)
            log_warn "Port $port ($name) kullanımda: $proc"
        else
            log_info "Port $port ($name) müsait"
        fi
    done
}

# ═══════════════════════════════════════════════════════════════════════════════
# KURULUM FONKSİYONLARI
# ═══════════════════════════════════════════════════════════════════════════════

install_base_packages() {
    log_step "Temel Paketler Kuruluyor"
    
    log_progress "Paket listesi güncelleniyor"
    run_silent apt-get update
    log_done "Paket listesi güncellendi"
    
    local base_packages=(
        curl wget git unzip tar
        software-properties-common
        apt-transport-https
        ca-certificates
        gnupg lsb-release
        net-tools
    )
    
    log_progress "Temel paketler kuruluyor"
    DEBIAN_FRONTEND=noninteractive apt-get install -y "${base_packages[@]}" > /dev/null 2>&1
    log_done "Temel paketler kuruldu"
    
    for pkg in curl wget git; do
        check_command $pkg
    done
}

install_build_tools() {
    log_step "Derleme Araçları Kuruluyor"
    
    log_progress "build-essential kuruluyor"
    DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential > /dev/null 2>&1
    log_done "build-essential kuruldu"
    
    check_command gcc
    check_command make
    
    log_info "GCC sürümü: $(gcc --version | head -1)"
}

install_go() {
    log_step "Go Programlama Dili Kuruluyor"
    
    local GO_VERSION="1.21.5"
    local GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
    
    if command -v go &> /dev/null; then
        local current_version=$(go version | awk '{print $3}')
        log_info "Go zaten kurulu: $current_version"
    else
        log_progress "Go ${GO_VERSION} indiriliyor"
        wget -q "$GO_URL" -O /tmp/go.tar.gz
        log_done "Go indirildi"
        
        log_progress "Go kuruluyor"
        rm -rf /usr/local/go
        tar -C /usr/local -xzf /tmp/go.tar.gz
        rm /tmp/go.tar.gz
        log_done "Go kuruldu"
    fi
    
    # PATH'e ekle
    export PATH=$PATH:/usr/local/go/bin
    
    if ! grep -q "/usr/local/go/bin" /etc/profile; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    fi
    
    log_info "Go sürümü: $(go version | awk '{print $3}')"
}

install_nodejs() {
    log_step "Node.js Kuruluyor"
    
    if command -v node &> /dev/null; then
        local node_version=$(node --version)
        log_info "Node.js zaten kurulu: $node_version"
    else
        log_progress "NodeSource repository ekleniyor"
        curl -fsSL https://deb.nodesource.com/setup_20.x | bash - > /dev/null 2>&1
        log_done "Repository eklendi"
        
        log_progress "Node.js kuruluyor"
        DEBIAN_FRONTEND=noninteractive apt-get install -y nodejs > /dev/null 2>&1
        log_done "Node.js kuruldu"
    fi
    
    log_info "Node.js sürümü: $(node --version)"
    log_info "NPM sürümü: $(npm --version)"
}

install_webserver() {
    log_step "Apache Web Sunucusu Kuruluyor"
    
    local apache_packages=(
        apache2
        libapache2-mod-fcgid
    )
    
    log_progress "Apache kuruluyor"
    DEBIAN_FRONTEND=noninteractive apt-get install -y "${apache_packages[@]}" > /dev/null 2>&1
    log_done "Apache kuruldu"
    
    # Modülleri aktif et
    log_progress "Apache modülleri aktifleştiriliyor"
    a2enmod proxy_fcgi setenvif rewrite headers ssl expires > /dev/null 2>&1
    log_done "Modüller aktif"
    
    # Varsayılan siteyi devre dışı bırak
    a2dissite 000-default > /dev/null 2>&1 || true
    
    systemctl enable apache2 > /dev/null 2>&1
    systemctl restart apache2
    
    log_info "Apache durumu: $(systemctl is-active apache2)"
}

install_php() {
    log_step "PHP-FPM Kuruluyor"
    
    # PHP sürümünü belirle
    source /etc/os-release
    if [[ "$VERSION_ID" == "24.04" ]]; then
        PHP_VERSION="8.3"
    elif [[ "$VERSION_ID" == "22.04" ]]; then
        PHP_VERSION="8.1"
    else
        PHP_VERSION="8.1"
    fi
    
    log_detail "PHP sürümü: $PHP_VERSION"
    
    local php_packages=(
        php${PHP_VERSION}-fpm
        php${PHP_VERSION}-cli
        php${PHP_VERSION}-mysql
        php${PHP_VERSION}-curl
        php${PHP_VERSION}-gd
        php${PHP_VERSION}-mbstring
        php${PHP_VERSION}-xml
        php${PHP_VERSION}-zip
        php${PHP_VERSION}-intl
        php${PHP_VERSION}-bcmath
    )
    
    log_progress "PHP paketleri kuruluyor"
    DEBIAN_FRONTEND=noninteractive apt-get install -y "${php_packages[@]}" > /dev/null 2>&1
    log_done "PHP paketleri kuruldu"
    
    # Apache'ye PHP-FPM entegrasyonu
    a2enconf php${PHP_VERSION}-fpm > /dev/null 2>&1 || true
    
    systemctl enable php${PHP_VERSION}-fpm > /dev/null 2>&1
    systemctl restart php${PHP_VERSION}-fpm
    systemctl restart apache2
    
    log_info "PHP-FPM durumu: $(systemctl is-active php${PHP_VERSION}-fpm)"
}

install_mysql() {
    log_step "MySQL Veritabanı Kuruluyor"
    
    log_progress "MySQL Server kuruluyor"
    DEBIAN_FRONTEND=noninteractive apt-get install -y mysql-server > /dev/null 2>&1
    log_done "MySQL kuruldu"
    
    systemctl enable mysql > /dev/null 2>&1
    systemctl start mysql
    
    # Root şifresi oluştur
    local MYSQL_ROOT_PASS=$(openssl rand -base64 24 | tr -dc 'a-zA-Z0-9' | head -c 16)
    
    # Root şifresini ayarla
    mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '${MYSQL_ROOT_PASS}';" 2>/dev/null || true
    mysql -e "FLUSH PRIVILEGES;" 2>/dev/null || true
    
    # Şifreyi kaydet
    mkdir -p /root/.serverpanel
    echo "MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASS}" > /root/.serverpanel/mysql.conf
    chmod 600 /root/.serverpanel/mysql.conf
    
    log_info "MySQL durumu: $(systemctl is-active mysql)"
    log_info "Root şifresi: /root/.serverpanel/mysql.conf"
}

install_dns() {
    log_step "BIND DNS Sunucusu Kuruluyor"
    
    local dns_packages=(
        bind9
        bind9-utils
        bind9-dnsutils
    )
    
    log_progress "BIND kuruluyor"
    DEBIAN_FRONTEND=noninteractive apt-get install -y "${dns_packages[@]}" > /dev/null 2>&1
    log_done "BIND kuruldu"
    
    # Zone dizini oluştur
    mkdir -p /etc/bind/zones
    chown bind:bind /etc/bind/zones
    
    systemctl enable bind9 > /dev/null 2>&1
    systemctl start bind9
    
    log_info "BIND durumu: $(systemctl is-active bind9)"
}

install_certbot() {
    log_step "Let's Encrypt (Certbot) Kuruluyor"
    
    log_progress "Certbot kuruluyor"
    DEBIAN_FRONTEND=noninteractive apt-get install -y certbot python3-certbot-apache > /dev/null 2>&1
    log_done "Certbot kuruldu"
    
    log_info "Certbot sürümü: $(certbot --version 2>&1 | head -1)"
}

install_serverpanel() {
    log_step "ServerPanel Kuruluyor"
    
    # Dizinleri oluştur
    mkdir -p $INSTALL_DIR
    mkdir -p $DATA_DIR
    mkdir -p $LOG_DIR
    mkdir -p $INSTALL_DIR/public
    
    # Projeyi indir
    log_progress "Kaynak kod indiriliyor"
    if [[ -d "$INSTALL_DIR/.git" ]]; then
        cd $INSTALL_DIR
        git pull > /dev/null 2>&1
    else
        rm -rf $INSTALL_DIR/*
        git clone https://github.com/${GITHUB_REPO}.git $INSTALL_DIR > /dev/null 2>&1
    fi
    log_done "Kaynak kod indirildi"
    
    cd $INSTALL_DIR
    
    # Frontend build
    log_progress "Frontend derleniyor"
    cd web
    npm install --silent > /dev/null 2>&1
    npm run build > /dev/null 2>&1
    cp -r dist/* $INSTALL_DIR/public/
    cd $INSTALL_DIR
    log_done "Frontend derlendi"
    
    # Backend build
    log_progress "Backend derleniyor"
    export PATH=$PATH:/usr/local/go/bin
    CGO_ENABLED=1 go build -o serverpanel ./cmd/panel > /dev/null 2>&1
    log_done "Backend derlendi"
    
    log_info "Binary: $INSTALL_DIR/serverpanel"
    log_info "Frontend: $INSTALL_DIR/public/"
}

create_service() {
    log_step "Sistem Servisi Oluşturuluyor"
    
    # Sunucu IP'sini al
    local SERVER_IP=$(curl -s ifconfig.me 2>/dev/null || hostname -I | awk '{print $1}')
    
    # MySQL şifresini oku
    local MYSQL_PASS=""
    if [[ -f /root/.serverpanel/mysql.conf ]]; then
        source /root/.serverpanel/mysql.conf
        MYSQL_PASS=$MYSQL_ROOT_PASSWORD
    fi
    
    # PHP sürümünü belirle
    source /etc/os-release
    if [[ "$VERSION_ID" == "24.04" ]]; then
        PHP_VERSION="8.3"
    else
        PHP_VERSION="8.1"
    fi
    
    # Systemd service dosyası
    cat > /etc/systemd/system/serverpanel.service << EOF
[Unit]
Description=ServerPanel - Web Hosting Control Panel
Documentation=https://github.com/${GITHUB_REPO}
After=network.target mysql.service apache2.service

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/serverpanel
Restart=always
RestartSec=5
StandardOutput=append:${LOG_DIR}/panel.log
StandardError=append:${LOG_DIR}/error.log

# Environment
Environment="ENVIRONMENT=production"
Environment="PORT=8443"
Environment="MYSQL_ROOT_PASSWORD=${MYSQL_PASS}"
Environment="SERVER_IP=${SERVER_IP}"
Environment="PHP_VERSION=${PHP_VERSION}"
Environment="WEB_SERVER=apache"

[Install]
WantedBy=multi-user.target
EOF

    log_info "Service dosyası oluşturuldu"
    
    # Servisi başlat
    systemctl daemon-reload
    systemctl enable serverpanel > /dev/null 2>&1
    systemctl start serverpanel
    
    sleep 2
    
    if systemctl is-active --quiet serverpanel; then
        log_info "ServerPanel durumu: ${GREEN}aktif${NC}"
    else
        log_error "ServerPanel başlatılamadı!"
        log_error "Hata için: journalctl -u serverpanel -n 20"
    fi
}

configure_firewall() {
    log_step "Güvenlik Duvarı Yapılandırılıyor"
    
    if command -v ufw &> /dev/null; then
        log_progress "UFW yapılandırılıyor"
        
        ufw allow 22/tcp > /dev/null 2>&1    # SSH
        ufw allow 80/tcp > /dev/null 2>&1    # HTTP
        ufw allow 443/tcp > /dev/null 2>&1   # HTTPS
        ufw allow 8443/tcp > /dev/null 2>&1  # Panel
        ufw allow 53/tcp > /dev/null 2>&1    # DNS
        ufw allow 53/udp > /dev/null 2>&1    # DNS
        
        log_done "UFW yapılandırıldı"
        
        log_info "Açık portlar: 22, 80, 443, 8443, 53"
    else
        log_warn "UFW bulunamadı, firewall manuel yapılandırın"
    fi
}

# ═══════════════════════════════════════════════════════════════════════════════
# KURULUM ÖZETİ
# ═══════════════════════════════════════════════════════════════════════════════

print_summary() {
    local SERVER_IP=$(curl -s ifconfig.me 2>/dev/null || hostname -I | awk '{print $1}')
    
    echo ""
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                                           ║${NC}"
    echo -e "${GREEN}║                    ${WHITE}${BOLD}KURULUM BAŞARIYLA TAMAMLANDI!${NC}${GREEN}                       ║${NC}"
    echo -e "${GREEN}║                                                                           ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${CYAN}┌─────────────────────────────────────────────────────────────────────────────┐${NC}"
    echo -e "${CYAN}│${NC} ${WHITE}${BOLD}Panel Erişimi${NC}                                                              ${CYAN}│${NC}"
    echo -e "${CYAN}├─────────────────────────────────────────────────────────────────────────────┤${NC}"
    echo -e "${CYAN}│${NC}                                                                             ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   ${YELLOW}URL:${NC}        ${WHITE}http://${SERVER_IP}:8443${NC}                              ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   ${YELLOW}Kullanıcı:${NC}  ${WHITE}admin${NC}                                                       ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   ${YELLOW}Şifre:${NC}      ${WHITE}admin123${NC}                                                    ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}                                                                             ${CYAN}│${NC}"
    echo -e "${CYAN}└─────────────────────────────────────────────────────────────────────────────┘${NC}"
    echo ""
    echo -e "${YELLOW}⚠️  GÜVENLİK: İlk girişte şifrenizi değiştirin!${NC}"
    echo ""
    echo -e "${CYAN}┌─────────────────────────────────────────────────────────────────────────────┐${NC}"
    echo -e "${CYAN}│${NC} ${WHITE}${BOLD}Önemli Dosyalar${NC}                                                            ${CYAN}│${NC}"
    echo -e "${CYAN}├─────────────────────────────────────────────────────────────────────────────┤${NC}"
    echo -e "${CYAN}│${NC}   Panel Binary:    ${WHITE}${INSTALL_DIR}/serverpanel${NC}                        ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   Panel Logları:   ${WHITE}${LOG_DIR}/${NC}                              ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   MySQL Şifresi:   ${WHITE}/root/.serverpanel/mysql.conf${NC}                     ${CYAN}│${NC}"
    echo -e "${CYAN}└─────────────────────────────────────────────────────────────────────────────┘${NC}"
    echo ""
    echo -e "${CYAN}┌─────────────────────────────────────────────────────────────────────────────┐${NC}"
    echo -e "${CYAN}│${NC} ${WHITE}${BOLD}Servis Komutları${NC}                                                           ${CYAN}│${NC}"
    echo -e "${CYAN}├─────────────────────────────────────────────────────────────────────────────┤${NC}"
    echo -e "${CYAN}│${NC}   Başlat:   ${GREEN}systemctl start serverpanel${NC}                               ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   Durdur:   ${GREEN}systemctl stop serverpanel${NC}                                ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   Durum:    ${GREEN}systemctl status serverpanel${NC}                              ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}   Loglar:   ${GREEN}journalctl -u serverpanel -f${NC}                              ${CYAN}│${NC}"
    echo -e "${CYAN}└─────────────────────────────────────────────────────────────────────────────┘${NC}"
    echo ""
    
    if [[ $ERRORS -gt 0 ]]; then
        echo -e "${RED}⚠️  Kurulum sırasında $ERRORS hata oluştu. Logları kontrol edin.${NC}"
    fi
    
    if [[ $WARNINGS -gt 0 ]]; then
        echo -e "${YELLOW}ℹ️  Kurulum sırasında $WARNINGS uyarı oluştu.${NC}"
    fi
    
    echo ""
    echo -e "${GREEN}Kurulum tamamlandı! Tarayıcınızda paneli açabilirsiniz.${NC}"
    echo ""
}

# ═══════════════════════════════════════════════════════════════════════════════
# ANA FONKSİYON
# ═══════════════════════════════════════════════════════════════════════════════

main() {
    print_banner
    
    echo -e "${WHITE}Kurulum başlatılıyor...${NC}"
    echo -e "${WHITE}Bu işlem birkaç dakika sürebilir.${NC}"
    
    # Kontroller
    check_root
    check_os
    check_resources
    check_ports
    
    # Kurulumlar
    install_base_packages
    install_build_tools
    install_go
    install_nodejs
    install_webserver
    install_php
    install_mysql
    install_dns
    install_certbot
    install_serverpanel
    create_service
    configure_firewall
    
    # Özet
    print_summary
}

# Scripti çalıştır
main "$@"
