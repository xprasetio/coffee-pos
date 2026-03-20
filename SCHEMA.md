# Database Schema Design
## Coffee Shop Point of Sales System

---

## 1. Daftar Tabel

| No | Nama Tabel | Deskripsi |
|----|------------|-----------|
| 1 | `users` | Menyimpan data user (Owner dan Cashier) |
| 2 | `categories` | Kategori produk untuk mengelompokkan produk |
| 3 | `products` | Katalog produk coffee shop |
| 4 | `product_images` | Foto produk (satu produk bisa punya banyak foto) |
| 5 | `tables` | Data meja untuk dine-in customers |
| 6 | `shifts` | Shift kerja cashier (buka/tutup shift) |
| 7 | `stock_movements` | Riwayat pergerakan stok (masuk/keluar/adjustment) |
| 8 | `promos` | Promo dan diskon yang tersedia |
| 9 | `transactions` | Transaksi penjualan |
| 10 | `transaction_items` | Detail item dalam satu transaksi |
| 11 | `payments` | Informasi pembayaran transaksi |
| 12 | `refunds` | Pengembalian dana/refund transaksi |

---

## 2. Detail Tabel

### 2.1 `users`

Menyimpan data user yang dapat mengakses sistem (Owner dan Cashier).

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `username` | VARCHAR(50) | UNIQUE, NOT NULL | Username untuk login |
| `password_hash` | VARCHAR(255) | NOT NULL | Hash password (bcrypt) |
| `full_name` | VARCHAR(100) | NOT NULL | Nama lengkap |
| `role` | ENUM('owner', 'cashier') | NOT NULL | Role user |
| `is_active` | TINYINT(1) | DEFAULT 1 | Status aktif (1=aktif, 0=nonaktif) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| `deleted_at` | TIMESTAMP | NULL | Soft delete timestamp |

**Index:**
- `idx_username` pada `username`
- `idx_role` pada `role`
- `idx_is_active` pada `is_active`

---

### 2.2 `categories`

Kategori untuk mengelompokkan produk (Coffee, Non-Coffee, Pastry, dll).

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `name` | VARCHAR(50) | NOT NULL | Nama kategori |
| `description` | TEXT | NULL | Deskripsi kategori |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| `deleted_at` | TIMESTAMP | NULL | Soft delete timestamp |

**Index:**
- `idx_name` pada `name`

---

### 2.3 `products`

Katalog produk yang dijual di coffee shop.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `sku` | VARCHAR(20) | UNIQUE, NOT NULL | Stock Keeping Unit |
| `name` | VARCHAR(100) | NOT NULL | Nama produk |
| `description` | TEXT | NULL | Deskripsi produk |
| `category_id` | VARCHAR(36) | FOREIGN KEY → categories(id) | Kategori produk |
| `price` | BIGINT | NOT NULL | Harga jual (dalam sen/rupiah) |
| `stock` | INT UNSIGNED | DEFAULT 0 | Stok saat ini |
| `min_stock` | INT UNSIGNED | DEFAULT 5 | Threshold alert stok minimum |
| `is_active` | TINYINT(1) | DEFAULT 1 | Status ketersediaan (1=aktif, 0=nonaktif) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| `deleted_at` | TIMESTAMP | NULL | Soft delete timestamp |

**Index:**
- `idx_sku` pada `sku`
- `idx_category_id` pada `category_id`
- `idx_is_active` pada `is_active`
- `idx_name` pada `name` (untuk pencarian)

---

### 2.4 `product_images`

Foto produk (satu produk dapat memiliki multiple foto).

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `product_id` | VARCHAR(36) | FOREIGN KEY → products(id) | Produk yang dimiliki |
| `file_path` | VARCHAR(255) | NOT NULL | Path file foto |
| `file_name` | VARCHAR(100) | NOT NULL | Nama file asli |
| `file_size` | BIGINT | NOT NULL | Ukuran file dalam bytes |
| `mime_type` | VARCHAR(50) | NOT NULL | MIME type (image/jpeg, image/png) |
| `is_primary` | TINYINT(1) | DEFAULT 0 | Foto utama (1=ya, 0=tidak) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |

**Index:**
- `idx_product_id` pada `product_id`
- `idx_is_primary` pada `is_primary`

---

### 2.5 `tables`

Data meja untuk dine-in customers.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `table_number` | VARCHAR(10) | UNIQUE, NOT NULL | Nomor/meja identifier |
| `capacity` | TINYINT UNSIGNED | NOT NULL | Kapasitas (jumlah orang) |
| `location` | VARCHAR(50) | NULL | Lokasi meja (indoor/outdoor/vip) |
| `status` | ENUM('available', 'occupied', 'reserved') | DEFAULT 'available' | Status meja saat ini |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| `deleted_at` | TIMESTAMP | NULL | Soft delete timestamp |

**Index:**
- `idx_table_number` pada `table_number`
- `idx_status` pada `status`

---

### 2.6 `shifts`

Shift kerja cashier dengan modal kas dan rekap tutup shift.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `user_id` | VARCHAR(36) | FOREIGN KEY → users(id) | Cashier yang membuka shift |
| `opening_balance` | BIGINT | NOT NULL | Modal kas awal (dalam sen/rupiah) |
| `closing_balance` | BIGINT | NULL | Kas akhir aktual (diisi saat tutup shift) |
| `expected_balance` | BIGINT | NULL | Kas yang seharusnya ada (dari sistem) |
| `balance_difference` | BIGINT | NULL | Selisih kas (positif=surplus, negatif=shortage) |
| `opened_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu buka shift |
| `closed_at` | TIMESTAMP | NULL | Waktu tutup shift |
| `closing_notes` | TEXT | NULL | Catatan saat tutup shift |
| `status` | ENUM('open', 'closed') | DEFAULT 'open' | Status shift |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |

**Index:**
- `idx_user_id` pada `user_id`
- `idx_status` pada `status`
- `idx_opened_at` pada `opened_at`

---

### 2.7 `stock_movements`

Riwayat pergerakan stok untuk audit trail.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `product_id` | VARCHAR(36) | FOREIGN KEY → products(id) | Produk yang bergerak |
| `user_id` | VARCHAR(36) | FOREIGN KEY → users(id) | User yang melakukan movement |
| `movement_type` | ENUM('in', 'out', 'adjustment') | NOT NULL | Tipe pergerakan |
| `source` | ENUM('purchase', 'return', 'adjustment', 'waste', 'expired', 'sample', 'sale') | NOT NULL | Sumber/alasan movement |
| `quantity` | INT | NOT NULL | Quantity (positif untuk in, negatif untuk out) |
| `previous_stock` | INT UNSIGNED | NOT NULL | Stok sebelum movement |
| `current_stock` | INT UNSIGNED | NOT NULL | Stok setelah movement |
| `notes` | TEXT | NULL | Catatan/keterangan |
| `reference_id` | VARCHAR(36) | NULL | ID referensi (contoh: transaction_id untuk sale) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |

**Index:**
- `idx_product_id` pada `product_id`
- `idx_user_id` pada `user_id`
- `idx_movement_type` pada `movement_type`
- `idx_source` pada `source`
- `idx_reference_id` pada `reference_id`
- `idx_created_at` pada `created_at` (untuk filter tanggal)

---

### 2.8 `promos`

Promo dan diskon yang dapat diterapkan ke transaksi.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `name` | VARCHAR(100) | NOT NULL | Nama promo |
| `description` | TEXT | NULL | Deskripsi promo |
| `type` | ENUM('percentage', 'nominal') | NOT NULL | Tipe diskon |
| `value` | BIGINT | NOT NULL | Nilai diskon (persentase dalam sen [1-10000 untuk 1%-100%], atau nominal dalam sen) |
| `min_transaction` | BIGINT | DEFAULT 0 | Minimum transaksi untuk menggunakan promo (dalam sen) |
| `max_discount` | BIGINT | NULL | Maksimum nominal diskon (untuk type=percentage) |
| `started_at` | TIMESTAMP | NOT NULL | Mulai periode aktif |
| `ended_at` | TIMESTAMP | NOT NULL | Akhir periode aktif |
| `is_active` | TINYINT(1) | DEFAULT 1 | Status aktif (1=aktif, 0=nonaktif) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| `deleted_at` | TIMESTAMP | NULL | Soft delete timestamp |

**Index:**
- `idx_type` pada `type`
- `idx_is_active` pada `is_active`
- `idx_period` pada `started_at`, `ended_at`

---

### 2.9 `transactions`

Transaksi penjualan.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `transaction_number` | VARCHAR(20) | UNIQUE, NOT NULL | Nomor transaksi (format: TRX-YYYYMMDD-XXXX) |
| `user_id` | VARCHAR(36) | FOREIGN KEY → users(id) | Cashier yang membuat transaksi |
| `shift_id` | VARCHAR(36) | FOREIGN KEY → shifts(id) | Shift saat transaksi dibuat |
| `table_id` | VARCHAR(36) | FOREIGN KEY → tables(id) | Meja (NULL untuk takeaway) |
| `is_takeaway` | TINYINT(1) | DEFAULT 0 | Flag takeaway (1=ya, 0=tidak) |
| `subtotal` | BIGINT | NOT NULL | Subtotal sebelum diskon |
| `discount_amount` | BIGINT | DEFAULT 0 | Nominal diskon |
| `discount_type` | ENUM('none', 'percentage', 'nominal') | DEFAULT 'none' | Tipe diskon yang diterapkan |
| `promo_id` | VARCHAR(36) | FOREIGN KEY → promos(id) | Promo yang digunakan (NULL jika tidak ada) |
| `total` | BIGINT | NOT NULL | Total setelah diskon |
| `status` | ENUM('pending', 'paid', 'cancelled', 'refunded') | DEFAULT 'pending' | Status transaksi |
| `customer_notes` | TEXT | NULL | Catatan dari customer |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |

**Index:**
- `idx_transaction_number` pada `transaction_number`
- `idx_user_id` pada `user_id`
- `idx_shift_id` pada `shift_id`
- `idx_table_id` pada `table_id`
- `idx_status` pada `status`
- `idx_created_at` pada `created_at` (untuk laporan)

---

### 2.10 `transaction_items`

Detail item dalam satu transaksi.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `transaction_id` | VARCHAR(36) | FOREIGN KEY → transactions(id) | Transaksi induk |
| `product_id` | VARCHAR(36) | FOREIGN KEY → products(id) | Produk yang dibeli |
| `product_name` | VARCHAR(100) | NOT NULL | Snapshot nama produk saat transaksi |
| `product_price` | BIGINT | NOT NULL | Snapshot harga produk saat transaksi |
| `quantity` | INT UNSIGNED | NOT NULL | Quantity item |
| `subtotal` | BIGINT | NOT NULL | Subtotal per item (price × quantity) |
| `notes` | VARCHAR(255) | NULL | Catatan per item (less sugar, extra ice, dll) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |

**Index:**
- `idx_transaction_id` pada `transaction_id`
- `idx_product_id` pada `product_id`

---

### 2.11 `payments`

Informasi pembayaran transaksi melalui Midtrans.

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `transaction_id` | VARCHAR(36) | FOREIGN KEY → transactions(id), UNIQUE | Transaksi yang dibayar |
| `midtrans_order_id` | VARCHAR(50) | UNIQUE, NOT NULL | Order ID di Midtrans |
| `midtrans_transaction_id` | VARCHAR(50) | NULL | Transaction ID dari Midtrans |
| `snap_token` | VARCHAR(255) | NULL | Snap token dari Midtrans |
| `payment_method` | ENUM('credit_card', 'debit_card', 'gopay', 'shopeepay', 'bank_transfer', 'qris') | NULL | Metode pembayaran yang dipilih customer |
| `payment_amount` | BIGINT | NOT NULL | Jumlah pembayaran (dalam sen/rupiah) |
| `payment_status` | ENUM('pending', 'success', 'failed', 'expired', 'cancelled') | DEFAULT 'pending' | Status pembayaran |
| `paid_at` | TIMESTAMP | NULL | Waktu pembayaran sukses |
| `midtrans_response` | JSON | NULL | Response lengkap dari Midtrans (untuk audit) |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |

**Index:**
- `idx_transaction_id` pada `transaction_id`
- `idx_midtrans_order_id` pada `midtrans_order_id`
- `idx_payment_status` pada `payment_status`
- `idx_paid_at` pada `paid_at`

---

### 2.12 `refunds`

Pengembalian dana/refund transaksi (hanya Owner yang dapat approve).

| Kolom | Tipe Data | Constraint | Keterangan |
|-------|-----------|------------|------------|
| `id` | VARCHAR(36) | PRIMARY KEY | UUID |
| `transaction_id` | VARCHAR(36) | FOREIGN KEY → transactions(id), UNIQUE | Transaksi yang di-refund |
| `user_id` | VARCHAR(36) | FOREIGN KEY → users(id) | Owner yang approve refund |
| `refund_amount` | BIGINT | NOT NULL | Jumlah refund (dalam sen/rupiah) |
| `reason` | TEXT | NOT NULL | Alasan refund |
| `status` | ENUM('pending', 'approved', 'rejected', 'processed', 'completed') | DEFAULT 'pending' | Status refund |
| `approved_at` | TIMESTAMP | NULL | Waktu approval |
| `processed_at` | TIMESTAMP | NULL | Waktu refund diproses |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| `updated_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |

**Index:**
- `idx_transaction_id` pada `transaction_id`
- `idx_user_id` pada `user_id`
- `idx_status` pada `status`

---

## 3. Relasi Antar Tabel

```
┌─────────────┐
│   users     │
└──────┬──────┘
       │
       ├─────────────────────────────────────┐
       │                                     │
       ▼                                     ▼
┌──────────────┐                    ┌─────────────┐
│  categories  │                    │   shifts    │
└──────┬───────┘                    └──────┬──────┘
       │                                   │
       │                                   │
       ▼                                   ▼
┌──────────────┐                    ┌─────────────┐
│   products   │                    │transactions │
└──────┬───────┘                    └──────┬──────┘
       │                                   │
       ├──────────────┐                    ├──────────────┐
       │              │                    │              │
       ▼              ▼                    ▼              ▼
┌──────────────┐ ┌──────────────┐  ┌─────────────┐ ┌─────────────┐
│product_images│ │stock_movement│  │trans_items  │ │  payments   │
└──────────────┘ └──────────────┘  └──────┬──────┘ └─────────────┘
                                          │
                                          ▼
                                   ┌─────────────┐
                                   │   refunds   │
                                   └─────────────┘
```

### Foreign Key Constraints

| Tabel | Kolom FK | Referensi | On Delete | On Update |
|-------|----------|-----------|-----------|-----------|
| `products` | `category_id` | `categories(id)` | SET NULL | CASCADE |
| `product_images` | `product_id` | `products(id)` | CASCADE | CASCADE |
| `stock_movements` | `product_id` | `products(id)` | RESTRICT | CASCADE |
| `stock_movements` | `user_id` | `users(id)` | SET NULL | CASCADE |
| `shifts` | `user_id` | `users(id)` | RESTRICT | CASCADE |
| `transactions` | `user_id` | `users(id)` | RESTRICT | CASCADE |
| `transactions` | `shift_id` | `shifts(id)` | RESTRICT | CASCADE |
| `transactions` | `table_id` | `tables(id)` | SET NULL | CASCADE |
| `transactions` | `promo_id` | `promos(id)` | SET NULL | CASCADE |
| `transaction_items` | `transaction_id` | `transactions(id)` | CASCADE | CASCADE |
| `transaction_items` | `product_id` | `products(id)` | RESTRICT | CASCADE |
| `payments` | `transaction_id` | `transactions(id)` | CASCADE | CASCADE |
| `refunds` | `transaction_id` | `transactions(id)` | RESTRICT | CASCADE |
| `refunds` | `user_id` | `users(id)` | SET NULL | CASCADE |

### Daftar Lengkap Relasi

| Relasi | Keterangan |
|--------|------------|
| `[products.category_id]` → `[categories.id]` | Produk milik kategori tertentu |
| `[product_images.product_id]` → `[products.id]` | Foto produk |
| `[stock_movements.product_id]` → `[products.id]` | Pergerakan stok untuk produk |
| `[stock_movements.user_id]` → `[users.id]` | User yang melakukan stock movement |
| `[shifts.user_id]` → `[users.id]` | Cashier yang membuka shift |
| `[transactions.user_id]` → `[users.id]` | Cashier yang membuat transaksi |
| `[transactions.shift_id]` → `[shifts.id]` | Shift saat transaksi dibuat |
| `[transactions.table_id]` → `[tables.id]` | Meja untuk transaksi dine-in |
| `[transactions.promo_id]` → `[promos.id]` | Promo yang diterapkan ke transaksi |
| `[transaction_items.transaction_id]` → `[transactions.id]` | Item detail dari transaksi |
| `[transaction_items.product_id]` → `[products.id]` | Produk yang dibeli dalam transaksi |
| `[payments.transaction_id]` → `[transactions.id]` | Pembayaran untuk transaksi |
| `[refunds.transaction_id]` → `[transactions.id]` | Refund untuk transaksi |
| `[refunds.user_id]` → `[users.id]` | Owner yang approve refund |

---

## 4. Keputusan Desain Penting

### 4.1 UUID sebagai Primary Key (VARCHAR 36)

**Keputusan:** Menggunakan UUID versi 4 untuk semua primary key.

**Alasan:**
- **Distributed-friendly:** Dapat generate ID di application layer tanpa perlu query database untuk AUTO_INCREMENT
- **Security:** ID tidak sequential, sulit ditebak oleh pihak luar
- **Scalability:** Memudahkan sharding di masa depan jika diperlukan
- **Data migration:** Lebih mudah saat merge data dari multiple databases

**Trade-off:** UUID lebih besar (16 bytes) dan lebih lambat untuk indexing dibanding INTEGER, namun untuk skala coffee shop menengah dampaknya negligible.

---

### 4.2 BIGINT untuk Uang (dalam Sen/Rupiah Terkecil)

**Keputusan:** Semua kolom yang berhubungan dengan uang menggunakan `BIGINT` dengan satuan sen (1 rupiah = 100 sen).

**Alasan:**
- **Presisi:** Menghindari floating-point rounding errors yang umum terjadi dengan DECIMAL/FLOAT
- **Konsistensi:** Perhitungan matematika lebih sederhana dan akurat
- **Midtrans compatibility:** Midtrans API juga menggunakan format integer dalam satuan terkecil

**Contoh:**
- Rp 25.000 = 2.500.000 (dalam sen)
- Diskon 15% = disimpan sebagai 1500 (dalam sen persentase)

**Trade-off:** Application layer harus melakukan konversi dari/to rupiah untuk display ke user.

---

### 4.3 ENUM untuk Status

**Keputusan:** Menggunakan ENUM MySQL untuk kolom dengan nilai terbatas dan fixed.

**Kolom yang menggunakan ENUM:**
- `users.role`: 'owner', 'cashier'
- `tables.status`: 'available', 'occupied', 'reserved'
- `shifts.status`: 'open', 'closed'
- `stock_movements.movement_type`: 'in', 'out', 'adjustment'
- `stock_movements.source`: 'purchase', 'return', 'adjustment', 'waste', 'expired', 'sample', 'sale'
- `promos.type`: 'percentage', 'nominal'
- `transactions.status`: 'pending', 'paid', 'cancelled', 'refunded'
- `transactions.discount_type`: 'none', 'percentage', 'nominal'
- `payments.payment_method`: 'credit_card', 'debit_card', 'gopay', 'shopeepay', 'bank_transfer', 'qris'
- `payments.payment_status`: 'pending', 'success', 'failed', 'expired', 'cancelled'
- `refunds.status`: 'pending', 'approved', 'rejected', 'processed', 'completed'

**Alasan:**
- **Data integrity:** Database enforce valid values
- **Performance:** ENUM lebih efisien storage (1-2 bytes) dibanding VARCHAR
- **Self-documenting:** Nilai yang valid terlihat di schema

**Trade-off:** Menambah nilai ENUM baru memerlukan ALTER TABLE.

---

### 4.4 Soft Delete dengan `deleted_at`

**Keputusan:** Tabel master data menggunakan soft delete dengan kolom `deleted_at TIMESTAMP NULL`.

**Tabel dengan soft delete:**
- `users` (kecuali hard delete untuk user yang tidak pernah transaksi)
- `categories`
- `products`
- `tables`
- `promos`

**Alasan:**
- **Audit trail:** Data historis tetap tersimpan
- **Recovery:** Dapat restore data yang terhapus tidak sengaja
- **Referential integrity:** Tidak break foreign key constraints

**Implementasi:** Semua query SELECT harus include `WHERE deleted_at IS NULL`.

---

### 4.5 Snapshot Data Produk di `transaction_items`

**Keputusan:** `transaction_items` menyimpan snapshot `product_name` dan `product_price` saat transaksi terjadi.

**Alasan:**
- **Price history:** Harga produk dapat berubah di masa depan, tapi transaksi historis harus tetap akurat
- **Data integrity:** Jika produk dihapus atau diubah namanya, transaksi lama tidak terpengaruh
- **Reporting:** Laporan revenue tetap akurat meskipun data produk berubah

**Trade-off:** Redundansi data, namun diperlukan untuk accuracy historis.

---

### 4.6 `created_at` dan `updated_at` di Semua Tabel

**Keputusan:** Setiap tabel memiliki kolom `created_at` dan `updated_at`.

**Alasan:**
- **Audit trail:** Mengetahui kapan data dibuat dan terakhir diubah
- **Debugging:** Membantu troubleshoot issue dengan timeline
- **Reporting:** `created_at` digunakan untuk filter periode laporan

**Implementasi:** Menggunakan `DEFAULT CURRENT_TIMESTAMP` dan `ON UPDATE CURRENT_TIMESTAMP`.

---

### 4.7 JSON untuk Midtrans Response

**Keputusan:** Kolom `midtrans_response` menggunakan tipe JSON.

**Alasan:**
- **Flexibility:** Response Midtrans dapat berubah tanpa perlu ALTER TABLE
- **Debugging:** Full response tersimpan untuk troubleshoot payment issues
- **Webhook verification:** Data lengkap tersedia untuk verifikasi webhook signature

**Trade-off:** JSON tidak dapat di-index secara efisien, namun hanya digunakan untuk audit.

---

### 4.8 Stok Tidak Berkurang Sebelum Payment Confirmed

**Keputusan:** Stok hanya berkurang setelah webhook Midtrans dengan status `payment_status = success` diterima.

**Implementasi:**
- Saat checkout: stok di-reserve (bisa implementasi `reserved_stock` kolom di `products`)
- Webhook success: kurangi `stock` dan catat di `stock_movements` dengan `source = 'sale'`
- Webhook failed/expired: release reservation

**Alasan:**
- **Race condition prevention:** Mencegah stok negatif jika ada multiple checkout bersamaan
- **Payment failure handling:** Jika payment gagal, stok kembali available
- **Accuracy:** Stok di database mencerminkan stok fisik yang benar-benar terjual

---

## 5. SQL Create Statements

Untuk implementasi lengkap, lihat file `database/migrations/001_initial_schema.sql`.

---

## Version History

| Version | Date | Author | Description |
|---------|------|--------|-------------|
| 1.0 | 20 Maret 2026 | Database Team | Initial schema design |
