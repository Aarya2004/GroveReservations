"use client"

import { use, useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { getResource } from "@/lib/api"
import { CalendarView } from "@/components/calendar-view"
import { BookingModal } from "@/components/booking-modal"
import { Button } from "@/components/ui/button"
import { ArrowLeft, Clock, MapPin, CalendarDays, Timer } from "lucide-react"
import Link from "next/link"

export default function ResourceDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)

  const {
    data: resource,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["resource", id],
    queryFn: () => getResource(id),
  })

  const [booking, setBooking] = useState<{
    startsAt: Date
    endsAt: Date
  } | null>(null)

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20 text-muted-foreground">
        Loading resource...
      </div>
    )
  }

  if (error || !resource) {
    return (
      <div className="py-20 text-center">
        <p className="text-destructive">Failed to load resource</p>
        <Button variant="outline" className="mt-4" asChild>
          <Link href="/resources">Back to resources</Link>
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Back link */}
      <Link
        href="/resources"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        <ArrowLeft className="size-3.5" />
        All resources
      </Link>

      {/* Resource info */}
      <div className="space-y-2">
        <h1 className="text-2xl font-semibold tracking-tight">{resource.name}</h1>
        <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
          <span className="flex items-center gap-1.5 capitalize">
            <CalendarDays className="size-3.5" />
            {resource.type.replace("_", " ")}
          </span>
          {resource.location && (
            <span className="flex items-center gap-1.5">
              <MapPin className="size-3.5" />
              {resource.location}
            </span>
          )}
          <span className="flex items-center gap-1.5">
            <Clock className="size-3.5" />
            {resource.slot_minutes} min slots
          </span>
          {resource.buffer_minutes > 0 && (
            <span className="flex items-center gap-1.5">
              <Timer className="size-3.5" />
              {resource.buffer_minutes} min buffer
            </span>
          )}
        </div>
      </div>

      {/* Calendar */}
      <div>
        <h2 className="mb-3 text-lg font-medium">Availability</h2>
        <CalendarView
          resourceId={id}
          onSelectSlot={(startsAt, endsAt) => setBooking({ startsAt, endsAt })}
        />
      </div>

      {/* Booking modal */}
      {booking && (
        <BookingModal
          open={true}
          onClose={() => setBooking(null)}
          resourceId={id}
          resourceName={resource.name}
          startsAt={booking.startsAt}
          endsAt={booking.endsAt}
        />
      )}
    </div>
  )
}
