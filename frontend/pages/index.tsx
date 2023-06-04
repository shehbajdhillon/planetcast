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

const HeroSection: React.FC = () => {
  return (
    <Stack
      spacing={"25px"}
      display={"flex"}
      alignItems={{ base: "center" }}
    >
      <Box>
        <Heading
          size={"4xl"}
          textAlign={"center"}
          fontWeight={"semibold"}
        >
          Cast Content Across the {' '}
          <Text
            as={"span"}
            bgGradient={'linear(to-tl, orange, red)'}
            bgClip='text'
          >
            Planet
          </Text>
        </Heading>
      </Box>
      <HStack>
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
            bgGradient: 'linear(to-tl, orange, red)',
          }}
        >
          Read More
        </Button>
      </HStack>
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
        p="10px"
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

