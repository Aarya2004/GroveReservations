"use client"

import { useQuery } from "@tanstack/react-query"
import { listResources } from "@/lib/api"
import Link from "next/link"
import { Clock, MapPin, ArrowRight } from "lucide-react"

export default function ResourcesPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["resources"],
    queryFn: listResources,
  })

  if (isLoading) return <p className="py-12 text-center text-muted-foreground">Loading...</p>
  if (error) return <p className="py-12 text-center text-destructive">{String(error)}</p>

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Resources</h1>
        <p className="text-sm text-muted-foreground">
          Select a resource to view availability and book a slot.
        </p>
      </div>

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {data!.map((r) => (
          <Link
            key={r.id}
            href={`/resources/${r.id}`}
            className="group flex flex-col justify-between rounded-xl border bg-card p-5 shadow-sm transition-all hover:border-primary/30 hover:shadow-md"
          >
            <div>
              <div className="font-semibold text-card-foreground group-hover:text-primary transition-colors">
                {r.name}
              </div>
              <div className="mt-1 text-sm capitalize text-muted-foreground">
                {r.type.replace("_", " ")}
              </div>
            </div>

            <div className="mt-4 flex items-center justify-between text-xs text-muted-foreground">
              <div className="flex items-center gap-3">
                {r.location && (
                  <span className="flex items-center gap-1">
                    <MapPin className="size-3" />
                    {r.location}
                  </span>
                )}
                <span className="flex items-center gap-1">
                  <Clock className="size-3" />
                  {r.slot_minutes}m
                </span>
              </div>
              <ArrowRight className="size-3.5 text-muted-foreground/50 group-hover:text-primary transition-colors" />
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
