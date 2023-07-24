import Navbar from '@/components/marketing_page_navbar'
import {
  Box,
  HStack,
  Stack,
  Heading,
  Button,
  Text,
  useColorModeValue,
  Center,
  VStack,
} from '@chakra-ui/react';
import Head from 'next/head'

import useWindowDimensions from '@/hooks/useWindowDimensions';
import Image from 'next/image';
import Link from 'next/link';

const HeroSection: React.FC = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      direction={{ base: "column-reverse", md: "row" }}
      maxW={"1920px"}
      mx={"10px"}
    >
      <Box
        w="full"
        maxW={{ md: "75%" }}
        alignItems={{ base: "center", md: "left" }}
        justifyContent={{ base: "center", md: "left" }}
        display={"flex"}
        flexDir={"column"}
      >
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"semibold"}
          w={{ md: "full" }}
        >
          Dub
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"semibold"}
          w={{ md: "full" }}
        >
          Translate
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"semibold"}
          w={{ md: "full" }}
        >
          Broadcast
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"semibold"}
          w={{ md: "full" }}
        >
          Content Across the {' '}
          <Text
            as={"span"}
            bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
            bgClip='text'
          >
            Planet
          </Text>
        </Heading>
        <HStack w={{ md: "full" }} pt="10px">
          {/**
          <Link
            href={'/dashboard'}
          >
            <Button
              size={"lg"}
              backgroundColor={useColorModeValue("black", "white")}
              textColor={useColorModeValue("white", "black")}
              borderColor={useColorModeValue("black", "white")}
              borderWidth={"1px"}
              _hover={{
                backgroundColor: useColorModeValue("white", "black"),
                textColor: useColorModeValue("black", "white")
              }}
            >
              Try for Free
            </Button>
          </Link>
          **/}
          <Button
            size={"lg"}
            variant={"outline"}
            backgroundColor={useColorModeValue("black", "white")}
            textColor={useColorModeValue("white", "black")}
            _hover={{
              backgroundColor: useColorModeValue("black", "white:"),
              textColor: useColorModeValue("white", "black"),
              bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
            }}
          >
            Request Access
          </Button>
        </HStack>
      </Box>
      <Box maxW={{ base: "200px", md: "25%" }}>
        <Image
          height={1000}
          width={400}
          src={useColorModeValue('/planetcastgradientlight.svg', '/planetcastgradientdark.svg')}
          alt='planet cast gradient logo'
        />
      </Box>
    </Stack>
  );
};

const BenefitsSection: React.FC = () => {
  return (
    <Stack
      w="full"
      direction={{ base: "column", md: "row" }}
      h="full"
      spacing={0}
    >
      <VStack w="full" h="full">
        <Center h="full">

          <Heading>Reach 100m+ People</Heading>

        </Center>
      </VStack>

      <VStack w="full" h="full">
        <Center h="full">
          Reach 10x platforms easily
        </Center>
      </VStack>

      <VStack w="full" h="full">
        <Center h="full">
          Reach 10x platforms easily
        </Center>
      </VStack>

    </Stack>
  );
};

export default function Home() {

  const { height } = useWindowDimensions();

  const bgColor = useColorModeValue("white", "black");

  return (
    <VStack>

      <Head>
        <title>PlanetCast</title>
        <meta name="description" content="Cast your Content Across the Planet" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Box position={"fixed"} top={0} left={0} w="full" p="10px" backgroundColor={bgColor} zIndex={100}>
        <Navbar marketing />
      </Box>

      <VStack
        w="full"
        height={height}
      >

        <Center
          h="full"
          maxW={"1920px"}
        >
          <HeroSection />
        </Center>
      </VStack>

    </VStack>
  )
}

