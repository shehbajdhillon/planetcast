import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { ChakraProvider } from '@chakra-ui/react'
import { mode } from '@chakra-ui/theme-tools'
import { extendTheme } from '@chakra-ui/react';

import { Inter } from 'next/font/google';

const inter = Inter({ subsets: ['latin'] });

const theme = extendTheme({
  fonts: {
    heading: inter.style.fontFamily,
    body: inter.style.fontFamily,
  },
  styles: {
    global: (props: any) => ({
      body: {
        bg: mode('white', 'black')(props),
        textColor: mode('black', 'white')(props),
      },
    }),
  },
});

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ChakraProvider theme={theme}>
      <main className={inter.className}>
        <Component {...pageProps} />
      </main>
    </ChakraProvider>
  );
}
