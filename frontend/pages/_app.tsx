import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { ChakraProvider } from '@chakra-ui/react'
import { mode } from '@chakra-ui/theme-tools'
import { extendTheme } from '@chakra-ui/react';

import { Inter } from 'next/font/google';
import { ClerkProvider, useAuth } from '@clerk/nextjs';
import { PropsWithChildren, useMemo } from 'react';
import { ApolloClient, ApolloProvider, InMemoryCache, from } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { createUploadLink } from 'apollo-upload-client'

const apiServer =
  process.env.NODE_ENV === 'production'
    ? 'https://api.withsync.ai'
    : 'http://localhost:8080';

const httpLink: any = createUploadLink({ uri: apiServer });

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

const ApolloProviderWrapper = ({ children }: PropsWithChildren) => {
  const { getToken } = useAuth();

  const client = useMemo(() => {
    const authMiddleware = setContext(async (_, { headers }) => {
      const token = await getToken({ template: 'PlanetCast_GQL_Backend' });
      return {
        headers: {
          ...headers,
          authorization: token ? `Bearer ${token}` : '',
        },
      };
    });
    return new ApolloClient({
      link: from([authMiddleware, httpLink]),
      cache: new InMemoryCache(),
    });
  }, [getToken]);

  return <ApolloProvider client={client}>{children}</ApolloProvider>;
};

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ChakraProvider theme={theme}>
      <ClerkProvider>
        <ApolloProviderWrapper>
          <main className={inter.className}>
            <Component {...pageProps} />
          </main>
        </ApolloProviderWrapper>
      </ClerkProvider>
    </ChakraProvider>
  );
}
