import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export type Role = 'admin' | 'affiliate';

export interface User {
    id: string;
    name: string;
    email: string;
    role: Role;
}

export interface Sale {
    id: string;
    date: string;
    amount: number;
    affiliateId: string;
    status: 'completed' | 'refunded';
}

export interface Commission {
    id: string;
    saleId: string;
    affiliateId: string;
    amount: number;
    date: string;
}

export interface Payment {
    id: string;
    affiliateId: string;
    amount: number;
    date: string;
    status: 'pending' | 'paid';
}

export interface Notification {
    id: number; // Ubah ke number sesuai backend
    userId: string;
    message: string;
    is_read: boolean; // Sesuaikan dengan backend is_read
    created_at: string; // Sesuaikan dengan backend created_at
}

interface AppState {
    currentUser: User | null;
    token: string | null;
    users: User[];
    sales: Sale[];
    commissions: Commission[];
    payments: Payment[];
    notifications: Notification[];
    globalCommissionRate: number;
    darkMode: boolean;
    login: (user: User, token: string) => void;
    logout: () => void;
    setNotifications: (notifications: Notification[]) => void;
    addSale: (sale: Omit<Sale, 'id'>) => Promise<void>;
    markPaymentPaid: (paymentId: string) => void;
    markNotificationRead: (id: number) => void;
    convertAllToIDR: () => void;
    toggleDarkMode: () => void;
}

export const mockUsers: User[] = [
    { id: '1', name: 'Admin User', email: 'admin@loopaffi.com', role: 'admin' },
    { id: '2', name: 'John Affiliate', email: 'john@example.com', role: 'affiliate' },
    { id: '3', name: 'Jane Marketer', email: 'jane@example.com', role: 'affiliate' },
];

export const useAppStore = create<AppState>()(
    persist(
        (set, get) => ({
            currentUser: null,
            token: null,
            users: mockUsers,
            sales: [],
            commissions: [],
            payments: [],
            notifications: [],
            globalCommissionRate: 0.1, // 10%
            darkMode: false,
            login: (user, token) => set({ currentUser: user, token, notifications: [] }),
            logout: () => set({ currentUser: null, notifications: [] }),
            setNotifications: (notifications) => set({ notifications: notifications || [] }),
            toggleDarkMode: () => set((state) => ({ darkMode: !state.darkMode })),
            convertAllToIDR: () => {
                set((state) => {
                    const rate = 16000;
                    // Auto conversion logic if not yet converted
                    // Heuristic: if any sale is < 100000, probably not IDR yet.
                    const needsConversion = state.sales.some(s => s.amount < 100000);
                    if (!needsConversion) return state;

                    return {
                        sales: state.sales.map(s => ({ ...s, amount: s.amount * rate })),
                        commissions: state.commissions.map(c => ({ ...c, amount: c.amount * rate })),
                        payments: state.payments.map(p => ({ ...p, amount: p.amount * rate }))
                    };
                });
            },
            addSale: async (saleData) => {
                try {
                    const token = localStorage.getItem("token");
                    const apiUrl = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";
                    const response = await fetch(`${apiUrl}/admin/sales`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${token}`
                        },
                        body: JSON.stringify({
                            amount: saleData.amount,
                            affiliate_id: saleData.affiliateId
                        }),
                    });

                    const result = await response.json();
                    if (result.status === "success") {
                        const newSale = result.data;
                        
                        // We still need to update commissions and payments in local state 
                        // if they don't have their own backend endpoints yet.
                        // But NOT notifications, as backend handles that now.

                        const commissionAmount = saleData.amount * get().globalCommissionRate;
                        const newCommission: Commission = {
                            id: Math.random().toString(36).substr(2, 9),
                            saleId: newSale.id,
                            affiliateId: saleData.affiliateId,
                            amount: commissionAmount,
                            date: newSale.date,
                        };

                        const newPayment: Payment = {
                            id: Math.random().toString(36).substr(2, 9),
                            affiliateId: saleData.affiliateId,
                            amount: commissionAmount,
                            date: newSale.date,
                            status: 'pending',
                        };

                        set((state) => ({
                            sales: [newSale, ...(state.sales || [])],
                            commissions: [newCommission, ...(state.commissions || [])],
                            payments: [newPayment, ...(state.payments || [])],
                        }));
                    }
                } catch (err) {
                    console.error("Gagal mencatat penjualan:", err);
                }
            },
            markPaymentPaid: (paymentId) => {
                set((state) => ({
                    payments: state.payments.map((p) =>
                        p.id === paymentId ? { ...p, status: 'paid' } : p
                    ),
                }));
            },
            markNotificationRead: (id) => {
                set((state) => ({
                    notifications: (state.notifications || []).map((n) =>
                        n ? (n.id === id ? { ...n, is_read: true } : n) : n
                    ),
                }));
            },
        }),
        {
            name: 'loopaffi-storage',
        }
    )
);
