"use client"

import { useEffect, useState } from "react"
import { supabase } from "@/lib/supabase"
import type { User, Session } from "@supabase/supabase-js"

type AuthState = {
  user: User | null
  session: Session | null
  loading: boolean
}

export function useAuth() {
  const [auth, setAuth] = useState<AuthState>({ user: null, session: null, loading: true })

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setAuth({ user: session?.user ?? null, session, loading: false })
    })

    const { data: { subscription } } = supabase.auth.onAuthStateChange(
      (_event, session) => {
        setAuth({ user: session?.user ?? null, session, loading: false })
      }
    )

    return () => subscription.unsubscribe()
  }, [])

  const signOut = async () => {
    await supabase.auth.signOut()
  }

  return { ...auth, signOut }
}
