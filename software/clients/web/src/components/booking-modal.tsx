"use client"

import { useState } from "react"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { api } from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { format } from "date-fns"
import { CalendarCheck, X, AlertCircle, CheckCircle2 } from "lucide-react"
import { cn } from "@/lib/utils"

interface BookingModalProps {
  open: boolean
  onClose: () => void
  resourceId: string
  resourceName?: string
  startsAt: Date
  endsAt: Date
}

export function BookingModal({
  open,
  onClose,
  resourceId,
  resourceName,
  startsAt,
  endsAt,
}: BookingModalProps) {
  const queryClient = useQueryClient()
  const [modal, setModal] = useState<{ state: "idle" | "success" | "error"; errorMsg: string }>({
    state: "idle",
    errorMsg: "",
  })

  const mutation = useMutation({
    mutationFn: () =>
      api("/reservations", {
        method: "POST",
        body: JSON.stringify({
          resource_id: resourceId,
          starts_at: startsAt.toISOString(),
          ends_at: endsAt.toISOString(),
        }),
      }),
    onSuccess: () => {
      setModal({ state: "success", errorMsg: "" })
      queryClient.invalidateQueries({ queryKey: ["availability", resourceId] })
    },
    onError: (err: Error) => {
      setModal({
        state: "error",
        errorMsg: err.message.includes("409") ? "This slot is no longer available." : err.message,
      })
    },
  })

  const handleConfirm = () => {
    setModal({ state: "idle", errorMsg: "" })
    mutation.mutate()
  }

  const handleClose = () => {
    setModal({ state: "idle", errorMsg: "" })
    onClose()
  }

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={handleClose} />

      {/* Modal */}
      <Card className="relative z-10 w-full max-w-md mx-4 shadow-xl animate-in fade-in zoom-in-95 duration-200">
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <CalendarCheck className="size-5 text-emerald-600" />
              {modal.state === "success" ? "Booking confirmed" : "Confirm booking"}
            </CardTitle>
            <button
              onClick={handleClose}
              className="rounded-md p-1 text-muted-foreground hover:bg-accent hover:text-foreground transition-colors"
            >
              <X className="size-4" />
            </button>
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          {modal.state === "success" ? (
            <div className="flex flex-col items-center gap-3 py-4 text-center">
              <CheckCircle2 className="size-10 text-emerald-600" />
              <p className="text-sm text-muted-foreground">
                Your reservation has been confirmed. You can view it in{" "}
                <a href="/reservations" className="underline underline-offset-2">
                  My Reservations
                </a>
                .
              </p>
            </div>
          ) : (
            <>
              {resourceName && (
                <div className="text-sm">
                  <span className="text-muted-foreground">Resource: </span>
                  <span className="font-medium">{resourceName}</span>
                </div>
              )}

              <div className="rounded-lg border bg-muted/30 p-4">
                <div className="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <div className="text-muted-foreground">Date</div>
                    <div className="font-medium">{format(startsAt, "EEEE, MMM d, yyyy")}</div>
                  </div>
                  <div>
                    <div className="text-muted-foreground">Time</div>
                    <div className="font-medium">
                      {format(startsAt, "HH:mm")} — {format(endsAt, "HH:mm")}
                    </div>
                  </div>
                </div>
              </div>

              {modal.state === "error" && (
                <div className="flex items-start gap-2 rounded-lg border border-destructive/30 bg-destructive/5 p-3 text-sm text-destructive">
                  <AlertCircle className="mt-0.5 size-4 shrink-0" />
                  {modal.errorMsg}
                </div>
              )}
            </>
          )}
        </CardContent>

        <CardFooter
          className={cn("gap-2", modal.state === "success" ? "justify-center" : "justify-end")}
        >
          {modal.state === "success" ? (
            <Button onClick={handleClose}>Done</Button>
          ) : (
            <>
              <Button variant="outline" onClick={handleClose} disabled={mutation.isPending}>
                Cancel
              </Button>
              <Button onClick={handleConfirm} disabled={mutation.isPending}>
                {mutation.isPending ? "Booking..." : "Confirm booking"}
              </Button>
            </>
          )}
        </CardFooter>
      </Card>
    </div>
  )
}
