# ServerPanel - Özellik Yol Haritası

Bu dosya WHM/cPanel özelliklerini analiz eder ve ServerPanel'e eklenecek özellikleri listeler.

---

## 📊 Mevcut Durum Özeti (Son Güncelleme: 16 Aralık 2025)

| Kategori | cPanel/WHM | ServerPanel | Tamamlanma |
|----------|------------|-------------|------------|
| Authentication | ✅ | ✅ JWT + Rol bazlı | %100 |
| Dashboard | ✅ | ✅ Sistem istatistikleri | %85 |
| Kullanıcı Yönetimi | ✅ | ✅ CRUD + Paket atama | %75 |
| Domain Yönetimi | ✅ | ✅ Domain + Subdomain + Silme seçenekleri | %95 |
| DNS Yönetimi | ✅ | ✅ BIND9 Zone Editor + SPF/DMARC | %95 |
| **E-posta Yönetimi** | ✅ | ✅ **Postfix + Dovecot + Roundcube + DKIM + Rate Limiting** | **%95** |
| Veritabanı Yönetimi | ✅ | ✅ phpMyAdmin SSO | %85 |
| Dosya Yönetimi | ✅ | ✅ Tam fonksiyonel | %95 |
| FTP Yönetimi | ✅ | ✅ Pure-FTPd | %90 |
| SSL/TLS | ✅ | ✅ Let's Encrypt + Otomatik Vhost | %98 |
| **PHP Yönetimi** | ✅ | ✅ **MultiPHP + Yazılım Yöneticisi** | **%95** |
| **Sunucu Yönetimi** | ✅ | ✅ **Sunucu Durumu + Yazılım Yöneticisi + Sistem Sağlığı** | **%95** |
| **Node.js Yönetimi** | ✅ | ✅ **NVM + PM2 + NPM + Kaynak İzleme** | **%95** |
| Backup | ✅ | ❌ | %0 |
| **Migration** | ✅ | ✅ **cPanel Import** | **%60** |
| **Cron Jobs** | ✅ | ✅ **Tam fonksiyonel** | **%95** |
| **Güvenlik** | ✅ | ✅ **Fail2ban + UFW + SSH Key + Malware + ModSecurity** | **%95** |
| Metrics/Logs | ✅ | ⚠️ Temel | %15 |
| Reseller Sistemi | ✅ | ⚠️ Rol var | %10 |
| **Kurulum Scripti** | ✅ | ✅ Tam otomatik + Migration + Mail + MultiPHP | %98 |
| **UI/UX** | ✅ | ✅ **Lottie Loading Animasyonları + Tema Uyumu** | **%90** |
| **Terminal** | ✅ | ✅ **WebSocket Terminal** | **%95** |

### 🆕 Son Eklenen Özellikler (23 Ocak 2026)
- ✅ **cPanel Migration** (YENİ!)
  - cPanel backup (.tar.gz) upload ve analiz
  - Otomatik hesap bilgisi çıkarma
  - Node.js/Passenger → PM2 dönüşümü
  - MySQL, DNS, FTP, DKIM, Cron import
  - Seçici import seçenekleri
  - Modern UI ile wizard akışı

### Önceki Özellikler (16 Aralık 2025)
- ✅ **Node.js Uygulama Yönetimi** (GÜNCELLENDİ!)
  - NVM (Node Version Manager) entegrasyonu
  - PM2 process manager ile uygulama yönetimi
  - Birden fazla Node.js sürümü desteği
  - Apache mod_proxy ile reverse proxy
  - Uygulama başlatma/durdurma/yeniden başlatma
  - Canlı log görüntüleme
  - Ortam değişkenleri yönetimi
  - Opsiyonel özellik (Sunucu Ayarları'ndan etkinleştirme)
  - **NPM Komutları** (YENİ!)
    - npm install, build, run script vb.
    - WebSocket ile real-time output (terminal deneyimi)
    - package.json scriptlerini otomatik tespit
    - Tehlikeli komut engelleme (güvenlik)
  - **PM2 Kaynak İzleme** (YENİ!)
    - CPU kullanımı (%)
    - RAM kullanımı (MB)
    - Uptime (ne kadar süredir çalışıyor)
    - Restart sayısı
- ✅ **Web Terminal** (YENİ!)
  - WebSocket tabanlı terminal erişimi
  - xterm.js ile tam terminal emülasyonu
  - Tam ekran modu
  - Keyboard shortcuts desteği

### Önceki Özellikler (5 Aralık 2025)
- ✅ **ModSecurity WAF** (YENİ!)
  - Web Application Firewall
  - OWASP Core Rule Set (CRS) entegrasyonu
  - Tespit/Engelleme modları
  - Audit log görüntüleme ve istatistikler
  - IP whitelist yönetimi
  - Kural listesi görüntüleme
  - **CMS Exclusion kuralları** (WordPress, Joomla, Drupal, PrestaShop, Magento)
  - **Manuel kural devre dışı bırakma** (ID ile)
  - **Detaylı bilgilendirme UI** (ModSecurity nedir, modlar, öneriler)
- ✅ **Malware Tarama Sistemi**
  - Arka planda tarama (sayfa kapatılabilir)
  - Canlı ilerleme gösterimi (progress bar, dosya adı)
  - Hızlı/Tam tarama seçenekleri
  - Tarama iptali
  - Tehdit tespiti ve karantina
  - Tarama geçmişi (veritabanında saklanır)
- ✅ **Güvenlik Bölümü**
  - Fail2ban Yönetimi (jail'ler, ban/unban IP, whitelist)
  - UFW Firewall Yönetimi (kurallar, varsayılan portlar, güvenli etkinleştirme)
  - SSH Güvenliği (port, root login, şifre/key authentication ayarları)
  - SSH Key Yönetimi (ED25519 key oluşturma, key ekleme/silme, fingerprint)
  - Güvenlik uyarıları (şifre kapatma, root kapatma için onay modalları)
- ✅ **Yazılım Yöneticisi Fail2ban Entegrasyonu**
  - Kurulumda otomatik jail yapılandırması (SSH, Apache, Postfix, Dovecot, FTP)
  - Log dosyaları otomatik oluşturma

### Önceki Özellikler (3 Aralık 2025)
- ✅ **Sistem Sağlığı Bölümü** (YENİ!)
  - Arka Plan İşlem Sonlandırıcı (tehlikeli işlemler, güvenilir kullanıcılar)
  - İşlem Yöneticisi (CPU/Memory kullanımı, kill, kullanıcı filtreleme)
  - Geçerli Disk Kullanımı (disk bilgisi, I/O istatistikleri)
  - Geçerli Çalışma İşlemleri (tüm işlemler listesi)
- ✅ **Cron Jobs Yönetimi** (YENİ!)
  - Cron işi oluşturma/düzenleme/silme
  - Zamanlama şablonları (dakikalık, saatlik, günlük, haftalık, aylık)
  - Özel cron ifadesi desteği
  - Manuel çalıştırma ve çıktı görüntüleme
  - Aktif/pasif durumu değiştirme
  - Sistem crontab senkronizasyonu
- ✅ **Spam Filtreleri Sayfası** (YENİ!)
  - SpamAssassin ayarları (spam skoru, otomatik silme)
  - ClamAV antivirüs durumu görüntüleme
  - Whitelist/Blacklist yönetimi
  - Veritabanı güncelleme tetikleme
- ✅ **Gelişmiş Yazılım Yönetimi** (YENİ!)
  - ClamAV tam kurulum/kaldırma (daemon + freshclam + temizlik)
  - ImageMagick tam kurulum/kaldırma (config temizliği dahil)
  - SpamAssassin/Fail2ban servis yönetimi
  - Kalıntısız kaldırma (paketler, config, kullanıcılar, gruplar)
- ✅ **Lottie Loading Animasyonları** (YENİ!)
  - Tema uyumlu loading animasyonu
  - Ortak LoadingAnimation bileşeni
  - Dark/Light mode desteği
- ✅ **Mail Rate Limiting & Kuyruk Sistemi**
  - Hesap bazlı saatlik/günlük mail limiti
  - Paket bazlı limit tanımlama (Admin)
  - Postfix Policy Daemon entegrasyonu
  - Limit aşıldığında otomatik kuyruğa alma
  - Kuyruk yönetimi (silme, yeniden deneme, temizleme)
  - Kullanıcı mail istatistikleri görüntüleme
  - Queue Processor daemon (otomatik gönderim)
- ✅ **Yazılım Yöneticisi** (Admin Panel)
  - PHP sürümleri kurma/kaldırma (7.4, 8.0, 8.1, 8.2, 8.3)
  - PHP eklentileri kurma/kaldırma
  - Apache modülleri etkinleştirme/devre dışı bırakma
  - Ek yazılımlar kurma/kaldırma (Redis, Memcached, ImageMagick vs.)
  - **Gerçek zamanlı log görüntüleme** (WebSocket)
- ✅ **Sunucu Ayarları** (Admin Panel)
  - MultiPHP aktif/pasif
  - Domain bazlı PHP aktif/pasif
  - Varsayılan PHP sürümü seçimi
  - İzin verilen PHP sürümlerini belirleme
- ✅ **Sunucu Özellikleri** (Müşteri Panel)
  - Kurulu PHP sürümlerini görüntüleme
  - Kurulu PHP eklentilerini görüntüleme
  - Aktif Apache modüllerini görüntüleme
  - Kurulu ek yazılımları görüntüleme
- ✅ **Ondrej PHP PPA** (install.sh)
  - Tüm PHP sürümleri için destek (7.4-8.3)
- ✅ **Sunucu Durumu Sayfaları** (Admin Panel)
  - Sunucu Bilgileri
  - Günlük İşlem Günlüğü
  - Top Processes
  - Task Queue (Postfix + Rate Limit Kuyruğu + Kullanıcı İstatistikleri)

- ✅ **Tam Mail Sistemi** (Postfix + Dovecot + Roundcube)
- ✅ **DKIM Otomatik Kurulum** (hesap oluşturulduğunda)
- ✅ **SPF/DMARC DNS Kayıtları** (otomatik eklenir)
- ✅ **OpenDKIM Entegrasyonu** (mail imzalama)
- ✅ **SpamAssassin** (spam filtreleme)
- ✅ **ClamAV** (virüs tarama)
- ✅ **webmail.domain.com** subdomain desteği
- ✅ **SSL Otomatik Vhost** (webmail, mail, ftp, www için)
- ✅ Subdomain SSL sertifikası alma (her FQDN için ayrı)
- ✅ SSL Status sayfası (cPanel benzeri tablo görünümü)
- ✅ Domain/Subdomain silme sırasında dosya silme seçeneği
- ✅ Subdomain için modern hoşgeldin sayfası
- ✅ DNS A kaydı otomatik ekleme (subdomain için)
- ✅ Veritabanı migration (mevcut kurulumlar için)

---

## 🔐 1. AUTHENTICATION & GÜVENLİK

### Mevcut ✅
- [x] JWT tabanlı authentication
- [x] Rol bazlı erişim (Admin/Reseller/User)
- [x] Login/Logout

### Eksik Özellikler
- [ ] **İki Faktörlü Kimlik Doğrulama (2FA)**
  - TOTP (Google Authenticator, Authy)
  - SMS doğrulama
  - Yedek kodlar
- [ ] **Şifre Politikaları**
  - Minimum uzunluk
  - Karmaşıklık gereksinimleri
  - Şifre geçmişi
  - Otomatik kilitleme
- [ ] **Session Yönetimi**
  - Aktif oturumları görme
  - Uzaktan oturum kapatma
  - Session timeout ayarları
- [ ] **IP Kısıtlamaları**
  - Beyaz liste
  - Kara liste
  - Ülke bazlı engelleme
- [ ] **API Token Yönetimi**
  - Token oluşturma/silme
  - İzin bazlı tokenlar
  - Token son kullanma tarihi
- [ ] **Güvenlik Logları**
  - Başarısız giriş denemeleri
  - Şüpheli aktiviteler
  - Brute-force koruması (fail2ban entegrasyonu)

---

## 👥 2. KULLANICI YÖNETİMİ

### Mevcut ✅
- [x] Kullanıcı listeleme
- [x] Kullanıcı oluşturma/güncelleme/silme
- [x] Rol atama (Admin/Reseller/User)

### Eksik Özellikler
- [ ] **Paket Atama**
  - Kullanıcıya hosting paketi atama
  - Kota yönetimi
  - Kaynak limitleri
- [ ] **Kullanıcı Detay Sayfası**
  - Kullanıcının tüm kaynaklarını görme
  - Disk kullanımı
  - Bandwidth kullanımı
- [ ] **Toplu İşlemler**
  - Çoklu kullanıcı askıya alma
  - Çoklu paket değiştirme
  - CSV import/export
- [ ] **Kullanıcı Arama & Filtreleme**
  - Domain'e göre arama
  - Duruma göre filtreleme
  - Pakete göre filtreleme
- [ ] **Reseller Hiyerarşisi**
  - Alt kullanıcıları görme
  - Reseller kota limitleri
  - Özel fiyatlandırma
- [ ] **Hesap Askıya Alma/Aktifleştirme**
  - Geçici askıya alma
  - Otomatik askıya alma (kota aşımı)
  - Ödeme gecikme entegrasyonu

---

## 🌐 3. DOMAİN YÖNETİMİ

### Mevcut ✅
- [x] Domain listeleme API
- [x] Domain ekleme/silme API
- [x] **Domain Yönetim Arayüzü** ✅
  - Domain listesi sayfası (tab görünümü)
  - Domain ekleme formu
  - Paket limitleri kontrolü
  - Domain tipi (primary, addon, alias)
- [x] **Addon Domains** ✅
  - Ana domain'e ek domain ekleme
  - Ayrı document root
  - Otomatik Apache vhost
  - Otomatik DNS zone
- [x] **Subdomain Yönetimi** ✅
  - Subdomain oluşturma/silme
  - Subdomain yönlendirme (301/302)
  - Otomatik DNS A kaydı
  - Otomatik Apache vhost

### Eksik Özellikler
- [ ] **Wildcard Subdomain**
  - *.domain.com desteği
- [ ] **Domain Alias (Parked Domains)**
  - Aynı içeriği farklı domain'de gösterme
- [ ] **Domain Yönlendirme**
  - 301/302 redirect
  - Wildcard redirect
  - Koşullu yönlendirme
- [ ] **Document Root Yönetimi**
  - Klasör seçimi
  - Otomatik klasör oluşturma
- [ ] **NGINX/Apache Konfigürasyonu**
  - Virtual host oluşturma
  - PHP sürüm seçimi
  - Custom direktifler

---

## 🔤 4. DNS YÖNETİMİ

### Mevcut ✅
- [x] **Zone Editor** (BIND9)
  - A, AAAA, CNAME, MX, TXT, NS, SRV, CAA kayıtları
  - TTL yönetimi (preset seçenekleri)
  - Kayıt ekleme/düzenleme/silme
  - Zone sıfırlama (varsayılana döndürme)
  - Kullanıcı izolasyonu (sadece kendi domainleri)
  - Admin tüm zone'ları yönetebilir
- [x] **Otomatik Zone Oluşturma**
  - Hesap oluşturulduğunda otomatik DNS zone
  - Varsayılan A, MX, TXT (SPF) kayıtları
- [x] **cPanel Benzeri UI**
  - Kayıt tipi filtreleme
  - Renkli tip badge'leri
  - Domain seçici sidebar

### Eksik Özellikler
- [ ] **DNS Şablonları**
  - Özel kayıt şablonları
  - Hızlı kurulum
- [ ] **DNS Cluster**
  - Birden fazla DNS sunucu desteği
  - Zone senkronizasyonu
- [ ] **DNSSEC**
  - DNSSEC aktivasyonu
  - Anahtar yönetimi
- [ ] **DNS Propagation Kontrolü**
  - Propagation durumu
  - DNS sorgu testi
- [ ] **Reverse DNS (PTR)**
  - PTR kayıt yönetimi
- [ ] **Dynamic DNS**
  - API ile DNS güncelleme
  - Dinamik IP desteği

---

## 📧 5. E-POSTA YÖNETİMİ

### Mevcut ✅
- [x] **Mail Server Kurulumu** (Postfix + Dovecot)
- [x] **E-posta Hesapları Arayüzü**
  - [x] Hesap listesi
  - [x] Hesap oluşturma/silme
  - [x] Kota yönetimi
  - [x] Şifre değiştirme
- [x] **Webmail Entegrasyonu**
  - [x] Roundcube (webmail.domain.com)
  - [x] Otomatik SSL vhost
- [x] **E-posta Yönlendirme (Forwarders)**
  - [x] Tek adrese yönlendirme
  - [x] Çoklu yönlendirme
- [x] **Otomatik Yanıtlayıcı (Autoresponder)**
  - [x] Tatil mesajı
  - [x] Zamanlı yanıtlar (başlangıç/bitiş tarihi)
- [x] **E-posta Filtreleri**
  - [x] SpamAssassin entegrasyonu
  - [x] ClamAV virüs tarama
- [x] **DKIM/SPF/DMARC**
  - [x] Otomatik DKIM key oluşturma (hesap oluşturulduğunda)
  - [x] SPF kaydı otomatik ekleme
  - [x] DMARC kaydı otomatik ekleme
  - [x] OpenDKIM entegrasyonu
- [x] **Rate Limiting**
  - [x] Saatlik mail limiti (varsayılan: 100)
  - [x] Günlük mail limiti (varsayılan: 500)
- [x] **TLS/SSL Güvenliği**
  - [x] SMTP TLS (port 587)
  - [x] SMTPS (port 465)
  - [x] IMAPS (port 993)

### Eksik Özellikler
- [ ] **Mailing Lists**
  - Liste oluşturma
  - Üye yönetimi
  - Mailman entegrasyonu
- [ ] **E-posta Routing**
  - Local/Remote mail exchanger
  - Backup MX
- [ ] **E-posta İstatistikleri**
  - Gönderim/alım sayıları
  - Bounce oranları
  - Queue durumu
- [ ] **BoxTrapper**
  - Challenge-response spam koruması
- [ ] **Track Delivery**
  - E-posta takibi
  - Log analizi
- [ ] **Catch-All Email**
  - Tüm mailleri tek adrese yönlendirme

---

## 🗄️ 6. VERİTABANI YÖNETİMİ

### Mevcut ⚠️
- [x] Veritabanı listeleme API
- [x] Veritabanı oluşturma/silme API

### Eksik Özellikler
- [ ] **Veritabanı Arayüzü**
  - Veritabanı listesi sayfası
  - Oluşturma formu
  - Boyut bilgisi
- [ ] **MySQL/MariaDB Yönetimi**
  - Veritabanı oluşturma
  - Kullanıcı oluşturma
  - Yetki yönetimi
  - Remote access
- [ ] **PostgreSQL Desteği**
  - Veritabanı oluşturma
  - Kullanıcı yönetimi
- [ ] **phpMyAdmin Entegrasyonu**
  - Tek tıkla erişim
  - SSO (Single Sign-On)
- [ ] **phpPgAdmin Entegrasyonu**
  - PostgreSQL için web arayüzü
- [ ] **Veritabanı Yedekleme**
  - Manuel backup
  - Zamanlanmış backup
  - Restore
- [ ] **Remote Database**
  - Uzak bağlantı izinleri
  - IP whitelist
- [ ] **Veritabanı İstatistikleri**
  - Boyut takibi
  - Sorgu istatistikleri

---

## 📁 7. DOSYA YÖNETİMİ

### Mevcut ✅
Web tabanlı dosya yöneticisi tam fonksiyonel çalışıyor.

### Tamamlanan Özellikler
- [x] **Web Tabanlı File Manager** ✅
  - ✅ Dosya/klasör listeleme
  - ✅ Dosya yükleme (drag & drop)
  - ✅ Dosya indirme
  - ✅ Dosya düzenleme (code editor)
  - ✅ Dosya kopyalama/taşıma
  - ✅ Dosya silme
  - ✅ Yeniden adlandırma
  - ✅ Zip/Unzip (Archive)
  - ✅ Dosya arama
  - ✅ Resim önizleme
  - ✅ Dark mode desteği
  - ✅ ESC ile modal kapatma
  - ✅ Kaydedilmemiş değişiklik uyarısı
- [ ] **Dosya İzinleri (Permissions)**
  - chmod arayüzü
  - chown desteği
  - Recursive izin değişikliği
- [ ] **Directory Privacy**
  - .htpasswd koruması
  - Klasör şifreleme
- [ ] **Disk Usage Analizi**
  - Klasör bazlı kullanım
  - En büyük dosyalar
  - Görsel grafik
- [ ] **Hotlink Protection**
  - Resim/dosya koruması
  - İzin verilen domainler
- [ ] **Index Ayarları**
  - Directory listing
  - Custom index sayfası
- [ ] **MIME Types**
  - Özel MIME tanımları
- [ ] **Image Manager**
  - Thumbnail oluşturma
  - Resim boyutlandırma
  - Format dönüştürme

---

## 📤 8. FTP YÖNETİMİ

### Mevcut ✅
- [x] **FTP Hesapları** (Pure-FTPd)
  - Hesap oluşturma/silme
  - Şifre yönetimi (güçlü şifre generator)
  - Directory kısıtlaması (chroot)
  - Kota belirleme (sınırsız seçeneği)
  - Hesap aktif/pasif yapma
  - Kullanıcı adı kopyalama
- [x] **FTP Sunucu Yönetimi** (Admin)
  - Sunucu durumu görüntüleme
  - TLS şifreleme ayarları
  - Bağlantı limitleri
  - Pasif port aralığı
- [x] **UI/UX**
  - cPanel benzeri form tasarımı
  - Autocomplete dizin seçimi
  - Şifre gücü göstergesi
  - Loading animasyonları

### Eksik Özellikler
- [ ] **FTP İstatistikleri**
  - Bağlantı logları
  - Transfer istatistikleri
- [ ] **Anonymous FTP**
  - Anonim erişim ayarları
- [ ] **SFTP Desteği**
  - SSH üzerinden FTP
- [ ] **FTP Session Yönetimi**
  - Aktif bağlantıları görme
  - Bağlantı sonlandırma

---

## 🔒 9. SSL/TLS YÖNETİMİ

### Mevcut ✅
- [x] Let's Encrypt entegrasyonu (certbot)
- [x] Tek tıkla SSL sertifikası alma
- [x] Otomatik yenileme (cron job)
- [x] SSL durumu görüntüleme
- [x] Sertifika yenileme
- [x] Sertifika silme/iptal

### Eksik Özellikler
- [ ] **Gelişmiş SSL Yönetimi**
  - Manuel sertifika yükleme
  - Private key yönetimi
  - CSR oluşturma
  - Wildcard SSL
- [ ] **Force HTTPS**
  - Otomatik yönlendirme
  - HSTS ayarları

---

## 💾 10. YEDEKLEME (BACKUP)

### Mevcut ❌
Henüz yok

### Eklenecek Özellikler
- [ ] **Manuel Backup**
  - Full backup
  - Home directory backup
  - Database backup
  - E-posta backup
- [ ] **Zamanlanmış Backup**
  - Günlük/Haftalık/Aylık
  - Retention policy
- [ ] **Backup Hedefleri**
  - Lokal disk
  - Remote FTP/SFTP
  - Amazon S3
  - Google Cloud Storage
  - Backblaze B2
- [ ] **Restore**
  - Full restore
  - Kısmi restore
  - Dosya bazlı restore
- [ ] **Backup İstatistikleri**
  - Backup geçmişi
  - Boyut bilgisi
  - Durum raporları

---

## ⏰ 11. CRON JOBS

### Mevcut ❌
Henüz yok

### Eklenecek Özellikler
- [ ] **Cron Job Yönetimi**
  - Job oluşturma
  - Zamanlama editörü
  - Komut girişi
- [ ] **Cron Şablonları**
  - Yaygın zamanlamalar
  - Kolay seçim
- [ ] **Cron Logları**
  - Çalışma geçmişi
  - Hata logları
  - E-posta bildirimi

---

## 📊 12. METRİKLER & LOGLAR

### Mevcut ⚠️
- [x] Temel sistem istatistikleri (CPU, RAM, Disk)

### Eksik Özellikler
- [ ] **Bandwidth İstatistikleri**
  - Günlük/Aylık kullanım
  - Domain bazlı
  - Grafikler
- [ ] **Ziyaretçi İstatistikleri**
  - AWStats entegrasyonu
  - Webalizer
  - Analog Stats
- [ ] **Error Logs**
  - Apache/Nginx hata logları
  - PHP hataları
  - Canlı log takibi
- [ ] **Access Logs**
  - Ham erişim logları
  - Log analizi
  - IP bazlı filtreleme
- [ ] **Resource Usage**
  - CPU kullanımı (process bazlı)
  - Memory kullanımı
  - I/O istatistikleri
- [ ] **Uptime Monitoring**
  - Servis durumu
  - Uptime geçmişi
  - Uyarı sistemi

---

## 🛡️ 13. GÜVENLİK ÖZELLİKLERİ

### Mevcut ✅
- [x] Temel authentication
- [x] **SpamAssassin Entegrasyonu**
  - Spam skoru ayarlama
  - Otomatik silme eşiği
  - Spam klasörüne taşıma
- [x] **ClamAV Antivirüs**
  - Virüs veritabanı durumu
  - Otomatik güncelleme
  - Daemon yönetimi
- [x] **Whitelist/Blacklist**
  - E-posta/domain bazlı filtreleme
  - Dinamik liste yönetimi
- [x] **Spam Filtreleri UI**
  - Ayarlar sayfası
  - Durum görüntüleme
  - İstatistikler

### Yeni Eklenen ✅
- [x] **Fail2ban Yönetimi**
  - Servis durumu görüntüleme
  - Jail listesi ve istatistikleri
  - IP ban/unban
  - Jail ayarları (bantime, findtime, maxretry)
  - Whitelist yönetimi
- [x] **UFW Firewall Yönetimi**
  - Firewall durumu görüntüleme
  - Kural ekleme/silme
  - Varsayılan portlar (SSH, HTTP, HTTPS, Panel, FTP, Mail, DNS, MySQL)
  - Güvenli etkinleştirme (portlar önce açılır)
- [x] **SSH Güvenliği**
  - SSH port değiştirme
  - Root login ayarları (izin ver, sadece key ile, yasakla)
  - Şifre/Key authentication ayarları
  - Max deneme sayısı ve giriş süresi
  - Güvenlik puanı hesaplama
- [x] **SSH Key Yönetimi**
  - ED25519 key çifti oluşturma
  - Private key tek seferlik indirme (sunucuda saklanmaz)
  - Mevcut public key ekleme
  - Key listeleme (fingerprint ile)
  - Key silme
- [x] **Güvenlik Uyarıları**
  - Şifre girişi kapatılırken SSH key kontrolü
  - Root girişi kapatılırken onay modalı

- [x] **Malware Tarama (ClamAV)**
  - Arka planda tarama (sayfa kapatılabilir)
  - Canlı ilerleme gösterimi (progress bar)
  - Taranan dosya adı gösterimi
  - Hızlı/Tam tarama seçenekleri
  - Tarama iptali
  - Tehdit tespiti ve listeleme
  - Karantinaya alma
  - Dosya silme
  - Tarama geçmişi (veritabanında saklanır)
  - Admin tüm kullanıcıların taramalarını görebilir

- [x] **ModSecurity WAF**
  - WAF aktivasyonu/deaktivasyonu
  - Tespit/Engelleme mod seçimi
  - OWASP CRS kural listesi
  - Audit log görüntüleme
  - İstatistikler (engellenen, loglanan)
  - IP whitelist yönetimi
  - CMS Exclusion kuralları (WordPress, Joomla, Drupal, PrestaShop, Magento)
  - Manuel kural devre dışı bırakma (ID ile)
  - Detaylı bilgilendirme UI

### Eksik Özellikler
- [ ] **ModSecurity Gelişmiş**
  - Domain bazlı ModSecurity yönetimi
  - Otomatik OWASP CRS güncelleme
  - Vendor kural seti seçimi (Comodo, Atomicorp)
- [ ] **Leech Protection**
  - Şifre sızıntı koruması

---

## 🔧 14. SUNUCU YÖNETİMİ (WHM)

### Mevcut ⚠️
- [x] Servis listesi API
- [x] Servis restart API

### Eksik Özellikler
- [ ] **Servis Yönetimi Arayüzü**
  - Servis durumları
  - Start/Stop/Restart
  - Otomatik başlatma
- [ ] **PHP Yönetimi**
  - Çoklu PHP sürümü
  - PHP-FPM yönetimi
  - php.ini editörü
  - PHP extension yönetimi
- [ ] **Apache/NGINX Yönetimi**
  - Konfigürasyon editörü
  - Module yönetimi
  - Virtual host yönetimi
- [ ] **MySQL/MariaDB Yönetimi**
  - my.cnf editörü
  - Performans ayarları
  - Slow query log
- [ ] **Mail Server Yönetimi**
  - Exim/Postfix konfigürasyonu
  - Queue yönetimi
  - Mail log analizi
- [ ] **Sistem Güncelleme**
  - OS güncellemeleri
  - Paket yönetimi
- [ ] **Server Bilgisi**
  - Donanım bilgisi
  - OS bilgisi
  - Yüklü yazılımlar

---

## 📦 15. PAKET YÖNETİMİ

### Mevcut ✅
- [x] Paket listeleme
- [x] Paket oluşturma/güncelleme/silme
- [x] **Paket Yönetimi Arayüzü** ✅
  - Paket listesi sayfası (grid görünümü)
  - Detaylı kota ayarları
  - PHP ayarları (memory, upload, execution time)
  - Kullanıcı sayısı gösterimi
  - Oluşturma/düzenleme/silme modal'ları

### Eksik Özellikler
- [ ] **Gelişmiş Kota Seçenekleri**
  - Inode limiti
  - MySQL veritabanı sayısı
  - PostgreSQL veritabanı sayısı
  - Email hesap sayısı
  - Mailing list sayısı
  - Subdomain sayısı
  - Addon domain sayısı
  - FTP hesap sayısı
  - Max email gönderimi/saat
- [ ] **Özellik Listeleri**
  - cPanel özellik seçimi
  - Modül bazlı erişim
- [ ] **Reseller Paketleri**
  - Reseller kotaları
  - Overselling ayarları

---

## 🔄 16. MİGRASYON

### Mevcut ✅
- [x] **cPanel Migration** (YENİ!)
  - cPanel backup (.tar.gz) upload ve analiz
  - Otomatik hesap bilgisi çıkarma (username, domain, email, PHP version)
  - Node.js/Passenger uygulaması tespiti ve PM2'ye dönüştürme
  - MySQL veritabanı import
  - DNS zone import (BIND formatı)
  - FTP hesapları import (Pure-FTPd uyumlu hash)
  - DKIM anahtarları import
  - Cron job'lar import
  - Seçici import (dosyalar, DB, DNS, email, FTP, Node.js, cron)
  - Paket atama ve yeni şifre belirleme
  - İlerleme takibi ve detaylı sonuç raporu

### Eksik Özellikler
- [ ] **Plesk Migration**
  - Plesk backup import
- [ ] **DirectAdmin Migration**
  - DirectAdmin backup import
- [ ] **Manuel Migration**
  - Dosya yükleme
  - Veritabanı import
  - DNS import

---

## 🎨 17. TEMA & ÖZELLEŞTİRME

### Mevcut ⚠️
- [x] Temel dashboard

### Eksik Özellikler
- [ ] **Tema Sistemi**
  - Açık/Koyu mod
  - Renk şemaları
- [ ] **Branding**
  - Logo değiştirme
  - Favicon
  - Şirket adı
- [ ] **Dil Desteği**
  - Çoklu dil
  - Türkçe
  - İngilizce
- [ ] **Dashboard Özelleştirme**
  - Widget düzeni
  - Hızlı erişim kısayolları

---

## 📱 18. API & ENTEGRASYONLAR

### Mevcut ⚠️
- [x] REST API (temel)

### Eksik Özellikler
- [ ] **API Dokümantasyonu**
  - Swagger/OpenAPI
  - Interaktif docs
- [ ] **Webhook Desteği**
  - Event bazlı bildirimler
  - Custom webhook URL
- [ ] **WHMCS Entegrasyonu**
  - Provisioning modülü
  - SSO desteği
- [ ] **Cloudflare Entegrasyonu**
  - DNS senkronizasyonu
  - Proxy ayarları
- [ ] **WordPress Toolkit**
  - WP kurulumu
  - WP yönetimi
  - Güvenlik taraması

---

## 🚀 GELİŞTİRME SIRASI (Gerçek Kullanım Öncelikli)

Bir hosting müşterisinin temel ihtiyaçlarına göre sıralandı:

### 🎯 Faz 1 - MVP (Minimum Viable Product)
> Müşteri website yayınlayabilmeli

| # | Özellik | Neden Gerekli? | Durum |
|---|---------|----------------|-------|
| 1 | ✅ Authentication & Dashboard | Panele giriş | ✅ Tamam |
| 2 | ✅ Hesap Yönetimi (CRUD) | Hosting hesabı | ✅ Tamam |
| 3 | ✅ Apache/PHP-FPM Entegrasyonu | Web sunucu | ✅ Tamam |
| 4 | ✅ DNS Zone (BIND9) | Domain yönlendirme | ✅ Tamam |
| 5 | ✅ Welcome Page | İlk açılış sayfası | ✅ Tamam |
| 6 | ✅ Dosya Yöneticisi | Site dosyalarını yükleme | ✅ Tamam |
| 7 | Veritabanı UI + phpMyAdmin | WordPress vb. kurulum | ✅ Tamam |
| 8 | SSL/Let's Encrypt | HTTPS zorunlu | ✅ Tamam |

### 🎯 Faz 2 - Temel Hosting
> Müşteri e-posta kullanabilmeli, yedek alabilmeli

| # | Özellik | Neden Gerekli? | Durum |
|---|---------|----------------|-------|
| 6 | E-posta Hesapları UI | info@domain.com | ⏳ Bekliyor |
| 7 | Webmail (Roundcube) | Tarayıcıdan mail okuma | ⏳ Bekliyor |
| 8 | FTP Hesapları | Büyük dosya yükleme | ✅ Tamam |
| 9 | Backup & Restore | Veri kaybını önleme | ⏳ Bekliyor |
| 10 | DNS Zone Editor | Mail/subdomain ayarları | ✅ Tamam |

### 🎯 Faz 3 - Profesyonel Hosting
> Gelişmiş müşteriler için

| # | Özellik | Neden Gerekli? | Durum |
|---|---------|----------------|-------|
| 11 | Cron Jobs | Zamanlanmış görevler | ⏳ Bekliyor |
| 12 | PHP Sürüm Seçimi | Farklı PHP versiyonları | ⏳ Bekliyor |
| 13 | SSH/Terminal Erişimi | Geliştiriciler için | ⏳ Bekliyor |
| 14 | Subdomain Yönetimi | blog.domain.com | ✅ Tamam |
| 15 | Error Logs | Hata ayıklama | ⏳ Bekliyor |

### 🎯 Faz 4 - Reseller & Enterprise
> Hosting satışı yapanlar için

| # | Özellik | Durum |
|---|---------|-------|
| 16 | Paket Yönetimi | ✅ Tam UI | %90 |
| 17 | Reseller Panel | ⏳ Bekliyor |
| 18 | WHMCS Entegrasyonu | ⏳ Bekliyor |
| 19 | Çoklu Sunucu | ⏳ Bekliyor |
| 20 | Migration Tools | ⏳ Bekliyor |

---

## 📈 İlerleme Durumu

- **Tamamlanan**: 55+ özellik
- **Devam Eden**: 1 özellik (Backup)
- **Bekleyen**: 105+ özellik
- **Toplam İlerleme**: ~%50

### ✅ Son Tamamlanan Özellikler (2 Aralık 2025)
- Tek komutla kurulum scripti (install.sh)
- Linux user yönetimi (useradd/userdel)
- Apache vhost yönetimi (a2ensite/a2dissite)
- PHP-FPM pool yönetimi
- BIND9 DNS zone yönetimi
- Home dizini izin yönetimi (711/755)
- Welcome page otomatik oluşturma
- Hesap CRUD (Create/Read/Update/Delete)
- Veritabanı yönetimi + phpMyAdmin SSO
- **SSL/Let's Encrypt entegrasyonu**
  - Tek tıkla SSL sertifikası
  - Otomatik yenileme
  - SSL durumu görüntüleme
- **Dosya Yöneticisi (Tam fonksiyonel)**
  - Çoklu dosya yükleme + progress bar
  - Drag & drop desteği
  - Dosya düzenleme (code editor)
  - Kopyalama/Taşıma/Silme
  - Zip/Unzip (Archive)
  - Resim önizleme
  - 512MB yükleme limiti
- **MultiPHP Yönetimi**
  - PHP versiyon seçimi (7.4, 8.0, 8.1, 8.2, 8.3)
  - PHP INI ayarları düzenleme
  - Paket bazlı PHP limitleri
  - memory_limit, upload_max_filesize, max_execution_time
- **FTP Yönetimi (Pure-FTPd)**
  - FTP hesabı oluşturma/silme/aktif-pasif
  - Dizin kısıtlaması (chroot)
  - Kota yönetimi (sınırsız seçeneği)
  - Şifre gücü göstergesi
  - Admin sunucu ayarları
- **UI/UX İyileştirmeleri**
  - Merkezi tema renk sistemi (CSS variables)
  - Light/Dark mode tutarlılığı
  - Tüm sayfalarda tutarlı başlık boyutları
  - Badge ve alert renkleri düzeltildi
  - phpMyAdmin blowfish_secret otomatik yapılandırma
- **DNS Zone Editor (BIND9)**
  - A, AAAA, CNAME, MX, TXT, NS, SRV, CAA kayıtları
  - TTL yönetimi (preset seçenekleri)
  - Kayıt ekleme/düzenleme/silme
  - Zone sıfırlama (varsayılana döndürme)
  - Kullanıcı izolasyonu
  - cPanel benzeri UI
  - **Kayıt arama çubuğu** (isim, içerik, tip filtreleme)
- **Paket Yönetimi UI**
  - Paket listesi (grid görünümü)
  - Paket oluşturma/düzenleme/silme
  - PHP ayarları (memory, upload, execution time)
  - Disk, bant genişliği, domain, veritabanı, e-posta, FTP limitleri
- **Domain & Subdomain Yönetimi**
  - Domain ekleme/silme (addon domain)
  - Subdomain ekleme/silme
  - Yönlendirme desteği (301/302)
  - Paket limitleri kontrolü
  - Otomatik Apache vhost oluşturma
  - Otomatik DNS zone/kayıt oluşturma
  - Kullanım limitleri gösterimi

---

*Son güncelleme: 16 Aralık 2025*
