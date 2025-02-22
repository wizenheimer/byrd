import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams;
    const code = searchParams.get('code');
    const state = searchParams.get('state');


    if (!code || !state) {
      return NextResponse.redirect(
        new URL('/error?message=missing_parameters', request.url)
      );
    }

    if (!process.env.BACKEND_ORIGIN) {
      return NextResponse.redirect(
        new URL('/error?message=configuration_error', request.url)
      );
    }

    // Construct backend URL with query parameters
    const backendURL = new URL(`${process.env.BACKEND_ORIGIN}/api/public/v1/integration/slack/oauth/callback`);
    backendURL.searchParams.set('code', code);
    backendURL.searchParams.set('state', state);

    // Forward to Go backend
    const response = await fetch(backendURL, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Backend error:', errorText);
      throw new Error(`Backend error: ${response.status}`);
    }

    const data = await response.json();

    if (!data.deep_link) {
      console.error('No deep_link in response:', data);
      throw new Error('Invalid backend response');
    }


    // Create redirect response with appropriate headers
    const redirectResponse = NextResponse.redirect(data.deep_link);

    // Add any necessary headers
    redirectResponse.headers.set('Cache-Control', 'no-store');

    return redirectResponse;

  } catch (error) {
    console.error('Failed to handle OAuth callback:', error);

    // Construct error URL with details
    const errorURL = new URL('/error', request.url);
    errorURL.searchParams.set('message', 'oauth_failed');
    errorURL.searchParams.set('details', error instanceof Error ? error.message : 'Unknown error');

    return NextResponse.redirect(errorURL);
  }
}
