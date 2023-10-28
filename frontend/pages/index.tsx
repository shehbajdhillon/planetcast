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
  Avatar,
  Center,
} from '@chakra-ui/react';
import Head from 'next/head'

import Image from 'next/image';
import Link from 'next/link';
import {
  ArrowUpFromDot,
  DollarSign,
  ExternalLink,
  GlobeIcon,
  TrendingDownIcon,
} from 'lucide-react';
import VideoPlayer from '@/components/video_player';
import { useState } from 'react';
import PricingComponent from '@/components/pricing_component';
import FooterComponent from '@/components/footer_component';

const HeroSection: React.FC = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      direction={{ base: "column-reverse", md: "row" }}
      maxW={"1400px"}
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
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Dub
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Translate
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
          w={{ md: "full" }}
        >
          Broadcast
        </Heading>
        <Heading
          size={{ base: '3xl', md: "4xl" }}
          textAlign={{ base: "center", md: "left" }}
          fontWeight={"medium"}
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
          <Link href={'/dashboard'}>
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
          <Link href={"#benefits"}>
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
          </Link>
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
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/english.mp4",
    },
    {
      language: "SPANISH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/spanish.mp4",
    },
    {
      language: "HINDI",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/hindi.mp4",
    },
    {
      language: "FRENCH",
      link: "https://planetcastpublic.s3.us-west-1.amazonaws.com/french.mp4",
    },
  ]

  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
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
              fontWeight={'medium'}
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
              fontWeight={'normal'}
              textAlign={{ base: "center", md: "left" }}
              size={{ base: "sm", sm: "lg" }}
              mt={{ md: "10px" }}
            >
              Engage listeners from every corner of the globe
            </Heading>
            <Heading
              w={{ md: "full" }}
              textAlign={{ base: "center", md: "left" }}
              fontWeight={'normal'}
              size={{ base: "sm", sm: "lg" }}
            >
              Save time and money over traditional dubbing
            </Heading>
            <Heading
              w={{ md: "full" }}
              textAlign={{ base: "center", md: "left" }}
              fontWeight={'normal'}
              size={{ base: "sm", sm: "lg" }}
            >
              Preserve original voices in every translation
            </Heading>
          </Box>
        </GridItem>
        <GridItem area={"video"} display={"grid"} placeItems={"center"}>
          <Box display={"flex"} h="full" w="full" px={{ base: "16px", sm: "0px" }} rounded={"sm"}>
            <VideoPlayer src={transformations[tfnIdx].link} />
          </Box>
          <HStack pt="10px" w="full">
            {transformations.map((tfn, idx) => (
              <Button
                key={idx}
                onClick={() => setTfnIdx(idx)}
                variant={idx == tfnIdx ? "solid" : "outline"}
                pointerEvents={idx === tfnIdx ? "none" : "auto"}
                background={idx === tfnIdx ? buttonBg : buttonColor }
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

interface TestimonialCardProps {
  name: string;
  title: string;
  src: string;
  text: string[];
  link: string;
};

const TestimonialCard: React.FC<TestimonialCardProps> = (props) => {
  const { name, title, src, text, link } = props;

  return (
    <Box
      w="full"
      borderWidth={"1px"}
      rounded={"lg"}
      p="10px"
      shadow="lg"
      borderRadius={"md"}
    >
      <HStack mb="30px">
        <Avatar src={src} />
        <Stack spacing={1}>
          <HStack>
            <Heading
              size={{ base: "md" }}
              fontWeight={'medium'}
            >
              { name }
            </Heading>
            <Link href={link} target='_blank'>
              <ExternalLink />
            </Link>
          </HStack>
          <Heading
            size={{ base:"sm" }}
            fontWeight={'normal'}
          >
            { title }
          </Heading>
        </Stack>
      </HStack>
      <Heading
        fontSize={"md"}
        fontWeight={"normal"}
        as={"em"}
      >
        {text.map((txt, idx) => (
          <Text pt="10px" key={idx}>{'"'}{txt}{'"'}</Text>
        ))}
      </Heading>
    </Box>
  );
};

const TestimonialSection = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >
      <Heading
        fontWeight={"medium"}
        size={{ base: '2xl', md: "3xl" }}
        textAlign={{ base: "center", md: "left" }}
        w={{ md: "full" }}
      >
        Welcome to efficient {' '}
        <Text
          as={"span"}
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          bgClip='text'
        >
          broadcasting
        </Text>
      </Heading>
      <Heading
        fontWeight={'normal'}
        size={{ base: "sm", sm: "lg" }}
        textAlign={{ base: "center", md: "left" }}
        w={{ md: "full" }}
      >
        Discover how our users are revolutionizing their content reach
      </Heading>

      <Stack
        direction={{ base: "column", md: "row" }}
        justifyContent={{ md: "space-between" }}
        w="full"
        mt="45px"
        px="16px"
        spacing={"80px"}
      >
        <TestimonialCard
          name='Rahul Pandey'
          title='Educator'
          src={'/rahulimg.jpeg'}
          text={['Amazed by this 🤯😮', 'The next few years for creators are going to be wild.']}
          link='https://www.youtube.com/@RahulInHindi'
        />
        <TestimonialCard
          name='Devin Estopinal'
          title='Social Media Content Strategist'
          src={'/devnimg.jpeg'}
          text={["I can't wait to see the results", "Just got done downloading the video after dubbing in spanish and it is absolutely flawless"]}
          link='https://x.com/NotDevn'
        />
        <TestimonialCard
          name='Dallon Asnes'
          title='Travel Content Creator'
          src={'/dallonimg.jpeg'}
          text={['first vid in progress', 'the hindi dubbing is excellent']}
          link='https://www.youtube.com/@dallonearth'
        />
      </Stack>
    </Stack>
  );
};


const PricingSection = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >
      <Heading
        fontWeight={"medium"}
        size={{ base: '2xl', md: "3xl" }}
        textAlign={{ base: "center", md: "left" }}
        w={{ md: "full" }}
      >
        Start dubbing {' '}
        <Text
          as={"span"}
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          bgClip='text'
        >
          today
        </Text>
      </Heading>
      <Heading
        fontWeight={'normal'}
        size={{ base: "sm", sm: "lg" }}
        textAlign={{ base: "center", md: "left" }}
        w={{ md: "full" }}
        mb="45px"
      >
        Select the perfect plan tailored to your needs
      </Heading>
      <PricingComponent />
    </Stack>
  );
};

export default function Home() {

  const bgColor = useColorModeValue("white", "black");
  const alternateBgColor = useColorModeValue("blackAlpha.50", "whiteAlpha.50");
  const bgGradient = useColorModeValue('radial(blackAlpha.200, white, white)', 'radial(whiteAlpha.200, black, black)')

  return (
    <VStack>
      <Head>
        <title>PlanetCast</title>
        <meta name="description" content="Cast your Content Across the Planet" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box position={"fixed"} top={0} left={0} w="full" px="10px" backgroundColor={bgColor} zIndex={100}>
        <Navbar marketing />
      </Box>
      <VStack w="full">

        <Center
          w="full"
          py={{ base:"110px", md: "250px" }}
        >
          <HeroSection />
        </Center>

        <Center
          w="full"
          py={{ base:"110px", md: "250px" }}
          borderTopWidth={"1px"}
          id="benefits"
        >
          <BenefitsSection />
        </Center>

        <Center
          w="full"
          py={{ base:"110px", md: "250px" }}
          borderTopWidth={"1px"}
          id="pricing"
        >
          <PricingSection />
        </Center>

        <FooterComponent />

      </VStack>
    </VStack>
  )
}

