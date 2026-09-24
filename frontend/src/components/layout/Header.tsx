"use client"

import { Bell, Search } from "lucide-react"
import { useState } from "react"
import { useRouter } from "next/navigation"

export function Header() {
  const [query, setQuery] = useState("")
  const router = useRouter()

  const handleSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && query.trim()) {
      // For now, assume the query is a transaction ID and route there.
      // If it doesn't exist, the transaction detail page handles the 404.
      router.push(`/transactions/${query.trim()}`)
    }
  }

  return (
    <header className="flex h-16 items-center justify-between border-b bg-card px-6 py-4 sticky top-0 z-50">
      <div className="flex flex-1 items-center gap-4">
        <div className="relative w-full max-w-md">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <input
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleSearch}
            placeholder="Search transaction ID... (Enter to search)"
            className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 pl-9 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
          />
        </div>
      </div>
      <div className="flex items-center gap-4">
        <button className="relative rounded-full p-2 text-muted-foreground hover:bg-muted hover:text-foreground">
          <Bell className="h-5 w-5" />
          <span className="absolute right-1.5 top-1.5 flex h-2 w-2 rounded-full bg-destructive" />
        </button>
      </div>
    </header>
  )
}
