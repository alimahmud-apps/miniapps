## Cara running test

1. didalam folder `test` running command :
    - `go test --race -v -timeout 30m`

## Jawaban Pertanyaan Tambahan

1. cara memastikan atomic yaitu :
    - implementasi `BEGIN TRANSACTION` dan `commit` ,dan jika mengalami error maka menggunakan `ROLLBACK`

2. Racing condition terjadi karena ada bebera request secara masive yang berjalan bersamaan, dan cara mengatasi nya yaitu dengan
    - implementasi locking dalam case ini didalam `service` menggunakan `mutex` agar mencegah thread lain mengakses dan mengubah data yang sama

3. Langkah-langkah rollback yaitu:
    - menggunakan `BEGIN TRANSACTION` untuk memulai transaksi.
    - jika ada error maka pastikan akan melakukan `ROLLBACK` untuk membatalkan semua perubahan yang sudah dilakukan.
    - jika success makan melakukan `COMMIT` untuk menyimpan semua perubahan secara permanen.
