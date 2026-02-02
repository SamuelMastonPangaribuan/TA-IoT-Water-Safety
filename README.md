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
[🛠️ **Installation Guide**](#-installation--setup-guide-from-scratch) •
[📊 Test Results](#-hasil-validasi-lapangan-field-test) •
[🔌 Wiring](#-pin-mapping-wiring-configuration) •
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

Panduan ini disusun secara berurutan mulai dari persiapan *hardware*, instalasi *software*, konfigurasi *server*, hingga alat siap digunakan.

### 📦 Tahap 1: Persiapan Hardware
Pastikan Anda memiliki komponen berikut sebelum memulai:
1.  **2x ESP32 DevKit V1** (Satu untuk Transmitter/Korban, satu untuk Receiver/Pos Pantau).
2.  **2x Modul LoRa E32-TTL-100** (Pastikan frekuensi sama, misal 433MHz atau 915MHz).
3.  **1x Modul GPS Neo-6M** (Beserta antena keramik).
4.  **1x Water Level Sensor** (Tipe resistif/garis-garis).
5.  **1x Push Button** (Untuk tombol SOS).
6.  **1x Active Buzzer** & LED Indikator.
7.  **Kabel Jumper** & Breadboard/PCB.

> **Note:** Lakukan perakitan sesuai dengan diagram **[🔌 Pin Mapping](#-pin-mapping-wiring-configuration)** di bawah.

### 💻 Tahap 2: Persiapan Environment (Laptop/PC)
1.  **Install Arduino IDE:** Unduh versi terbaru di [arduino.cc](https://www.arduino.cc/en/software).
2.  **Install Driver USB:**
    * Jika ESP32 tidak terbaca, install driver **CP210x** atau **CH340** (sesuai chip USB di board ESP32 Anda).
3.  **Setup Board ESP32:**
    * Buka Arduino IDE -> `File` -> `Preferences`.
    * Isi *Additional Boards Manager URLs*: `https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json`
    * Buka `Tools` -> `Board` -> `Boards Manager`, cari **"esp32"** dan klik **Install**.

### 📚 Tahap 3: Instalasi Library (Wajib)
Tanpa library ini, kode tidak akan bisa di-compile. Install melalui `Sketch` -> `Include Library` -> `Manage Libraries`:

| Nama Library | Penulis | Fungsi |
| :--- | :--- | :--- |
| **TinyGPSPlus** | Mikal Hart | Parsing data NMEA dari modul GPS |
| **LoRa_E32** | KrisKasprzak | Komunikasi UART dengan modul E32 |
| **PubSubClient** | Nick O'Leary | Protokol MQTT (untuk Receiver ke Server) |

### 🗄️ Tahap 4: Konfigurasi Database & Backend
Sistem ini membutuhkan server lokal/cloud untuk dashboard.
1.  **MySQL (XAMPP/MariaDB):**
    * Buka phpMyAdmin.
    * Buat database baru bernama `db_watersafety`.
    * Import file SQL yang tersedia di folder `database/db_watersafety.sql`.
2.  **InfluxDB (Time Series):**
    * Install InfluxDB.
    * Buat Organization: `watersafety_org`.
    * Buat Bucket: `sensor_data`.
    * Salin **API Token** untuk dimasukkan ke konfigurasi Node-RED/Backend.

### 🚀 Tahap 5: Kompilasi & Upload Firmware
Lakukan langkah ini dua kali (sekali untuk alat korban, sekali untuk alat pos pantau).

**A. Upload ke Transmitter (Alat Korban):**
1.  Buka file `Transmitter_Final.ino`.
2.  Sambungkan ESP32 Transmitter ke PC.
3.  Pilih Board: `DOIT ESP32 DEVKIT V1`.
4.  Cek Port: `Tools` -> `Port` (Pilih COM yang aktif).
5.  Klik **Upload (➡️)**. Tunggu sampai "Done uploading".

**B. Upload ke Receiver (Gateway Pos Pantau):**
1.  Buka file `Receiver_Gateway.ino`.
2.  **Edit Konfigurasi WiFi:** Cari baris `const char* ssid = "..."` dan ubah sesuai hotspot/WiFi Anda.
3.  Sambungkan ESP32 Receiver ke PC.
4.  Klik **Upload (➡️)**.

### ✅ Tahap 6: Pengujian Sistem (Running Test)
1.  **Nyalakan Serial Monitor** di Arduino IDE (Baudrate **115200**).
2.  Pastikan modul LoRa sudah terhubung (Pesan: *"LoRa Init Success"*).
3.  Bawa alat Transmitter ke luar ruangan agar GPS mendapat sinyal (*GPS Lock*).
4.  Tekan tombol SOS atau celupkan sensor air.
5.  Cek apakah data muncul di Dashboard/Database.

---

## 📊 Hasil Validasi Lapangan (Field Test)

Perangkat telah diuji di 3 lokasi berbeda untuk memvalidasi akurasi modul GPS Neo-6M. Berikut adalah ringkasan hasil pengujian:

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
