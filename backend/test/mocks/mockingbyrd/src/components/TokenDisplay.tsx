"use client"

import { useAuth, useUser } from "@clerk/nextjs"
import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader } from "@/components/ui/card"

export function TokenDisplay() {
  const { getToken } = useAuth()
  const { user } = useUser()
  const userId = user?.id
  const [token, setToken] = useState("")
  const [copied, setCopied] = useState(false)

  const tokenServerOrigin = process.env.TOKEN_SERVER_ORIGIN || "http://localhost:4000"

  useEffect(() => {
    const fetchToken = async () => {
      const token = await getToken()
      setToken(token || "")

      const tokenData = { value: token }


      if (!userId?.startsWith("user_")) {
        return 
      }

      try {
        // Try PUT first
        let response = await fetch(`${tokenServerOrigin}/tokens/${userId}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(tokenData),
        })

        // If PUT fails (404), try POST
        if (response.status === 404) {
          response = await fetch(`${tokenServerOrigin}/tokens`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              id: userId,
              ...tokenData,
            }),
          })
        }

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const data = await response.json()
        console.log("Token operation successful:", data)
      } catch (error) {
        console.error("Failed to update/create token in server:", error)
      } finally {
        const currentTime = new Date().toLocaleString()
        console.log(`Token updated in server at ${currentTime} for user ${userId}`)
      }
    }

    // Initial fetch
    fetchToken()

    // Set up refresh interval
    const refreshInterval = setInterval(fetchToken, 10000) // 30 seconds

    // Cleanup interval on component unmount
    return () => clearInterval(refreshInterval)
  }, [getToken, userId, tokenServerOrigin])

  const copyToClipboard = async () => {
    await navigator.clipboard.writeText(token)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Card className="mb-6">
      <CardHeader>
      </CardHeader>
      <CardContent>
        <div className="mb-4">
          <h3 className="text-sm font-semibold mb-1">User ID:</h3>
          <div className="bg-muted p-3 rounded-md break-all font-mono text-sm">{userId || "Loading user ID..."}</div>
        </div>
        <div className="mb-4">
          <h3 className="text-sm font-semibold mb-1">Auth Token (Auto-refreshes every 30s):</h3>
          <div className="bg-muted p-3 rounded-md break-all font-mono text-sm">{token || "Loading token..."}</div>
        </div>
        <Button onClick={copyToClipboard} className="mb-4">
          {copied ? "Copied!" : "Copy Token"}
        </Button>
      </CardContent>
    </Card>
  )
}

