# go-rajaongkir-location

**go-rajaongkir-location** is a lightweight and high-performance API service built in Go, designed to replicate the location-based endpoints of [RajaOngkir](https://rajaongkir.com) — including **provinces**, **cities**, and **subdistricts** — with improved performance and flexibility.

## 🚀 Why This Exists

This project was meant to replace the location service provided by RajaOngkir, which can sometimes be slow or less responsive.  
By self-hosting the static location data (SQLite), this service mimics the same response format while significantly boosting performance and reliability.

It’s ideal for internal use or as a drop-in replacement where RajaOngkir’s full shipping cost API isn’t needed.

---

## 📦 Features

- ⚡ **Fast & Lightweight** – Built with [Gin](https://github.com/gin-gonic/gin) and optimized for blazing fast response times.
- 🗃️ **SQLite-Powered** – Uses SQLite for lightweight and portable local data storage — no external database needed.
- 🔁 **RajaOngkir-Compatible** – Mirrors the structure and response format of RajaOngkir for seamless integration.
- 📌 **Essential Location Data Only** – Focuses purely on **provinces**, **cities**, and **subdistricts** — no shipping costs, just what you need.
- 🔐 **Optional API Key Authentication** – Secure the API by setting an `API_KEY`, or leave it unset for open access.
- 🧩 **Easy to Integrate** – Acts as a drop-in replacement for RajaOngkir location APIs in your frontend or backend.

---

## 📚 API Endpoints

All responses are wrapped in a `rajaongkir` object just like the original API.

### 🔍 Get All Provinces (with optional filter by ID)
```

GET /province?id={province\_id}

```

### 🔍 Get All Cities (with optional province filter)
```

GET /city?province={province\_id}

```

### 🔍 Get All Subdistricts (with optional city filter)
```

GET /subdistrict?city={city\_id}

```

> 🔐 If `API_KEY` is set in your `.env`, include this header in requests:

```

X-API-KEY: your-secret-key

````

---

## 📦 Response Format

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
````

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

4. **Run the server**

```bash
go run main.go
```

---

## 🗃️ Data Source

The data used in this project was originally sourced from RajaOngkir for **provinces**, **cities**, and **subdistricts**, and is stored locally for performance.

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

* [RajaOngkir](https://rajaongkir.com) for providing the original location data.
* [Gin](https://github.com/gin-gonic/gin) for powering the HTTP server.