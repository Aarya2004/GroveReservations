"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/hooks/useAuth"

export default function Home() {
  const { user, loading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (loading) return
    if (user) {
      router.replace("/resources")
    } else {
      router.replace("/login")
    }
  }, [user, loading, router])

  return (
    <div className="flex min-h-svh items-center justify-center">
      <p className="text-muted-foreground">Loading...</p>
    </div>
  )
}
