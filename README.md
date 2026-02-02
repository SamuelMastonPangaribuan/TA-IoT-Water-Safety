<div align="center">

# 🌊 IoT Water Safety & Rescue System
### **LoRa E32 UART Implementation for Blank Spot Areas**

![Status](https://img.shields.io/badge/Status-FINAL%20FIX-success?style=for-the-badge&logo=checkbox)
![Hardware](https://img.shields.io/badge/Hardware-ESP32%20%7C%20GPS%20Neo--6M-blue?style=for-the-badge&logo=arduino)
![LoRa](https://img.shields.io/badge/Comms-LoRa%20E32%20(4KM)-orange?style=for-the-badge&logo=lora)
![Accuracy](https://img.shields.io/badge/GPS%20Accuracy-~1.7%20Meter-green?style=for-the-badge&logo=google-maps)

<p align="center">
  <b>Sistem keselamatan cerdas untuk perairan tanpa sinyal seluler.</b><br>
  Mengirimkan titik koordinat korban secara <i>real-time</i> dengan latensi rendah dan jangkauan jauh.
</p>

[📖 About](#-deskripsi-proyek) •
[🚀 Features](#-fitur-unggulan) •
[🛠️ **Installation**](#-installation--setup-guide-from-scratch) •
[📊 Results](#-hasil-validasi-lapangan-field-test) •
[⚙️ Usage](#-user-guide-panduan-operasional)

</div>

---

## 📖 Deskripsi Proyek

**IoT Water Safety System** adalah perangkat *wearable* keselamatan yang dirancang khusus untuk area *blank spot* (laut lepas/hutan). Sistem ini mengatasi keterbatasan sinyal GSM dengan menggunakan teknologi **LoRa (Long Range)**.

Modul **E32-TTL-100** memungkinkan transmisi data telemetri (Status Sensor & GPS) hingga jarak **4 KM** (Line of Sight) dari korban ke pos pantau. Data kemudian diproses menggunakan arsitektur **Dual Database** untuk visualisasi yang cepat dan penyimpanan log yang aman.

---

## 🚀 Fitur Unggulan

| Fitur | Deskripsi |
| :--- | :--- |
| 🚨 **Hybrid Triggering** | Aktivasi alarm ganda: **Otomatis** (Sensor Air) saat tenggelam, atau **Manual** (Tombol SOS). |
| 🛰️ **Smart GPS Lock** | Algoritma cerdas yang mengirimkan koordinat presisi, atau data terakhir jika satelit belum terkunci. |
| 📡 **Long Range (LoRa)** | Komunikasi radio mandiri (433/915 MHz) yang stabil hingga 4KM tanpa pulsa/internet. |
| 🛡️ **Safety Reset** | Mencegah reset alarm yang tidak disengaja dengan mekanisme *Hold Button* selama 3 detik. |
| ⚡ **Real-time Dashboard** | Visualisasi posisi korban pada peta digital dengan update per detik. |

---

## 🛠️ Installation & Setup Guide (From Scratch)

Panduan ini disusun secara berurutan mulai dari perakitan kabel (*wiring*), instalasi *software*, hingga alat siap digunakan.

### 📦 Tahap 1: Persiapan Hardware & Wiring
Siapkan komponen utama: 2x ESP32, 2x LoRa E32, 1x GPS Neo-6M, Sensor Air, dan Tombol.
Rakit komponen mengikuti tabel pin di bawah ini (Sesuai kode `FINAL FIX`):

#### **A. Koneksi LoRa E32 ke ESP32**
*Catatan: Pastikan mode jumper M0 & M1 terhubung ke Ground untuk Mode Normal.*
| Pin LoRa | Pin ESP32 | Keterangan |
| :---: | :---: | :--- |
| `VCC` | 3.3V / 5V | Cek spesifikasi modul |
| `GND` | GND | Ground Common |
| `TX` | **GPIO 16** | Masuk ke RX2 ESP32 |
| `RX` | **GPIO 17** | Masuk ke TX2 ESP32 |
| `M0` | GND | Mode Normal (0) |
| `M1` | GND | Mode Normal (0) |

#### **B. Koneksi GPS Neo-6M ke ESP32**
| Pin GPS | Pin ESP32 | Keterangan |
| :---: | :---: | :--- |
| `VCC` | 3.3V | Power Supply |
| `TX` | **GPIO 34** | *Input Only Pin* (Aman untuk RX) |
| `RX` | **GPIO 12** | Serial TX |

#### **C. Sensor & Aktuator**
| Komponen | Pin ESP32 | Mode Pin |
| :--- | :---: | :--- |
| **Water Sensor** | **GPIO 32** | `INPUT` (Analog ADC1) |
| **Tombol SOS** | **GPIO 4** | `INPUT_PULLUP` (Aktif LOW) |
| **Buzzer** | **GPIO 13** | `OUTPUT` (Active High) |
| **LED Status (TX)** | **GPIO 26** | `OUTPUT` (Kedip saat kirim) |
| **LED Bahaya (SOS)**| **GPIO 27** | `OUTPUT` (Nyala saat bahaya) |

---

### 💻 Tahap 2: Persiapan Environment
1.  **Install Arduino IDE:** Unduh di [arduino.cc](https://www.arduino.cc/en/software).
2.  **Install Driver USB:** Pastikan driver **CP210x** atau **CH340** sudah terinstall agar ESP32 terbaca.
3.  **Setup Board:**
    * Buka `File` -> `Preferences`.
    * Tambahkan URL: `https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json`
    * Buka `Tools` -> `Board` -> `Boards Manager`, install **"esp32"**.

### 📚 Tahap 3: Instalasi Library
Install library berikut via `Sketch` -> `Include Library` -> `Manage Libraries`:
* `TinyGPSPlus` (by Mikal Hart)
* `LoRa_E32` (by KrisKasprzak)
* `PubSubClient` (by Nick O'Leary)

### 🗄️ Tahap 4: Konfigurasi Database & Backend
1.  **MySQL:** Import file `database/db_watersafety.sql` ke phpMyAdmin.
2.  **InfluxDB:** Buat bucket `sensor_data` dan salin API Token.
3.  **Konfigurasi Kode:**
    * Buka `Receiver_Gateway.ino`.
    * Sesuaikan `SSID`, `PASSWORD`, dan kredensial Database Anda.

### 🚀 Tahap 5: Upload Firmware
1.  **Transmitter (Alat Korban):** Upload file `Transmitter_Final.ino`.
2.  **Receiver (Pos Pantau):** Upload file `Receiver_Gateway.ino`.
    * *Setting:* Board "DOIT ESP32 DEVKIT V1", Upload Speed "921600".

---

## 📊 Hasil Validasi Lapangan (Field Test)

Perangkat telah diuji di 3 lokasi berbeda untuk memvalidasi akurasi modul GPS Neo-6M.

### 📍 Ringkasan Akurasi GPS
| Lokasi Uji | Koordinat Acuan (*Ground Truth*) | Rata-rata Error | Status |
| :--- | :--- | :---: | :---: |
| **Lokasi 1** (Lapangan Terbuka) | `2.385994, 99.148044` | **1.71 Meter** | ✅ Sangat Valid |
| **Lokasi 2** (Area Taman) | `2.386593, 99.147932` | **2.04 Meter** | ✅ Valid |
| **Lokasi 3** (Bundaran) | `2.385112, 99.147781` | **1.26 Meter** | 🌟 **Terbaik** |

> **Analisis:** Dengan rata-rata error keseluruhan di bawah **2.5 Meter**, alat ini memenuhi standar keselamatan (toleransi GPS sipil umumnya 2.5 - 5 meter).

---

## 💾 Arsitektur Sistem & Database

Sistem ini menggunakan strategi **Dual Database** untuk memisahkan beban kerja data sensor yang berat dengan data manajemen aplikasi.

```mermaid
graph TD
    A[ESP32 Transmitter] -->|LoRa RF| B[ESP32 Receiver/Gateway]
    B -->|MQTT/WiFi| C[Node-RED / Backend]
    
    C -->|Telemetry Data| D[(InfluxDB)]
    C -->|User & Logs| E[(MySQL / MariaDB)]
    
    D --> F[Real-time Map Dashboard]
    E --> F
    
    style D fill:#ee0,stroke:#333,stroke-width:2px
    style E fill:#0dd,stroke:#333,stroke-width:2px
