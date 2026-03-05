import { createBrowserClient, createServerClient } from "@supabase/ssr"
import type { NextRequest } from "next/server"
import type { NextResponse } from "next/server"

// Browser client (used in client components and hooks)
// Uses @supabase/ssr so PKCE code_verifier is stored in cookies, not localStorage
export function createClient() {
  return createBrowserClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
  )
}

// Singleton for client components that import `supabase` directly
export const supabase = createClient()

// Server client (used in route handlers and middleware)
export function createServerSupabase(request: NextRequest, response: NextResponse) {
  return createServerClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        getAll() {
          return request.cookies.getAll()
        },
        setAll(cookiesToSet) {
          for (const { name, value, options } of cookiesToSet) {
            response.cookies.set(name, value, options)
          }
        },
      },
    },
  )
}
