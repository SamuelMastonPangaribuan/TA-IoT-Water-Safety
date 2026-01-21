package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	client "github.com/influxdata/influxdb1-client/v2"
)

// --- KONFIGURASI ---
const (
	DB_USER     = "root"
	DB_PASS     = ""
	DB_NAME     = "iot_safety"
	INFLUX_URL  = "http://localhost:8086"
	INFLUX_DB   = "iot_safety"
	MQTT_BROKER = "tcp://broker.emqx.io:1883"
	MQTT_TOPIC  = "samuel/project/ta"
)

var db *sql.DB
var influxClient client.Client

// --- STRUKTUR DATA ---
type SensorPayload struct {
	ID    interface{} `json:"id"`
	Lat   interface{} `json:"lat"`
	Lon   interface{} `json:"lon"`
	Water interface{} `json:"w"`
	SOS   interface{} `json:"s"`
}

type DeviceData struct {
	ID         int     `json:"id"`
	DeviceID   string  `json:"device_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	WaterLevel float64 `json:"water_level"`
	Status     string  `json:"status"`
	Timestamp  string  `json:"timestamp"`
	Owner      string  `json:"owner_snapshot"`
}

type AuthData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AssignData struct {
	DeviceID  string `json:"device_id"`
	OwnerName string `json:"owner_name"`
}

// --- MAIN FUNCTION ---
func main() {
	fmt.Println("--- MEMULAI SISTEM BACKEND IOT ---")

	// 1. KONEKSI MYSQL
	var err error
	// Pastikan port sesuai XAMPP (biasanya 3306 atau 3307)
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3307)/%s?parseTime=true", DB_USER, DB_PASS, DB_NAME)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Config MySQL Salah:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Gagal Konek MySQL:", err)
	}
	fmt.Println("✅ MySQL Terhubung!")

	// 2. KONEKSI INFLUXDB
	influxClient, err = client.NewHTTPClient(client.HTTPConfig{Addr: INFLUX_URL})
	if err != nil {
		log.Fatal("Gagal Config Influx:", err)
	}
	_, _, err = influxClient.Ping(0)
	if err != nil {
		fmt.Println("⚠️ InfluxDB tidak merespon (Pastikan service jalan)")
	} else {
		fmt.Println("✅ InfluxDB Terhubung!")
		q := client.NewQuery("CREATE DATABASE "+INFLUX_DB, "", "")
		influxClient.Query(q)
	}

	// 3. BACKGROUND TASKS
	go startMQTT()
	go startSimulation()

	// 4. API SERVER (GIN)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	// --- ENDPOINTS ---

	// Login
	r.POST("/api/login", func(c *gin.Context) {
		var u AuthData
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(400, gin.H{"status": "fail"})
			return
		}
		var dbPass string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", u.Username).Scan(&dbPass)
		if err != nil || dbPass != u.Password {
			c.JSON(200, gin.H{"status": "fail", "message": "Login Gagal"})
		} else {
			c.JSON(200, gin.H{"status": "success", "user": u.Username})
		}
	})

	// Assign Device (Daftar Pemilik)
	r.POST("/api/assign-device", func(c *gin.Context) {
		var d AssignData
		if err := c.ShouldBindJSON(&d); err != nil {
			c.JSON(400, gin.H{"status": "fail"})
			return
		}
		// Insert atau Update jika sudah ada
		_, err := db.Exec("INSERT INTO device_owners (device_id, owner_name) VALUES (?, ?) ON DUPLICATE KEY UPDATE owner_name = VALUES(owner_name)", d.DeviceID, d.OwnerName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": "success"})
		}
	})

	// === API HAPUS PEMILIK (UNASSIGN) ===
	r.POST("/api/unassign-device", func(c *gin.Context) {
		var d struct {
			DeviceID string `json:"device_id"`
		}
		if err := c.ShouldBindJSON(&d); err != nil {
			c.JSON(400, gin.H{"status": "fail"})
			return
		}

		// 1. Hapus dari tabel pemilik
		_, err := db.Exec("DELETE FROM device_owners WHERE device_id = ?", d.DeviceID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Gagal Hapus DB"})
			return
		}

		// 2. Catat Log ke History bahwa alat ini di-reset
		_, err = db.Exec("INSERT INTO sensor_data (device_id, latitude, longitude, water_level, status, owner_snapshot) VALUES (?, 0, 0, 0, 'INFO: UNREGISTERED', 'SYSTEM')", d.DeviceID)
		if err != nil {
			fmt.Println("Gagal catat log unregister:", err)
		}

		c.JSON(200, gin.H{"status": "success"})
	})

	// Get Owner (Cek Pemilik Live)
	r.GET("/api/get-owner", func(c *gin.Context) {
		id := c.Query("id")
		var owner string
		err := db.QueryRow("SELECT owner_name FROM device_owners WHERE device_id = ?", id).Scan(&owner)
		if err != nil {
			c.JSON(200, gin.H{"owner": "Belum Terdaftar"})
		} else {
			c.JSON(200, gin.H{"owner": owner})
		}
	})

	// Latest Data
	r.GET("/api/latest-data", func(c *gin.Context) {
		rows, err := db.Query(`SELECT t1.id, t1.device_id, t1.latitude, t1.longitude, t1.water_level, t1.status, t1.timestamp, COALESCE(t1.owner_snapshot, 'Unknown') FROM sensor_data t1 JOIN (SELECT device_id, MAX(id) as max_id FROM sensor_data GROUP BY device_id) t2 ON t1.id = t2.max_id`)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var results []DeviceData
		for rows.Next() {
			var d DeviceData
			var ts []uint8
			rows.Scan(&d.ID, &d.DeviceID, &d.Latitude, &d.Longitude, &d.WaterLevel, &d.Status, &ts, &d.Owner)
			d.Timestamp = string(ts)
			results = append(results, d)
		}
		if results == nil {
			results = []DeviceData{}
		}
		c.JSON(200, results)
	})

	// History Data (Tanpa LIMIT agar semua muncul)
	r.GET("/api/history", func(c *gin.Context) {
		id := c.Query("id")
		rows, err := db.Query("SELECT id, device_id, latitude, longitude, water_level, status, timestamp, owner_snapshot FROM sensor_data WHERE device_id = ? ORDER BY id DESC", id)
		if err != nil {
			c.JSON(500, []DeviceData{})
			return
		}
		defer rows.Close()
		var results []DeviceData
		for rows.Next() {
			var d DeviceData
			var ts []uint8
			var owner sql.NullString
			rows.Scan(&d.ID, &d.DeviceID, &d.Latitude, &d.Longitude, &d.WaterLevel, &d.Status, &ts, &owner)
			d.Timestamp = string(ts)
			if owner.Valid {
				d.Owner = owner.String
			} else {
				d.Owner = "Unknown"
			}
			results = append(results, d)
		}
		if results == nil {
			results = []DeviceData{}
		}
		c.JSON(200, results)
	})

	fmt.Println("🚀 Backend Golang Berjalan di Port 1880...")
	r.Run(":1880")
}

// --- FUNGSI HYBRID STORAGE ---
func saveDataHybrid(id string, lat, lon float64, w, s int) {
	status := "AMAN"
	if s == 1 {
		status = "BAHAYA (SOS)"
	} else if w == 1 {
		status = "WASPADA (BASAH)"
	}

	query := `INSERT INTO sensor_data (device_id, latitude, longitude, water_level, status, owner_snapshot) 
              SELECT ?, ?, ?, ?, ?, COALESCE(MAX(owner_name), 'Belum Terdaftar') 
              FROM device_owners WHERE device_id = ?`
	_, err := db.Exec(query, id, lat, lon, w, status, id)
	if err != nil {
		fmt.Println("❌ MySQL Error:", err)
	}

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: INFLUX_DB, Precision: "s"})
	tags := map[string]string{"deviceId": id}
	fields := map[string]interface{}{"lat": lat, "lon": lon, "water_val": w, "sos_val": s, "status_desc": status}
	pt, err := client.NewPoint("water_safety", tags, fields, time.Now())
	if err == nil {
		bp.AddPoint(pt)
		influxClient.Write(bp)
	}
	fmt.Printf("💾 [Dual-DB] %s -> Data Saved\n", id)
}

// --- SIMULASI ---
func startSimulation() {
	fmt.Println("🤖 Simulasi Aktif...")
	for range time.Tick(5 * time.Second) {
		offLat := (rand.Float64() * 0.005) - 0.0025
		offLon := (rand.Float64() * 0.005) - 0.0025
		saveDataHybrid("NODE02", 2.383+offLat, 99.148+offLon, 0, 0)
		saveDataHybrid("NODE03", 2.386+offLat, 99.145+offLon, 1, 1)
	}
}

// --- MQTT ---
func startMQTT() {
	opts := mqtt.NewClientOptions().AddBroker(MQTT_BROKER)
	opts.SetClientID("Go-Backend-" + fmt.Sprint(rand.Int()))
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		payloadStr := string(msg.Payload())
		fmt.Printf("📡 MQTT: %s\n", payloadStr)
		var raw SensorPayload
		if err := json.Unmarshal([]byte(payloadStr), &raw); err == nil {
			id := fmt.Sprintf("%v", raw.ID)
			if id == "" || id == "<nil>" {
				id = "Unknown"
			}
			saveDataHybrid(id, toFloat(raw.Lat), toFloat(raw.Lon), toInt(raw.Water), toInt(raw.SOS))
		}
	})
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("❌ MQTT Error:", token.Error())
		return
	}
	fmt.Println("✅ MQTT Terhubung!")
	client.Subscribe(MQTT_TOPIC, 0, nil)
}

func toFloat(v interface{}) float64 {
	switch i := v.(type) {
	case float64:
		return i
	case int:
		return float64(i)
	case string:
		f, _ := strconv.ParseFloat(i, 64)
		return f
	default:
		return 0
	}
}

func toInt(v interface{}) int {
	switch i := v.(type) {
	case float64:
		return int(i)
	case int:
		return i
	case string:
		v, _ := strconv.Atoi(i)
		return v
	default:
		return 0
	}
}
