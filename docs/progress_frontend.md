# Progress Tracker — Frontend (React + TypeScript SPA)

Target arsitektur: **Separation of Concerns (SoC)** untuk antarmuka Asisten Laboratorium Hardware / IoT.

## Checklist Implementasi

- [x] **Setup**: Struktur folder `frontend/src/`, `package.json`, `package-lock.json`
- [ ] **Scaffold**: `vite.config.ts`, `tsconfig.json`, `index.html`, `src/main.tsx`, `.env.example`
- [ ] **Types** (`src/types/component.ts`): TypeScript interfaces (`Component`, `Category`, `Condition`, `CreateComponentPayload`, `UpdateComponentPayload`)
- [ ] **API Client** (`src/api/client.ts`): Fetch client untuk `POST /items`, `GET /items`, `GET /items/{id}`, `PUT /items/{id}`, `DELETE /items/{id}`
- [ ] **Utils** (`src/utils/inventory.ts`): Helper `isLowStock(stock_qty <= min_stock)`, formatter ISO 8601, badge mapping
- [ ] **Hooks** (`src/hooks/useInventory.ts`): State management, filtering, loading/error handling, CRUD actions
- [ ] **Styling** (`src/index.css`): Desain dashboard laboratorium yang bersih & responsif
- [ ] **Components** (`src/components/`):
  - [ ] `StatsBar.tsx` (KPI ringkasan stok & peringatan)
  - [ ] `InventoryTable.tsx` (Tabel komponen + indikator stok menipis)
  - [ ] `ComponentFormDialog.tsx` (Modal tambah & edit komponen)
  - [ ] `DeleteConfirmDialog.tsx` (Modal konfirmasi hapus)
- [ ] **Composition** (`src/App.tsx`): Integrasi halaman utama
- [ ] **Verifikasi**: `npm run build` (`tsc -b && vite build`) lulus tanpa error
