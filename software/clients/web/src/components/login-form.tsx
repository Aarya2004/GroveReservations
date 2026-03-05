"use client"

import { useReducer } from "react"
import { useRouter } from "next/navigation"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { supabase } from "@/lib/supabase"

const loginSchema = z.object({
  email: z.string().email("Invalid email"),
  password: z.string().min(1, "Password is required"),
})

type LoginValues = z.infer<typeof loginSchema>

type AsyncState = { status: "idle" | "loading"; error: string | null }
type AsyncAction =
  | { type: "loading" }
  | { type: "error"; message: string }
  | { type: "reset" }

function asyncReducer(_: AsyncState, action: AsyncAction): AsyncState {
  switch (action.type) {
    case "loading": return { status: "loading", error: null }
    case "error":   return { status: "idle", error: action.message }
    case "reset":   return { status: "idle", error: null }
  }
}

type SubmitFn = (e: React.FormEvent<HTMLFormElement>) => Promise<void> | void

interface LoginProps extends React.HTMLAttributes<HTMLDivElement> {
  className?: string
  handleSubmit?: SubmitFn | null
}

export function LoginForm({ className, handleSubmit = null, ...props }: LoginProps) {
  const form = useForm<LoginValues>({ resolver: zodResolver(loginSchema) })
  const [async, dispatch] = useReducer(asyncReducer, { status: "idle", error: null })
  const router = useRouter()

  const onSubmit = async (values: LoginValues) => {
    dispatch({ type: "loading" })

    const { data, error } = await supabase.auth.signInWithPassword({
      email: values.email,
      password: values.password,
    })

    if (error) {
      dispatch({ type: "error", message: error.message })
      return
    }

    console.log("Signed in:", data)
    router.push("/resources")
  }

  const signInWithGoogle = async () => {
    dispatch({ type: "loading" })
    const { error } = await supabase.auth.signInWithOAuth({
      provider: "google",
      options: {
        redirectTo: `${window.location.origin}/auth/callback`,
        queryParams: { prompt: "select_account" }
      }
    })
    if (error) dispatch({ type: "error", message: error.message })
  }

  const sendResetLink = async () => {
    const email = form.getValues("email")
    if (!email) {
      dispatch({ type: "error", message: "Enter your email to receive a reset link." })
      return
    }
    dispatch({ type: "loading" })
    const { error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo: `${window.location.origin}/auth/update-password`
    })
    if (error) dispatch({ type: "error", message: error.message })
    else dispatch({ type: "reset" })
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle>Login to your account</CardTitle>
          <CardDescription>Enter your email below to login to your account</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit ?? form.handleSubmit(onSubmit)}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  {...form.register("email")}
                />
                {form.formState.errors.email && (
                  <p className="text-sm text-red-600">{form.formState.errors.email.message}</p>
                )}
              </div>

              <div className="grid gap-3">
                <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                  <button
                    type="button"
                    onClick={sendResetLink}
                    className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                  >
                    Forgot your password?
                  </button>
                </div>
                <Input
                  id="password"
                  type="password"
                  {...form.register("password")}
                />
                {form.formState.errors.password && (
                  <p className="text-sm text-red-600">{form.formState.errors.password.message}</p>
                )}
              </div>

              {async.error && (
                <p className="text-sm text-red-600">
                  {async.error}
                </p>
              )}

              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full" disabled={async.status === "loading"}>
                  {async.status === "loading" ? "Logging in..." : "Login"}
                </Button>
                <Button type="button" variant="outline" className="w-full" onClick={signInWithGoogle} disabled={async.status === "loading"}>
                  {async.status === "loading" ? "Please wait..." : "Login with Google"}
                </Button>
              </div>
            </div>

            <div className="mt-4 text-center text-sm">
              Don&apos;t have an account?{" "}
              <a href="/signup" className="underline underline-offset-4">
                Sign up
              </a>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
