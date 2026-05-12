"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAppStore } from "@/lib/store";
import { Sidebar } from "./Sidebar";
import { Topbar } from "./Topbar";

// Daftar halaman publik yang bisa diakses TANPA login
const PUBLIC_ROUTES = ["/login", "/register"];

export function ClientWrapper({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const pathname = usePathname();
    const { currentUser } = useAppStore();
    const [mounted, setMounted] = useState(false);

    useEffect(() => {
        setMounted(true);
    }, []);

    useEffect(() => {
        if (!mounted) return;

        const isPublicRoute = PUBLIC_ROUTES.includes(pathname) || pathname === "/";

        // User belum login & bukan di halaman publik → tendang ke login
        if (!currentUser && !isPublicRoute) {
            router.push("/login");
            return;
        }

        // User sudah login → atur routing berdasarkan role
        if (currentUser) {
            // Kalau user sudah login tapi masih di halaman login/register/root → ke dashboard
            if (isPublicRoute || pathname === "/") {
                router.push(`/${currentUser.role}/dashboard`);
                return;
            }

            // Proteksi role: admin tidak bisa akses halaman affiliate, dan sebaliknya
            const adminPath = pathname.startsWith("/admin");
            const affiliatePath = pathname.startsWith("/affiliate");

            if (currentUser.role === "admin" && affiliatePath) {
                router.push("/admin/dashboard");
            } else if (currentUser.role === "affiliate" && adminPath) {
                router.push("/affiliate/dashboard");
            }
        }
    }, [currentUser, pathname, router, mounted]);

    if (!mounted) return null;

    const isPublicRoute = PUBLIC_ROUTES.includes(pathname) || pathname === "/";

    if (isPublicRoute) {
        return <>{children}</>;
    }

    // Halaman terproteksi → tampilkan DENGAN Sidebar & Topbar
    return (
        <div className="flex h-screen bg-slate-50 text-slate-900">
            <Sidebar />
            <div className="flex-1 flex flex-col overflow-hidden">
                <Topbar />
                <main className="flex-1 overflow-y-auto p-8 bg-white">
                    {children}
                </main>
            </div>
        </div>
    );
}
