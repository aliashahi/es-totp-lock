#include <WiFi.h>
#include <PubSubClient.h>
#include <WiFiClientSecure.h>
#include <Keypad.h>
#include <LiquidCrystal.h>

//////////PINS/////////////////////
// #define CLK 12
// #define DIO 13

#define ROW1 13
#define ROW2 12
#define ROW3 14
#define ROW4 27

#define COL1 26
#define COL2 25
#define COL3 33

#define LCD_BRIGHTNESS 15

#define LCD_RS 2
#define LCD_EN 4
#define LCD_D0 5
#define LCD_D1 18
#define LCD_D2 19
#define LCD_D3 21

// #define BUZZER 23
#define SLOCK 22
#define RED_LIGHT 23
///////////////////////////////////

#define ROW_NUM 4     // four rows
#define COLUMN_NUM 3  // three columns
char keys[ROW_NUM][COLUMN_NUM] = {
  { '1', '2', '3' },
  { '4', '5', '6' },
  { '7', '8', '9' },
  { '*', '0', '#' }
};

byte pin_rows[ROW_NUM] = { ROW1, ROW2, ROW3, ROW4 };
byte pin_column[COLUMN_NUM] = { COL1, COL2, COL3 };

Keypad keypad = Keypad(makeKeymap(keys), pin_rows, pin_column, ROW_NUM, COLUMN_NUM);

// LCD
LiquidCrystal lcd(LCD_RS, LCD_EN, LCD_D0, LCD_D1, LCD_D2, LCD_D3);

// WiFi
const char *ssid = "ali";               // Enter your WiFi name
const char *password = "123456123456";  // Enter WiFi password

// MQTT Broker
const char *mqtt_broker = "f25aeeaa.ala.us-east-1.emqxsl.com";  // broker address
const char *pub_topic = "es-project/esp/pub";                   // define topic
const char *sub_topic = "es-project/server/pub";                // define topic
const char *mqtt_username = "test";                             // username for authentication
const char *mqtt_password = "test";                             // password for authentication
const int mqtt_port = 8883;
const int qos = 0;

// load DigiCert Global Root CA ca_cert
const char *ca_cert =
  "-----BEGIN CERTIFICATE-----\n"
  "MIIDrzCCApegAwIBAgIQCDvgVpBCRrGhdWrJWZHHSjANBgkqhkiG9w0BAQUFADBh\n"
  "MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\n"
  "d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD\n"
  "QTAeFw0wNjExMTAwMDAwMDBaFw0zMTExMTAwMDAwMDBaMGExCzAJBgNVBAYTAlVT\n"
  "MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j\n"
  "b20xIDAeBgNVBAMTF0RpZ2lDZXJ0IEdsb2JhbCBSb290IENBMIIBIjANBgkqhkiG\n"
  "9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4jvhEXLeqKTTo1eqUKKPC3eQyaKl7hLOllsB\n"
  "CSDMAZOnTjC3U/dDxGkAV53ijSLdhwZAAIEJzs4bg7/fzTtxRuLWZscFs3YnFo97\n"
  "nh6Vfe63SKMI2tavegw5BmV/Sl0fvBf4q77uKNd0f3p4mVmFaG5cIzJLv07A6Fpt\n"
  "43C/dxC//AH2hdmoRBBYMql1GNXRor5H4idq9Joz+EkIYIvUX7Q6hL+hqkpMfT7P\n"
  "T19sdl6gSzeRntwi5m3OFBqOasv+zbMUZBfHWymeMr/y7vrTC0LUq7dBMtoM1O/4\n"
  "gdW7jVg/tRvoSSiicNoxBN33shbyTApOB6jtSj1etX+jkMOvJwIDAQABo2MwYTAO\n"
  "BgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUA95QNVbR\n"
  "TLtm8KPiGxvDl7I90VUwHwYDVR0jBBgwFoAUA95QNVbRTLtm8KPiGxvDl7I90VUw\n"
  "DQYJKoZIhvcNAQEFBQADggEBAMucN6pIExIK+t1EnE9SsPTfrgT1eXkIoyQY/Esr\n"
  "hMAtudXH/vTBH1jLuG2cenTnmCmrEbXjcKChzUyImZOMkXDiqw8cvpOp/2PV5Adg\n"
  "06O/nVsJ8dWO41P0jmP6P6fbtGbfYmbW0W5BjfIttep3Sp+dWOIrWcBAI+0tKIJF\n"
  "PnlUkiaY4IBIqDfv8NZ5YBberOgOzW6sRBc4L0na4UU+Krk2U886UAb3LujEV0ls\n"
  "YSEY1QSteDwsOoBrp+uvFRTp2InBuThs4pFsiv9kuXclVzDAGySj4dzp30d8tbQk\n"
  "CAUw7C29C79Fv1C5qfPrmAESrciIxpg0X40KPMbp1ZWVbd4="
  "-----END CERTIFICATE-----\n";

// init secure wifi client
WiFiClientSecure espClient;
// use wifi client to init mqtt client
PubSubClient client(espClient);

// varables
char num[7] = "000000";
int digitCount = 0;

int state = 0;
int watchdog = 0;

void connectWIFI() {
  int i = 0;
  if (WiFi.status() == WL_CONNECTED)
    return;

  lcd.clear();
  lcd.setCursor(2, 0);
  lcd.print("Connecting to");
  lcd.setCursor(5, 1);
  lcd.print("WiFi ");

  Serial.println("Connecting to WiFi...");
  while (WiFi.status() != WL_CONNECTED) {
    if (i == 3) {
      i = 0;
      lcd.setCursor(5, 1);
      lcd.print("WiFi    ");
    } else {
      lcd.setCursor(9 + i, 1);
      lcd.print(".");
      i += 1;
    }
    delay(1000);
  }
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("WiFi Connected");
  delay(1000);
  Serial.println("Connected to the WiFi");
}

void connectMQTT() {
  int i = 0;

  connectWIFI();

  if (client.connected()) {
    client.subscribe(sub_topic, qos);
    return;
  }

  lcd.clear();
  lcd.setCursor(2, 0);
  lcd.print("Connecting to");
  lcd.setCursor(5, 1);
  lcd.print("MQTT");

  Serial.println("connecting to MQTT broker...");
  while (!client.connected()) {
    String client_id = "esp32-client-";
    client_id += String(WiFi.macAddress());
    if (client.connect(client_id.c_str(), mqtt_username, mqtt_password)) {
      lcd.clear();
      lcd.setCursor(0, 0);
      lcd.print("MQTT Connected");
      Serial.println("connected to MQTT broker.");
      delay(2000);
      break;
    } else {
      Serial.print("Failed to connect to MQTT broker, rc=");
      Serial.print(client.state());
      Serial.println("Retrying in 2 seconds.");
    }

    if (i == 3) {
      i = 0;
      lcd.setCursor(5, 1);
      lcd.print("MQTT   ");
      continue;
    }
    lcd.setCursor(9 + i, 1);
    lcd.print(".");
    i += 1;
    delay(2000);
  }

  client.subscribe(sub_topic, qos);
}

void callback(char *topic, byte *payload, unsigned int length) {
  if(state != 3)return;

  Serial.print("Message arrived in \n topic: ");
  Serial.println(topic);
  Serial.print("Message:");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();
  Serial.println("--------------------------------");
  if (payload[0] == '1') {
    lcd.clear();
    lcd.setCursor(0, 0);
    lcd.print("Wellcome");
    lcd.setCursor(0, 1);
    for (int i = 1; i < length && i < 16; i++) {
      lcd.print((char)payload[i]);
    }
    state = 4;
  } else {
    state = 5;
  }
}


void clean() {
  for (int i = 0; i < 6; i++) {
    num[i] = '0';
  }
  digitCount = 0;
  state = 1;
}

void readingKeypad() {
  char key = keypad.getKey();
  if (key) {
    if (key == '*') {
      clean();
      return;
    } else if (key == '#') {
      return;
    }
    num[digitCount] = key;
    digitCount += 1;
    lcd.print(key);
    Serial.print("input: ");
    Serial.println(num);
  }
  if (digitCount == 6) {
    delay(500);
    state = 2;
  }
}

void openDoor() {
  // digitalWrite(BUZZER, HIGH);
  analogWrite(SLOCK, 100);
  delay(500);
  // digitalWrite(BUZZER, LOW);
  delay(3000);
  analogWrite(SLOCK, 0);
  clean();
  state = 1;
}

void sendingCode() {
  Serial.println("sending input ...");

  if (!client.connected()) {
    connectMQTT();
  }

  lcd.clear();
  lcd.setCursor(4, 0);
  lcd.print("sending...");

  while (!client.publish(pub_topic, num)) {
    delay(1000);
  }

  client.subscribe(sub_topic, qos);
  Serial.println("published");
  watchdog = 0;
  state = 3;
}

void waitForResponse() {
  if(watchdog == 20){
    state = 6;
    return;
  }
  watchdog+= 1;
  // real wait is in callback function;
  lcd.clear();
  lcd.setCursor(5, 0);
  lcd.print("waiting");
  
  connectMQTT();
  client.loop();
  delay(500);
}

void showError() {
  // digitalWrite(BUZZER,HIGH);
  digitalWrite(RED_LIGHT,HIGH);
  lcd.clear();
  lcd.setCursor(2, 0);
  lcd.print("unauthorized!");
  delay(2000);
  // digitalWrite(BUZZER,LOW);
  digitalWrite(RED_LIGHT,LOW);
  clean();
  state = 1;
}

void timeOut() {
  // digitalWrite(BUZZER,HIGH);
  digitalWrite(RED_LIGHT,HIGH);
  lcd.clear();
  lcd.setCursor(4, 0);
  lcd.print("timeout!");
  delay(2000);
  clean();
  // digitalWrite(BUZZER,LOW);
  digitalWrite(RED_LIGHT,LOW);
  state = 1;
}

void setup() {
  pinMode(RED_LIGHT, OUTPUT);
  pinMode(SLOCK, OUTPUT);
  digitalWrite(SLOCK, 0);

  pinMode(LCD_BRIGHTNESS, OUTPUT);
  analogWrite(LCD_BRIGHTNESS, 100);

  lcd.begin(16, 2);
  lcd.setCursor(4, 0);
  lcd.print("TOTP Lock");
  lcd.setCursor(5,1);
  lcd.print("v1.0.0");
  // Set software serial baud to 115200;
  Serial.begin(115200);
  // connecting to a WiFi network
  WiFi.begin(ssid, password);
  // connectWIFI();
  // setup mqtt broker
  espClient.setCACert(ca_cert);
  client.setServer(mqtt_broker, mqtt_port);
  client.setCallback(callback);
  // connectMQTT();
  delay(3000);
  lcd.clear();
  state = 1;
}

void loop() {
  switch (state) {
    case 1:
      lcd.clear();
      lcd.setCursor(0, 0);
      lcd.print("enter passcode:");
      lcd.setCursor(5, 1);
      state = 11;
      break;
    case 11:
      readingKeypad();
      break;
    case 2:
      sendingCode();
      break;
    case 3:
      waitForResponse();
      break;
    case 4:
      openDoor();
      break;
    case 5:
      showError();
      break;
    case 6:
      timeOut();
      break;
  }
}