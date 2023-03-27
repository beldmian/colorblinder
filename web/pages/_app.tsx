import { AppWrapper } from '@/context/state'
import '@/styles/globals.css'
import '@/styles/player.css'
import type { AppProps } from 'next/app'

export default function App({ Component, pageProps }: AppProps) {
  return (
    <AppWrapper>
      <Component {...pageProps} />
    </AppWrapper>
  )
}
