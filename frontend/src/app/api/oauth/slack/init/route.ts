// frontend/src/app/api/oauth/slack/init/route.ts
import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const data = await request.json();

    // Your backend should construct the proper Slack OAuth URL with all parameters
    const response = await fetch(
      `${process.env.BACKEND_ORIGIN}/api/public/v1/integration/slack/oauth/init`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      }
    );

    if (!response.ok) {
      throw new Error('Failed to initialize OAuth');
    }

    // Return the OAuth URL as JSON instead of redirecting
    const { oauth_url } = await response.json();
    return NextResponse.json({ oauth_url });

  } catch (error) {
    console.error('Failed to initiate Slack OAuth:', error);
    return NextResponse.json(
      { error: 'Failed to initialize OAuth' },
      { status: 500 }
    );
  }
}
