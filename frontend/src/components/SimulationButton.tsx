"use client"

import { useState, useRef, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { runSimulation } from "@/lib/api"
import { PlayCircle, Loader2, X, CheckCircle2 } from "lucide-react"
import { useRouter } from "next/navigation"

export function SimulationButton() {
  const [loading, setLoading] = useState(false)
  const [expanded, setExpanded] = useState(false)
  const [count, setCount] = useState("5")
  const router = useRouter()
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (expanded && inputRef.current) {
      inputRef.current.focus()
      inputRef.current.select()
    }
  }, [expanded])

  async function handleSimulate() {
    const simCount = parseInt(count, 10)
    if (isNaN(simCount) || simCount <= 0) return

    setLoading(true)
    try {
      await runSimulation(simCount)
      router.refresh()
      setExpanded(false)
    } catch (e) {
      console.error(e)
      alert(e instanceof Error ? e.message : "Failed to run simulation")
    } finally {
      setLoading(false)
    }
  }

  if (expanded) {
    return (
      <div className="flex items-center gap-2 animate-in fade-in slide-in-from-right-4 duration-300">
        <Input
          ref={inputRef}
          type="number"
          min="1"
          max="1000"
          value={count}
          onChange={(e) => setCount(e.target.value)}
          className="w-20 h-9"
          disabled={loading}
          onKeyDown={(e) => {
            if (e.key === 'Enter') handleSimulate()
            if (e.key === 'Escape') setExpanded(false)
          }}
        />
        <Button onClick={handleSimulate} disabled={loading} size="sm" className="gap-2 h-9">
          {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <CheckCircle2 className="h-4 w-4" />}
          Run
        </Button>
        <Button variant="ghost" size="icon" className="h-9 w-9 text-muted-foreground" onClick={() => setExpanded(false)} disabled={loading}>
          <X className="h-4 w-4" />
        </Button>
      </div>
    )
  }

  return (
    <Button onClick={() => setExpanded(true)} className="gap-2">
      <PlayCircle className="h-4 w-4" />
      Simulate Traffic
    </Button>
  )
}
