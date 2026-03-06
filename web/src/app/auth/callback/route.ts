import { NextRequest, NextResponse } from "next/server"
import { createServerSupabase } from "@/lib/supabase"

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url)
  const code = searchParams.get("code")
  const next = searchParams.get("next") ?? "/resources"

  if (code) {
    const response = NextResponse.redirect(new URL(next, request.url))
    const supabase = createServerSupabase(request, response)

    const { error } = await supabase.auth.exchangeCodeForSession(code)
    if (error) {
      console.error("Auth callback error:", error.message)
    } else {
      return response
    }
  }

  return NextResponse.redirect(new URL("/login", request.url))
}
