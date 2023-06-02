import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { ChakraProvider } from '@chakra-ui/react'
import { mode } from '@chakra-ui/theme-tools'
import { extendTheme } from '@chakra-ui/react';

import { Inter } from 'next/font/google';
import { ClerkProvider } from '@clerk/nextjs';

const inter = Inter({ subsets: ['latin'] });

const components = {
  Modal: {
    // setup light/dark mode component defaults
    baseStyle: (props: any) => ({
      dialog: {
        bg: mode('white', 'black')(props),
        padding: 2,
        borderColor: mode("blackAlpha.300", "whiteAlpha.300")(props),
        borderWidth: '1px'
      },
    }),
  },

  Drawer: {
    // setup light/dark mode component defaults
    baseStyle: (props: any) => ({
      dialog: {
        bg: mode('white', 'black')(props),
      },
    }),
  },
};

const config = {
  initialColorMode: 'dark',
};

const theme = extendTheme({
  config,
  components,
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
      <ClerkProvider>
        <main className={inter.className}>
          <Component {...pageProps} />
        </main>
      </ClerkProvider>
    </ChakraProvider>
  );
}
