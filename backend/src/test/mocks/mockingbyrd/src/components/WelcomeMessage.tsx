import { SignInButton } from "@clerk/nextjs"
import { Button } from "@/components/ui/button"

export function WelcomeMessage() {
  return (
    <div className="text-center my-10">
      <h2 className="text-2xl font-semibold mb-4">Welcome to the Auth Demo</h2>
      <p className="mb-4">Sign in to view your authentication token</p>
      <SignInButton mode="modal">
        <Button>Sign In</Button>
      </SignInButton>
    </div>
  )
}

