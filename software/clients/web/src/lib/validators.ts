import { z } from "zod"

export const createReservationSchema = z.object({
  resource_id: z.string().uuid(),
  starts_at: z.string().datetime(),
  ends_at: z.string().datetime(),
})

export type CreateReservationInput = z.infer<typeof createReservationSchema>

export const updateUserSchema = z.object({
  name: z.string().min(1).optional(),
  role: z.enum(["ADMIN", "MEMBER", "GUEST"]).optional(),
  villa_number: z.number().int().positive().optional(),
  phone_number: z.string().optional(),
})
export type UpdateUserInput = z.infer<typeof updateUserSchema>
