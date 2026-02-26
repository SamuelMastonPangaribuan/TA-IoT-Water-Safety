<div align="center">

# 🌊 IoT Water Safety & Rescue System
### **LoRa E32 UART Implementation for Blank Spot Areas**

![Status](https://img.shields.io/badge/Status-FINAL%20FIX-success?style=for-the-badge&logo=checkbox)
![Hardware](https://img.shields.io/badge/Hardware-ESP32%20%7C%20GPS%20Neo--6M-blue?style=for-the-badge&logo=arduino)
![LoRa](https://img.shields.io/badge/Comms-LoRa%20E32%20(4KM)-orange?style=for-the-badge&logo=lora)
![Backend](https://img.shields.io/badge/Backend-Golang-cyan?style=for-the-badge&logo=go)

<p align="center">
  <b>Sistem keselamatan cerdas untuk perairan tanpa sinyal seluler.</b><br>
  Mengirimkan titik koordinat korban secara <i>real-time</i> dengan latensi rendah dan jangkauan jauh.
</p>

[📖 Deskripsi](#-deskripsi-proyek) •
[🚀 Fitur](#-fitur-unggulan) •
[🛠️ Installation](#-installation--setup-guide) •
[📊 Hasil Uji](#-hasil-pengujian--validasi-lapangan) •
[⚙️ Cara Pakai](#️-panduan-pemakaian-alat-untuk-klien)

</div>

---

## 📖 Deskripsi Proyek

**IoT Water Safety System** adalah perangkat *wearable* keselamatan yang dirancang khusus untuk area *blank spot* (laut lepas, danau, atau hutan). Sistem ini mengatasi keterbatasan sinyal GSM dengan menggunakan teknologi komunikasi **LoRa (Long Range)**.

Modul **E32-TTL-100** memungkinkan transmisi data telemetri (Status Sensor Air, Tombol SOS, & Koordinat GPS) hingga jarak maksimal **4 KM** (Line of Sight) dari korban ke pos pantau (Gateway). Data kemudian diproses menggunakan **Golang Backend** dan disimpan dalam arsitektur **Dual Database** untuk visualisasi peta *real-time*.

---

## 🚀 Fitur Unggulan

| Fitur | Deskripsi |
| :--- | :--- |
| 🚨 **Hybrid Triggering** | Aktivasi alarm ganda: **Otomatis** (saat sensor air tenggelam) atau **Manual** (tekan tombol SOS). |
| 🛰️ **Smart GPS Lock** | Algoritma cerdas yang mengirimkan koordinat presisi. Jika satelit belum terkunci, mengirim lokasi terakhir. |
| 📡 **Long Range (LoRa)** | Komunikasi radio frekuensi mandiri (433MHz/915MHz) tanpa memerlukan pulsa, kuota, atau internet di area korban. |
| 🛡️ **Safety Reset** | Mencegah *false alarm* (reset tidak disengaja) dengan mekanisme *Hold Button* selama 3 detik. |
| ⚡ **Dual DB Dashboard** | Visualisasi posisi korban yang cepat dan log historis yang aman menggunakan InfluxDB dan MySQL. |

---

## ⚙️ Panduan Pemakaian Alat (Untuk Klien)

Penggunaan alat ini dirancang sangat mudah dan otomatis bagi pengguna di lapangan. Cukup ikuti 3 langkah berikut:

### 1️⃣ Pasang & Nyalakan Alat
* Pasang perangkat *Water Safety* ini di lengan, pelampung, atau badan Anda.
* Colokkan baterai atau nyalakan tombol *Power*.
* **Selesai!** Anda tidak perlu melakukan apa-apa lagi. Saat berada di luar ruangan, alat akan otomatis melacak koordinat GPS Anda dan langsung mengirimkannya ke Server. Pos pantau (tim SAR/Pengawas) sudah bisa melihat titik lokasi Anda bergerak secara *real-time* di layar Web Dashboard mereka.

### 2️⃣ Jika Terjadi Kondisi Darurat (Butuh Bantuan)
Sistem memiliki dua cara untuk meminta tolong:
* **Cara Manual (Tekan Tombol):** Jika Anda tiba-tiba mengalami kram, kelelahan, atau merasa dalam bahaya, cukup **Tekan Tombol SOS 1x**.
* **Cara Otomatis (Tenggelam):** Jika Anda jatuh dan tidak sempat menekan tombol, tenang saja. Sensor pada alat akan langsung mendeteksi air saat terendam penuh dan otomatis memicu alarm.
> *Begitu mode darurat aktif, alat di lengan Anda akan berbunyi beep terus-menerus. Di saat yang sama, Web Dashboard di Pos Pantau akan berkedip merah dan membunyikan alarm agar tim penyelamat segera menuju lokasi koordinat Anda.*

### 3️⃣ Jika Sudah Aman (Matikan Alarm)
* Jika bantuan sudah datang, atau jika alarm tidak sengaja menyala (misal terkena cipratan ombak besar padahal Anda aman), Anda bisa mematikan alarmnya.
* Caranya: **Tekan dan TAHAN Tombol SOS selama 3 Detik.**
* Bunyi alat akan mati, dan Web di Pos Pantau akan kembali ke status "AMAN".

---

## 🛠️ Installation & Setup Guide (Teknis)

*(Bagian ini khusus untuk developer/teknisi yang ingin merakit atau memodifikasi ulang sistem).*

### 📦 Tahap 1: Perakitan Hardware (Wiring)
Rakit komponen Node/Transmitter mengikuti tabel pin di bawah ini (Sesuai kode `FINAL FIX`):

**A. Koneksi LoRa E32 ke ESP32**
| Pin LoRa | Pin ESP32 | Keterangan |
| :---: | :---: | :--- |
| `VCC` | 3.3V / 5V | Cek spesifikasi tegangan modul |
| `GND` | GND | Ground Common |
| `TX` | **GPIO 16** | Masuk ke RX2 ESP32 (Cross-Wiring) |
| `RX` | **GPIO 17** | Masuk ke TX2 ESP32 (Cross-Wiring) |
| `M0 & M1`| GND | Hubungkan ke Ground untuk Mode Normal |

**B. Koneksi Modul GPS & Sensor ke ESP32**
| Komponen | Pin ESP32 | Mode / Keterangan |
| :--- | :---: | :--- |
| **GPS Neo-6M (TX)** | **GPIO 34** | `Input Only Pin` (Menerima data dari satelit) |
| **Water Sensor (S)** | **GPIO 32** | `INPUT` (Membaca Analog ADC1) |
| **Tombol SOS** | **GPIO 4** | `INPUT_PULLUP` (Aktif LOW saat ditekan) |
| **Active Buzzer** | **GPIO 13** | `OUTPUT` (Alarm suara saat bahaya) |

### 💻 Tahap 2: Persiapan Software
1.  **Visual Studio Code (VS Code):** Install dari [code.visualstudio.com](https://code.visualstudio.com/). Buka menu *Extensions*, install **`Go`** (by Google).
2.  **Go (Golang) Compiler:** Download dan install dari [go.dev/dl](https://go.dev/dl/).
3.  **Arduino IDE:** Download dari [arduino.cc](https://www.arduino.cc/en/software). Tambahkan URL Board ESP32: `https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json`. Install **"esp32"** via Boards Manager.

### 📚 Tahap 3: Instalasi Library (Arduino IDE)
Masuk ke `Sketch` -> `Include Library` -> `Manage Libraries`, lalu install: `TinyGPSPlus`, `LoRa_E32`, dan `PubSubClient`.

### 🗄️ Tahap 4: Konfigurasi Database
1.  **MySQL:** Buka phpMyAdmin, buat DB `db_watersafety`, lalu *Import* file `database/db_watersafety.sql`.
2.  **InfluxDB:** Buat bucket bernama `sensor_data` dan simpan **API Token** Anda.

### 🚀 Tahap 5: Upload Firmware ke ESP32
1.  **Transmitter:** Buka `Transmitter_Final.ino`, lalu *Upload* ke ESP32 Korban.
2.  **Receiver:** Buka `Receiver_Gateway.ino`. Edit `ssid` dan `password` WiFi, lalu *Upload* ke ESP32 Pos Pantau.

### 🖥️ Tahap 6: Menjalankan Backend Golang (Server)
1.  Buka folder **`backend`** di VS Code.
2.  Buka terminal (`Ctrl + J`), download dependensi: `go mod tidy`.
3.  Konfigurasi koneksi MySQL dan InfluxDB di `config/config.go`.
4.  Jalankan server: `go run main.go`. Akses Web di `http://localhost:8080`.

---

## 📊 Hasil Pengujian & Validasi Lapangan

### 📍 1. Akurasi GPS Neo-6M (Static Test)
| Lokasi Uji | Kondisi Lingkungan | Rata-rata Error | Kesimpulan |
| :--- | :--- | :---: | :---: |
| **Area Lapangan** | Minim halangan (*Open Sky*) | **1.71 Meter** | Sangat Valid |
| **Area Taman** | Semi-terbuka (Pohon & Bangunan) | **2.04 Meter** | Valid |
| **Area Bundaran** | Terbuka penuh, minim interferensi | **1.26 Meter** | **Sangat Presisi** |

### ⏱️ 2. Stabilitas Pengiriman Data LoRa (Interval Target: 5 Detik)
* **Jarak 200m:** Rata-rata interval **5.00 detik** (Delay 0s). 100% Packet Delivery.
* **Jarak 300m - 1 KM:** Rata-rata interval **~5.7 detik**. Terdapat delay minor (< 1s) untuk *retransmission*.
* **Jarak 1.2 KM:** Rata-rata interval **6.42 detik**. Delay mencapai 1.4 detik per paket. Alat berfungsi, tapi *real-time* menurun.

---

## 💾 Arsitektur Sistem

```mermaid
graph TD
    A[ESP32 Transmitter] -->|LoRa RF 433MHz| B[ESP32 Gateway]
    B -->|MQTT / TCP/IP| C[Golang Backend Server]
    
    C -->|Time-Series Data| D[(InfluxDB)]
    C -->|User, Auth, & Logs| E[(MySQL / MariaDB)]
    
    D --> F[Real-time Web Dashboard]
    E --> F
    
    style C fill:#0ff,stroke:#333,stroke-width:2px
