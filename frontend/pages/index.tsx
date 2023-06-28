import Navbar from '@/components/marketing_page_navbar'
import {
  Box,
  HStack,
  Stack,
  Heading,
  Button,
  Text,
  useColorModeValue,
} from '@chakra-ui/react';
import Head from 'next/head'

import useWindowDimensions from '@/hooks/useWindowDimensions';
import Image from 'next/image';

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
          <Button
            size={"lg"}
            variant={"outline"}
            _hover={{
              backgroundColor: useColorModeValue("black", "white:"),
              textColor: useColorModeValue("white", "black"),
              bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
            }}
          >
            Read More
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

export default function Home() {

  const { height } = useWindowDimensions();

  return (
    <>
      <Head>
        <title>PlanetCast</title>
        <meta name="description" content="Cast your Content Across the Planet" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box position={"fixed"} top={0} left={0} w="full" p="10px">
        <Navbar marketing />
      </Box>
      <Box
        display={"flex"}
        justifyContent={"center"}
        flexDir={"column"}
        height={height}
      >
        <Box
          h="full"
          display={'grid'}
          placeItems={"center"}
        >
          <HeroSection />
        </Box>
      </Box>
    </>
  )
}

