import { supabase } from "./supabase"
import type { Resource, Reservation } from "./types"

export type Slot = {
  starts_at: string
  ends_at: string
  available: boolean
  booked_by_me: boolean
}

type DayHours = {
  open: string
  close: string
}

const DAY_NAMES = ["sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"]

function parseHHMM(hhmm: string): { hour: number; minute: number } | null {
  const [h, m] = hhmm.split(":").map(Number)
  if (isNaN(h) || isNaN(m)) return null
  return { hour: h, minute: m }
}

export async function getAvailability(
  resourceId: string,
  from: Date,
  to: Date,
): Promise<{ resource_id: string; from: string; to: string; slots: Slot[] }> {
  // Fetch resource, overlapping reservations, and current user in parallel
  const [resourceResult, reservationsResult, userResult] = await Promise.all([
    supabase.from("resources").select("*").eq("id", resourceId).single(),
    supabase
      .from("reservations")
      .select("starts_at, ends_at, user_id")
      .eq("resource_id", resourceId)
      .in("status", ["HELD", "CONFIRMED"])
      .lt("starts_at", to.toISOString())
      .gt("ends_at", from.toISOString()),
    supabase.auth.getUser(),
  ])

  if (resourceResult.error) throw resourceResult.error
  const resource = resourceResult.data as Resource
  const reservations = (reservationsResult.data ?? []) as Pick<Reservation, "starts_at" | "ends_at" | "user_id">[]
  const currentUserId = userResult.data.user?.id

  // Parse open_hours
  const openHours: Record<string, DayHours> = typeof resource.open_hours === "object" && resource.open_hours !== null
    ? (resource.open_hours as unknown as Record<string, DayHours>)
    : {}

  const slotMs = resource.slot_minutes * 60_000
  const bufferMs = resource.buffer_minutes * 60_000

  const slots: Slot[] = []

  // Generate slots day by day
  const day = new Date(from)
  while (day < to) {
    const dayName = DAY_NAMES[day.getDay()]
    const hours = openHours[dayName] ?? { open: "06:00", close: "22:00" }

    const openTime = parseHHMM(hours.open)
    const closeTime = parseHHMM(hours.close)
    if (!openTime || !closeTime) {
      day.setDate(day.getDate() + 1)
      continue
    }

    const dayStart = new Date(day)
    dayStart.setHours(openTime.hour, openTime.minute, 0, 0)

    const dayEnd = new Date(day)
    dayEnd.setHours(closeTime.hour, closeTime.minute, 0, 0)

    let slotStart = dayStart.getTime()
    while (slotStart + slotMs <= dayEnd.getTime()) {
      const slotEnd = slotStart + slotMs

      let available = true
      let booked_by_me = false
      for (const r of reservations) {
        const rStart = new Date(r.starts_at).getTime()
        const rEnd = new Date(r.ends_at).getTime()
        if (rStart < slotEnd && rEnd > slotStart) {
          available = false
          if (r.user_id === currentUserId) booked_by_me = true
          break
        }
      }

      slots.push({
        starts_at: new Date(slotStart).toISOString(),
        ends_at: new Date(slotEnd).toISOString(),
        available,
        booked_by_me,
      })

      slotStart += slotMs + bufferMs
    }

    day.setDate(day.getDate() + 1)
  }

  slots.sort((a, b) => a.starts_at.localeCompare(b.starts_at))

  return {
    resource_id: resourceId,
    from: from.toISOString(),
    to: to.toISOString(),
    slots,
  }
}
