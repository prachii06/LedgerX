"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { Activity, LayoutDashboard, List, ActivitySquare, Server, ServerCrash, LineChart, ExternalLink } from "lucide-react"
import { cn } from "@/lib/utils"

const navigation = [
  { name: "Overview", href: "/overview", icon: LayoutDashboard },
  { name: "Transactions", href: "/transactions", icon: List },
  { name: "Reconciliation", href: "/reconciliation", icon: ActivitySquare },
  { name: "Events", href: "/events", icon: Activity },
]

const systemNavigation = [
  { name: "System Health", href: "/system", icon: ServerCrash },
  { name: "Metrics", href: "/system/metrics", icon: Server },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <div className="flex h-full w-64 flex-col border-r bg-card px-4 py-6 text-card-foreground">
      <div className="flex items-center gap-2 px-2 pb-6">
        <span className="text-xl font-bold tracking-tight">LedgerX</span>
      </div>

      <nav className="flex-1 space-y-1 overflow-y-auto">
        {navigation.map((item) => {
          const isActive = pathname === item.href || pathname.startsWith(item.href + "/")
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "group flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-secondary text-secondary-foreground"
                  : "text-muted-foreground hover:bg-muted hover:text-foreground"
              )}
            >
              <item.icon className="h-4 w-4" />
              {item.name}
            </Link>
          )
        })}

        <div className="pt-6 pb-2">
          <p className="px-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            System
          </p>
        </div>

        {systemNavigation.map((item) => {
          const isActive = pathname === item.href || pathname.startsWith(item.href + "/")
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "group flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-secondary text-secondary-foreground"
                  : "text-muted-foreground hover:bg-muted hover:text-foreground"
              )}
            >
              <item.icon className="h-4 w-4" />
              {item.name}
            </Link>
          )
        })}

        <div className="pt-6 pb-2">
          <p className="px-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Observability
          </p>
        </div>

        <a
          href="http://localhost:3001/d/1d1a72f3-9719-4fa2-8ad4-d6092fd3df9c/ledgerx-overview?orgId=1&from=now-6h&to=now&timezone=browser&refresh=5s"
          target="_blank"
          rel="noopener noreferrer"
          className={cn(
            "group flex items-center justify-between rounded-md px-3 py-2 text-sm font-medium transition-colors text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
        >
          <div className="flex items-center gap-3">
            <LineChart className="h-4 w-4" />
            Grafana
          </div>
          <ExternalLink className="h-3 w-3 text-muted-foreground/70" />
        </a>
      </nav>

      <div className="mt-auto border-t pt-4">
        <div className="flex items-center justify-between px-3 py-2">
          <div className="flex items-center gap-2 text-sm">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex h-2 w-2 rounded-full bg-green-500"></span>
            </span>
            <span className="font-medium">Production</span>
          </div>
        </div>
        <div className="flex items-center gap-3 px-3 py-3">
          <div className="h-8 w-8 rounded-full bg-secondary flex items-center justify-center font-semibold text-sm">
            P
          </div>
          <div className="flex flex-col">
            <span className="text-sm font-medium">Prachi</span>
            <span className="text-xs text-muted-foreground">Ops Team</span>
          </div>
        </div>
      </div>
    </div>
  )
}
