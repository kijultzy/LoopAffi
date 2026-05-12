"use client";

import { useState, useEffect } from "react";
import { useAppStore } from "@/lib/store";
import { fetchMySales, fetchMyCommissions, fetchMyPayments, DBSale, DBCommission, DBPayment } from "@/lib/api";
import { TrendingUp, Coins, ShoppingCart, Clock, CheckCircle2, Loader2 } from "lucide-react";
import { formatIDR } from "@/lib/utils";

export default function AffiliateDashboard() {
    const { currentUser } = useAppStore();
    const [sales, setSales] = useState<DBSale[]>([]);
    const [commissions, setCommissions] = useState<DBCommission[]>([]);
    const [payments, setPayments] = useState<DBPayment[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            setIsLoading(true);
            try {
                const [salesData, commissionsData, paymentsData] = await Promise.all([
                    fetchMySales(),
                    fetchMyCommissions(),
                    fetchMyPayments(),
                ]);
                setSales(salesData);
                setCommissions(commissionsData);
                setPayments(paymentsData);
            } catch (err) {
                console.error("Gagal memuat data dashboard afiliasi:", err);
            } finally {
                setIsLoading(false);
            }
        };
        loadData();
    }, []);

    if (!currentUser) return null;
    if (isLoading) {
        return (
            <div className="max-w-6xl mx-auto flex items-center justify-center py-20">
                <div className="flex items-center gap-2 text-slate-500">
                    <Loader2 className="w-5 h-5 animate-spin" />
                    Memuat data dashboard...
                </div>
            </div>
        );
    }

    const totalSalesAmount = sales.reduce((acc, s) => acc + s.amount, 0);
    const totalEarned = commissions.reduce((acc, c) => acc + c.amount, 0);
    const pendingPayout = payments.filter((p) => p.status === "pending").reduce((acc, p) => acc + p.amount, 0);
    const paidOut = payments.filter((p) => p.status === "paid").reduce((acc, p) => acc + p.amount, 0);

    const recentSales = [...sales]
        .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
        .slice(0, 5);


    return (
        <div className="max-w-6xl mx-auto space-y-8">
            <div>
                <h2 className="text-3xl font-bold text-slate-900 tracking-tight">Dasbor</h2>
                <p className="text-slate-500 mt-1">Selamat datang kembali, {currentUser.name}. Berikut adalah rincian afiliasi Anda.</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {[
                    { title: "Penjualan Rujukan", value: formatIDR(totalSalesAmount), sub: `${sales.length} konversi`, icon: ShoppingCart },
                    { title: "Total Pendapatan", value: formatIDR(totalEarned), sub: "Komisi seumur hidup", icon: Coins },
                    { title: "Pembayaran Tertunda", value: formatIDR(pendingPayout), sub: "Akan diproses", icon: Clock },
                    { title: "Telah Dibayar", value: formatIDR(paidOut), sub: "Diterima sejauh ini", icon: TrendingUp },
                ].map((card, i) => {
                    const Icon = card.icon;
                    return (
                        <div key={i} className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm flex flex-col gap-4 relative overflow-hidden">
                            <div className="absolute top-0 right-0 p-4 opacity-5">
                                <Icon className="w-24 h-24" />
                            </div>
                            <div className="flex items-center gap-3 relative z-10">
                                <div className="w-10 h-10 rounded-full bg-red-50 text-red-600 flex items-center justify-center">
                                    <Icon className="w-5 h-5" />
                                </div>
                                <h3 className="text-sm font-semibold text-slate-600">{card.title}</h3>
                            </div>
                            <div className="relative z-10">
                                <p className="text-3xl font-bold text-slate-900">{card.value}</p>
                                <p className="text-sm text-slate-500 mt-1">{card.sub}</p>
                            </div>
                        </div>
                    );
                })}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
                    <div className="px-6 py-5 border-b border-slate-100 flex items-center justify-between">
                        <h3 className="font-semibold text-slate-900">Konversi Terbaru</h3>
                    </div>
                    <div className="p-6">
                        {recentSales.length === 0 ? (
                            <p className="text-slate-500 text-center text-sm py-4">Belum ada penjualan rujukan.</p>
                        ) : (
                            <div className="space-y-4">
                                {recentSales.map((sale) => (
                                    <div key={sale.id} className="flex items-center justify-between">
                                        <div className="flex items-center gap-3">
                                            <div className="w-8 h-8 rounded-full bg-slate-50 flex items-center justify-center border border-slate-100">
                                                <CheckCircle2 className="w-4 h-4 text-green-600" />
                                            </div>
                                            <div>
                                                <p className="text-sm font-medium text-slate-900">Penjualan #{sale.id}</p>
                                                <p className="text-xs text-slate-500">{new Date(sale.date).toLocaleDateString()}</p>
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <p className="text-sm font-bold text-slate-900">{formatIDR(sale.amount)}</p>
                                            <p className="text-xs text-red-600 font-medium">+{formatIDR(sale.amount * 0.1)}</p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>

                <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm flex flex-col gap-6">
                    <h3 className="font-semibold text-slate-900">Info Komisi</h3>
                    <div className="bg-red-50 p-4 rounded-xl border border-red-100 text-red-900">
                        <h4 className="font-semibold">Anda mendapatkan 10% untuk setiap penjualan!</h4>
                        <p className="text-sm mt-1 text-red-700">Ketika pelanggan membeli melalui rujukan Anda, Anda secara otomatis menerima komisi sebesar 10%. Pembayaran diproses secara berkala oleh tim admin.</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
