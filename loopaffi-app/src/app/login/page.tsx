"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAppStore } from "@/lib/store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Infinity, Loader2 } from "lucide-react";
import Link from "next/link";


export default function LoginPage() {
    const router = useRouter();
    const { login } = useAppStore();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        setIsLoading(true);

        try {
            const apiUrl = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080") + "/api/v1";
            const response = await fetch(`${apiUrl}/auth/login`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password }),
            });

            const result = await response.json();

            if (response.ok && result.status === "success") {
                const { token, user } = result.data;
                // Simpan token ke localStorage atau cookie jika diperlukan
                localStorage.setItem("token", token);
                
                // Login ke store (pastikan store menerima struktur user yang benar)
                login(user, token);

                // Redirect berdasarkan role_id (role_admin -> /admin, role_affiliate -> /affiliate)
                const targetPath = user.role === "role_admin" ? "admin" : "affiliate";
                router.push(`/${targetPath}/dashboard`);
            } else {
                setError(result.message || "Email atau password salah.");
            }
        } catch (err) {
            setError("Gagal terhubung ke server. Pastikan backend sudah menyala.");
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-slate-50">
            {/* Background decoration */}
            <div className="absolute inset-0 overflow-hidden pointer-events-none">
                <div className="absolute -top-40 -right-40 w-80 h-80 bg-red-100 rounded-full blur-3xl opacity-40" />
                <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-orange-100 rounded-full blur-3xl opacity-40" />
            </div>

            <div className="relative w-full max-w-md bg-white p-8 rounded-2xl shadow-lg shadow-slate-200/50 border border-slate-100">
                <div className="flex flex-col items-center mb-8">
                    <div className="w-14 h-14 bg-gradient-to-br from-red-500 to-red-700 text-white rounded-2xl flex items-center justify-center mb-4 shadow-lg shadow-red-500/25">
                        <Infinity className="w-8 h-8" />
                    </div>
                    <h1 className="text-2xl font-bold text-slate-900">Selamat Datang di LoopAffi</h1>
                    <p className="text-slate-500 text-sm mt-1">Masuk untuk mengelola afiliasi Anda</p>
                </div>

                <form onSubmit={handleLogin} className="space-y-5">
                    {error && (
                        <div className="p-3 text-sm text-red-600 bg-red-50 rounded-lg border border-red-100 flex items-start gap-2">
                            <span className="shrink-0 mt-0.5">⚠️</span>
                            <span>{error}</span>
                        </div>
                    )}

                    <div className="space-y-2">
                        <Label htmlFor="input-email" className="text-sm font-medium text-slate-700">Alamat Email</Label>
                        <Input
                            id="input-email"
                            data-testid="input-email"
                            type="email"
                            placeholder="admin@loopaffi.com"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            className="h-11"
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="input-password" className="text-sm font-medium text-slate-700">Kata Sandi</Label>
                        <Input
                            id="input-password"
                            data-testid="input-password"
                            type="password"
                            placeholder="••••••••"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            className="h-11"
                        />
                    </div>

                    <Button
                        id="btn-login"
                        data-testid="btn-login"
                        type="submit"
                        disabled={isLoading}
                        className="w-full h-11 bg-gradient-to-r from-red-600 to-red-700 hover:from-red-700 hover:to-red-800 text-white text-sm font-semibold shadow-lg shadow-red-500/25 transition-all duration-200"
                    >
                        {isLoading ? (
                            <span className="flex items-center gap-2">
                                <Loader2 className="w-4 h-4 animate-spin" />
                                Sedang masuk...
                            </span>
                        ) : "Masuk"}
                    </Button>
                </form>

                <div className="mt-6 text-center">
                    <p className="text-sm text-slate-500 font-medium">
                        Belum punya akun affiliator?{" "}
                        <Link href="/register" className="text-red-600 font-bold hover:text-red-700 underline underline-offset-4 decoration-2 decoration-red-200 hover:decoration-red-600 transition-all">
                            Daftar di sini
                        </Link>
                    </p>
                </div>


                <div className="mt-8 pt-6 border-t border-slate-100 flex gap-4 text-center justify-center">
                    <button
                        id="btn-quickfill-admin"
                        data-testid="btn-quickfill-admin"
                        onClick={() => { setEmail("admin@loopaffi.com"); setPassword("password"); }}
                        className="text-xs text-slate-500 hover:text-red-600 underline underline-offset-4 transition-colors"
                    >
                        Gunakan Admin
                    </button>
                    <button
                        id="btn-quickfill-affiliate"
                        data-testid="btn-quickfill-affiliate"
                        onClick={() => { setEmail("john@example.com"); setPassword("password"); }}
                        className="text-xs text-slate-500 hover:text-red-600 underline underline-offset-4 transition-colors"
                    >
                        Gunakan Afiliasi
                    </button>
                </div>
            </div>
        </div>
    );
}
