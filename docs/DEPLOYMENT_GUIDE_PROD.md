# Tutorial Deployment LoopAffi ke Render

Projek Anda sekarang sudah terhubung ke GitHub dan konfigurasi `render.yaml` serta `Dockerfile` telah diperbarui untuk mendukung deployment otomatis. Berikut adalah langkah-langkah untuk mendeploy semuanya secara benar.

## 1. Sinkronisasi Awal (Blueprint)
Karena kita sudah mengupdate `render.yaml`, Render akan mendeteksi perubahan ini jika Anda sudah menghubungkan repo ke Render.

1.  Buka [Dashboard Render](https://dashboard.render.com/).
2.  Cari layanan yang bernama **Blueprint** (biasanya ada di menu 'Blueprints' di navigasi atas).
3.  Pilih Blueprint untuk projek `LOOPAAFFI`.
4.  Klik **Manual Sync** untuk memicu penarikan data terbaru dari GitHub.
5.  Render akan otomatis mulai membuat 3 hal:
    *   `loopaffi-db` (PostgreSQL)
    *   `loopaffi-backend` (Go Web Service)
    *   `loopaffi-frontend` (Next.js Web Service)

## 2. Inisialisasi Database (PENTING)
Database Anda baru saja dibuat, namun tabel-tabelnya mungkin belum memiliki data awal (dummy data) yang diperlukan untuk login.

1.  Di Dashboard Render, buka layanan **loopaffi-db**.
2.  Pilih tab **Info**.
3.  Cari **External Connection String**. Copy URL tersebut.
4.  Buka **pgAdmin 4** di laptop Anda:
    *   Klik kanan pada 'Servers' -> Register -> Server.
    *   Beri nama (misal: 'Render Prod').
    *   Di tab 'Connection', isi bagian 'Connection String' dengan URL yang tadi Anda copy.
5.  Setelah terhubung, klik kanan pada database `loopaffi_db` -> **Query Tool**.
6.  Buka file `full_setup.sql` dari folder projek Anda, copy isinya, paste ke Query Tool pgAdmin, lalu tekan **F5** (Execute).
    *   *Catatan: Backend sebenarnya melakukan AutoMigrate, tapi file SQL ini diperlukan untuk mengisi data Roles dan User Admin.*

## 3. Konfigurasi Environment Variables
Sebagian besar sudah diatur otomatis oleh `render.yaml`, tapi ada satu hal yang perlu Anda perhatikan:

*   **ALLOWED_ORIGIN**: Di `render.yaml`, saya mengaturnya ke `*`. Jika nanti sudah stabil, sebaiknya ubah ini menjadi URL frontend Anda (contoh: `https://loopaffi-frontend.onrender.com`) untuk keamanan.
*   **NEXT_PUBLIC_API_URL**: Jika URL backend Anda di Render ternyata berbeda (misal Render menambahkan angka unik di belakangnya), Anda harus memperbarui variabel ini di dashboard Render layanan frontend agar Next.js bisa memanggil API yang benar.

## 4. Cara Verifikasi
Setelah semua status menjadi **Live** (Hijau):

1.  Buka URL **loopaffi-frontend** yang diberikan oleh Render.
2.  Coba login menggunakan data dummy dari `full_setup.sql`:
    *   **Email**: `admin@loopaffi.com`
    *   **Password**: `password123`
3.  Jika berhasil masuk ke Dashboard, berarti integrasi Fullstack Anda sudah 100% sukses!

> [!TIP]
> Jika deployment gagal lagi, cek **Logs** pada masing-masing layanan di Render. Biasanya logs akan memberitahu detail error jika ada dependency yang kurang atau port yang masih bentrok.

---
**Status Terkini:** Perubahan sudah saya push ke repo GitHub Anda di commit `bd90e81`. Anda bisa langsung cek di Dashboard Render sekarang.
