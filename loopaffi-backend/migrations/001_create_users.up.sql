-- 1. Hapus tabel lama jika ada
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS commissions CASCADE;
DROP TABLE IF EXISTS sales CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS affiliates CASCADE;

-- 2. Buat tabel Roles
CREATE TABLE roles (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- 3. Buat tabel Users (Menggabungkan info user dan affiliate)
CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id VARCHAR(50) NOT NULL REFERENCES roles(id),
    phone VARCHAR(20),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Buat tabel Sales
CREATE TABLE sales (
    id VARCHAR(50) PRIMARY KEY,
    date DATE NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    affiliate_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'completed'
);

-- 5. Buat tabel Commissions
CREATE TABLE commissions (
    id VARCHAR(50) PRIMARY KEY,
    sale_id VARCHAR(50) NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    affiliate_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL,
    date DATE NOT NULL
);

-- 6. Buat tabel Payments
CREATE TABLE payments (
    id VARCHAR(50) PRIMARY KEY,
    affiliate_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL,
    date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'paid'
);

-- 6.5. Buat tabel Notifications
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 7. ISI DATA AWAL
INSERT INTO roles (id, name) VALUES
('admin', 'Admin'),
('affiliate', 'Affiliate');

INSERT INTO users (id, name, email, password_hash, role_id, phone, status) VALUES
('USR-ADMIN-01', 'Super Admin', 'admin@loopaffi.com', '$2a$10$xOgNPFIr0e0wWA56HfQ2k.KrLrMTCuNE0/mJDBg46ewaFSN3IzfJ.', 'admin', '081100000000', 'active'),
('USR-AFF-01', 'Rizky Dzulfikar Ahmad', 'rizky@example.com', '$2a$10$vx6htugaV4KRG2ucXc8iHOo/Ch4FRfM7aa6Tpc79j9ecPo9U6APsu', 'affiliate', '081234567890', 'active');

INSERT INTO sales (id, date, amount, affiliate_id, status) VALUES
('SALE-001', '2026-05-10', 500000, 'USR-AFF-01', 'completed');

INSERT INTO commissions (id, sale_id, affiliate_id, amount, date) VALUES
('COM-001', 'SALE-001', 'USR-AFF-01', 50000, '2026-05-10');

INSERT INTO payments (id, affiliate_id, amount, date, status) VALUES
('PAY-001', 'USR-AFF-01', 50000, '2026-05-12', 'paid');
