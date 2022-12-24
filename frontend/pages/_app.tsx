import '../styles/globals.css'
import { SessionProvider } from "next-auth/react"
import { SWRConfig } from 'swr';
import type { AppProps } from 'next/app'

async function fetcher<JSON = any>(input: RequestInfo, init?: RequestInit): Promise<JSON> {
  const res = await fetch(input, init);
  return res.json();
}
export default function App({ Component, pageProps: { session, ...pageProps } }: AppProps) {
  return (
    <SessionProvider session={session}>
      <SWRConfig value={{ fetcher }}>
        <Component {...pageProps} />
      </SWRConfig>
    </SessionProvider>
  )
}
