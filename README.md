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

[📖 About](#-deskripsi-proyek) •
[🚀 Features](#-fitur-unggulan) •
[🛠️ **Installation Guide**](#-installation--setup-guide-from-scratch) •
[📊 Results](#-hasil-validasi-lapangan-field-test) •
[⚙️ Usage](#-user-guide-panduan-operasional)

</div>

---

## 📖 Deskripsi Proyek

**IoT Water Safety System** adalah perangkat *wearable* keselamatan yang dirancang khusus untuk area *blank spot* (laut lepas/hutan). Sistem ini mengatasi keterbatasan sinyal GSM dengan menggunakan teknologi **LoRa (Long Range)**.

Modul **E32-TTL-100** memungkinkan transmisi data telemetri (Status Sensor & GPS) hingga jarak **4 KM** (Line of Sight) dari korban ke pos pantau. Data kemudian diproses menggunakan **Golang Backend** dan disimpan dalam arsitektur **Dual Database**.

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

Panduan ini disusun secara berurutan mulai dari perakitan, instalasi tools (VS Code & Arduino), hingga menjalankan server Backend.

### 📦 Tahap 1: Persiapan Hardware & Wiring
Rakit komponen mengikuti tabel pin di bawah ini (Sesuai kode `FINAL FIX`):

**A. Koneksi LoRa E32 ke ESP32**
*Catatan: Pastikan jumper M0 & M1 terhubung ke Ground.*
| Pin LoRa | Pin ESP32 | Keterangan |
| :---: | :---: | :--- |
| `VCC` | 3.3V / 5V | Cek spesifikasi modul |
| `GND` | GND | Ground Common |
| `TX` | **GPIO 16** | Masuk ke RX2 ESP32 |
| `RX` | **GPIO 17** | Masuk ke TX2 ESP32 |
| `M0` | GND | Mode Normal (0) |
| `M1` | GND | Mode Normal (0) |

**B. Koneksi GPS Neo-6M ke ESP32**
| Pin GPS | Pin ESP32 | Keterangan |
| :---: | :---: | :--- |
| `VCC` | 3.3V | Power Supply |
| `TX` | **GPIO 34** | *Input Only Pin* (Aman untuk RX) |
| `RX` | **GPIO 12** | Serial TX |

**C. Sensor & Aktuator**
| Komponen | Pin ESP32 | Mode Pin |
| :--- | :---: | :--- |
| **Water Sensor** | **GPIO 32** | `INPUT` (Analog ADC1) |
| **Tombol SOS** | **GPIO 4** | `INPUT_PULLUP` (Aktif LOW) |
| **Buzzer** | **GPIO 13** | `OUTPUT` (Active High) |
| **LED Status** | **GPIO 26** | `OUTPUT` (Indikator TX) |

---

### 💻 Tahap 2: Persiapan Software (VS Code & Tools)
Sebelum coding, install software berikut secara berurutan:

1.  **Visual Studio Code (VS Code):**
    * Download dan install dari [code.visualstudio.com](https://code.visualstudio.com/).
    * Buka VS Code, pilih menu **Extensions** (kotak kiri), cari dan install:
        * `Go` (by Go Team at Google).
        * `Arduino` (optional, jika ingin coding Arduino di VS Code).

2.  **Go (Golang) Compiler:**
    * Download dan install dari [go.dev/dl](https://go.dev/dl/).
    * Cek instalasi via CMD/Terminal: `go version`.

3.  **Arduino IDE:**
    * Download dari [arduino.cc](https://www.arduino.cc/en/software).
    * **Install Driver USB:** Pastikan driver **CP210x** atau **CH340** terinstall.
    * **Setup Board:** Buka Arduino IDE -> `File` -> `Preferences`. Tambahkan URL:
        `https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json`
    * Install **"esp32"** via Boards Manager.

---

### 📚 Tahap 3: Instalasi Library Arduino
Install library berikut via Arduino IDE (`Sketch` -> `Include Library` -> `Manage Libraries`):
* `TinyGPSPlus` (by Mikal Hart)
* `LoRa_E32` (by KrisKasprzak)
* `PubSubClient` (by Nick O'Leary) - *Khusus Receiver*

---

### 🗄️ Tahap 4: Konfigurasi Database
1.  **MySQL:**
    * Buka phpMyAdmin.
    * Import file `database/db_watersafety.sql`.
2.  **InfluxDB:**
    * Buat bucket `sensor_data` dan salin **API Token**.

---

### 🚀 Tahap 5: Upload Firmware (Hardware)
1.  **Transmitter (Alat Korban):** Buka `Transmitter_Final.ino`, upload ke ESP32 Korban.
2.  **Receiver (Pos Pantau):** Buka `Receiver_Gateway.ino`.
    * Edit bagian `ssid` dan `password` WiFi.
    * Upload ke ESP32 Gateway.

---

### 🖥️ Tahap 6: Setup Backend Golang (VS Code)
Bagian ini menjelaskan cara menjalankan server backend untuk Dashboard.

1.  **Buka Project di VS Code:**
    * Buka aplikasi **Visual Studio Code**.
    * Klik `File` -> `Open Folder`.
    * Pilih folder **`backend`** yang ada di dalam folder proyek ini.
    * *Struktur folder biasanya berisi: `main.go`, `go.mod`, `controllers/`, dll.*

2.  **Install Dependencies:**
    * Di VS Code, buka Terminal (`Ctrl + J` atau `Terminal` -> `New Terminal`).
    * Ketik perintah berikut untuk mengunduh library yang dibutuhkan:
        ```bash
        go mod tidy
        ```

3.  **Konfigurasi Koneksi:**
    * Cari file `config/config.go` atau `.env` (tergantung struktur).
    * Sesuaikan user/password database MySQL dan InfluxDB Token.

4.  **Jalankan Server:**
    * Di Terminal VS Code, ketik:
        ```bash
        go run main.go
        ```
    * Jika berhasil, akan muncul pesan: `Server running on port :8080`.
    * Buka browser dan akses: `http://localhost:8080` untuk melihat Dashboard.

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

```mermaid
graph TD
    A[ESP32 Transmitter] -->|LoRa RF| B[ESP32 Receiver/Gateway]
    B -->|MQTT/WiFi| C[Golang Backend]
    
    C -->|Telemetry Data| D[(InfluxDB)]
    C -->|User & Logs| E[(MySQL / MariaDB)]
    
    D --> F[Web Dashboard]
    E --> F
    
    style C fill:#0ff,stroke:#333,stroke-width:2px
