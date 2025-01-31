import { SignedIn, SignedOut } from "@clerk/nextjs"
import { TokenDisplay } from "@/components/TokenDisplay"
import { Header } from "@/components/Header"
import { WelcomeMessage } from "@/components/WelcomeMessage"

export default function Home() {
  return (
    <div className="max-w-4xl mx-auto p-6">
      <Header />

      <main>
        <SignedOut>
          <WelcomeMessage />
        </SignedOut>

        <SignedIn>
          <TokenDisplay />
        </SignedIn>
      </main>
    </div>
  )
}

