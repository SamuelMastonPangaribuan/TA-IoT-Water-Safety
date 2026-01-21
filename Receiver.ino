  #include <WiFi.h>
  #include <PubSubClient.h>
  #include "LoRa_E32.h"

  // === WIFI CONFIG ===
  const char* ssid = "Sipangalo";
  const char* password = "00001111";

  // === MQTT CONFIG ===
  const char* mqtt_server = "broker.emqx.io";
  const int mqtt_port = 1883;
  const char* mqtt_topic = "samuel/project/ta";

  // === LoRa E32 ===
  #define E32_RX 16
  #define E32_TX 17
  LoRa_E32 e32ttl100(E32_RX, E32_TX, &Serial2, UART_BPS_RATE_9600);

  // === [TAMBAHAN] KONFIGURASI BUZZER ===
  #define BUZZER_PIN 4  // Sambungkan Positif Buzzer ke Pin D4, Negatif ke GND

  WiFiClient espClient;
  PubSubClient client(espClient);

  // -------------------------------------
  // [TAMBAHAN] FUNGSI BUNYI ALARM
  // -------------------------------------
  void bunyikanAlarm() {
    Serial.println("!!! PERINGATAN: SOS DITERIMA - BUZZER BUNYI !!!");
    // Pola bunyi: Tit-Tit-Tit (3x Cepat)
    for (int i = 0; i < 3; i++) {
      digitalWrite(BUZZER_PIN, HIGH);
      delay(100);
      digitalWrite(BUZZER_PIN, LOW);
      delay(100);
    }
  }

  // -------------------------------------
  // WIFI Connect
  // -------------------------------------
  void setupWiFi() {
    Serial.print("Menghubungkan ke WiFi: ");
    Serial.println(ssid);
    WiFi.begin(ssid, password);
    while (WiFi.status() != WL_CONNECTED) {
      delay(300);
      Serial.print(".");
    }
    Serial.println("\nWiFi Terhubung!");
    Serial.print("IP Address: ");
    Serial.println(WiFi.localIP());
  }

  // -------------------------------------
  // MQTT Reconnect
  // -------------------------------------
  boolean reconnectMQTT() {
    if (client.connect("ESP32-GATEWAY")) {
      Serial.println("Terhubung ke MQTT Broker!");
    }
    return client.connected();
  }

  // -------------------------------------
  // SETUP
  // -------------------------------------
  void setup() {
    Serial.begin(115200);
    e32ttl100.begin();

    // === [TAMBAHAN] SETUP BUZZER ===
    pinMode(BUZZER_PIN, OUTPUT);
    digitalWrite(BUZZER_PIN, LOW); // Pastikan mati di awal

    setupWiFi();
    client.setServer(mqtt_server, mqtt_port);

    Serial.println("====================================");
    Serial.println("        GATEWAY RECEIVER LoRa       ");
    Serial.println("====================================");
    
    // Test Buzzer sebentar (Tanda nyala)
    digitalWrite(BUZZER_PIN, HIGH); delay(100); digitalWrite(BUZZER_PIN, LOW);

    Serial.println("Menunggu data JSON dari Transmitter...");
    Serial.println("------------------------------------");
  }

  // -------------------------------------
  // LOOP
  // -------------------------------------
  void loop() {
    // MQTT Loop
    if (!client.connected()) reconnectMQTT();
    client.loop();

    // Cek data LoRa
    if (e32ttl100.available() > 1) {
      ResponseContainer rc = e32ttl100.receiveMessage();
      if (rc.status.code != 1) return;

      String data = rc.data;
      data.trim(); // Hapus spasi/enter di awal/akhir

      // Hapus baris ini jika ingin serial monitor bersih
      // Serial.println("Data Masuk: " + data);

      // === VALIDASI JSON ===
      if (!data.startsWith("{") || !data.endsWith("}")) {
        Serial.println("Format data salah (Bukan JSON), diabaikan.");
        return;
      }

      // === PARSING JSON MANUAL ===
      int idStart = data.indexOf("\"id\":\"");
      int latStart = data.indexOf("\"lat\":");
      int lonStart = data.indexOf("\"lon\":");
      int wStart = data.indexOf("\"w\":");
      int sStart = data.indexOf("\"s\":");

      if (idStart < 0 || latStart < 0 || lonStart < 0 || wStart < 0 || sStart < 0) {
        Serial.println("JSON tidak lengkap!");
        return;
      }

      // Ekstrak Nilai
      String id = data.substring(idStart + 6, data.indexOf("\"", idStart + 6));
      float lat = data.substring(latStart + 6, data.indexOf(",", latStart)).toFloat();
      float lon = data.substring(lonStart + 6, data.indexOf(",", lonStart)).toFloat();
      int w = data.substring(wStart + 4, data.indexOf(",", wStart)).toInt();
      int s = data.substring(sStart + 4, data.indexOf("}", sStart)).toInt();

      // === TAMPILKAN DEBUG ===
      Serial.println(">> PARSED DATA:");
      Serial.println("ID    : " + id);  
      Serial.printf("Lat   : %.6f\n", lat);
      Serial.printf("Lon   : %.6f\n", lon);
      Serial.printf("Water : %d\n", w);
      Serial.printf("SOS   : %d", s);
      
      // === [TAMBAHAN] LOGIKA BUZZER ===
      if (s == 1) {
        Serial.println(" [BAHAYA - ALARM BUNYI]");
        bunyikanAlarm(); // Panggil fungsi alarm
      } else {
        Serial.println(" [AMAN]");
      }
      Serial.println("----------------------");

      // === KIRIM KE MQTT ===
      String mqttPayload =
        "{\"id\":\"" + id +
        "\",\"lat\":" + String(lat, 6) +
        ",\"lon\":" + String(lon, 6) +
        ",\"w\":" + String(w) +
        ",\"s\":" + String(s) + "}";

      client.publish(mqtt_topic, mqttPayload.c_str());
      Serial.println(">> Sukses Publish ke MQTT");
    }
  }