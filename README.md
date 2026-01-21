# 🌊 IoT Water Safety & Rescue System (LoRa E32 UART)

![Status](https://img.shields.io/badge/Status-Final%20Fix-success)
![Hardware](https://img.shields.io/badge/Hardware-ESP32%20%7C%20E32--TTL--100-blue)
![Architecture](https://img.shields.io/badge/Architecture-Dual%20Database-purple)

## 📖 Deskripsi Proyek
**IoT Water Safety System** adalah perangkat keselamatan cerdas yang dirancang untuk memantau kondisi korban di area perairan *blank spot* (tanpa sinyal seluler). Sistem ini menggunakan modul komunikasi **LoRa E32-TTL-100 (UART)** untuk mengirimkan data telemetri dan koordinat GPS secara *real-time* ke pos pantau hingga jarak **4 KM** (Line of Sight).

Sistem dilengkapi dengan mekanisme **Hybrid Triggering** (Sensor Air Otomatis & Tombol SOS Manual) serta fitur pencegahan *false alarm* melalui algoritma *debounce* pada tombol.

---

## 🚀 Fitur Unggulan

1.  **Hybrid Triggering:** Deteksi otomatis saat tenggelam (Sensor Air) dan pemicu manual (Tombol SOS).
2.  **Smart GPS Lock:** Mengirimkan koordinat presisi. Jika GPS belum *lock*, sistem tetap mengirim status darurat dengan koordinat terakhir/kosong.
3.  **Reliable Communication:** Menggunakan LoRa E32 (Ebyte) dengan antarmuka UART yang stabil dan jangkauan jauh.
4.  **Safety Reset Mechanism:** Fitur reset alarm membutuhkan penekanan tombol 3 detik untuk mencegah reset tidak sengaja.

---

## 💾 Arsitektur Dual Database
Sistem ini menerapkan strategi penyimpanan ganda (*Dual Database*) untuk memisahkan beban kerja antara data sensor yang cepat (*high-speed*) dengan data administratif aplikasi.

### 1. InfluxDB (Time-Series Database)
* **Fungsi:** Menyimpan data mentah sensor (*raw telemetry*) yang masuk secara *real-time* dan terus-menerus.
* **Data yang Disimpan:** Koordinat (Lat, Lon), Nilai Sensor Air, Status SOS, RSSI (Kekuatan Sinyal), dan Timestamp.
* **Tujuan:** Mengoptimalkan performa *query* untuk menampilkan grafik pergerakan korban dan rute perjalanan tanpa membebani database utama.

### 2. MySQL / MariaDB (Relational Database)
* **Fungsi:** Menyimpan data terstruktur yang berkaitan dengan manajemen aplikasi dan pengguna.
* **Data yang Disimpan:** Profil User (Login), Daftar ID Perangkat (*Whitelisting*), dan Log Riwayat Insiden (History Alarm).
* **Tujuan:** Menjamin integritas relasi data dan keamanan akses pengguna ke Dashboard.

---

## 🔌 Pin Mapping (Wiring Configuration)

Konfigurasi kabel ini disesuaikan dengan kode `FINAL FIX`. Harap perhatikan **Cross-Wiring** (RX ketemu TX).

### 1. Modul LoRa E32 (UART) ke ESP32
| Pin LoRa E32 | Pin ESP32 | Fungsi | Keterangan |
| :--- | :--- | :--- | :--- |
| **VCC** | 3.3V / 5V | Power | Cek spek modul |
| **GND** | GND | Ground | - |
| **TX** | **GPIO 16** | Serial RX | Data Masuk ke ESP32 |
| **RX** | **GPIO 17** | Serial TX | Data Keluar dari ESP32 |
| **M0** | GND | Mode | Mode Normal (Transparent) |
| **M1** | GND | Mode | Mode Normal (Transparent) |

### 2. Modul GPS Neo-6M ke ESP32
| Pin GPS | Pin ESP32 | Fungsi | Keterangan |
| :--- | :--- | :--- | :--- |
| **VCC** | 3.3V | Power | - |
| **GND** | GND | Ground | - |
| **TX** | **GPIO 34** | Serial RX | Masuk ke ESP32 (Input Only Pin) |
| **RX** | **GPIO 12** | Serial TX | Keluar dari ESP32 |

### 3. Sensor & Aktuator
| Komponen | Pin ESP32 | Mode |
| :--- | :--- | :--- |
| **Water Level Sensor** | **GPIO 32** | INPUT Analog |
| **Tombol SOS** | **GPIO 4** | INPUT_PULLUP (Aktif LOW) |
| **Buzzer** | **GPIO 13** | OUTPUT (Active High) |
| **LED Indikator TX** | **GPIO 26** | OUTPUT (Nyala saat kirim data) |
| **LED Indikator SOS** | **GPIO 27** | OUTPUT (Nyala saat Bahaya) |

---

## 📖 Tata Cara Pemakaian Alat (User Guide)

Berikut adalah panduan operasional perangkat Transmitter bagi pengguna/korban:

### A. Persiapan Awal (Power On)
1.  **Nyalakan Alat:** Hubungkan baterai ke perangkat.
2.  **Tunggu GPS Lock:** Bawa alat ke area terbuka (langit terlihat). Tunggu hingga LED pada modul GPS berkedip (sekitar 1-3 menit).
3.  **Mode Standby:** Jika LED TX (GPIO 26) berkedip singkat setiap 5 detik, berarti alat sudah aktif dan mengirim sinyal "Heartbeat" (Status: AMAN).

### B. Kondisi Darurat Manual (Tombol SOS)
1.  **Aktivasi:** Tekan **Tombol SOS** (GPIO 4) satu kali.
2.  **Indikator:**
    * Buzzer akan berbunyi putus-putus (Beep-Beep).
    * LED SOS (GPIO 27) akan menyala merah terus-menerus.
    * Sistem mengirim data status `s:1` (BAHAYA) ke penerima secara instan.

### C. Kondisi Darurat Otomatis (Jatuh ke Air)
1.  **Aktivasi:** Saat sensor air (GPIO 32) terendam air sepenuhnya (Nilai analog > 800).
2.  **Indikator:** Sistem otomatis beralih ke mode SOS tanpa perlu menekan tombol. Buzzer dan LED SOS akan aktif.

### D. Cara Reset (Mematikan Alarm)
Untuk mematikan mode SOS dan kembali ke mode Normal:
1.  Pastikan tombol SOS sudah **dilepas**.
2.  Tekan dan **TAHAN Tombol SOS selama 3 Detik**.
3.  **Indikator Reset:**
    * Buzzer akan berbunyi panjang (1 detik).
    * LED SOS mati.
    * Sistem mengirim status `s:0` (AMAN) ke penerima.

---

## ⚠️ Disclaimer
Sistem ini menggunakan modul LoRa E32-TTL-100. Pastikan jumper M0 dan M1 pada modul LoRa terhubung ke GND (Ground) agar modul bekerja pada **Mode 0 (Normal Mode)**. Jangkauan 4 KM tercapai pada pengujian *Line of Sight* dengan antena eksternal.

---
### Author
**[Samuel Maston Pangaribuan]**
