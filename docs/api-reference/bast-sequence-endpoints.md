# BAST Sequence Endpoints

Endpoint untuk memeriksa dan me-reset nomor urut (running number) BAST per format & periode. Sumber kode: [`bast_sequence_handler.go`](../../internal/handlers/bast_sequence_handler.go), [`bast_sequence_service.go`](../../internal/services/bast_sequence_service.go), [`bast_sequence_repository.go`](../../internal/repositories/bast_sequence_repository.go).

> 🔑 **Terlindungi** — semua endpoint butuh token.

> ⚠️ Endpoint ini untuk **monitoring & koreksi manual**. Nomor urut di-generate otomatis saat `POST /bast-requests` — Anda tidak perlu memanggil ini untuk operasi normal.

---

## Skema: BastSequence

```json
{
  "sequence_id": "seq-123-uuid",
  "format_id": "fmt-...",
  "format": { ... },              // nested (jika di-preload)
  "year": 2026,
  "month": 6,
  "last_number": 15,
  "created_at": "...",
  "updated_at": "..."
}
```

| Field | Tipe | Keterangan |
|---|---|---|
| `sequence_id` | UUID | Primary key |
| `format_id` | UUID | FK ke format |
| `year` | int | Tahun periode |
| `month` | int | Bulan periode (1-12) |
| `last_number` | int | Nomor urut terakhir yang dikeluarkan |

**Unique index komposit:** `(format_id, year, month)` — tiap format per bulan hanya punya satu baris.

---

## Get Sequence
<a id="get"></a>

Cek nomor urut terakhir untuk format + periode tertentu.

### Request
```http
GET /api/bast-sequences?format_id=<uuid>&year=2026&month=6
Authorization: Bearer <token>
```

**Query param:**
| Param | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `format_id` | string | ✅ | UUID format |
| `year` | int | ✅ | Tahun, mis. `2026` |
| `month` | int | ✅ | Bulan, `1`-`12` |

### Response — `200 OK`
```json
{
  "sequence_id": "seq-...",
  "format_id": "fmt-...",
  "year": 2026,
  "month": 6,
  "last_number": 15,
  "created_at": "...",
  "updated_at": "..."
}
```

### `404 Not Found`
Jika belum ada sequence untuk periode itu (belum ada BAST dibuat bulan itu).
```json
{ "error": "Sequence not found" }
```

### Contoh curl
```bash
curl -H "Authorization: Bearer <TOKEN>" \
  "http://localhost:8080/api/bast-sequences?format_id=fmt-...&year=2026&month=6"
```

---

## Reset Sequence
<a id="reset"></a>

Set / timpa `last_number` untuk format + periode tertentu. Berguna untuk koreksi manual, migrasi, atau testing.

### Request
```http
POST /api/bast-sequences/reset
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "format_id": "fmt-...",
  "year": 2026,
  "month": 6,
  "last_number": 0
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `format_id` | string | ✅ | UUID format |
| `year` | int | ✅ | Tahun |
| `month` | int | ✅ | Bulan |
| `last_number` | int | ✅ | Nomor baru (akan jadi `last_number`) |

> Perilaku: jika baris sequence belum ada → dibuat baru. Jika sudah ada → `last_number` di-timpa.

### Response — `200 OK`
```json
{
  "sequence_id": "seq-...",
  "format_id": "fmt-...",
  "year": 2026,
  "month": 6,
  "last_number": 0,
  ...
}
```

### `400 Bad Request`
```json
{ "error": "Key: 'ResetSequenceRequest.LastNumber' ..." }
```

### Contoh curl
```bash
curl -X POST http://localhost:8080/api/bast-sequences/reset \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"format_id":"fmt-...","year":2026,"month":6,"last_number":0}'
```

---

## ⚠️ Peringatan Reset

| Skenario | Konsekuensi |
|---|---|
| Set `last_number` **lebih rendah** dari nomor BAST yang sudah terpakai | Bisa menyebabkan **nomor ganda** saat BAST berikutnya dibuat |
| Set ke `0` | BAST berikutnya akan dapat nomor `1` lagi (jika belum dipakai) |
| Set ke angka besar (mis. `999`) | Berguna untuk "lompat" nomor |

**Selalu cek** nomor BAST yang sudah ada sebelum reset. Hanya admin yang seharusnya melakukan ini.

---

## Hubungan dengan BAST Request

Flow normal (tanpa perlu panggil sequence manual):
1. User `POST /bast-requests` dengan `tipe_nomor: Internal`.
2. Service otomatis:
   - Cek sequence (format, year, month).
   - Jika ada → increment +1.
   - Jika belum ada → buat baru mulai dari 1.
3. Sequence terupdate otomatis.

Endpoint di file ini hanya untuk **lihat** atau **koreksi** hasilnya.

---

← **[Daftar Referensi API](README.md)** · 🤿 **[Deep Dive Penomoran](../guides/bast-numbering-deep-dive.md)**
