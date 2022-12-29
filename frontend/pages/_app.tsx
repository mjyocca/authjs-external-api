import '../styles/globals.css'
import { SessionProvider } from "next-auth/react"
import { SWRConfig } from 'swr';
import type { AppProps } from 'next/app'
import fetcher from '../lib/fetch'

export default function App({ Component, pageProps: { session, ...pageProps } }: AppProps) {
  return (
    <SessionProvider session={session}>
      <SWRConfig value={{ fetcher }}>
        <Component {...pageProps} />
      </SWRConfig>
    </SessionProvider>
  )
}
