# API E-Commerce Rakamin-Evermos (Golang)

*Backend* API untuk layanan *e-commerce* yang dibuat dengan bahasa Go (Golang). API ini menerapkan *Clean Architecture* dan memiliki fitur seperti autentikasi (Register, Login), manajemen pengguna, pengelolaan toko, produk, kategori (Admin), alamat, dan sistem transaksi penuh dengan *database transaction*.

## Daftar Isi

  * Tugas Proyek
  * Aturan & Implementasi
  * Struktur Folder
  * Alur Kerja (Request Lifecycle)
  * Cara Menjalankan Proyek

-----

## Fungsionalitas Proyek

API ini dibuat untuk 8 fungsionalitas utama berikut:

1.  Service **Login** dan **Register** Pengguna (*User*).
2.  Pembuatan **Toko** otomatis saat *register* *User* berhasil.
3.  Service untuk mengelola **Akun** (*User* Profile).
4.  Service untuk mengelola **Toko** (Update, Upload Foto).
5.  Service untuk mengelola **Alamat** (CRUD Alamat).
6.  Service untuk mengelola **Kategori** (Hanya Admin).
7.  Service untuk mengelola **Produk** (CRUD, Pagination, Filtering).
8.  Service untuk **Transaksi** (Checkout, Riwayat).

-----

## Aturan & Implementasi

Proyek ini menerapkan 18 aturan yang telah ditentukan. Berikut adalah rincian semua aturan dan di diimplementasikan dalam kode:

 
1.  **Harus memiliki *routing* sesuai koleksi Postman.**

      **Link :** https://github.com/Fajar-Islami/go-example-cruid/blob/master/Rakamin%20Evermos%20Virtual%20Internship.postman_collection.json

2.  **Boleh menambahkan API tapi tidak boleh mengurangi.**


3.  **Email dan No Telepon *user* tidak boleh sama.**

      * **Kode:** `model/user.go`
        ```go
        Email        string `gorm:"size:255;unique"`
        NoTelp       string `gorm:"size:255;unique"`
        ```
      * Validasi tambahan juga ada di `usecase/auth_usecase.go` (fungsi `Register`) untuk memberikan pesan *error* yang jelas.

4.  **Menggunakan JWT.**

      * **Pembuatan Token:** `utils/token.go` (fungsi `GenerateToken`), dipanggil oleh `usecase/auth_usecase.go` saat `Login`.
      * **Validasi Token:** `middleware/auth_middleware.go` (fungsi `AuthMiddleware`), yang membaca `Authorization` header.

5.  **Harus terdapat API yang meng-upload file.**

    *  **Upload Foto Toko:** `handler/toko_handler.go` (fungsi `UploadTokoPhoto`) yang menerima `form-data` dan menyimpan ke folder `/uploads`.
    *  **Upload Foto Produk:** `handler/produk_handler.go` (fungsi `UploadFotoProduk`) yang menerima `form-data`.

6.  **Toko otomatis terbuat ketika user mendaftar.**

      * **Kode:** `usecase/auth_usecase.go` (fungsi `Register`). Setelah `userRepo.Save(user)` berhasil, kode akan langsung memanggil `tokoRepo.Save(newToko)`.

7.  **Alamat diperlukan untuk alamat kirim produk.**

      * **Kode:** `usecase/transaksi_usecase.go` (fungsi `CreateTransaksi`), yang pertama kali memanggil `addressRepo.FindByIDAndUserID(alamatID, userID)` untuk memvalidasi kepemilikan alamat. Alamat divalidasi dan digunakan saat membuat transaksi.

8.  **Yang dapat mengelola kategori hanyalah admin.**

      * **Kode:** `middleware/auth_middleware.go` (fungsi `AdminOnlyMiddleware`) yang memeriksa `c.Get("currentUserIsAdmin")`.
      * **Penerapan:** `router/router.go` di mana semua rute `/categories` ditempatkan di dalam grup `admin.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())`.

9.  **Menerapkan *pagination*.**

      * **Kode:**
          * `utils/pagination.go`: *Helper* untuk mengambil `?page=` dan `?limit=` dari URL dan untuk memformat hasil balasan.
          * `repository/produk_repository.go`: Menggunakan `query.Scopes(utils.Paginate(pagination.Page, pagination.Limit))` untuk memodifikasi *query* GORM. Diterapkan pada `GET /produk` (publik) dan `GET /my-produk` (penjual).

10. **Menerapkan *filtering data*.**

      * **Kode:** `handler/produk_handler.go` (fungsi `parseFilterAndPagination`) membaca *query* `?search=` dan `?category_id=`.
      * `repository/produk_repository.go` (fungsi `buildFilterQuery`) menerapkan filter ini ke *query* GORM menggunakan `query.Where(...)`. Diterapkan pada `GET /produk` dan `GET /my-produk`.

11. **User tidak dapat mendapatkan/meng-update data *user* lain.**

      * **Kode:** `handler/user_handler.go` (fungsi `GetProfile` dan `UpdateProfile`) selalu mengambil `userID` dari `c.Get("currentUserID")`. Diterapkan dengan **tidak** menggunakan ID dari URL (`/users/:id`), tapi menggunakan `currentUserID` yang didapat dari token JWT.

12. **User tidak dapat mengelola alamat data *user* lain.**

      * **Kode:** `usecase/address_usecase.go` selalu memanggil `repository.FindByIDAndUserID(addressID, userID)` untuk memvalidasi kepemilikan sebelum melakukan *update* atau *delete*. Diterapkan dengan menyertakan `userID` yang sedang login dalam *query* GORM.

13. **User tidak dapat mengelola data toko *user* lain.**

      * **Kode:** `usecase/toko_usecase.go` selalu menggunakan `repository.FindByUserID(userID)` untuk menemukan toko berdasarkan `currentUserID` dari token.

14. **User tidak dapat mengelola data produk *user* lain.**

      * **Kode:** `usecase/produk_usecase.go`. Sebelum *update* atau *delete*, *usecase* memanggil `getTokoByUserID(userID)` untuk mendapatkan toko milik *seller*, kemudian memanggil `repository.FindByTokoIDAndProdukID(toko.ID, produkID)` untuk memastikan produk tersebut milik tokonya.

15. **User tidak dapat mengelola data transaksi *user* lain.**

      * **Kode:** `usecase/transaksi_usecase.go` memanggil `repository.FindByUserAndTrxID(userID, trxID)` untuk memastikan `userID` dari token cocok dengan `id_user` di transaksi. Diterapkan saat mengambil riwayat/detail transaksi.

16. **Tabel `log_product` diisi ketika melakukan transaksi.**

      * **Kode:** `usecase/transaksi_usecase.go` (fungsi `CreateTransaksi`). Di dalam *loop* item, *usecase* membuat `logProduk` dari produk asli, lalu memanggil `logProdukRepo.Save(tx, logProduk)`. Bagian inti dari GORM Transaction saat *checkout*.

17. **Tabel `log_produk` digunakan untuk menyimpan data produk di transaksi.**

      * **Kode:** `model/detail_trx.go` memiliki `IDLogProduk` (bukan `IDProduk`), yang menunjuk ke `log_produk` yang dibuat pada aturan \#16. Ini mengunci harga dan nama produk pada saat pembelian.

18. **Menerapkan *Clean Architecture*.**

      * **Implementasi:** Seluruh struktur proyek mengikuti prinsip *Clean Architecture* dengan pemisahan yang jelas antara lapisan.

-----

## Struktur Folder

Proyek ini menggunakan struktur folder berbasis *Clean Architecture* untuk memastikan pemisahan tanggung jawab (*Separation of Concerns*).

```
/rakamin-evermos
  ├── /config       # (config.go) Koneksi database GORM & .env loader.
  ├── /handler      # (toko_handler.go, user_handler.go, dll.) Lapisan HTTP (Pelayan). Menerima request, validasi, dan kirim response.
  ├── /middleware   # (auth_middleware.go) Pos pemeriksaan keamanan (cek JWT, cek Admin).
  ├── /model        # (user.go, produk.go, dll.) Blueprint tabel database (GORM structs).
  ├── /repository   # (user_repository.go, produk_repository.go, dll.) Penjaga Gudang Data. Yang berbicara langsung dengan GORM/DB.
  ├── /router       # (router.go) Peta API. Mendaftarkan semua endpoint dan menghubungkannya ke Handler.
  ├── /usecase      # (auth_usecase.go, produk_usecase.go, dll.) Otak/Koki Utama. Berisi semua logika bisnis.
  ├── /utils        # (hash.go, token.go, response.go, dll.) Alat bantu (helpers) yang digunakan di semua lapisan.
  ├── /uploads      # (Folder .gitignore) Tempat foto toko/produk yang di-upload disimpan.
  ├── .env          # (Rahasia) Menyimpan password DB dan JWT Secret.
  ├── .gitignore    # Mengabaikan file .env, /uploads, dan hasil build.
  ├── go.mod        # Daftar dependensi proyek.
  └── main.go       # Titik masuk dan tempat untuk merakit semua lapisan (Dependency Injection) dan menyalakan server.
```

-----

## Alur Kerja (Request Lifecycle)

Berikut adalah alur kerja sebuah *request* dari *user* hingga kembali lagi ke *user*:

1.  **User (Client)**: Mengirim *request* (misal: `POST /api/v1/my-produk` dengan data JSON dan Token Auth).
2.  **Gin Engine (`main.go`)**: Menerima *request*.
3.  **Router (`router/router.go`)**: Mencocokkan *endpoint* `/my-produk` dan melihat bahwa rute ini ada di grup `authenticated`.
4.  **Middleware (`middleware/auth_middleware.go`)**: "Pos Pemeriksaan" pertama. `AuthMiddleware()` berjalan, memvalidasi Token JWT. Jika valid, *middleware* mengambil `userID` dari token dan menyimpannya di *context* (`c.Set("currentUserID", ...)`) lalu melanjutkan (`c.Next()`).
5.  **Handler (`handler/produk_handler.go`)**: *Handler* `CreateProduk` dipanggil.
      * Ia mengambil `userID` dari *context* (`c.Get("currentUserID")`).
      * Ia memvalidasi JSON *body* (`c.ShouldBindJSON(&input)`).
6.  **Usecase (`usecase/produk_usecase.go`)**: *Handler* memanggil *usecase* `CreateProduk(userID, input)`.
      * "Otak" berjalan: Ia memanggil `tokoRepo.FindByUserID(userID)` untuk mendapatkan `toko.ID` (Aturan \#14).
      * Ia mengatur `input.IDToko = toko.ID`.
7.  **Repository (`repository/produk_repository.go`)**: *Usecase* memanggil `produkRepo.Save(input)`.
      * "Penjaga Gudang" menerjemahkan *struct* Go menjadi *query* GORM.
8.  **GORM & Database**: GORM menjalankan `INSERT INTO produk (...) VALUES (...)`.
9.  **Kembali ke Repository**: Database mengembalikan data produk yang baru dibuat (dengan `ID` baru).
10. **Kembali ke Usecase**: Repository mengembalikan `savedProduk` ke *usecase*.
11. **Kembali ke Handler**: Usecase mengembalikan `savedProduk` ke *handler*.
12. **Handler (`utils/response.go`)**: *Handler* menggunakan `utils.SendCreatedResponse()` untuk memformat balasan JSON yang rapi (dengan `meta` dan `data`).
13. **User (Client)**: Menerima balasan `201 Created` dengan data produk baru.

-----

## Cara Menjalankan Proyek (Tutorial Clone)

Berikut adalah langkah-langkah untuk menjalankan proyek ini di komputer lokal Anda.

### 1\. Prasyarat

  * **Go:** Pastikan Anda telah menginstal [Go (Golang)](https://go.dev/doc/install) versi 1.18 atau lebih baru.
  * **Git:** Terinstal [Git](https://www.google.com/search?q=https://git-scm.com/downloads).
  * **Database:** Server MySQL yang sedang berjalan. Cara termudah adalah menggunakan [Laragon](https://laragon.org/download/) atau XAMPP (yang sudah termasuk phpMyAdmin).

### 2\. Instalasi

1.  **Clone Proyek**
    Buka terminal Anda dan *clone* repositori ini:

    ```bash
    git clone https://github.com/Obiwannn11/golang-mini-project-api.git
    cd golang-mini-project-api
    ```

2.  **Siapkan Database**

      * Nyalakan Laragon/XAMPP (pastikan MySQL berjalan).
      * Buka **phpMyAdmin** (biasanya di `http://localhost/phpmyadmin`).
      * Buat database baru dengan nama persis: **`rakamin_evermos`**.
      * Biarkan database tersebut kosong.

3.  **Buat File `.env`**

      * Di folder *root* proyek, buat file baru bernama `.env`.
      * Salin-tempel konten dari `.env.example` atau gunakan template di bawah ini.
      * **PENTING:** Sesuaikan `DB_PASSWORD` dengan *password* MySQL Anda (di Laragon, *password default* biasanya kosong `""` atau `root`).

    **Template `.env`:**

    ```.env
    # === Konfigurasi Database ===
    DB_USER=root
    DB_PASSWORD=
    DB_NAME=database-backend-golang-evermos
    DB_HOST=127.0.0.1
    DB_PORT=3306

    # === Kunci Rahasia JWT ===
    JWT_SECRET=INI_STRING_ACAK_YANG_SANGAT_PANJANG_DAN_RAHASIA
    ```

4.  **Install Dependensi**
    Di terminal, jalankan `go mod tidy`. Perintah ini akan membaca `go.mod` dan `go.sum` lalu mengunduh semua dependensi yang dibutuhkan secara otomatis (termasuk Gin, GORM, JWT, UUID, dll).

    ```bash
    go mod tidy
    ```

### 3\. Menjalankan Server

1.  Setelah dependensi terinstal, jalankan server:

    ```bash
    go run main.go
    ```

2.  Jika berhasil, Anda akan melihat output di terminal:

    ```
    Koneksi Database Berhasil!
    Menjalankan Migrasi Database...
    Migrasi Database Selesai.
    Server berjalan di http://localhost:Sesuai_PORT_ENV
    ```

      * Pada saat ini, `AutoMigrate` telah membuat semua 9 tabel di database `rakamin_evermos` Anda.

### 4\. Menggunakan API

Buka **Postman** (atau *API Client* lain) dan mulai lakukan *request* ke *endpoint* Anda:

  * **Register:** `POST http://localhost:8080/api/v1/register`
  * **Login:** `POST http://localhost:8080/api/v1/login`
  * ...dan seterusnya.