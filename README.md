# 🤖 WhatsAppBot

**WhatsAppBot** adalah aplikasi yang memungkinkan Anda menghubungkan akun WhatsApp melalui pemindaian kode QR dan mengirim pesan langsung dari antarmuka web. Aplikasi ini mendukung pengiriman pesan tunggal dan pesan massal menggunakan unggahan file CSV.

---

## 🚀 Fitur

- 🔗 **Koneksi mudah** melalui pemindaian Kode QR.
- 🌐 **Antarmuka web** yang ramah pengguna.
- 💬 **Mengirim pesan teks** ke kontak individual.
- 📤 **Mengirim pesan massal** dengan mengunggah file CSV berisi nomor telepon.
- 🔄 **Pembaruan status koneksi** secara real-time (Terhubung, Terputus, Butuh Scan QR).
- 🪵 **Pencatatan log kesalahan** untuk pengiriman pesan massal yang gagal.

---

## 🛠️ Cara Menjalankan Aplikasi

Aplikasi ini dapat dijalankan di berbagai sistem operasi menggunakan file eksekusi yang sesuai. Setelah dijalankan, aplikasi akan secara otomatis memulai server web pada port `9000` dan membuka antarmuka web di browser default Anda.

### 🪟 Windows

1. Jalankan file `bot-whatsapp-windows.exe` dengan **double-click** atau **Run as Administrator**.
2. Aplikasi akan membuka terminal dan menjalankan server secara otomatis.
3. Antarmuka web akan terbuka di browser default Anda.

### 🍎 macOS (Apple Silicon)

```bash
# Buka Terminal
cd /path/to/file
chmod +x bot-whatsapp-mac-arm64
./bot-whatsapp-mac-arm64
```

### 🍏 macOS (Intel)

```bash
# Buka Terminal
cd /path/to/file
chmod +x bot-whatsapp-mac-intel
./bot-whatsapp-mac-intel
```

### 🐧 Linux (Ubuntu)

```bash
# Buka Terminal
cd /path/to/file
chmod +x bot-whatsapp-ubuntu
./bot-whatsapp-ubuntu
```

---

## 🌐 Cara Menggunakan Aplikasi Web

Setelah aplikasi berjalan, antarmuka web akan terbuka di `http://localhost:9000`.

### 1️⃣ Menghubungkan Akun WhatsApp Anda

1. Pada halaman web, status awal akan terlihat seperti "Memeriksa Status..." atau "Tidak Terhubung".
2. Klik tombol **Connect**.
3. Status akan berubah menjadi **Butuh Scan QR Code**.
4. Sebuah **Kode QR** akan muncul di halaman.
5. Buka WhatsApp di ponsel Anda, masuk ke menu **Perangkat Tertaut (Linked Devices)**, dan pindai Kode QR yang muncul.
6. Setelah berhasil dipindai, status akan berubah menjadi **Terhubung dan Login**. Tombol **Connect** akan dinonaktifkan, dan tombol **Disconnect** akan diaktifkan.

### 2️⃣ Mengirim Pesan

Setelah berhasil terhubung, bagian pengiriman pesan akan muncul. Anda dapat memilih untuk mengirim **pesan tunggal** atau **pesan massal**.

#### 📩 Kirim Pesan Tunggal

1. Klik tombol **Kirim Pesan Tunggal**.
2. Masukkan **Nomor Telepon** penerima (gunakan format internasional, misalnya `628...`).
3. Tulis pesan Anda di kolom **Tulis Pesan**.
4. Klik **Kirim Pesan**.
5. Anda akan mendapatkan notifikasi yang mengonfirmasi apakah pesan berhasil terkirim.

#### 📊 Kirim Pesan Massal (via CSV)

1. Klik tombol **Kirim Pesan Massal**.
2. Unggah file CSV yang berisi daftar nomor telepon (pastikan hanya ada satu kolom nomor telepon).
3. Setelah diunggah, pratinjau nomor telepon yang akan dikirimkan akan muncul.
4. Tulis pesan Anda di kolom **Tulis Pesan Massal**.
5. Tombol kirim hanya akan aktif setelah Anda menulis pesan.
6. Klik tombol **Kirim Pesan Massal**.
7. Aplikasi akan mengirimkan pesan ke semua nomor telepon. Setelah selesai, sebuah notifikasi akan menunjukkan ringkasan jumlah pesan yang berhasil dan gagal.
8. Jika ada kegagalan, tombol **Download Log Error** akan muncul, yang memungkinkan Anda mengunduh file teks berisi detail nomor yang gagal dikirimi pesan.

### 3️⃣ Memutuskan Koneksi

- **Disconnect**: Klik untuk memutuskan sesi dari server. Status akan berubah menjadi **Tidak Terhubung**. Anda dapat menghubungkan kembali dengan menekan tombol **Connect**.
- **Logout**: Untuk keluar sepenuhnya, putuskan koneksi menggunakan tombol **Disconnect**, lalu gunakan ponsel Anda untuk keluar dari sesi perangkat tertaut melalui menu WhatsApp.

---

## 📁 Contoh File CSV

File CSV yang digunakan untuk pengiriman pesan massal harus berisi daftar nomor telepon dalam satu kolom, contohnya:

```
6281234567890
6289876543210
6281122334455
```

---

## 📞 Kontak

irvanirawanit@gmail.com
