# ServerPanel Mimari Dökümanı

## Mevcut Durum vs Hedef

### ❌ Mevcut (Yanlış)
```
Admin giriş yapar
├── Domain ekler → Sadece DB'ye yazılır
├── Veritabanı ekler → Sadece DB'ye yazılır
└── Hiçbir şey gerçekte oluşmaz
```

### ✅ Hedef (Doğru)
```
Admin (WHM benzeri):
├── Kullanıcı/Hesap oluşturur
│   ├── Linux user: useradd -m -d /home/username -s /bin/bash username
│   ├── Dizin yapısı: /home/username/{public_html, mail, logs, tmp}
│   ├── İzinler: chown -R username:username /home/username
│   ├── Quota: setquota veya disk limit
│   └── PHP-FPM pool: /etc/php/8.x/fpm/pool.d/username.conf
│
├── Domain atar
│   ├── Nginx vhost: /etc/nginx/sites-available/domain.com
│   ├── Document root: /home/username/public_html/domain.com
│   ├── SSL config: Let's Encrypt için hazırlık
│   └── DNS zone: BIND veya PowerDNS

Kullanıcı (cPanel benzeri):
├── Kendi hesabına giriş yapar
├── Sadece kendi kaynaklarını görür
├── Kendi domainlerini yönetir
├── Kendi veritabanlarını yönetir
└── Kendi mail hesaplarını yönetir
```

---

## Gerekli Sistem Servisleri

### 1. Web Server (Nginx veya Apache)
```bash
# Nginx kurulumu
apt install nginx

# Her domain için vhost
/etc/nginx/sites-available/domain.com
/etc/nginx/sites-enabled/domain.com -> symlink
```

### 2. PHP-FPM (Her kullanıcı için ayrı pool)
```bash
# PHP-FPM kurulumu
apt install php8.2-fpm

# Her kullanıcı için pool
/etc/php/8.2/fpm/pool.d/username.conf
```

### 3. MySQL/MariaDB
```bash
# MariaDB kurulumu
apt install mariadb-server

# Veritabanı ve kullanıcı oluşturma
CREATE DATABASE username_dbname;
CREATE USER 'username_dbuser'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON username_dbname.* TO 'username_dbuser'@'localhost';
```

### 4. Mail Server (Postfix + Dovecot)
```bash
# Mail kurulumu
apt install postfix dovecot-core dovecot-imapd
```

### 5. FTP Server (ProFTPD veya Pure-FTPd)
```bash
apt install proftpd
```

### 6. DNS Server (BIND9 veya PowerDNS)
```bash
apt install bind9
```

---

## Kullanıcı İzolasyonu

### Yöntem 1: Linux Kullanıcıları + PHP-FPM Pools
```
/home/
├── user1/
│   ├── public_html/
│   ├── mail/
│   ├── logs/
│   └── tmp/
├── user2/
│   ├── public_html/
│   └── ...
```

Her kullanıcı:
- Kendi Linux kullanıcısı
- Kendi PHP-FPM pool'u (farklı uid/gid ile çalışır)
- Kendi dizin izinleri (700 veya 750)

### Yöntem 2: Docker Containerization (Gelişmiş)
Her kullanıcı ayrı container'da çalışır.

---

## Güvenlik Kontrolleri

### Domain Ekleme Güvenliği
```go
// 1. Domain formatı kontrolü
func isValidDomain(domain string) bool {
    // Regex ile kontrol
    // Sadece alfanumerik, tire ve nokta
}

// 2. Kullanıcı yetkisi kontrolü
func canUserAddDomain(userID int, domain string) bool {
    // Kullanıcının paketinde domain hakkı var mı?
    // Limit aşılmış mı?
}

// 3. Path traversal koruması
func sanitizePath(path string) string {
    // ../../../etc/passwd gibi saldırıları engelle
    // Sadece /home/username/ altına izin ver
}
```

---

## Uygulama Akışı

### Hesap Oluşturma (Admin)
```
1. Admin "Hesap Oluştur" der
2. Form: kullanıcı adı, email, şifre, paket seçimi
3. Backend:
   a. Kullanıcı adı uygun mu? (sistemde var mı, geçerli mi)
   b. Linux user oluştur
   c. Home dizini oluştur
   d. PHP-FPM pool oluştur
   e. Nginx default config oluştur
   f. Veritabanına kaydet
   g. Hoşgeldin emaili gönder
```

### Domain Ekleme (Kullanıcı)
```
1. Kullanıcı kendi panelinde "Domain Ekle" der
2. Backend:
   a. Bu kullanıcı domain ekleyebilir mi? (paket limiti)
   b. Domain geçerli mi?
   c. /home/username/public_html/domain.com oluştur
   d. Nginx vhost oluştur
   e. Nginx reload
   f. Veritabanına kaydet
```

---

## Dosya Yapısı (Revize)

```
/internal/
├── api/           # HTTP handlers
├── auth/          # Authentication
├── database/      # SQLite (panel verisi)
├── models/        # Data models
├── services/      # İş mantığı
│   ├── account/   # Hesap oluşturma/silme
│   ├── domain/    # Domain yönetimi
│   ├── database/  # MySQL veritabanı yönetimi
│   ├── email/     # Mail hesap yönetimi
│   └── ssl/       # Let's Encrypt
└── system/        # Linux komutları
    ├── user.go    # useradd, userdel
    ├── nginx.go   # vhost yönetimi
    ├── php.go     # PHP-FPM pool
    ├── mysql.go   # MySQL yönetimi
    └── dns.go     # DNS zone yönetimi
```

---

## Öncelik Sırası (Güncellendi)

### Faz 0 - Temel Altyapı (ÖNCELİK!)
1. [ ] Linux user yönetimi (useradd/userdel)
2. [ ] Dizin yapısı oluşturma
3. [ ] Nginx vhost yönetimi
4. [ ] Hesap oluşturma akışı

### Faz 1 - MVP
1. [ ] Hesap yönetimi UI (Admin)
2. [ ] Kullanıcının kendi paneli
3. [ ] Domain ekleme (gerçek)
4. [ ] Dosya yöneticisi

### Faz 2 - Temel Hosting
1. [ ] MySQL veritabanı yönetimi
2. [ ] SSL/Let's Encrypt
3. [ ] Backup
