
# 🚚 Delivery App — Golang Backend

A full-featured delivery platform built with **Golang (Gin)**, supporting real-time order tracking, chat, and admin management.

## 🧩 Features

### 👤 Customer
- Register & login with **JWT authentication**
- **Email OTP** verification for account activation & password recovery
- Browse products and **place orders**
- View order status (**real-time tracking**)
- **Chat with shipper** (via WebSocket)
- **Rate & review** products (with image upload)

### 🚴 Shipper
- View available orders & **accept deliveries**
- View order info on **interactive map**
- **Share live location** via WebSocket
- **Chat** with customers
- **Mark orders as completed**

### 🛠️ Admin
- **Manage products**: add, edit, delete
- **Confirm orders**
- **Ban / unban** users (customer or shipper)
- Dashboard overview of platform activity

---

## 🏗️ Tech Stack

| Layer | Technology |
|-------|-------------|
| **Backend** | Golang (Gin Framework) |
| **Database** | MySQL |
| **Authentication** | JWT + Email OTP |
| **Real-time** | WebSocket (chat + location updates) |
| **File storage** | Cloudinary |
| **Email service** | SMTP (for OTP verification) |
| **Deployment** | (optional) Docker / Render / Railway / etc. |

---

## ⚙️ Project Structure

