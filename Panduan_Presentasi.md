# Panduan Presentasi & Demo Postman "Sistem Komisi LoopAffi"

## Persiapan Environment Variables (.env)
Sebelum deploy ke cloud (Vercel & Render/Railway), pastikan variabel lingkungan berikut diset:
### Frontend (Vercel)
- `NEXT_PUBLIC_API_URL` = `https://<backend-app-url>/api`

### Backend (Render/Railway)
- `PORT` = `8080` (Railway/Render akan otomatis mengisinya jika dinamis)
- `DATABASE_URL` = `postgres://user:password@host:port/dbname`
- `JWT_SECRET` = `rahasia_super_aman_anda`

---

## Script Presentasi & Skenario Postman

### 1. Skenario RBAC (Role-Based Access Control)
**Narasi Presentasi:**
> "Pertama, kita akan mendemonstrasikan sistem keamanan RBAC dari LoopAffi. Di sini, sistem membedakan hak akses secara ketat antara Admin dan Affiliate. Ketika seorang Affiliate mencoba mengakses endpoint khusus Admin, seperti menghitung komisi, sistem middleware kita akan otomatis memblokir request tersebut."

**Eksekusi Postman:**
- **Endpoint:** `POST {{BASE_URL}}/api/commissions/calculate`
- **Auth:** Bearer Token (Gunakan token milik Affiliate)
- **Ekspektasi:** Status `403 Forbidden` dengan pesan "Akses diblokir! Hanya admin yang diizinkan mengakses rute ini."

### 2. Use Case #1: Input Data Penjualan
**Narasi Presentasi:**
> "Sekarang kita masuk ke skenario utama. Seorang Admin mencatat data penjualan baru yang berhasil didapatkan oleh seorang Affiliate. Sistem akan menyimpan data ini dan secara otomatis men-trigger kalkulasi awal."

**Eksekusi Postman:**
- **Endpoint:** `POST {{BASE_URL}}/api/sales`
- **Auth:** Bearer Token (Gunakan token Admin)
- **Body (JSON):**
  ```json
  {
      "amount": 500000,
      "affiliateId": "USER-AFFILIATE-1",
      "status": "completed"
  }
  ```
- **Ekspektasi:** Status `200 OK`. Data sale berhasil dibuat.

### 3. Use Case #2: Hitung Komisi
**Narasi Presentasi:**
> "Setelah penjualan tercatat, Admin dapat memverifikasi dan menghitung komisi secara spesifik (Implementasi FR-02). Sistem menarik setting persentase aktif saat ini dan menghitung nominal komisi secara presisi. Algoritma kami juga menggunakan keyword `break` pada iterasi item penjualan untuk memfilter anomali nilai negatif demi efisiensi logika."

**Eksekusi Postman:**
- **Endpoint:** `POST {{BASE_URL}}/api/commissions/calculate`
- **Auth:** Bearer Token (Admin)
- **Body (JSON):**
  ```json
  {
      "id_sale": "SALE-<ID-DARI-UC1>"
  }
  ```
- **Ekspektasi:** Status `200 OK`. Komisi berstatus "Pending".

### 4. Use Case #3: Kelola Pembayaran
**Narasi Presentasi:**
> "Komisi yang berstatus Pending kemudian akan dibayarkan oleh Admin. Begitu proses pembayaran dieksekusi, status komisi akan ter-update menjadi 'Lunas', dan sistem akan men-trigger notifikasi otomatis kepada Affiliate yang bersangkutan bahwa komisi mereka telah cair."

**Eksekusi Postman:**
- **Endpoint:** `POST {{BASE_URL}}/api/payments/process`
- **Auth:** Bearer Token (Admin)
- **Body (JSON):**
  ```json
  {
      "id_commission": "COM-<ID-DARI-UC2>",
      "id_payment_method": "PAY-METHOD-1"
  }
  ```
- **Ekspektasi:** Status `200 OK`. Data pembayaran tersimpan dan status lunas.

### 5. Use Case #4: Buat Laporan
**Narasi Presentasi:**
> "Sebagai penutup, seluruh alur transaksi tadi terekam dengan baik. Admin dapat meng-generate laporan komprehensif yang menampilkan keseluruhan performa penjualan dan komisi."

**Eksekusi Postman:**
- **Endpoint:** `GET {{BASE_URL}}/api/reports`
- **Auth:** Bearer Token (Admin)
- **Ekspektasi:** Status `200 OK`. Mengembalikan rekap data.

---

## Tabel Rekapitulasi Pengujian (Table of Method)
*Silakan salin tabel Markdown di bawah ini ke dalam Excel.*

| Sequence Diagram | Method | Status Progress | Due Date | URL Test | Hasil Uji Coba | Bukti Uji Coba |
| --- | --- | --- | --- | --- | --- | --- |
| Registrasi Affiliate | POST | Done | 05-May-2026 | {{BASE_URL}}/api/users | 200 OK (User baru dibuat, Role dipaksa Affiliate) | UC_Register_Success.png |
| Login Pengguna | POST | Done | 05-May-2026 | {{BASE_URL}}/api/auth/login | 200 OK (Return token JWT) | UC_Login_Success.png |
| Validasi RBAC | POST | Done | 05-May-2026 | {{BASE_URL}}/api/commissions/calculate | 403 Forbidden (Ditolak oleh AdminOnly middleware) | UC_RBAC_Block.png |
| UC #1 (Input Sales) | POST | Done | 05-May-2026 | {{BASE_URL}}/api/sales | 200 OK (Data penjualan tercatat di DB) | UC1_Input_Sales.png |
| UC #2 (Hitung Komisi) | POST | Done | 05-May-2026 | {{BASE_URL}}/api/commissions/calculate | 200 OK (Status = Pending, hitung %) | UC2_Hitung_Komisi.png |
| UC #3 (Kelola Bayar) | POST | Done | 05-May-2026 | {{BASE_URL}}/api/payments/process | 200 OK (Status = Lunas, trigger Notif) | UC3_Kelola_Pembayaran.png |
| UC #4 (Buat Laporan) | GET | Done | 05-May-2026 | {{BASE_URL}}/api/reports | 200 OK (Array data laporan JSON) | UC4_Buat_Laporan.png |
