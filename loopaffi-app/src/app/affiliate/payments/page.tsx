"use client";

import { useState, useEffect } from "react";
import { useAppStore } from "@/lib/store";
import { fetchMyPayments, DBPayment } from "@/lib/api";
import { formatIDR } from "@/lib/utils";
import { Loader2 } from "lucide-react";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";

export default function AffiliatePaymentsPage() {
    const { currentUser } = useAppStore();
    const [payments, setPayments] = useState<DBPayment[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadPayments = async () => {
            setIsLoading(true);
            try {
                const data = await fetchMyPayments();
                setPayments(data);
            } catch (err) {
                console.error("Gagal memuat data pembayaran:", err);
            } finally {
                setIsLoading(false);
            }
        };
        loadPayments();
    }, []);

    if (!currentUser) return null;

    return (
        <div className="max-w-6xl mx-auto space-y-8">
            <div>
                <h2 className="text-3xl font-bold text-slate-900 tracking-tight">Pembayaran Saya</h2>
                <p className="text-slate-500 mt-1">
                    Lacak riwayat pembayaran komisi Anda dan pembayaran yang tertunda.
                </p>
            </div>

            <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
                <Table>
                    <TableHeader className="bg-slate-50">
                        <TableRow>
                            <TableHead className="font-semibold text-slate-900">Tanggal</TableHead>
                            <TableHead className="font-semibold text-slate-900">ID Pembayaran</TableHead>
                            <TableHead className="font-semibold text-slate-900">Jumlah</TableHead>
                            <TableHead className="font-semibold text-slate-900 text-right">Status</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            <TableRow>
                                <TableCell colSpan={4} className="text-center py-8">
                                    <div className="flex items-center justify-center gap-2 text-slate-500">
                                        <Loader2 className="w-4 h-4 animate-spin" />
                                        Memuat data pembayaran...
                                    </div>
                                </TableCell>
                            </TableRow>
                        ) : payments.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={4} className="text-center text-slate-500 py-8">
                                    Belum ada pembayaran.
                                </TableCell>
                            </TableRow>
                        ) : (
                            payments.map((payment) => (
                                <TableRow key={payment.id}>
                                    <TableCell>{new Date(payment.date).toLocaleDateString()}</TableCell>
                                    <TableCell className="font-mono text-xs text-slate-500">{payment.id}</TableCell>
                                    <TableCell className="font-bold text-slate-900">{formatIDR(payment.amount)}</TableCell>
                                    <TableCell className="text-right">
                                        {payment.status === "paid" ? (
                                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-green-100 text-green-800">
                                                Lunas
                                            </span>
                                        ) : (
                                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-amber-100 text-amber-800">
                                                Tertunda
                                            </span>
                                        )}
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>
        </div>
    );
}
