export type Resource = {
  id: string
  name: string
  type: string
  location?: string | null
  slot_minutes: number
  buffer_minutes: number
  max_advance_days: number
  open_hours: Record<string, string[]>
}
export type Reservation = {
  id: string
  resource_id: string
  user_id: string
  starts_at: string
  ends_at: string
  status: string
}
