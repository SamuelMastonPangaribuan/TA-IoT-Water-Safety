/*
 * SISTEM MONITORING KESELAMATAN PERAIRAN (TRANSMITTER - FINAL FIX)
 * Platform: Arduino IDE
 * Board: DOIT ESP32 DEVKIT V1
 */

#include <TinyGPSPlus.h>
#include "LoRa_E32.h"

// ================= KONFIGURASI PIN =================
#define E32_RX 16  // Ke TX LoRa
#define E32_TX 17  // Ke RX LoRa
LoRa_E32 e32ttl100(E32_RX, E32_TX, &Serial2, UART_BPS_RATE_9600);

#define GPS_RX 34  // Ke TX GPS
#define GPS_TX 12  // Ke RX GPS
HardwareSerial GPSserial(1);
TinyGPSPlus gps;

#define WATER_PIN   32
#define BUTTON_PIN  4
#define BUZZER_PIN  13
#define LED_TX      26
#define LED_SOS     27

// ================= VARIABEL GLOBAL =================
String NODE_ID = "NODE01"; 
bool isEmergency = false;

// Timer Variables
unsigned long lastSend = 0;
unsigned long lastBuzzerTime = 0;
unsigned long buttonPressStart = 0;

// [PENTING] Variabel Pengaman Tombol (Biar ga stuck)
bool hasReleasedButton = true; 

// ================= FUNGSI BANTUAN =================
int readWaterWet() {
  int v = analogRead(WATER_PIN);  
  return (v > 800) ? 1 : 0; 
}

void beepSOS() {
  if (millis() - lastBuzzerTime >= 500) { 
    digitalWrite(BUZZER_PIN, !digitalRead(BUZZER_PIN)); 
    lastBuzzerTime = millis();
  }
}

void kirimData() {
  digitalWrite(LED_TX, HIGH); 

  float lat = gps.location.isValid() ? gps.location.lat() : 0.0;
  float lon = gps.location.isValid() ? gps.location.lng() : 0.0;
  int w = readWaterWet();
  int s = isEmergency ? 1 : 0;

  // === [PERBAIKAN DI SINI] ===
  // Sebelumnya tertulis "drfs", sekarang sudah kembali jadi "s"
  String payload = "{\"id\":\"" + NODE_ID +
                   "\",\"lat\":" + String(lat, 6) +
                   ",\"lon\":" + String(lon, 6) +
                   ",\"w\":" + String(w) +
                   ",\"s\":" + String(s) + "}";

  ResponseStatus rs = e32ttl100.sendMessage(payload);
  Serial.println("[LORA TX] " + payload);

  delay(100); 
  digitalWrite(LED_TX, LOW);
}

// ================= SETUP =================
void setup() {
  Serial.begin(115200);
  GPSserial.begin(9600, SERIAL_8N1, GPS_RX, GPS_TX);
  e32ttl100.begin();

  pinMode(WATER_PIN, INPUT);
  pinMode(BUTTON_PIN, INPUT_PULLUP);
  pinMode(LED_TX, OUTPUT);
  pinMode(LED_SOS, OUTPUT);
  pinMode(BUZZER_PIN, OUTPUT);

  digitalWrite(LED_TX, LOW);
  digitalWrite(LED_SOS, LOW);
  digitalWrite(BUZZER_PIN, LOW);

  Serial.println("=== SYSTEM READY (ARDUINO IDE) ===");
}

// ================= MAIN LOOP =================
void loop() {
  while (GPSserial.available()) gps.encode(GPSserial.read());

  // Baca Tombol (LOW = Ditekan)
  int btnState = digitalRead(BUTTON_PIN);

  // 1. LOGIKA TOMBOL (ANTI-MACET & WAJIB LEPAS DULU)
  if (btnState == LOW) { 
    // === TOMBOL SEDANG DITEKAN ===
    
    if (buttonPressStart == 0) buttonPressStart = millis();
    unsigned long duration = millis() - buttonPressStart;

    // SKENARIO A: NYALAKAN SOS
    if (!isEmergency) {
      isEmergency = true;
      hasReleasedButton = false; // Kunci biar ga langsung reset
      
      Serial.println(">>> SOS AKTIF! <<<");
      digitalWrite(LED_SOS, HIGH);
      kirimData(); 
      delay(200); 
    }

    // SKENARIO B: MATIKAN SOS (RESET)
    // Syarat: Sedang SOS + Tahan 3 Detik + SUDAH LEPAS TOMBOL DULU
    else if (isEmergency && duration > 3000 && hasReleasedButton) {
      isEmergency = false;
      Serial.println(">>> RESET KE NORMAL <<<");
      
      // Bunyi Reset
      digitalWrite(BUZZER_PIN, HIGH); delay(1000); digitalWrite(BUZZER_PIN, LOW);
      digitalWrite(LED_SOS, LOW);
      
      kirimData();
      
      // Kunci lagi biar ga looping
      hasReleasedButton = false; 
    }

  } else {
    // === TOMBOL DILEPAS (HIGH) ===
    buttonPressStart = 0;
    
    // Izinkan fitur reset bekerja untuk penekanan berikutnya
    hasReleasedButton = true; 
  }

  // 2. OUTPUT HANDLING
  if (isEmergency) {
    beepSOS();
    digitalWrite(LED_SOS, HIGH);
  } else {
    digitalWrite(BUZZER_PIN, LOW);
    digitalWrite(LED_SOS, LOW);
  }

  // 3. JADWAL KIRIM DATA (Rutin 5 Detik)
  if (millis() - lastSend > 5000) {
    kirimData();
    lastSend = millis();
  }
}