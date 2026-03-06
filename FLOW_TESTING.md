# Panduan Flow Testing - Go Inventory System

Dokumen ini berisi panduan langkah demi langkah untuk menguji fitur-fitur utama pada sistem Go Inventory, memastikan kesesuaian dengan aturan bisnis dan integritas data stok.

**Prasyarat:**
1. **Jalankan Database & Restore Schema:**
   *   Buka terminal di root folder project (`C:/laragon/www/Go-inventory`) dan jalankan:
       ```bash
       docker-compose up -d
       ```
   *   **Restore Database (via DBeaver):**
       1.  Buka aplikasi **DBeaver** dan koneksikan ke database PostgreSQL (Host: `localhost`, Port: `5432`, Database: `inventory_db`, User: `postgres`, Pass: `postgres`).
       2.  Klik kanan pada koneksi database tersebut, pilih **SQL Editor** -> **New SQL Script**.
       3.  Buka file `backend/database/schema.sql` di text editor Anda, lalu **Copy** semua isinya.
       4.  **Paste** isi file tersebut ke dalam SQL Editor di DBeaver.
       5.  Jalankan script dengan menekan tombol **Execute SQL Script** (atau shortcut `Alt + X`).

2. **Jalankan Backend Service:**
   Buka terminal baru, masuk ke folder backend, dan jalankan aplikasi:
   ```bash
   cd backend
   go run main.go
   ```
   Pastikan tidak ada error dan server berjalan di port 8080.

3. **Siapkan Postman:**
   *   Buka aplikasi Postman.
   *   Import file `backend/postman_collection.json` yang sudah disediakan di dalam project ini.
   *   Pastikan environment variable `base_url` di Postman diset ke `http://localhost:8080`.

---

## A. Stock In (Barang Masuk)
**Tujuan:** Memastikan stok fisik hanya bertambah saat status `DONE` dan log tercatat.

### Skenario 1: Buat Stock In Baru (Status: CREATED)
*   **Action:** `POST /api/stock-in`
*   **Body:**
    ```json
    {
        "product_id": 1,
        "quantity": 100,
        "notes": "Restock Awal"
    }
    ```
*   **Ekspektasi:**
    *   Response Status: `201 Created`.
    *   Response Body: `status` = `"CREATED"`.
    *   **Verifikasi:** Cek Produk (`GET /api/products`). `physical_stock` & `available_stock` **BELUM** bertambah.

### Skenario 2: Update Status ke IN_PROGRESS
*   **Action:** `PUT /api/stock-in/{id}/status` (Ganti `{id}` dengan ID dari langkah 1)
*   **Body:**
    ```json
    {
        "status": "IN_PROGRESS"
    }
    ```
*   **Ekspektasi:**
    *   Response Status: `200 OK`.
    *   Response Body: `status` = `"IN_PROGRESS"`.
    *   **Verifikasi:** Cek Produk. Stok **MASIH BELUM** bertambah.

### Skenario 3: Selesaikan Stock In (Status: DONE)
*   **Action:** `PUT /api/stock-in/{id}/status`
*   **Body:**
    ```json
    {
        "status": "DONE"
    }
    ```
*   **Ekspektasi:**
    *   Response Status: `200 OK`.
    *   Response Body: `status` = `"DONE"`.
    *   **Verifikasi:**
        *   Cek Produk: `physical_stock` & `available_stock` **BERTAMBAH** 100.
        *   Cek Report (`GET /api/reports/transactions?type=stock_in`): Transaksi ini muncul di list.

### Skenario 4: Coba Cancel Stock In yang sudah DONE (Negative Test)
*   **Action:** `PUT /api/stock-in/{id}/status`
*   **Body:**
    ```json
    {
        "status": "CANCELLED"
    }
    ```
*   **Ekspektasi:**
    *   Response Status: `500 Internal Server Error` (Error: "transaction is already final").
    *   Status tidak berubah.

---

## B. Inventory (Cek Stok & Adjustment)
**Tujuan:** Memastikan pemisahan Physical vs Available Stock dan fitur Adjustment.

### Skenario 1: Cek Daftar Produk
*   **Action:** `GET /api/products?name=Laptop`
*   **Ekspektasi:**
    *   Response Status: `200 OK`.
    *   Menampilkan produk dengan `physical_stock` dan `available_stock` yang sesuai.

### Skenario 2: Stock Adjustment (Barang Hilang/Rusak)
*   **Kondisi Awal:** Misal Stok Fisik = 100, Available = 100.
*   **Action:** `PUT /api/products/{id}/adjust`
*   **Request Body:**
    ```json
    {
        "quantity": -5,
        "type": "BOTH",
        "notes": "Barang Rusak ditemukan saat opname"
    }
    ```
*   **Ekspektasi Response:**
    ```json
    {
        "ID": 1,
        "CreatedAt": "...",
        "UpdatedAt": "...",
        "DeletedAt": null,
        "name": "Laptop Gaming",
        "sku": "LPT-001",
        "customer": "Toko Komputer Jaya",
        "physical_stock": 95,
        "available_stock": 95
    }
    ```
*   **Verifikasi:**
    *   Cek Produk: `physical_stock` = 95, `available_stock` = 95.
    *   Cek Report (`GET /api/reports/transactions?type=adjustment`): Muncul log adjustment -5.

### Skenario 3: Stock Adjustment (Koreksi Fisik Saja)
*   **Action:** `PUT /api/products/{id}/adjust`
*   **Request Body:**
    ```json
    {
        "quantity": 2,
        "type": "PHYSICAL",
        "notes": "Salah Hitung Fisik"
    }
    ```
*   **Ekspektasi Response:**
    ```json
    {
        "ID": 1,
        "CreatedAt": "...",
        "UpdatedAt": "...",
        "DeletedAt": null,
        "name": "Laptop Gaming",
        "sku": "LPT-001",
        "customer": "Toko Komputer Jaya",
        "physical_stock": 97,
        "available_stock": 95
    }
    ```
*   **Verifikasi:** Cek Produk: `physical_stock` = 97, `available_stock` = 95 (Hanya fisik berubah).

---

## C. Stock Out (Barang Keluar - Two Phase Commit)
**Tujuan:** Memastikan alokasi stok mengurangi `Available Stock` tapi `Physical Stock` tetap sampai `DONE`.

### Skenario 1: Allocation (Stage 1 - DRAFT/ALLOCATED)
*   **Kondisi Awal:** Fisik = 97, Available = 95.
*   **Action:** `POST /api/stock-out`
*   **Body:**
    ```json
    {
        "product_id": 1,
        "quantity": 10,
        "notes": "Order #123"
    }
    ```
*   **Ekspektasi:**
    *   Response Status: `201 Created`.
    *   Response Body: `status` = `"ALLOCATED"`.
    *   **Verifikasi:**
        *   Cek Produk:
            *   `physical_stock` = 97 (TETAP).
            *   `available_stock` = 85 (BERKURANG 10). **(Reservasi Berhasil)**

### Skenario 2: Execution (Stage 2 - IN_PROGRESS)
*   **Action:** `PUT /api/stock-out/{id}/status`
*   **Body:**
    ```json
    {
        "status": "IN_PROGRESS"
    }
    ```
*   **Ekspektasi:**
    *   Response Status: `200 OK`.
    *   Response Body: `status` = `"IN_PROGRESS"`.
    *   **Verifikasi:** Cek Produk: Stok tidak berubah dari langkah 1 (Fisik 97, Available 85).

### Skenario 3: Rollback / Cancel (Jika Batal)
*   **Action:** Buat Stock Out baru (ID baru), status ALLOCATED -> IN_PROGRESS. Lalu `PUT` status ke `CANCELLED`.
*   **Ekspektasi:**
    *   Response Status: `200 OK`.
    *   **Verifikasi:** Cek Produk: `available_stock` **KEMBALI** bertambah (Rollback).

### Skenario 4: Completion (Stage 2 - DONE)
*   **Action:** `PUT /api/stock-out/{id}/status` (Gunakan ID dari Skenario 1)
*   **Body:**
    ```json
    {
        "status": "DONE"
    }
    ```
*   **Ekspektasi:**
    *   Response Status: `200 OK`.
    *   Response Body: `status` = `"DONE"`.
    *   **Verifikasi:**
        *   Cek Produk:
            *   `physical_stock` = 87 (BERKURANG 10). **(Barang Fisik Keluar)**
            *   `available_stock` = 85 (TETAP, karena sudah dikurangi saat alokasi).
        *   Cek Report (`GET /api/reports/transactions?type=stock_out`): Transaksi muncul.

---

## D. Fitur Report
**Tujuan:** Memastikan laporan hanya menampilkan data yang valid (DONE).

### Skenario 1: Cek Report Stock In
*   **Action:** `GET /api/reports/transactions?type=stock_in`
*   **Ekspektasi:** Hanya menampilkan transaksi Stock In yang statusnya `DONE`. Transaksi yang masih `CREATED` atau `IN_PROGRESS` **TIDAK** boleh muncul.

### Skenario 2: Cek Report Stock Out
*   **Action:** `GET /api/reports/transactions?type=stock_out`
*   **Ekspektasi:** Hanya menampilkan transaksi Stock Out yang statusnya `DONE`. Transaksi `ALLOCATED` atau `IN_PROGRESS` **TIDAK** boleh muncul.

### Skenario 3: Filter Tanggal
*   **Action:** `GET /api/reports/transactions?type=stock_in&start_date=2023-01-01&end_date=2023-12-31`
*   **Ekspektasi:** Menampilkan data sesuai rentang tanggal.

---

## Ringkasan Status Stok

| Aksi | Status Transaksi | Physical Stock | Available Stock | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| **Stock In** | CREATED / IN_PROGRESS | Tetap | Tetap | Barang belum dihitung masuk |
| **Stock In** | DONE | **Bertambah** | **Bertambah** | Barang resmi masuk gudang |
| **Stock Out** | ALLOCATED (Draft) | Tetap | **Berkurang** | Stok di-booking (Reservasi) |
| **Stock Out** | IN_PROGRESS | Tetap | Berkurang | Proses packing |
| **Stock Out** | DONE | **Berkurang** | Berkurang | Barang fisik keluar gudang |
| **Stock Out** | CANCELLED | Tetap | **Bertambah** | Rollback (Batal booking) |

---

## Arsitektur Sistem

Sistem ini dibangun menggunakan **Clean Architecture** dengan prinsip **SOLID** dan **DRY**, memisahkan kode menjadi tiga lapisan utama untuk skalabilitas dan kemudahan pemeliharaan:

1.  **Repository Layer (`repositories/`)**:
    *   Bertanggung jawab langsung untuk akses ke database (PostgreSQL menggunakan GORM).
    *   Hanya berisi query database (CRUD).
    *   Tidak mengandung logika bisnis.

2.  **Service Layer (`services/`)**:
    *   Berisi seluruh logika bisnis (Business Logic) dan aturan validasi.
    *   Menangani kalkulasi stok, validasi status transaksi, dan manajemen transaksi database (ACID).
    *   Menjadi jembatan antara Controller dan Repository.

3.  **Controller Layer (`controllers/`)**:
    *   Bertanggung jawab menangani HTTP Request dan Response (menggunakan Go Fiber).
    *   Memparsing input dari client dan memanggil Service yang sesuai.
    *   Mengembalikan response JSON yang standar.

**Dependency Injection** diterapkan di `main.go` untuk menghubungkan ketiga lapisan ini, membuat sistem mudah diuji (testable) dan modular.
