"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "@/lib/api"
import type { Reservation } from "@/lib/types"
import { Button } from "@/components/ui/button"
import { format, parseISO } from "date-fns"
import { Calendar, ExternalLink, Trash2 } from "lucide-react"
import Link from "next/link"
import { cn } from "@/lib/utils"

const statusConfig: Record<string, { label: string; className: string }> = {
  CONFIRMED: {
    label: "Confirmed",
    className: "bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-400",
  },
  HELD: {
    label: "Held",
    className: "bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400",
  },
  CANCELLED: {
    label: "Cancelled",
    className: "bg-red-50 text-red-600 dark:bg-red-950/40 dark:text-red-400",
  },
  NOSHOW: {
    label: "No-show",
    className: "bg-muted text-muted-foreground",
  },
  COMPLETED: {
    label: "Completed",
    className: "bg-blue-50 text-blue-700 dark:bg-blue-950/40 dark:text-blue-400",
  },
}

export default function ReservationsPage() {
  const queryClient = useQueryClient()

  const { data, isLoading, error } = useQuery({
    queryKey: ["my-reservations"],
    queryFn: () =>
      api<{ reservations: Reservation[] }>("/reservations/me").then(
        (d) => d.reservations
      ),
  })

  const cancelMutation = useMutation({
    mutationFn: (id: string) =>
      api(`/reservations/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["my-reservations"] })
    },
  })

  if (isLoading) {
    return (
      <p className="py-12 text-center text-muted-foreground">
        Loading reservations...
      </p>
    )
  }

  if (error) {
    return (
      <p className="py-12 text-center text-destructive">{String(error)}</p>
    )
  }

  const reservations = data ?? []

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">
          My Reservations
        </h1>
        <p className="text-sm text-muted-foreground">
          View and manage your upcoming bookings.
        </p>
      </div>

      {reservations.length === 0 ? (
        <div className="flex flex-col items-center gap-3 rounded-xl border border-dashed py-16 text-center">
          <Calendar className="size-8 text-muted-foreground/40" />
          <p className="text-muted-foreground">No reservations yet.</p>
          <Button asChild variant="outline" size="sm">
            <Link href="/resources">Browse resources</Link>
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          {reservations.map((r) => {
            const start = parseISO(r.starts_at)
            const end = parseISO(r.ends_at)
            const status = statusConfig[r.status] ?? statusConfig.CONFIRMED
            const canCancel =
              r.status === "CONFIRMED" || r.status === "HELD"

            return (
              <div
                key={r.id}
                className="flex items-center gap-4 rounded-lg border bg-card p-4 shadow-sm"
              >
                {/* Date badge */}
                <div className="hidden sm:flex flex-col items-center rounded-lg border bg-muted/30 px-3 py-2 text-center">
                  <span className="text-xs font-medium uppercase text-muted-foreground">
                    {format(start, "MMM")}
                  </span>
                  <span className="text-xl font-bold leading-tight">
                    {format(start, "d")}
                  </span>
                </div>

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="font-medium truncate">
                      {format(start, "EEE, MMM d")}
                    </span>
                    <span
                      className={cn(
                        "inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium",
                        status.className
                      )}
                    >
                      {status.label}
                    </span>
                  </div>
                  <div className="mt-0.5 text-sm text-muted-foreground">
                    {format(start, "HH:mm")} — {format(end, "HH:mm")}
                  </div>
                </div>

                {/* Actions */}
                <div className="flex items-center gap-2">
                  <Button variant="ghost" size="icon" asChild>
                    <Link href={`/resources/${r.resource_id}`}>
                      <ExternalLink className="size-4" />
                    </Link>
                  </Button>
                  {canCancel && (
                    <Button
                      variant="ghost"
                      size="icon"
                      className="text-destructive hover:text-destructive hover:bg-destructive/10"
                      disabled={cancelMutation.isPending}
                      onClick={() => cancelMutation.mutate(r.id)}
                    >
                      <Trash2 className="size-4" />
                    </Button>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
