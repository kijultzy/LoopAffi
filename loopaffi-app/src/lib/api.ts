const API_BASE = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";

function getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem("token");
    return {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
    };
}

export function setToken(token: string) {
    localStorage.setItem("token", token);
}

export function removeToken() {
    localStorage.removeItem("token");
}

export interface DBUser {
    id: string;
    name: string;
    email: string;
    role_id: string;
    phone?: string;
    status?: string;
}

export interface DBSale {
    id: string;
    date: string;
    amount: number;
    affiliate_id: string;
    status: string;
}

export interface DBCommission {
    id: string;
    sale_id: string;
    affiliate_id: string;
    amount: number;
    date: string;
}

export interface DBPayment {
    id: string;
    affiliate_id: string;
    amount: number;
    date: string;
    status: string;
}

export interface DBReportRow {
    affiliate_id: string;
    affiliate_name: string;
    affiliate_email: string;
    total_sales: number;
    sales_count: number;
    total_commission: number;
    paid_commission: number;
    pending_commission: number;
}

// ============ Admin API calls ============

export async function fetchAdminUsers(): Promise<DBUser[]> {
    const res = await fetch(`${API_BASE}/admin/users`, { headers: getAuthHeaders() });
    const json = await res.json();
    if (json.status === "success") return json.data || [];
    throw new Error(json.message || "Gagal mengambil data user");
}

export async function fetchAdminSales(): Promise<DBSale[]> {
    const res = await fetch(`${API_BASE}/admin/sales`, { headers: getAuthHeaders() });
    const json = await res.json();
    if (json.status === "success") return json.data || [];
    throw new Error(json.message || "Gagal mengambil data penjualan");
}

export async function fetchAdminCommissions(): Promise<DBCommission[]> {
    const res = await fetch(`${API_BASE}/admin/commissions`, { headers: getAuthHeaders() });
    const json = await res.json();
    if (json.status === "success") return json.data || [];
    throw new Error(json.message || "Gagal mengambil data komisi");
}

export async function fetchAdminPayments(): Promise<DBPayment[]> {
    const res = await fetch(`${API_BASE}/admin/payments`, { headers: getAuthHeaders() });
    const json = await res.json();
    if (json.status === "success") return json.data || [];
    throw new Error(json.message || "Gagal mengambil data pembayaran");
}

export async function fetchAdminReport(): Promise<DBReportRow[]> {
    const res = await fetch(`${API_BASE}/admin/reports`, { headers: getAuthHeaders() });
    const json = await res.json();
    if (json.status === "success") return json.data || [];
    throw new Error(json.message || "Gagal mengambil data laporan");
}

export async function createSale(amount: number, affiliateId: string): Promise<DBSale> {
    const res = await fetch(`${API_BASE}/admin/sales`, {
        method: "POST",
        headers: getAuthHeaders(),
        body: JSON.stringify({ amount, affiliate_id: affiliateId }),
    });
    const json = await res.json();
    if (json.status === "success") return json.data;
    throw new Error(json.message || "Gagal mencatat penjualan");
}

export async function markPaymentPaid(paymentId: string): Promise<void> {
    const res = await fetch(`${API_BASE}/admin/payments/${paymentId}/pay`, {
        method: "PUT",
        headers: getAuthHeaders(),
    });
    const json = await res.json();
    if (json.status !== "success") {
        throw new Error(json.message || "Gagal memperbarui status pembayaran");
    }
}
