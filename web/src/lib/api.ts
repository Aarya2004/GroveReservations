import { supabase } from "./supabase"
import type { Resource, Reservation } from "./types"

// ── Resources ──────────────────────────────────────────────

export async function listResources(): Promise<Resource[]> {
  const { data, error } = await supabase
    .from("resources")
    .select("*")
    .order("name")

  if (error) throw error
  return data as Resource[]
}

export async function getResource(id: string): Promise<Resource> {
  const { data, error } = await supabase
    .from("resources")
    .select("*")
    .eq("id", id)
    .single()

  if (error) throw error
  return data as Resource
}

// ── Reservations ───────────────────────────────────────────

export async function listMyReservations(): Promise<Reservation[]> {
  const { data, error } = await supabase
    .from("reservations")
    .select("*")
    .order("starts_at", { ascending: false })

  if (error) throw error
  return data as Reservation[]
}

export async function createReservation(input: {
  resource_id: string
  starts_at: string
  ends_at: string
}): Promise<Reservation> {
  const {
    data: { user },
  } = await supabase.auth.getUser()
  if (!user) throw new Error("Not authenticated")

  const { data, error } = await supabase
    .from("reservations")
    .insert({
      resource_id: input.resource_id,
      user_id: user.id,
      starts_at: input.starts_at,
      ends_at: input.ends_at,
      status: "CONFIRMED",
    })
    .select()
    .single()

  if (error) {
    // PostgreSQL exclusion constraint violation
    if (error.code === "23P01") {
      throw new Error("409: This slot is no longer available.")
    }
    throw error
  }
  return data as Reservation
}

export async function cancelReservation(id: string): Promise<void> {
  const { error } = await supabase
    .from("reservations")
    .update({ status: "CANCELLED" })
    .eq("id", id)

  if (error) throw error
}
