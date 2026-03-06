# AI Usage Report

## 1. AI Tools yang Digunakan
*   **Google Gemini** (Model: gemini-3.0-pro / gemini-1.5-pro via IDE Integration)
    *   Digunakan untuk generate boilerplate code, refactoring arsitektur, pembuatan dokumentasi (Markdown), dan pembuatan Postman Collection.

## 2. Prompt Paling Kompleks
Prompt berikut digunakan untuk melakukan **Refactoring Code** secara menyeluruh agar kode lebih terstruktur dan modular:

> "tolong refactoring semua nya ke struktur yang lebih rapi tanpa mengubah flow bisnis, dan jika ada query ke database maka letakkan ke folder repository"

Prompt ini kompleks karena AI harus:
1.  Memahami struktur kode yang ada (Controllers, Models, Routes).
2.  Memecah kode yang sebelumnya menyatu menjadi 3 layer terpisah: `Repository` (Database), `Service` (Logika Bisnis), dan `Controller` (HTTP Handler).
3.  Menerapkan Dependency Injection di `main.go` untuk menghubungkan antar layer.
4.  Memastikan fitur transaksional database tetap berjalan dengan baik di struktur yang baru.

## 3. Modifikasi Manual demi Kepatuhan Aturan Bisnis
Meskipun AI membantu membuat struktur kode, beberapa logika bisnis spesifik memerlukan penyesuaian manual agar sesuai dengan kebutuhan sistem:

1.  **Logika Stock Adjustment (Pemisahan Stok):**
    AI awalnya menyarankan update stok sederhana. Kode dimodifikasi manual untuk mematuhi aturan: *"Harus memisahkan antara Physical Stock dan Available Stock"*. Kami mengubah input menjadi `Type` (PHYSICAL/AVAILABLE/BOTH) agar user bisa mengoreksi stok fisik (Stock Opname) tanpa mengganggu alokasi penjualan, atau sebaliknya.

2.  **Proses Rollback pada Stock Out:**
    Pada fitur Stock Out yang menggunakan skema *Two-Phase Commit* (Allocation -> Execution), logika rollback dipertegas secara manual. Ketika transaksi dibatalkan (`CANCELLED`) setelah status `IN_PROGRESS` atau `ALLOCATED`, sistem harus secara eksplisit mengembalikan `AvailableStock` yang sebelumnya sudah direservasi. AI menyediakan wrapper transaksi database, namun logika bisnis spesifik mengenai kapan stok harus dikembalikan (rollback) ditambahkan manual untuk menjamin integritas data stok.
