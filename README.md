# go-rajaongkir-location

**go-rajaongkir-location** is a lightweight and high-performance API service built in Go, designed to replicate the location-based endpoints of [RajaOngkir](https://rajaongkir.com) — including **provinces**, **cities**, and **subdistricts** — with improved performance and flexibility. Additionally, it provides a **cost calculation endpoint** that bridges RajaOngkir's v1 and v2 APIs for backward compatibility.

## 🚀 Why This Exists

This project was meant to replace the location service provided by RajaOngkir, which can sometimes be slow or less responsive. By self-hosting the static location data (SQLite), this service mimics the same response format while significantly boosting performance and reliability.

**New in this version**: With RajaOngkir deprecating their v1 API, this service now includes a `/cost` endpoint that accepts v1 API format requests and seamlessly translates them to v2 API calls, ensuring your legacy applications continue to work without modification.

It's ideal for internal use or as a drop-in replacement where you need both location data and shipping cost calculations.

---

## 📦 Features

- ⚡ **Fast & Lightweight** – Built with [Gin](https://github.com/gin-gonic/gin) and optimized for blazing fast response times.
- 🗃️ **SQLite-Powered** – Uses SQLite for lightweight and portable local data storage — no external database needed.
- 🔁 **RajaOngkir-Compatible** – Mirrors the structure and response format of RajaOngkir for seamless integration.
- 📌 **Complete Location & Cost Data** – Provides **provinces**, **cities**, **subdistricts**, and **shipping cost calculations**.
- 🔄 **API Version Bridge** – Translates v1 API requests to v2 format automatically for backward compatibility.
- 🚚 **Multi-Courier Support** – Supports all major couriers: JNE, POS, TIKI, Wahana, SiCepat, J&T.
- 🔐 **Flexible Authentication** – Secure endpoints with a primary `API_KEY` while supporting both modern `Bearer` tokens and legacy single-key headers for backward compatibility.
- 🧩 **Easy to Integrate** – Acts as a drop-in replacement for RajaOngkir APIs in your frontend or backend.

---

## 📚 API Endpoints

All responses are wrapped in a `rajaongkir` object just like the original API.

### 🔍 Get All Provinces (with optional filter by ID)
```
GET /province?id={province_id}
```

### 🔍 Get All Cities (with optional province filter)
```
GET /city?province={province_id}
```

### 🔍 Get All Subdistricts (with optional city filter)
```
GET /subdistrict?city={city_id}
```

### 💰 Calculate Shipping Cost (v1 API Compatible)

```
POST /cost
Content-Type: application/x-www-form-urlencoded

origin=2089 (subdistrict id)
destination=2088 (subdistrict id)
weight=1000 (gram)
courier=jne
originType=subdistrict
destinationType=subdistrict
```

**Supported Couriers:**
The supported couriers was based on RajaOngkir's v2 API
- `jne`  -  Jalur Nugraha Ekakurir (JNE)
- `sicepat`  -  SiCepat Express (SICEPAT)
- `ide`  -  ID Express
- `sap`  -  SAP Express
- `jnt`  -  J&T Express (J&T)
- `ninja`  -  Ninja Expres
- `tiki`  -  Citra Van Titipan Kilat (TIKI)
- `lion`  -  Lion Parcel
- `anteraja`  -  Anteraja
- `pos`  -  POS Indonesia (POS)
- `ncs`  -  Nusantara Card Semesta
- `rex`  -  REX Express,
- `rpx`  -  RPX One Stop Logistics,
- `sentral`  -  Sentral Cargo,
- `star`  -  STAR Cargo,
- `wahana`  -  Wahana Prestasi Logistik (WAHANA)
- `dse`  -  DSE Logistic,

## 🔐 Authentication

If an `API_KEY` is set in your `.env` file, all requests to the service must be authenticated. The service supports two methods to handle both new and legacy applications.

### Method 1: Bearer Token (Recommended)

This is the standard method for new applications. It requires **two headers**:

1.  `Authorization`: The `API_KEY` you set in the `.env` file for your service, prefixed with `Bearer `.
2.  `rajaongkir-key`: Your actual API key from RajaOngkir, which is needed for the `/cost` endpoint.

**Example Request:**
```http
Authorization: Bearer your-secret-service-key
rajaongkir-key: your-personal-rajaongkir-api-key

---

## 🔄 How the Cost API Works

The `/cost` endpoint acts as a bridge between RajaOngkir's v1 and v2 APIs:

1. **Receives v1 format requests** with subdistrict IDs
2. **Looks up postal codes** from the local SQLite database
3. **Translates to v2 format** and calls RajaOngkir's v2 API
4. **Transforms the response** back to v1 format
5. **Returns familiar v1 response** to your application

This ensures your legacy applications continue working even after RajaOngkir discontinues their v1 API.

---

## 📦 Response Format

### Location Endpoints Response
```json
{
  "rajaongkir": {
    "status": {
      "code": 200,
      "description": "OK"
    },
    "results": [
      {
        "province_id": "1",
        "province": "Bali"
      },
      {
        "province_id": "2",
        "province": "Bangka Belitung"
      }
    ]
  }
}
```

### Cost Endpoint Response
```json
{
  "rajaongkir": {
    "status": {
      "code": 200,
      "description": "OK"
    },
    "results": [
      {
        "code": "jne",
        "name": "Jalur Nugraha Ekakurir (JNE)",
        "costs": [
          {
            "service": "CTC",
            "description": "JNE City Courier",
            "cost": [
              {
                "value": 28000,
                "etd": "5 day",
                "note": ""
              }
            ]
          }
        ]
      }
    ]
  }
}
```

---

## 🛠️ Setup

1. **Clone the repo**

```bash
git clone https://github.com/your-username/go-rajaongkir-location.git
cd go-rajaongkir-location
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Setup environment variables**

Copy the example environment file and adjust values as needed:

```bash
cp .env.example .env
```

Default rate limiting values:

```env
RATE_REQUEST=10
RATE_MINUTE=1
```

**(Optional) Enable API authentication:**

To protect your API with an API key, set:

```env
API_KEY=your-secret-key
```

If `API_KEY` is not set, authentication will be disabled and all requests will be allowed.

**For cost calculations, you'll need a RajaOngkir API key:**

The `/cost` endpoint requires a valid API key from RajaOngkir. This key must be passed in the request headers. See the [Authentication](#-authentication) section for details.

4. **Run the server**

```bash
go run main.go
```

---

## 🗃️ Data Source

The data used in this project was originally sourced from RajaOngkir for **provinces**, **cities**, and **subdistricts**, and is stored locally for performance. The cost calculations are fetched in real-time from RajaOngkir's v2 API.

---

## 🌐 Demo

A live demo is available at:
👉 **[https://go-rajaongkir-location.ardie.web.id](https://go-rajaongkir-location.ardie.web.id)**

Use it to test real responses or as a reference implementation.

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

## 🙏 Acknowledgements

* [RajaOngkir](https://rajaongkir.com) for providing the original location data and API services.
* [Gin](https://github.com/gin-gonic/gin) for powering the HTTP server.