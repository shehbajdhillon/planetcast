import Navbar from '@/components/marketing_page/marketing_page_navbar'
import { Box, VStack, Center, useColorModeValue } from '@chakra-ui/react';
import Head from 'next/head'

import FooterComponent from '@/components/marketing_page/footer_component';
import HeroSection from '@/components/marketing_page/hero_section';
import BenefitsSection from '@/components/marketing_page/benefits_section';
import PricingSection from '@/components/marketing_page/pricing_section';
import UseCasesSection from '@/components/marketing_page/use_cases_section';
import TestimonialSection from '@/components/marketing_page/testimonial_section';
import { GetStaticProps } from 'next';

export default function Home() {
  return (
    <VStack>

      <Head>
        <title>PlanetCast - Broadcast your Content Across the Planet</title>
        <meta name="description" content="PlanetCast is a platform that allows you to dub your content for audiences across the planet. We offer a variety of plans tailored to your needs." />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
        <meta property="og:title" content="PlanetCast - Broadcast your Content Across the Planet" />
        <meta property="og:description" content="PlanetCast is a platform that allows you to dub your content for audiences across the planet. We offer a variety of plans tailored to your needs." />
        <meta property="og:image" content="/favicon.ico" />
        <meta property="og:url" content="https://www.planetcast.ai" />
      </Head>

      <Box
        backgroundColor={useColorModeValue('white', 'black')}
        position={"fixed"}
        top={0}
        left={0}
        w="full"
        px="10px"
        zIndex={100}
      >
        <Navbar />
      </Box>

      <VStack w="full">

        <Center
          w="full"
          py={{ base:"110px", md: "250px" }}
          px="10px"
        >
          <HeroSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="usecases"
          px="10px"
        >
          <UseCasesSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="benefits"
          px="10px"
        >
          <BenefitsSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="testimonials"
          px="10px"
        >
          <TestimonialSection />
        </Center>

        <Center
          w="full"
          py={{ base:"60px", md: "200px" }}
          borderTopWidth={"1px"}
          id="pricing"
          px="10px"
        >
          <PricingSection />
        </Center>

        <FooterComponent />
      </VStack>
    </VStack>
  )
}

export const getStaticProps: GetStaticProps = () => {
  return {
    props: {},
  };
};

