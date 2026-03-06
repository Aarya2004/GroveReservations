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

const signupSchema = z.object({
  name: z.string().min(1, "Name is required"),
  email: z.string().email("Invalid email"),
  password: z.string().min(6, "Password must be at least 6 characters"),
})

type SignupValues = z.infer<typeof signupSchema>

type AsyncState = { status: "idle" | "loading" | "success"; error: string | null }
type AsyncAction =
  | { type: "loading" }
  | { type: "success" }
  | { type: "error"; message: string }
  | { type: "reset" }

function asyncReducer(_: AsyncState, action: AsyncAction): AsyncState {
  switch (action.type) {
    case "loading":
      return { status: "loading", error: null }
    case "success":
      return { status: "success", error: null }
    case "error":
      return { status: "idle", error: action.message }
    case "reset":
      return { status: "idle", error: null }
  }
}

type SubmitFn = (e: React.FormEvent<HTMLFormElement>) => Promise<void> | void

interface SignupProps extends React.HTMLAttributes<HTMLDivElement> {
  className?: string
  handleSubmit?: SubmitFn | null
}

export function SignupForm({ className, handleSubmit = null, ...props }: SignupProps) {
  const form = useForm<SignupValues>({ resolver: zodResolver(signupSchema) })
  const [async, dispatch] = useReducer(asyncReducer, { status: "idle", error: null })
  const router = useRouter()

  const onSubmit = async (values: SignupValues) => {
    dispatch({ type: "loading" })

    const { data, error } = await supabase.auth.signUp({
      email: values.email,
      password: values.password,
      options: { data: { name: values.name } },
    })

    if (error) {
      dispatch({ type: "error", message: error.message })
      return
    }

    if (data.user?.identities?.length === 0) {
      dispatch({ type: "error", message: "An account with this email already exists." })
      return
    }

    dispatch({ type: "success" })
  }

  const signUpWithGoogle = async () => {
    dispatch({ type: "loading" })
    const { error } = await supabase.auth.signInWithOAuth({
      provider: "google",
      options: {
        redirectTo: `${window.location.origin}/auth/callback`,
        queryParams: { prompt: "select_account" },
      },
    })
    if (error) dispatch({ type: "error", message: error.message })
  }

  if (async.status === "success") {
    return (
      <div className={cn("flex flex-col gap-6", className)} {...props}>
        <Card>
          <CardHeader>
            <CardTitle>Check your email</CardTitle>
            <CardDescription>
              We sent a confirmation link to <strong>{form.getValues("email")}</strong>. Click the
              link to activate your account.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button className="w-full" onClick={() => router.push("/login")}>
              Back to login
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle>Create an account</CardTitle>
          <CardDescription>Enter your details below to create your account</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit ?? form.handleSubmit(onSubmit)}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="name">Name</Label>
                <Input
                  id="name"
                  type="text"
                  placeholder="Your full name"
                  {...form.register("name")}
                />
                {form.formState.errors.name && (
                  <p className="text-sm text-red-600">{form.formState.errors.name.message}</p>
                )}
              </div>

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
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="At least 6 characters"
                  {...form.register("password")}
                />
                {form.formState.errors.password && (
                  <p className="text-sm text-red-600">{form.formState.errors.password.message}</p>
                )}
              </div>

              {async.error && <p className="text-sm text-red-600">{async.error}</p>}

              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full" disabled={async.status === "loading"}>
                  {async.status === "loading" ? "Creating account..." : "Sign up"}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  className="w-full"
                  onClick={signUpWithGoogle}
                  disabled={async.status === "loading"}
                >
                  {async.status === "loading" ? "Please wait..." : "Sign up with Google"}
                </Button>
              </div>
            </div>

            <div className="mt-4 text-center text-sm">
              Already have an account?{" "}
              <a href="/login" className="underline underline-offset-4">
                Log in
              </a>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
