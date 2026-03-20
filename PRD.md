# Product Requirements Document (PRD)
## Coffee Shop Point of Sales System

---

## 1. Deskripsi Produk

**Coffee POS** adalah sistem Point of Sales (POS) yang dirancang khusus untuk coffee shop skala menengah. Sistem ini membantu operasional harian coffee shop dengan mengelola transaksi penjualan, manajemen produk, stok bahan baku, dan pelaporan bisnis secara real-time.

**Target Pengguna:**
- Owner coffee shop yang membutuhkan visibilitas penuh terhadap bisnis mereka
- Cashier/barista yang menangani transaksi harian dengan pelanggan

---

## 2. User Roles

### 2.1 Owner
Owner memiliki akses penuh ke seluruh fitur sistem, termasuk:
- Manajemen master data (produk, kategori, meja, user)
- Manajemen stok dan inventory
- Akses ke semua laporan dan dashboard
- Konfigurasi promo dan diskon
- Export data laporan

### 2.2 Cashier
Cashier memiliki akses terbatas untuk operasional harian:
- Manajemen shift kerja (buka/tutup shift)
- Pembuatan transaksi penjualan
- Apply promo saat checkout
- Melihat riwayat transaksi shift mereka sendiri

---

## 3. Fitur Owner

### 3.1 Manajemen Produk

**Deskripsi:** Owner dapat mengelola katalog produk coffee shop termasuk nama, harga, kategori, foto, dan status ketersediaan.

**Fitur:**
- Create: Tambah produk baru dengan nama, deskripsi, harga, kategori, foto, dan SKU
- Read: Lihat daftar semua produk dengan filter dan pencarian
- Update: Edit informasi produk yang sudah ada
- Delete: Hapus produk dengan soft delete (data tetap tersimpan dengan flag `deleted_at`)
- Toggle status aktif/nonaktif produk

**Acceptance Criteria:**
- [ ] Owner dapat menambah produk baru dengan semua field required terisi
- [ ] Owner dapat mengedit produk yang sudah ada
- [ ] Produk yang dihapus tidak tampil di daftar produk aktif tetapi masih tersimpan di database
- [ ] Owner dapat mengaktifkan/menonaktifkan produk tanpa menghapus data
- [ ] Produk dapat diunggah dengan foto (max 2MB, format JPG/PNG)
- [ ] Produk dengan status nonaktif tidak muncul di halaman transaksi cashier

---

### 3.2 Manajemen Kategori Produk

**Deskripsi:** Owner dapat mengelola kategori untuk mengelompokkan produk (contoh: Coffee, Non-Coffee, Pastry, Snack).

**Fitur:**
- Create: Tambah kategori baru
- Read: Lihat daftar semua kategori
- Update: Edit nama kategori
- Delete: Hapus kategori (soft delete)

**Acceptance Criteria:**
- [ ] Owner dapat menambah kategori baru
- [ ] Owner dapat mengedit nama kategori
- [ ] Owner dapat menghapus kategori
- [ ] Kategori yang masih digunakan oleh produk tidak dapat dihapus (validasi)
- [ ] Kategori ditampilkan dengan urutan alfabetis

---

### 3.3 Manajemen Stok

**Deskripsi:** Owner dapat mengelola stok bahan baku dan produk dengan pelacakan riwayat pergerakan stok (masuk, keluar, adjustment).

**Fitur:**
- Lihat stok saat ini untuk semua produk/bahan
- Tambah stok (stock in) dengan sumber (purchase, return, adjustment)
- Kurangi stok (stock out) dengan alasan (waste, expired, sample)
- Lihat riwayat pergerakan stok dengan filter tanggal
- Stock alert untuk produk dengan stok di bawah minimum

**Acceptance Criteria:**
- [ ] Owner dapat melihat stok real-time semua produk
- [ ] Owner dapat menambah stok dengan input quantity dan sumber
- [ ] Owner dapat mengurangi stok dengan input quantity dan alasan
- [ ] Setiap pergerakan stok tercatat dengan timestamp, user, dan keterangan
- [ ] Sistem menampilkan alert visual untuk produk dengan stok di bawah threshold
- [ ] Riwayat stok dapat difilter berdasarkan tanggal dan jenis pergerakan

---

### 3.4 Manajemen Meja

**Deskripsi:** Owner dapat mengelola data meja untuk dine-in customers termasuk status dan kapasitas.

**Fitur:**
- Create: Tambah meja baru dengan nomor dan kapasitas
- Read: Lihat daftar semua meja dengan status (available/occupied)
- Update: Edit informasi meja
- Delete: Hapus meja (soft delete)

**Acceptance Criteria:**
- [ ] Owner dapat menambah meja baru dengan nomor dan kapasitas
- [ ] Owner dapat mengedit informasi meja
- [ ] Owner dapat menghapus meja
- [ ] Status meja (available/occupied) ditampilkan dengan indikator visual
- [ ] Meja yang sedang digunakan tidak dapat dihapus

---

### 3.5 Manajemen User Cashier

**Deskripsi:** Owner dapat mengelola akun cashier yang dapat mengakses sistem untuk transaksi.

**Fitur:**
- Create: Tambah user cashier baru dengan username, password, dan nama lengkap
- Read: Lihat daftar semua cashier
- Update: Edit informasi cashier
- Delete: Hapus/Hard disable user cashier
- Reset password cashier

**Acceptance Criteria:**
- [ ] Owner dapat menambah user cashier baru
- [ ] Owner dapat mengedit informasi cashier (nama, status aktif)
- [ ] Owner dapat menonaktifkan akun cashier
- [ ] Owner dapat mereset password cashier
- [ ] Password disimpan dengan enkripsi bcrypt
- [ ] Cashier yang sedang dalam shift aktif tidak dapat dinonaktifkan

---

### 3.6 Dashboard Laporan

**Deskripsi:** Owner dapat melihat dashboard dengan berbagai laporan bisnis untuk monitoring performa.

**Fitur:**
- Revenue harian, mingguan, bulanan dengan grafik tren
- Produk terlaris (best seller) berdasarkan quantity dan revenue
- Transaksi per cashier dengan total revenue
- Filter laporan berdasarkan periode tanggal

**Acceptance Criteria:**
- [ ] Dashboard menampilkan total revenue hari ini dengan perbandingan hari sebelumnya
- [ ] Dashboard menampilkan grafik revenue 7 hari terakhir
- [ ] Dashboard menampilkan top 5 produk terlaris berdasarkan quantity
- [ ] Dashboard menampilkan revenue per cashier untuk periode terpilih
- [ ] Owner dapat filter laporan berdasarkan range tanggal
- [ ] Data laporan di-refresh secara real-time

---

### 3.7 Manajemen Promo dan Diskon

**Deskripsi:** Owner dapat membuat dan mengelola promo/diskon yang dapat diterapkan saat transaksi.

**Fitur:**
- Create: Buat promo baru dengan nama, tipe (persentase/nominal), nilai, periode aktif
- Read: Lihat daftar semua promo dengan status (aktif/kadaluarsa)
- Update: Edit promo yang belum aktif
- Delete: Hapus promo (soft delete)
- Set periode aktif promo (tanggal mulai dan selesai)

**Acceptance Criteria:**
- [ ] Owner dapat membuat promo dengan tipe persentase (1-100%)
- [ ] Owner dapat membuat promo dengan tipe nominal (Rupiah)
- [ ] Owner dapat设定 periode aktif promo
- [ ] Promo yang sudah kadaluarsa otomatis tidak tampil di halaman cashier
- [ ] Owner dapat melihat daftar promo aktif dan tidak aktif
- [ ] Promo dapat dihapus dengan soft delete

---

### 3.8 Export Laporan ke CSV

**Deskripsi:** Owner dapat export laporan transaksi dan stok ke format CSV untuk analisis lebih lanjut.

**Fitur:**
- Export laporan transaksi ke CSV
- Export laporan stok ke CSV
- Export laporan revenue ke CSV
- Filter tanggal sebelum export

**Acceptance Criteria:**
- [ ] Owner dapat export laporan transaksi dengan filter periode
- [ ] Owner dapat export laporan stok saat ini
- [ ] Owner dapat export laporan revenue
- [ ] File CSV ter-download dengan nama file yang deskriptif (contoh: `laporan_transaksi_2026-03-20.csv`)
- [ ] Format CSV dapat dibuka di Excel/Google Sheets tanpa error

---

## 4. Fitur Cashier

### 4.1 Manajemen Shift

**Deskripsi:** Cashier harus membuka shift di awal kerja dan menutup shift di akhir kerja dengan rekap kas.

**Fitur:**
- Buka Shift: Input modal kas awal untuk memulai shift
- Tutup Shift: Input kas akhir, sistem hitung selisih dengan transaksi
- Lihat status shift saat ini

**Acceptance Criteria:**
- [ ] Cashier tidak dapat membuat transaksi sebelum membuka shift
- [ ] Cashier input modal kas awal saat buka shift
- [ ] Cashier dapat melihat status shift (aktif/selesai)
- [ ] Saat tutup shift, sistem menampilkan rekap: modal awal, total transaksi, kas yang seharusnya ada
- [ ] Cashier input kas akhir aktual, sistem hitung selisih (surplus/shortage)
- [ ] Hanya satu shift aktif per cashier dalam satu waktu

---

### 4.2 Buat Transaksi Baru

**Deskripsi:** Cashier dapat membuat transaksi penjualan baru untuk pelanggan.

**Fitur:**
- Pilih meja untuk dine-in atau takeaway
- Pilih produk dari katalog dengan pencarian
- Atur quantity per produk
- Lihat total harga real-time
- Tambah catatan per produk (contoh: less sugar, extra ice)

**Acceptance Criteria:**
- [ ] Cashier dapat memilih meja atau opsi takeaway
- [ ] Cashier dapat menambah multiple produk ke dalam satu transaksi
- [ ] Cashier dapat mengatur quantity per produk (min 1)
- [ ] Total harga dihitung otomatis dan ditampilkan real-time
- [ ] Cashier dapat menambah catatan opsional per produk
- [ ] Produk yang ditampilkan hanya produk dengan status aktif

---

### 4.3 Apply Promo Aktif

**Deskripsi:** Cashier dapat menerapkan promo yang sedang aktif ke transaksi.

**Fitur:**
- Lihat daftar promo yang sedang aktif
- Pilih satu promo untuk diterapkan ke transaksi
- Lihat discount amount sebelum checkout
- Validasi eligibility promo (jika ada)

**Acceptance Criteria:**
- [ ] Cashier dapat melihat daftar promo yang sedang aktif
- [ ] Cashier dapat memilih satu promo per transaksi
- [ ] Discount ditampilkan secara terpisah di total transaksi
- [ ] Promo hanya dapat diterapkan jika masih dalam periode aktif
- [ ] Cashier dapat menghapus promo dari transaksi sebelum checkout

---

### 4.4 Checkout dengan Midtrans Snap

**Deskripsi:** Transaksi dibayar menggunakan payment gateway Midtrans Snap.

**Fitur:**
- Integrasi dengan Midtrans Snap API
- Generate Snap token untuk transaksi
- Redirect customer ke Midtrans payment page
- Handle callback dari Midtrans (success, pending, failed)
- Update status transaksi berdasarkan payment status

**Acceptance Criteria:**
- [ ] Sistem generate Snap token dari Midtrans untuk setiap transaksi
- [ ] Customer di-redirect ke halaman pembayaran Midtrans
- [ ] Transaksi status menjadi `pending` setelah redirect
- [ ] Sistem handle Midtrans callback dan update status ke `paid` jika sukses
- [ ] Transaksi dengan status `paid` tercatat di laporan
- [ ] Struk/nota dapat dicetak setelah pembayaran sukses

---

### 4.5 Lihat Riwayat Transaksi Shift

**Deskripsi:** Cashier dapat melihat riwayat transaksi yang mereka buat selama shift saat ini.

**Fitur:**
- Lihat daftar transaksi shift hari ini
- Filter berdasarkan status pembayaran
- Lihat detail transaksi

**Acceptance Criteria:**
- [ ] Cashier dapat melihat daftar transaksi yang mereka buat di shift saat ini
- [ ] Transaksi ditampilkan dengan informasi: waktu, total, status pembayaran
- [ ] Cashier dapat melihat detail transaksi (produk, quantity, promo)
- [ ] Riwayat transaksi hanya menampilkan transaksi dari shift aktif

---

## 5. Business Rules

1. **Shift Management:**
   - Cashier wajib membuka shift sebelum dapat membuat transaksi
   - Satu cashier hanya boleh memiliki satu shift aktif pada satu waktu
   - Shift harus ditutup sebelum cashier logout

2. **Transaksi:**
   - Transaksi hanya dapat dibuat oleh cashier dengan shift aktif
   - Transaksi yang sudah dibayar tidak dapat dihapus, hanya dapat di-refund dengan approval Owner
   - Minimal satu produk harus ada dalam transaksi
   - Order yang sudah checkout tidak bisa diubah itemnya

3. **Stok:**
   - Produk dengan stok 0 tidak dapat ditambahkan ke transaksi
   - Stok hanya berkurang setelah webhook Midtrans confirmed diterima
   - Stock adjustment harus disertai keterangan/alasan

4. **Promo:**
   - Satu transaksi hanya bisa menggunakan satu promo
   - Promo tidak dapat digabung dengan promo lain
   - Promo yang sudah kadaluarsa otomatis tidak tersedia

5. **Payment:**
   - Semua pembayaran diproses melalui Midtrans Snap
   - Transaksi dengan status pending lebih dari 24 jam otomatis dibatalkan
   - Refund hanya dapat diproses oleh Owner

6. **Data Integrity:**
   - Soft delete digunakan untuk semua master data (produk, kategori, meja, user, promo)
   - Semua transaksi dan pergerakan stok tidak dapat dihapus (hard delete)
   - Password user di-hash menggunakan bcrypt dengan cost factor 10

---

## 6. Out of Scope

Fitur berikut **tidak** termasuk dalam versi ini:

1. **Multi-outlet support:** Sistem hanya mendukung satu lokasi coffee shop
2. **Inventory purchase order:** PO ke supplier tidak dikelola dalam sistem
3. **Recipe management:** Resipi dan bill of materials untuk produk tidak dikelola
4. **Customer loyalty program:** Member dan loyalty points tidak termasuk
5. **Mobile app:** Sistem hanya tersedia sebagai web application
6. **Kitchen display system:** Order tidak ditampilkan di kitchen display
7. **Reservation management:** Reservasi meja tidak dikelola
8. **Employee scheduling:** Jadwal shift cashier tidak dikelola dalam sistem
9. **Accounting integration:** Integrasi dengan software akuntansi tidak termasuk
10. **Multi-currency:** Sistem hanya mendukung Rupiah (IDR)

---

## Version History

| Version | Date | Author | Description |
|---------|------|--------|-------------|
| 1.0 | 20 Maret 2026 | Product Team | Initial PRD |

---

## Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Product Owner | | | |
| Tech Lead | | | |
| Project Manager | | | |
