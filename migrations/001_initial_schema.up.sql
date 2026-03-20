-- +migrate Up
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role ENUM('owner', 'cashier') NOT NULL,
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    sku VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    category_id VARCHAR(36) NULL,
    price BIGINT NOT NULL,
    stock INT UNSIGNED DEFAULT 0,
    min_stock INT UNSIGNED DEFAULT 5,
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL ON UPDATE CASCADE,
    INDEX idx_sku (sku),
    INDEX idx_category_id (category_id),
    INDEX idx_is_active (is_active),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(50) NOT NULL,
    is_primary TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_product_id (product_id),
    INDEX idx_is_primary (is_primary)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create tables table
CREATE TABLE IF NOT EXISTS `tables` (
    id VARCHAR(36) PRIMARY KEY,
    table_number VARCHAR(10) UNIQUE NOT NULL,
    capacity TINYINT UNSIGNED NOT NULL,
    location VARCHAR(50) NULL,
    status ENUM('available', 'occupied', 'reserved') DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_table_number (table_number),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create shifts table
CREATE TABLE IF NOT EXISTS shifts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    opening_balance BIGINT NOT NULL,
    closing_balance BIGINT NULL,
    expected_balance BIGINT NULL,
    balance_difference BIGINT NULL,
    opened_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP NULL,
    closing_notes TEXT NULL,
    status ENUM('open', 'closed') DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_opened_at (opened_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create stock_movements table
CREATE TABLE IF NOT EXISTS stock_movements (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NULL,
    movement_type ENUM('in', 'out', 'adjustment') NOT NULL,
    source ENUM('purchase', 'return', 'adjustment', 'waste', 'expired', 'sample', 'sale') NOT NULL,
    quantity INT NOT NULL,
    previous_stock INT UNSIGNED NOT NULL,
    current_stock INT UNSIGNED NOT NULL,
    notes TEXT NULL,
    reference_id VARCHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    INDEX idx_product_id (product_id),
    INDEX idx_user_id (user_id),
    INDEX idx_movement_type (movement_type),
    INDEX idx_source (source),
    INDEX idx_reference_id (reference_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create promos table
CREATE TABLE IF NOT EXISTS promos (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    type ENUM('percentage', 'nominal') NOT NULL,
    value BIGINT NOT NULL,
    min_transaction BIGINT DEFAULT 0,
    max_discount BIGINT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP NOT NULL,
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_type (type),
    INDEX idx_is_active (is_active),
    INDEX idx_period (started_at, ended_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY,
    transaction_number VARCHAR(20) UNIQUE NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    shift_id VARCHAR(36) NOT NULL,
    table_id VARCHAR(36) NULL,
    is_takeaway TINYINT(1) DEFAULT 0,
    subtotal BIGINT NOT NULL,
    discount_amount BIGINT DEFAULT 0,
    discount_type ENUM('none', 'percentage', 'nominal') DEFAULT 'none',
    promo_id VARCHAR(36) NULL,
    total BIGINT NOT NULL,
    status ENUM('pending', 'paid', 'cancelled', 'refunded') DEFAULT 'pending',
    customer_notes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (shift_id) REFERENCES shifts(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (table_id) REFERENCES `tables`(id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (promo_id) REFERENCES promos(id) ON DELETE SET NULL ON UPDATE CASCADE,
    INDEX idx_transaction_number (transaction_number),
    INDEX idx_user_id (user_id),
    INDEX idx_shift_id (shift_id),
    INDEX idx_table_id (table_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create transaction_items table
CREATE TABLE IF NOT EXISTS transaction_items (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    product_price BIGINT NOT NULL,
    quantity INT UNSIGNED NOT NULL,
    subtotal BIGINT NOT NULL,
    notes VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL UNIQUE,
    midtrans_order_id VARCHAR(50) UNIQUE NOT NULL,
    midtrans_transaction_id VARCHAR(50) NULL,
    snap_token VARCHAR(255) NULL,
    payment_method ENUM('credit_card', 'debit_card', 'gopay', 'shopeepay', 'bank_transfer', 'qris') NULL,
    payment_amount BIGINT NOT NULL,
    payment_status ENUM('pending', 'success', 'failed', 'expired', 'cancelled') DEFAULT 'pending',
    paid_at TIMESTAMP NULL,
    midtrans_response JSON NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_midtrans_order_id (midtrans_order_id),
    INDEX idx_payment_status (payment_status),
    INDEX idx_paid_at (paid_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create refunds table
CREATE TABLE IF NOT EXISTS refunds (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL UNIQUE,
    user_id VARCHAR(36) NULL,
    refund_amount BIGINT NOT NULL,
    reason TEXT NOT NULL,
    status ENUM('pending', 'approved', 'rejected', 'processed', 'completed') DEFAULT 'pending',
    approved_at TIMESTAMP NULL,
    processed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
