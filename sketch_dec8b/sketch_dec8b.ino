#include <WiFi.h>
#include <PubSubClient.h>
#include <WiFiClientSecure.h>
#include <Keypad.h>
#include <TM1637Display.h>
// KEYPAD
#define CLK  12 // The ESP32 pin GPIO22 connected to CLK
#define DIO  13 // The ESP32 pin GPIO23 connected to DIO

#define ROW_NUM     4 // four rows
#define COLUMN_NUM  3 // three columns

char keys[ROW_NUM][COLUMN_NUM] = {
  {'1', '2', '3'},
  {'4', '5', '6'},
  {'7', '8', '9'},
  {'*', '0', '#'}
};

byte pin_rows[ROW_NUM] = {5, 18, 19, 21}; 
byte pin_column[COLUMN_NUM] = {22, 4, 15};

const uint8_t waitLine[] = { SEG_G };

// create a display object of type TM1637Display
TM1637Display display = TM1637Display(CLK, DIO);

Keypad keypad = Keypad( makeKeymap(keys), pin_rows, pin_column, ROW_NUM, COLUMN_NUM );

// an array that sets individual segments per digit to display the word "dOnE"
const uint8_t done[] = {
  SEG_B | SEG_C | SEG_D | SEG_E | SEG_G,         // d
  SEG_A | SEG_B | SEG_C | SEG_D | SEG_E | SEG_F, // O
  SEG_C | SEG_E | SEG_G,                         // n
  SEG_A | SEG_D | SEG_E | SEG_F | SEG_G,         // E
};

const uint8_t conn[] = {
  SEG_A | SEG_D | SEG_E | SEG_F , // C
  SEG_C | SEG_D | SEG_E | SEG_G , // o
  SEG_C | SEG_E | SEG_G,          // n
  SEG_C | SEG_E | SEG_G,          // n
};

const uint8_t _open[] = {
  SEG_C | SEG_D | SEG_E | SEG_G ,                 // o
  SEG_A | SEG_B | SEG_D | SEG_E | SEG_F | SEG_G , // P
  SEG_A | SEG_D | SEG_E | SEG_F | SEG_G,          // E
  SEG_C | SEG_E | SEG_G,                          // n
};

const uint8_t eror[] = {
  SEG_A | SEG_D | SEG_E | SEG_F | SEG_G ,         // E
  SEG_A | SEG_E | SEG_F ,                         // r
  SEG_A | SEG_B | SEG_C | SEG_D | SEG_E | SEG_F , // O
  SEG_A | SEG_E | SEG_F ,                         // r
};

// BUZZER
#define BUZZER 23

// WiFi
const char *ssid = "ali"; // Enter your WiFi name
const char *password = "123456123456";  // Enter WiFi password

// MQTT Broker
const char *mqtt_broker = "f25aeeaa.ala.us-east-1.emqxsl.com";// broker address
const char *pub_topic = "es-project/esp/pub"; // define topic 
const char *sub_topic = "es-project/server/pub"; // define topic 
const char *mqtt_username = "test"; // username for authentication
const char *mqtt_password = "test";// password for authentication
const int mqtt_port = 8883;// port of MQTT over TLS/SSL

// load DigiCert Global Root CA ca_cert
const char* ca_cert= \
"-----BEGIN CERTIFICATE-----\n" \
"MIIDrzCCApegAwIBAgIQCDvgVpBCRrGhdWrJWZHHSjANBgkqhkiG9w0BAQUFADBh\n" \
"MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\n" \
"d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD\n" \
"QTAeFw0wNjExMTAwMDAwMDBaFw0zMTExMTAwMDAwMDBaMGExCzAJBgNVBAYTAlVT\n" \
"MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j\n" \
"b20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IENBMIIBIjANBgkqhkiG\n" \
"9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4jvhEXLeqKTTo1eqUKKPC3eQyaKl7hLOllsB\n" \
"CSDMAZOnTjC3U/dDxGkAV53ijSLdhwZAAIEJzs4bg7/fzTtxRuLWZscFs3YnFo97\n" \
"nh6Vfe63SKMI2tavegw5BmV/Sl0fvBf4q77uKNd0f3p4mVmFaG5cIzJLv07A6Fpt\n" \
"43C/dxC//AH2hdmoRBBYMql1GNXRor5H4idq9Joz+EkIYIvUX7Q6hL+hqkpMfT7P\n" \
"T19sdl6gSzeRntwi5m3OFBqOasv+zbMUZBfHWymeMr/y7vrTC0LUq7dBMtoM1O/4\n" \
"gdW7jVg/tRvoSSiicNoxBN33shbyTApOB6jtSj1etX+jkMOvJwIDAQABo2MwYTAO\n" \
"BgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUA95QNVbR\n" \
"TLtm8KPiGxvDl7I90VUwHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUw\n" \
"DQYJKoZIhvcNAQEFBQADggEBAMucN6pIExIK+t1EnE9SsPTfrgT1eXkIoyQY/Esr\n" \
"hMAtudXH/vTBH1jLuG2cenTnmCmrEbXjcKChzUyImZOMkXDiqw8cvpOp/2PV5Adg\n" \
"06O/nVsJ8dWO41P0jmP6P6fbtGbfYmbW0W5BjfIttep3Sp+dWOIrWcBAI+0tKIJF\n" \
"PnlUkiaY4IBIqDfv8NZ5YBberOgOzW6sRBc4L0na4UU+Krk2U886UAb3LujEV0ls\n" \
"YSEY1QSteDwsOoBrp+uvFRTp2InBuThs4pFsiv9kuXclVzDAGySj4dzp30d8tbQk\n" \
"CAUw7C29C79Fv1C5qfPrmAESrciIxpg0X40KPMbp1ZWVbd4=" \
"-----END CERTIFICATE-----\n";


// init secure wifi client
WiFiClientSecure espClient;
// use wifi client to init mqtt client
PubSubClient client(espClient); 

// varables
int num = 0;
int digitCount = 0;
int watchdog = 0; // end when is 100 

void setup() {
  display.clear();
  display.setBrightness(4); // set the brightness to 7 (0:dimmest, 7:brightest)
  display.setSegments(conn);
  pinMode(BUZZER, OUTPUT);
  // Set software serial baud to 115200;
  Serial.begin(115200);
  // connecting to a WiFi network
  connectWIFI();
  // set root ca cert
  espClient.setCACert(ca_cert);
  // setup mqtt broker
  client.setServer(mqtt_broker, mqtt_port);
  client.setCallback(callback);
  connectMQTT();
  display.setSegments(done);
  delay(1000);
  display.clear();
}

void connectWIFI(){
  WiFi.begin(ssid, password);
  if (WiFi.status() == WL_CONNECTED)
    return;

  Serial.println("Connecting to WiFi..");
  while (WiFi.status() != WL_CONNECTED) {
      delay(1000);
  }
  Serial.println("Connected to the WiFi network");
}

void connectMQTT() {
  if (client.connected())
    return;
    Serial.println("connecting to MQTT broker...");
  while (!client.connected()) {
    String client_id = "esp32-client-";
    client_id += String(WiFi.macAddress());
    if (client.connect(client_id.c_str(), mqtt_username, mqtt_password)) {
        Serial.println("connected to MQTT broker.");
    } else {
        Serial.print("Failed to connect to MQTT broker, rc=");
        Serial.print(client.state());
        Serial.println("Retrying in 2 seconds.");
        delay(2000);
    }
  }
  client.subscribe(sub_topic,2);
}

void callback(char* topic, byte* payload, unsigned int length) {
    Serial.print("Message arrived in topic: ");
    Serial.println(topic);
    Serial.print("Message:");
    for (int i = 0; i < length; i++) {
        Serial.print((char) payload[i]);
    }
    Serial.println();
    Serial.println("-----------------------");
    if (payload[0] == '1'){
      display.setSegments(_open);
    }else {
      display.setSegments(eror);
    }
    delay(2000);
      display.clear();
}

void clean(){
  num = 0;
  digitCount = 0;
  display.clear();
  return;
}

void handleKeypad(){
  char key = keypad.getKey();
  if (key) {
    if( key == '*') {
      clean();
      return;
    }
    else if(key == '#') {
      Serial.println("sending input ...");
      char temp[6];
      sprintf(temp, "%d", num);
      if (!client.connected()) {
          connectMQTT();
      }
      client.publish(pub_topic, temp);
      client.subscribe(sub_topic,0);
      Serial.println("published");
      clean();
      return;
    }
    else if( digitCount == 6){
      clean();
    }
    num*= 10;
    digitCount += 1;
    num+= key - 48;
    display.clear();
    display.showNumberDec(num);
    delay(500);
    Serial.print("input: ");
    Serial.println(num);
  }
}

void loop() {
  handleKeypad();
  client.loop();
}