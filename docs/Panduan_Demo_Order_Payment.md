# 📝 Panduan Presentasi LoopAffi: Demo UI, Database & Use Case Order & Payment

Dokumen ini berisi panduan alur (flow) sekaligus **teks naskah/skrip** yang bisa kamu baca saat melakukan presentasi. Karena aplikasinya sudah di-deploy ke Render & Vercel, pastikan kamu sudah membuka URL aplikasi live kamu di dua tab browser yang berbeda sebelum presentasi dimulai.

## ⚙️ Persiapan Sebelum Presentasi (Behind The Scenes)
1. **Siapkan 2 Tab Browser:**
   - **Tab 1:** Login sebagai **Admin** (`admin@loopaffi.com`). Buka halaman `/admin/sales`.
   - **Tab 2:** Login sebagai **Affiliate** (`rizky@example.com` atau akun affiliate lainnya). Buka halaman `/affiliate/dashboard`.
   - *(Tips: Gunakan mode penyamaran/Incognito untuk tab 2 agar sesi login tidak bertabrakan).*
2. **Siapkan ERD / Skema Database:** Jika kamu punya gambar relasi tabel (ERD), siapkan untuk ditampilkan saat bagian penjelasan database.
3. Pastikan koneksi internet stabil karena akan mengakses aplikasi yang di-deploy di cloud.

---

## 🎬 Skenario & Skrip Presentasi (Step-by-Step)

### 1. Pembukaan & Pengenalan Use Case (1 Menit)
> **Skrip Kamu:**
> "Selamat pagi/siang Bapak/Ibu dosen dan teman-teman. Pada kesempatan ini, saya akan mendemokan sistem **LoopAffi**. Fokus presentasi saya hari ini adalah memperlihatkan hasil UI, struktur Database, serta implementasi logika bisnis untuk **Use Case Order & Payment**."
> 
> "Secara singkat, use case ini mencakup bagaimana sebuah sistem mencatat penjualan (order), menghitung komisi secara otomatis untuk affiliator, hingga proses admin membayarkan (payment) komisi tersebut."

### 2. Penjelasan Database (2 Menit)
*(Tampilkan ERD atau sebutkan nama-nama tabel)*
> **Skrip Kamu:**
> "Sebelum masuk ke demo aplikasi, saya akan menjelaskan sedikit tentang struktur databasenya. Kami menggunakan PostgreSQL yang di-deploy di backend. Terdapat tiga tabel utama yang saling berelasi untuk use case ini:"
> 1. **Tabel `sales`**: Menyimpan data transaksi / order pembelian barang. Kolom pentingnya ada nominal belanja (`amount`) dan ID Affiliator (`affiliate_id`).
> 2. **Tabel `commissions`**: Tabel yang menampung hak atau komisi affiliator. Tabel ini terhubung dengan tabel `sales`.
> 3. **Tabel `payments`**: Menyimpan data tagihan yang harus dibayar oleh sistem ke rekening Affiliator.

### 3. Demo Use Case: Membuat Order/Sale Baru (POV Admin)
*(Buka Tab 1 - Tampilan Admin di halaman Penjualan)*
> **Skrip Kamu:**
> "Sekarang kita masuk ke aplikasinya. Karena LoopAffi adalah sistem komisi dan tidak terhubung langsung dengan toko online nyata, simulasi masuknya **Order Baru** dilakukan melalui dashboard Admin."
> 
> *(Sambil berbicara, klik tombol **Tambah Penjualan**, masukkan nominal misalnya Rp 500.000, lalu pilih nama affiliator, dan klik **Simpan**)*
> 
> "Di sini, admin mencatat ada penjualan sebesar 500 ribu menggunakan link referal milik Affiliator X. Saat saya klik simpan, backend otomatis bekerja. Sistem langsung melakukan kalkulasi komisi sebesar 10% (Rp 50.000) dan menyimpannya ke tabel `commissions` serta menyiapkan antrean di tabel `payments`."

### 4. Demo Use Case: Mengecek Komisi (POV Affiliator)
*(Pindah ke Tab 2 - Tampilan Affiliator)*
> **Skrip Kamu:**
> "Mari kita lihat dari sisi Affiliator-nya. Saat Affiliator login, mereka akan langsung mendapatkan notifikasi otomatis dari sistem."
> 
> *(Klik ikon lonceng notifikasi)*
> "Di sini terlihat ada pesan masuk bahwa mereka mendapatkan komisi dari order yang baru saja terjadi. Kemudian..."
> *(Buka menu Penjualan / Sales)*
> "...pada halaman Penjualan, affiliator bisa melihat riwayat transaksinya dengan transparan."
> 
> *(Buka menu Pembayaran / Payments)*
> "Lalu, komisi tersebut akan masuk ke halaman Pembayaran dengan status masih **Tertunda (Pending)** karena uangnya belum ditransfer oleh perusahaan."

### 5. Demo Use Case: Mencairkan Pembayaran (POV Admin)
*(Pindah kembali ke Tab 1 - Tampilan Admin)*
> **Skrip Kamu:**
> "Untuk menyelesaikan Use Case Payment ini, mari kita kembali ke sisi Admin."
> *(Buka menu Pembayaran / Payments di Admin)*
> 
> "Di halaman ini, admin bisa melihat semua tagihan pencairan dana ke affiliator. Terlihat ada tagihan baru sebesar Rp 50.000 dengan status **Tertunda**."
> 
> *(Klik tombol **Tandai Lunas / Mark as Paid**)*
> "Setelah bagian keuangan perusahaan mentransfer uang ke rekening affiliator, admin akan menekan tombol ini. Sistem akan mengirim API *PUT request* ke backend untuk mengubah status di tabel `payments` menjadi **Lunas (Paid)**."

### 6. Hasil Akhir (POV Affiliator)
*(Pindah lagi ke Tab 2 - Tampilan Affiliator)*
> **Skrip Kamu:**
> "Sebagai penutup, jika kita *refresh* halaman Pembayaran di akun Affiliator..."
> *(Refresh halaman)*
> "...statusnya sekarang sudah berubah menjadi **Lunas / Paid**. Artinya seluruh proses dari Order hingga Payment sudah tereksekusi dengan baik dalam sistem LoopAffi."
> 
> "Sekian demo dari saya. Dari hasil ini dapat dibuktikan bahwa sinkronisasi antara Next.js UI, logic di Go Backend, dan PostgreSQL database berjalan dengan sukses."

---

## 🌟 Tips Tambahan Saat Presentasi:
- **Jangan terburu-buru.** Saat menekan tombol "Simpan" atau "Bayar", diam sekitar 1-2 detik agar penonton bisa melihat UI *loading* dan menyadari bahwa ada proses yang berjalan ke cloud (Render).
- Jika dosen menanyakan: *"Dimana proses kalkulasi komisinya terjadi?"*
  **Jawaban:** *"Proses kalkulasinya ditangani sepenuhnya oleh backend (Go). Jadi ketika frontend menembak API `/sales`, backend yang menghitung 10% lalu menyimpan datanya ke 3 tabel sekaligus (sales, commissions, payments) menggunakan trigger di logic code."*
