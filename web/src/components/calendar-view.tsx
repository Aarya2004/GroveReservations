"use client"

import { useMemo } from "react"
import { useQuery } from "@tanstack/react-query"
import { getAvailability, type Slot } from "@/lib/availability"
import {
  startOfWeek,
  endOfWeek,
  addDays,
  format,
  isSameDay,
  parseISO,
  addWeeks,
  subWeeks,
} from "date-fns"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { cn } from "@/lib/utils"

interface CalendarViewProps {
  resourceId: string
  onSelectSlot: (startsAt: Date, endsAt: Date) => void
}

export function CalendarView({ resourceId, onSelectSlot }: CalendarViewProps) {
  const [weekStart, setWeekStart] = useState(() => startOfWeek(new Date(), { weekStartsOn: 1 }))

  const fromDate = weekStart
  const toDate = endOfWeek(weekStart, { weekStartsOn: 1 })

  const { data, isLoading } = useQuery({
    queryKey: ["availability", resourceId, fromDate.toISOString()],
    queryFn: () => getAvailability(resourceId, fromDate, toDate),
  })

  const days = useMemo(() => {
    const d = []
    for (let i = 0; i < 7; i++) d.push(addDays(weekStart, i))
    return d
  }, [weekStart])

  const slotsByDay = useMemo(() => {
    if (!data?.slots) return {}
    const map: Record<string, Slot[]> = {}
    for (const slot of data.slots) {
      const day = format(parseISO(slot.starts_at), "yyyy-MM-dd")
      if (!map[day]) map[day] = []
      map[day].push(slot)
    }
    return map
  }, [data])

  const now = new Date()

  return (
    <div className="space-y-4">
      {/* Week navigation */}
      <div className="flex items-center justify-between">
        <Button variant="outline" size="sm" onClick={() => setWeekStart((w) => subWeeks(w, 1))}>
          <ChevronLeft className="size-4" />
          Previous
        </Button>
        <span className="text-sm font-medium text-muted-foreground">
          {format(weekStart, "MMM d")} — {format(addDays(weekStart, 6), "MMM d, yyyy")}
        </span>
        <Button variant="outline" size="sm" onClick={() => setWeekStart((w) => addWeeks(w, 1))}>
          Next
          <ChevronRight className="size-4" />
        </Button>
      </div>

      {isLoading && (
        <div className="flex items-center justify-center py-12 text-muted-foreground">
          Loading availability...
        </div>
      )}

      {!isLoading && (
        <div className="grid grid-cols-7 gap-px overflow-hidden rounded-lg border bg-border">
          {days.map((day) => {
            const dayKey = format(day, "yyyy-MM-dd")
            const slots = slotsByDay[dayKey] || []
            const isToday = isSameDay(day, now)

            return (
              <div key={dayKey} className="bg-background">
                {/* Day header */}
                <div
                  className={cn(
                    "border-b px-2 py-2 text-center text-xs font-medium",
                    isToday ? "bg-primary/10 text-primary" : "text-muted-foreground",
                  )}
                >
                  <div>{format(day, "EEE")}</div>
                  <div className={cn("text-lg", isToday && "font-bold")}>{format(day, "d")}</div>
                </div>

                {/* Slots */}
                <div className="flex flex-col gap-px p-1" style={{ minHeight: 200 }}>
                  {slots.length === 0 && (
                    <div className="flex flex-1 items-center justify-center text-xs text-muted-foreground/50">
                      —
                    </div>
                  )}
                  {slots.map((slot) => {
                    const start = parseISO(slot.starts_at)
                    const isPast = start < now
                    const clickable = slot.available && !isPast

                    return (
                      <button
                        key={slot.starts_at}
                        disabled={!clickable}
                        onClick={() =>
                          clickable &&
                          onSelectSlot(parseISO(slot.starts_at), parseISO(slot.ends_at))
                        }
                        className={cn(
                          "rounded px-1.5 py-1 text-[11px] leading-tight transition-colors",
                          clickable
                            ? "cursor-pointer bg-emerald-50 text-emerald-700 hover:bg-emerald-100 dark:bg-emerald-950/40 dark:text-emerald-400 dark:hover:bg-emerald-950/60"
                            : isPast
                              ? "cursor-default bg-muted/40 text-muted-foreground/40 line-through"
                              : slot.booked_by_me
                                ? "cursor-default bg-blue-50 text-blue-600 dark:bg-blue-950/40 dark:text-blue-400"
                                : "cursor-default bg-red-50 text-red-400 dark:bg-red-950/30 dark:text-red-400/70",
                        )}
                      >
                        {format(start, "HH:mm")}
                      </button>
                    )
                  })}
                </div>
              </div>
            )
          })}
        </div>
      )}

      {/* Legend */}
      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        <span className="flex items-center gap-1.5">
          <span className="size-2.5 rounded-sm bg-emerald-100 dark:bg-emerald-950/40" />
          Available
        </span>
        <span className="flex items-center gap-1.5">
          <span className="size-2.5 rounded-sm bg-blue-50 dark:bg-blue-950/40" />
          Your booking
        </span>
        <span className="flex items-center gap-1.5">
          <span className="size-2.5 rounded-sm bg-red-50 dark:bg-red-950/30" />
          Booked
        </span>
        <span className="flex items-center gap-1.5">
          <span className="size-2.5 rounded-sm bg-muted/40" />
          Past
        </span>
      </div>
    </div>
  )
}
