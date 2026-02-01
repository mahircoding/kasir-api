-- Migration: Create tables for Kasir API
-- Run this SQL in your Supabase SQL Editor

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW())
);

-- Insert sample categories
INSERT INTO categories (name, description) VALUES
    ('Buah', 'Kategori untuk berbagai jenis buah-buahan segar'),
    ('Sayuran', 'Kategori untuk berbagai jenis sayuran segar'),
    ('Minuman', 'Kategori untuk berbagai jenis minuman')
ON CONFLICT DO NOTHING;

-- Insert sample products
INSERT INTO products (name, price, stock, category_id) VALUES
    ('Apple', 90000, 100, 1),
    ('Banana', 30000, 150, 1),
    ('Orange', 70000, 200, 1)
ON CONFLICT DO NOTHING;

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
