import Navbar from '@/components/marketing_page_navbar'
import {
  Box,
  HStack,
  Stack,
  Heading,
  Button,
  Text,
  useColorModeValue,
  VStack,
  useBreakpointValue,
  Grid,
  GridItem,
} from '@chakra-ui/react';
import Head from 'next/head'

import Image from 'next/image';
import Link from 'next/link';
import {
  ArrowUpFromDot,
  DollarSign,
  GlobeIcon,
  TrendingDownIcon,
} from 'lucide-react';
import VideoPlayer from '@/components/video_player';
import { useState } from 'react';
import Marquee from '@/components/Marquee';

const HeroSection: React.FC = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      direction={{ base: "column-reverse", md: "row" }}
      maxW={"1400px"}
      mt={{ base:"100px", md: "250px" }}
      w="full"
    >
      <Box
        mb={{ base: "auto", md: "0px" }}
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
          <Link
            href={'/dashboard'}
          >
            <Button
              size={"lg"}
              backgroundColor={useColorModeValue("black", "white")}
              textColor={useColorModeValue("white", "black")}
              borderColor={useColorModeValue("white", "black")}
              borderWidth={"1px"}
              _hover={{
                backgroundColor: useColorModeValue("white", "black"),
                textColor: useColorModeValue("white", "black"),
                bgGradient: 'linear(to-tl, #007CF0, #01DFD8)'
              }}
            >
              Start for Free
            </Button>
          </Link>
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
      <Box maxW={{ base: "200px", md: "25%" }} mt={{ base:"auto", md: "0px" }}>
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

  const iconSize = useBreakpointValue({ base: '40px', md: '60px' })

  const buttonBg = useColorModeValue("black", "white");
  const buttonColor = useColorModeValue("white", "black");
  const [tfnIdx, setTfnIdx] = useState(0);

  const transformations = [
    {
      language: "ENGLISH",
      link: "",
    },
    {
      language: "SPANISH",
      link: "",
    },
    {
      language: "HINDI",
      link: "",
    },
    {
      language: "FRENCH",
      link: "",
    },
  ]

  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      mt={{ base:"110px", md: "250px" }}
      w={"full"}
    >
      <Grid
        templateAreas={{
          base: `
            "info"
            "video"
          `,
          lg: `"info video"`
        }}
        gridTemplateColumns={{ base: "1fr", lg: "3fr 2fr" }}
        h="full"
        gap={{ base: "15px", lg: "50px" }}
        w={{ lg: "full" }}
      >
        <GridItem
          area={"info"}
          placeItems={"center"}
          display={"grid"}
        >
          <Box
            mb={{ base: "auto", md: "0px" }}
            w="full"
            alignItems={{ base: "center", md: "left" }}
            justifyContent={{ base: "center", md: "left" }}
            display={"flex"}
            flexDir={"column"}
          >
            <Heading
              size={{ base: '2xl', md: '3xl' }}
              fontWeight={'semibold'}
              textAlign={{ base: "center", md: "left" }}
              w={{ md: "full" }}
            >
              <HStack>
                <Text>10x your {' '}
                  <Text
                    as={"span"}
                    bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
                    bgClip='text'
                  >
                    reach
                  </Text>
                </Text>
                <ArrowUpFromDot size={iconSize} />
                <GlobeIcon size={iconSize} />
              </HStack>
              <HStack mt={{ md: "8px" }}>
                <Text>1/10th the {' '}
                  <Text
                    as={"span"}
                    bgGradient={'linear(to-tr, #01CF00, #90DD00)'}
                    bgClip='text'
                  >
                    cost
                  </Text>
                </Text>
                <TrendingDownIcon size={iconSize} />
                <DollarSign size={iconSize} />
              </HStack>
            </Heading>
            <Heading
              w={{ md: "full" }}
              fontWeight={'semibold'}
              textAlign={{ base: "center", md: "left" }}
              size={{ base: "sm", sm: "lg" }}
              mt={{ md: "10px" }}
            >
              Engage listeners from every corner of the globe
            </Heading>
            <Heading
              w={{ md: "full" }}
              textAlign={{ base: "center", md: "left" }}
              fontWeight={'semibold'}
              size={{ base: "sm", sm: "lg" }}
            >
              Save time and money over traditional dubbing
            </Heading>
            <Heading
              w={{ md: "full" }}
              textAlign={{ base: "center", md: "left" }}
              fontWeight={'semibold'}
              size={{ base: "sm", sm: "lg" }}
            >
              Replicate authentic voices in every translation
            </Heading>
          </Box>
        </GridItem>
        <GridItem area={"video"} display={"grid"} placeItems={"center"}>
          <Box display={"flex"} h="full" w="full" px={{ base: "16px", sm: "0px" }}>
            <VideoPlayer src={transformations[tfnIdx].link} />
          </Box>
          <HStack pt="10px" w="full">
            {transformations.map((tfn, idx) => (
              <Button
                key={idx}
                onClick={() => setTfnIdx(idx)}
                variant={idx == tfnIdx ? "solid" : "outline"}
                pointerEvents={idx === tfnIdx ? "none" : "auto"}
                background={idx === tfnIdx ? buttonBg : '' }
                color={idx === tfnIdx ? buttonColor : '' }
              >
                {tfn.language}
              </Button>
            ))}
          </HStack>
        </GridItem>
      </Grid>
    </Stack>
  );
};

export default function Home() {

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
      <VStack w="full">
        <HeroSection />
        <BenefitsSection />
        <Heading
          mt={{ base:"110px", md: "250px" }}
          size={{ base: 'xl', md: "2xl" }}
          textAlign={{ base: "center" }}
          fontWeight={"semibold"}
          w={{ md: "full" }}
          backgroundColor={"red"}
          py="30px"
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          textColor={bgColor}
        >
          <Marquee>
            {new Array(8).fill(0).map((_, idx) => (
              <HStack px="100px" spacing={5} key={idx}>
                <Text h="60px">
                  Welcome to efficient broadcasting
                </Text>
              </HStack>
            ))}
          </Marquee>
        </Heading>
      </VStack>
    </VStack>
  )
}

