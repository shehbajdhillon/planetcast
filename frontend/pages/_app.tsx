import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { ChakraProvider } from '@chakra-ui/react'
import { mode } from '@chakra-ui/theme-tools'
import { extendTheme } from '@chakra-ui/react';

import { Inter } from 'next/font/google';
import { ClerkProvider, useAuth } from '@clerk/nextjs';
import { PropsWithChildren, useMemo } from 'react';
import { ApolloProvider } from '@apollo/client';
import { GetApolloClient } from '@/apollo-client';

import { Analytics } from '@vercel/analytics/react';

import '../styles/nprogress.css';

import { useEffect } from 'react';
import { useRouter } from 'next/router';
import posthog from 'posthog-js';
import { PostHogProvider } from 'posthog-js/react';

import NProgress from 'nprogress';

import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';

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
  initialColorMode: 'light',
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

if (typeof window !== 'undefined') {
  posthog.init(process.env.NEXT_PUBLIC_POSTHOG_KEY || '', {
    api_host: 'https://app.posthog.com',
    // Enable debug mode in development
    loaded: (posthog) => {
      if (process.env.NODE_ENV === 'development') posthog.debug()
    },
    capture_pageview: true,
  })
}

const ApolloProviderWrapper = ({ children }: PropsWithChildren) => {
  const { getToken } = useAuth();

  const client = useMemo(() => {
    return GetApolloClient(false, getToken)
  }, [getToken]);

  return <ApolloProvider client={client}>{children}</ApolloProvider>;
};

const stripePromise = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || "")

export default function App({ Component, pageProps }: AppProps) {

  const router = useRouter();

  useEffect(() => {
    const handleRouteChange = () => posthog?.capture('$pageview');
    router.events.on('routeChangeComplete', handleRouteChange);
    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    }
  }, [router.events]);


  useEffect(() => {
    const handleRouteStart = () => NProgress.start();
    const handleRouteChange = () => {
      NProgress.done();
    }
    const handleRouteError = () => NProgress.done();

    router.events.on('routeChangeStart', handleRouteStart);
    router.events.on('routeChangeComplete', handleRouteChange);
    router.events.on('routeChangeError', handleRouteError);

    return () => {
      NProgress.done();
      router.events.off('routeChangeStart', handleRouteStart);
      router.events.off('routeChangeComplete', handleRouteChange);
      router.events.off('routeChangeError', handleRouteError);
    };
  }, [router.events]);

  return (
    <ChakraProvider theme={theme}>
      <ClerkProvider>
        <ApolloProviderWrapper>
          <PostHogProvider client={posthog}>
            <Elements stripe={stripePromise}>
              <main className={inter.className}>
                <Component {...pageProps} />
                <Analytics />
              </main>
            </Elements>
          </PostHogProvider>
        </ApolloProviderWrapper>
      </ClerkProvider>
    </ChakraProvider>
  );
}
